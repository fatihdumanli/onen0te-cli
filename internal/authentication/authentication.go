package authentication

import (
	"errors"
	"log"

	"github.com/fatihdumanli/cnote/internal/config"
	"github.com/fatihdumanli/cnote/internal/storage"
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

func AuthenticateUser(opts config.AppOptions, storer storage.Storer) (oauthv2.OAuthToken, error) {

	//If the user confirms to setup an account now we trigger the authentication process.
	t, err := oauthv2.Authorize(opts.OAuthParams, opts.Out)
	if err != nil {
		log.Fatalf("An error occured while trying to authenticate you. %s", err.Error())
	}

	//Save the token on local storage
	err = storer.Set(TOKEN_KEY, t)
	if err != nil {
		log.Fatalf("An error occured while trying to save the token. %s", err.Error())
		return t, storage.CouldntSaveTheKey
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
