package webapp

import (
	"net/http"

	"github.com/authgear/authgear-server/pkg/auth/dependency/newinteraction"
	"github.com/authgear/authgear-server/pkg/auth/dependency/webapp"
)

// nolint:golint
type WebAppService interface {
	GetIntent(webappIntent *webapp.Intent, stateID string) (*webapp.State, *newinteraction.Graph, []newinteraction.Edge, error)
	Get(stateID string) (*webapp.State, *newinteraction.Graph, []newinteraction.Edge, error)
	PostIntent(webappIntent *webapp.Intent, inputer func() (interface{}, error)) (*webapp.Result, error)
	PostInput(stateID string, inputer func() (interface{}, error)) (*webapp.Result, error)
}

func StateID(r *http.Request) string {
	return r.Form.Get("x_sid")
}