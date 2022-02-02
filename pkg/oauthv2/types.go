package oauthv2

import (
	"errors"
	"time"
)

var CannotGetAuthCode = errors.New("Couldn't get authorization_code")
var FailedToGetToken = errors.New("/token response was not 200")

type AuthorizationCode string
type OAuthParams struct {
	ClientId             string
	RedirectUri          string
	Scope                []string
	OAuthEndpoint        string
	RefreshTokenEndpoint string
	State                string
}

type getTokenParams struct {
	OAuthParams
	AuthCode AuthorizationCode
}

type OAuthToken struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	ExpiresIn    int    `json:"expires_in"`
	//will be used when saving the token on local storage
	ExpiresAt time.Time `json:"expires_at"`
}
