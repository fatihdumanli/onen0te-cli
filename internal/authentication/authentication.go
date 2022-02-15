package authentication

import (
	"github.com/fatihdumanli/onenote/internal/config"
	"github.com/fatihdumanli/onenote/internal/storage"
	"github.com/fatihdumanli/onenote/pkg/oauthv2"
	errors "github.com/pkg/errors"
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

func AuthenticateUser(opts config.AppOptions, storer storage.Storer) (oauthv2.OAuthToken, error) {

	//If the user confirms to setup an account now we trigger the authentication process.
	t, err := oauthv2.Authenticate(opts.OAuthParams, opts.Out)
	if err != nil {
		return oauthv2.OAuthToken{}, errors.Wrap(err, "couldn't authenticate the user")
	}

	//Save the token on local storage
	err = storer.Set(TOKEN_KEY, t)
	if err != nil {
		return t, errors.Wrap(err, "couldn't save the token")
	}

	return t, nil
}

func RefreshToken(opts config.AppOptions, token oauthv2.OAuthToken, storer storage.Storer) (oauthv2.OAuthToken, error) {
	newToken, err := oauthv2.RefreshToken(opts.OAuthParams, token.RefreshToken)
	if err != nil {
		panic(err)
	}

	err = storer.Set(TOKEN_KEY, newToken)
	if err != nil {
		panic(err)
	}

	return newToken, nil
}
