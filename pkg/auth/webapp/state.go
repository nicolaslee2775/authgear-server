package webapp

import (
	"net/url"

	"github.com/authgear/authgear-server/pkg/api/apierrors"
)

type State struct {
	ID              string                 `json:"id"`
	PrevID          string                 `json:"prev_id"`
	Error           *apierrors.APIError    `json:"error"`
	RedirectURI     string                 `json:"redirect_uri,omitempty"`
	KeepState       bool                   `json:"keep_state,omitempty"`
	GraphInstanceID string                 `json:"graph_instance_id,omitempty"`
	Extra           map[string]interface{} `json:"extra,omitempty"`
	UserAgentToken  string                 `json:"user_agent_token"`
	UILocales       string                 `json:"ui_locales,omitempty"`
}

func NewState(intent *Intent) *State {
	return &State{
		ID:          NewID(),
		RedirectURI: intent.RedirectURI,
		KeepState:   intent.KeepState,
		Extra:       make(map[string]interface{}),
		UILocales:   intent.UILocales,
	}
}

func (s *State) SetID(id string) {
	s.PrevID = s.ID
	s.ID = id
}

func (s *State) RestoreFrom(s2 *State) {
	s.RedirectURI = s2.RedirectURI
	s.KeepState = s2.KeepState
	s.UILocales = s2.UILocales
	s.UserAgentToken = s2.UserAgentToken
	s.Error = s2.Error
}

func AttachStateID(id string, input *url.URL) *url.URL {
	u := *input

	q := u.Query()
	q.Set("x_sid", id)

	u.Scheme = ""
	u.Opaque = ""
	u.Host = ""
	u.User = nil
	u.RawQuery = q.Encode()

	return &u
}
