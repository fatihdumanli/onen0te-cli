package authentication

import (
	"errors"

	"github.com/fatihdumanli/cnote/pkg/oauthv2"
)

type TokenStatus int

var (
	InvalidTokenType  = errors.New("Token type is invalid")
	TokenStorageError = errors.New("Stored token is corrupted")
)

const (
	DoesntExist TokenStatus = iota
	Expired
	Valid
	TOKEN_KEY = "msgraphtoken"
)

type Authenticator interface {
	GetToken() (oauthv2.OAuthToken, TokenStatus)
	StoreToken() error
	RefreshToken() error
}
