package webapp

import (
	"net/http"

	"github.com/authgear/authgear-server/pkg/auth/config"
	"github.com/authgear/authgear-server/pkg/auth/dependency/auth"
	"github.com/authgear/authgear-server/pkg/db"
	"github.com/authgear/authgear-server/pkg/httproute"
	"github.com/authgear/authgear-server/pkg/template"
)

const (
	TemplateItemTypeAuthUILogoutHTML config.TemplateItemType = "auth_ui_logout.html"
)

var TemplateAuthUILogoutHTML = template.Spec{
	Type:        TemplateItemTypeAuthUILogoutHTML,
	IsHTML:      true,
	Translation: TemplateItemTypeAuthUITranslationJSON,
	Defines:     defines,
	Components:  components,
	Default: `<!DOCTYPE html>
<html>
{{ template "auth_ui_html_head.html" . }}
<body class="page">
<div class="content">

{{ template "auth_ui_header.html" . }}

<form class="logout-form" method="post" novalidate>
  {{ $.csrfField }}
  <p class="primary-txt">{{ localize "logout-button-hint" }}</p>
  <button class="btn primary-btn align-self-center" type="submit" name="x_action" value="logout">{{ localize "logout-button-label" }}</button>
</form>

{{ template "auth_ui_footer.html" . }}

</div>
</body>
</html>
`,
}

func ConfigureLogoutRoute(route httproute.Route) httproute.Route {
	return route.
		WithMethods("OPTIONS", "POST", "GET").
		WithPathPattern("/logout")
}

type LogoutSessionManager interface {
	Logout(auth.AuthSession, http.ResponseWriter) error
}

type LogoutHandler struct {
	Database       *db.Handle
	ServerConfig   *config.ServerConfig
	SessionManager LogoutSessionManager
	BaseViewModel  *BaseViewModeler
	Renderer       Renderer
}

func (h *LogoutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		baseViewModel := h.BaseViewModel.ViewModel(r, nil)

		data := map[string]interface{}{}

		Embed(data, baseViewModel)

		h.Renderer.Render(w, r, TemplateItemTypeAuthUILogoutHTML, data)
		return
	}

	// FIXME(webapp): logout
	// h.Database.WithTx(func() error {
	// 	if r.Method == "POST" && r.Form.Get("x_action") == "logout" {
	// 		sess := auth.GetSession(r.Context())
	// 		h.SessionManager.Logout(sess, w)
	// 		webapp.RedirectToRedirectURI(w, r, h.ServerConfig.TrustProxy)
	// 	}
	// 	return nil
	// })
}
