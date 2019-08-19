package password

import (
	"github.com/skygeario/skygear-server/pkg/auth/dependency/audit"
	"github.com/skygeario/skygear-server/pkg/core/skydb"
	"github.com/skygeario/skygear-server/pkg/core/skyerr"
)

type ResetPasswordRequestContext struct {
	PasswordChecker      *audit.PasswordChecker
	PasswordAuthProvider Provider
}

func (r *ResetPasswordRequestContext) ExecuteWithPrincipals(newPassword string, principals []*Principal) (err error) {
	if err = r.PasswordChecker.ValidatePassword(audit.ValidatePasswordPayload{
		PlainPassword: newPassword,
	}); err != nil {
		return
	}

	for _, p := range principals {
		err = r.PasswordAuthProvider.UpdatePassword(p, newPassword)
		if err != nil {
			return
		}
	}

	return
}

func (r *ResetPasswordRequestContext) ExecuteWithUserID(newPassword string, userID string) (err error) {
	principals, err := r.PasswordAuthProvider.GetPrincipalsByUserID(userID)
	if err != nil {
		if err == skydb.ErrUserNotFound {
			err = skyerr.NewError(skyerr.ResourceNotFound, "user not found")
			return
		}

		return
	}

	err = r.ExecuteWithPrincipals(newPassword, principals)
	return
}