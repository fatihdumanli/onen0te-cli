package authentication

import (
	"github.com/fatihdumanli/onen0te-cli/internal/config"
	"github.com/fatihdumanli/onen0te-cli/internal/storage"
	"github.com/fatihdumanli/onen0te-cli/pkg/oauthv2"
	errors "github.com/pkg/errors"
)

type TokenStatus int

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

func AuthenticateUser(oauthClient *oauthv2.OAuthClient, opts config.AppOptions, storer storage.Storer) (oauthv2.OAuthToken, error) {
	//If the user confirms to setup an account now we trigger the authentication process.
	t, err := oauthClient.Authenticate(opts.OAuthParams)
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

func RefreshToken(oauthClient *oauthv2.OAuthClient, opts config.AppOptions, token oauthv2.OAuthToken, storer storage.Storer) (oauthv2.OAuthToken, error) {
	newToken, err := oauthClient.RefreshToken(opts.OAuthParams, token.RefreshToken)
	if err != nil {
		return oauthv2.OAuthToken{}, errors.Wrap(err, "couldn't refresh the token\n")
	}

	err = storer.Set(TOKEN_KEY, newToken)
	if err != nil {
		return oauthv2.OAuthToken{}, errors.Wrap(err, "couldn't save the token\n")
	}

	return newToken, nil
}
