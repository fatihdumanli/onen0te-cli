package authentiation

import (
	"errors"

	"github.com/fatihdumanli/cnote/pkg/oauthv2"
)

type TokenStatus int

const TOKEN_KEY = "msgraphtoken"

var InvalidTokenType = errors.New("Token type is invalid")

const (
	DoesntExist TokenStatus = iota
	Expired
	Valid
)

type Authenticator interface {
	CheckToken() (oauthv2.OAuthToken, TokenStatus)
	StoreToken() error
}
