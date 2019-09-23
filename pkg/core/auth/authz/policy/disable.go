package policy

import (
	"net/http"

	"github.com/skygeario/skygear-server/pkg/core/auth"
	"github.com/skygeario/skygear-server/pkg/core/auth/authz"
	"github.com/skygeario/skygear-server/pkg/core/skyerr"
)

func DenyDisabledUser(r *http.Request, ctx auth.ContextGetter) error {
	authInfo, _ := ctx.AuthInfo()
	if authInfo == nil {
		return skyerr.NewError(skyerr.UnexpectedAuthInfoNotFound, "user authentication info not found")
	}

	if authInfo.Disabled {
		// TODO:
		// return proper error code
		return skyerr.NewError(skyerr.UnexpectedError, "user disabled")
	}

	return nil
}

// this ensures that our structure conform to certain interfaces.
var (
	_ authz.PolicyFunc = DenyDisabledUser
)
