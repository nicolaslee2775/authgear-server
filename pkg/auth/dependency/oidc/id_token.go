package oidc

import (
	gotime "time"

	"github.com/dgrijalva/jwt-go"
	"github.com/skygeario/skygear-server/pkg/auth/dependency/urlprefix"
	"github.com/skygeario/skygear-server/pkg/core/config"
	"github.com/skygeario/skygear-server/pkg/core/time"
)

type IDToken struct {
	jwt.StandardClaims
	Nonce string `json:"nonce,omitempty"`
}

type IDTokenIssuer struct {
	OIDCConfig config.OIDCConfiguration
	URLPrefix  urlprefix.Provider
	Time       time.Provider
}

const CodeGrantValidDuration = 5 * gotime.Minute

func (ti *IDTokenIssuer) IssueIDToken(client config.OAuthClientConfiguration, userID string, nonce string) (string, error) {
	now := ti.Time.NowUTC()
	token := &IDToken{
		StandardClaims: jwt.StandardClaims{
			Issuer:    ti.URLPrefix.Value().String(),
			Audience:  client.ClientID(),
			IssuedAt:  now.Unix(),
			ExpiresAt: now.Add(CodeGrantValidDuration).Unix(),
			Subject:   userID,
		},
		Nonce: nonce,
	}

	key, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(ti.OIDCConfig.Keys[0].PrivateKey))
	if err != nil {
		return "", err
	}
	jwt := jwt.NewWithClaims(jwt.SigningMethodRS256, token)
	return jwt.SignedString(key)
}
