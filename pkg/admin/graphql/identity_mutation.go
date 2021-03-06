package graphql

import (
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/relay"

	"github.com/authgear/authgear-server/pkg/api/apierrors"
	"github.com/authgear/authgear-server/pkg/lib/authn/identity"
)

var deleteIdentityInput = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "DeleteIdentityInput",
	Fields: graphql.InputObjectConfigFieldMap{
		"identityID": &graphql.InputObjectFieldConfig{
			Type:        graphql.NewNonNull(graphql.ID),
			Description: "Target identity ID.",
		},
	},
})

var deleteIdentityPayload = graphql.NewObject(graphql.ObjectConfig{
	Name: "DeleteIdentityPayload",
	Fields: graphql.Fields{
		"success": &graphql.Field{
			Type: graphql.NewNonNull(graphql.Boolean),
		},
	},
})

var _ = registerMutationField(
	"deleteIdentity",
	&graphql.Field{
		Description: "Delete identity of user",
		Type:        graphql.NewNonNull(deleteIdentityPayload),
		Args: graphql.FieldConfigArgument{
			"input": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(deleteIdentityInput),
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			input := p.Args["input"].(map[string]interface{})
			identityNodeID := input["identityID"].(string)

			resolvedNodeID := relay.FromGlobalID(identityNodeID)
			if resolvedNodeID == nil || resolvedNodeID.Type != typeIdentity {
				return nil, apierrors.NewInvalid("invalid identity ID")
			}
			identityRef, err := decodeIdentityID(resolvedNodeID.ID)
			if err != nil {
				return nil, apierrors.NewInvalid("invalid identity ID")
			}

			gqlCtx := GQLContext(p.Context)
			lazy := gqlCtx.Identities.Get(identityRef)
			return lazy.
				Map(func(value interface{}) (interface{}, error) {
					i := value.(*identity.Info)
					if i == nil {
						return nil, apierrors.NewNotFound("identity not found")
					}
					return gqlCtx.Identities.Remove(i), nil
				}).
				Map(func(value interface{}) (interface{}, error) {
					return map[string]bool{"success": true}, nil
				}).Value, nil
		},
	},
)
