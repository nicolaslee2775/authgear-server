package newinteraction

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/authgear/authgear-server/pkg/auth/dependency/authenticator"
	"github.com/authgear/authgear-server/pkg/auth/dependency/identity"
)

const GraphLifetime = 5 * time.Minute

var ErrInputRequired = errors.New("new input is required")

type Graph struct {
	// GraphID is the unique ID for a graph.
	// It is a constant value through out a graph.
	// It is used to keep track of which instances belong to a particular graph.
	// When one graph is committed, any other instances sharing the same GraphID become invalid.
	GraphID string

	// InstanceID is a unique ID for a particular instance of a graph.
	InstanceID string

	// Intent is the intent (i.e. flow type) of the graph
	Intent Intent

	// Nodes are nodes in a specific path from intent of the interaction graph.
	Nodes []Node
}

func newGraph(intent Intent) *Graph {
	return &Graph{
		GraphID:    "",
		InstanceID: "",
		Intent:     intent,
		Nodes:      nil,
	}
}

func (g *Graph) CurrentNode() Node {
	return g.Nodes[len(g.Nodes)-1]
}

func (g *Graph) appendingNode(n Node) *Graph {
	nodes := make([]Node, len(g.Nodes)+1)
	copy(nodes, g.Nodes)
	nodes[len(nodes)-1] = n

	return &Graph{
		GraphID:    g.GraphID,
		InstanceID: "",
		Intent:     g.Intent,
		Nodes:      nodes,
	}
}

func (g *Graph) MarshalJSON() ([]byte, error) {
	var err error

	intent := ifaceJSON{Kind: IntentKind(g.Intent)}
	if intent.Data, err = json.Marshal(g.Intent); err != nil {
		return nil, err
	}

	nodes := make([]ifaceJSON, len(g.Nodes))
	for i, node := range g.Nodes {
		nodes[i].Kind = NodeKind(node)
		if nodes[i].Data, err = json.Marshal(node); err != nil {
			return nil, err
		}
	}

	graph := &graphJSON{
		GraphID:    g.GraphID,
		InstanceID: g.InstanceID,
		Intent:     intent,
		Nodes:      nodes,
	}
	return json.Marshal(graph)
}

func (g *Graph) UnmarshalJSON(d []byte) error {
	graph := &graphJSON{}
	if err := json.Unmarshal(d, graph); err != nil {
		return err
	}

	intent := InstantiateIntent(graph.Intent.Kind)
	if err := json.Unmarshal(graph.Intent.Data, intent); err != nil {
		return err
	}

	nodes := make([]Node, len(graph.Nodes))
	for i, node := range graph.Nodes {
		nodes[i] = InstantiateNode(node.Kind)
		if err := json.Unmarshal(node.Data, nodes[i]); err != nil {
			return err
		}
	}

	g.GraphID = graph.GraphID
	g.InstanceID = graph.InstanceID
	g.Intent = intent
	g.Nodes = nodes
	return nil
}

func (g *Graph) MustGetUserID() string {
	for i := len(g.Nodes) - 1; i >= 0; i-- {
		if n, ok := g.Nodes[i].(interface{ UserID() string }); ok {
			return n.UserID()
		}
	}
	panic("interaction: expect user ID presents")
}

func (g *Graph) MustGetUserIdentity() *identity.Info {
	for i := len(g.Nodes) - 1; i >= 0; i-- {
		if n, ok := g.Nodes[i].(interface{ UserIdentity() *identity.Info }); ok {
			return n.UserIdentity()
		}
	}
	panic("interaction: expect user identity presents")
}

func (g *Graph) GetAuthenticator(stage AuthenticationStage) (*authenticator.Info, bool) {
	for i := len(g.Nodes) - 1; i >= 0; i-- {
		if n, ok := g.Nodes[i].(interface {
			UserAuthenticator() (AuthenticationStage, *authenticator.Info)
		}); ok {
			s, authenticator := n.UserAuthenticator()
			if s == stage {
				return authenticator, true
			}
		}
	}
	return nil, false
}

// Apply applies the effect the the graph nodes into the context.
func (g *Graph) Apply(ctx *Context) error {
	for _, node := range g.Nodes {
		if err := node.Apply(ctx.perform, g); err != nil {
			return err
		}
	}
	return nil
}

// Accept run the graph to the deepest node using the input
func (g *Graph) Accept(ctx *Context, input interface{}) (*Graph, []Edge, error) {
	graph := g
	for {
		node := graph.CurrentNode()
		edges, err := node.DeriveEdges(ctx, graph)
		if err != nil {
			return nil, nil, err
		}

		if len(edges) == 0 {
			// No more edges, reached the end of the graph
			return graph, edges, nil
		}

		var nextNode Node
		for _, edge := range edges {
			nextNode, err = edge.Instantiate(ctx, graph, input)
			if errors.Is(err, ErrIncompatibleInput) {
				// Continue to check next edges
				continue
			} else if errors.Is(err, ErrSameNode) {
				// The next node is the same current node,
				// so no need to update the graph.
				// Continuing would keep traversing the same edge,
				// so stop and request new input.
				return graph, edges, ErrInputRequired
			} else if err != nil {
				return nil, nil, err
			}
			break
		}

		// No edges are followed, input is required
		if nextNode == nil {
			return graph, edges, ErrInputRequired
		}

		// Follow the edge to nextNode
		graph = graph.appendingNode(nextNode)
		err = nextNode.Apply(ctx.perform, graph)
		if err != nil {
			return nil, nil, err
		}
	}
}

type ifaceJSON struct {
	Kind string          `json:"kind"`
	Data json.RawMessage `json:"data"`
}

type graphJSON struct {
	GraphID    string      `json:"graph_id"`
	InstanceID string      `json:"instance_id"`
	Intent     ifaceJSON   `json:"intent"`
	Nodes      []ifaceJSON `json:"nodes"`
}