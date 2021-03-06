package oauthv2

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/fatihdumanli/onen0te-cli/internal/util/process"
	"github.com/fatihdumanli/onen0te-cli/pkg/rest"
	errors "github.com/pkg/errors"
)

func (t *OAuthToken) IsExpired() bool {
	return time.Now().After(t.ExpiresAt)
}

type OAuthClient struct {
	restClient rest.Requester
}

func NewOAuthClient(restClient rest.Requester) *OAuthClient {
	return &OAuthClient{
		restClient: restClient,
	}
}

func (o *OAuthClient) RefreshToken(p OAuthParams, refreshToken string) (OAuthToken, error) {
	var newToken OAuthToken

	var data = url.Values{}
	data.Set("client_id", p.ClientId)
	data.Set("scope", strings.Join(p.Scope, " "))
	data.Set("refresh_token", refreshToken)
	data.Set("redirect_uri", p.RedirectUri)
	data.Set("grant_type", "refresh_token")

	tokenPath := fmt.Sprintf("%s/token", p.OAuthEndpoint)

	var headers map[string]string = make(map[string]string)
	headers["Content-Type"] = "application/x-www-form-urlencoded"

	res, statusCode, err := o.restClient.Post(tokenPath, headers, strings.NewReader(data.Encode()))
	if statusCode != http.StatusOK {
		return newToken, fmt.Errorf("status code is %d %s", int(statusCode), string(res))
	}

	err = json.Unmarshal(res, &newToken)
	if err != nil {
		return newToken, errors.Wrap(err, "couldn't unmarshal the json while refreshing token.")
	}

	//set expiredAt proprety so that we can check if the token has expired
	t := time.Duration(newToken.ExpiresIn)
	newToken.ExpiresAt = time.Now().Add(time.Second * t)

	return newToken, nil
}

func (o *OAuthClient) Authenticate(p OAuthParams) (OAuthToken, error) {

	var token *OAuthToken
	authCode, err := o.getAuthCode(p)

	if err != nil {
		return OAuthToken{}, errors.Wrap(err, "couldn't authenticate the user")
	}

	var getTokenParams = getTokenParams{
		OAuthParams: p,
		AuthCode:    authCode,
	}

	token, err = o.getToken(getTokenParams)
	if err != nil {
		return OAuthToken{}, errors.Wrap(err, "couldn't get the oauth token")
	}

	//set expiredAt proprety so that we can check if the token has expired
	t := time.Duration(token.ExpiresIn)
	token.ExpiresAt = time.Now().Add(time.Second * t)

	return *token, nil
}

func (o *OAuthClient) getAuthCode(p OAuthParams) (AuthorizationCode, error) {

	authCodeUrl := fmt.Sprintf("%s/authorize?client_id=%s&response_type=code&redirect_uri=%s&response_mode=query&scope=%s&state=%s", p.OAuthEndpoint, p.ClientId, p.RedirectUri, strings.Join(p.Scope, " "), p.State)

	var authCode AuthorizationCode
	var srv *http.Server

	var fnCallback = func(w http.ResponseWriter, r *http.Request) {
		values := r.URL.Query()
		authCode = AuthorizationCode(values.Get("code"))
		fmt.Fprintln(w, "You can close this window now.")
		time.Sleep(1 * time.Second)
		srv.Shutdown(context.TODO())
	}

	httpServerExitDone := &sync.WaitGroup{}
	httpServerExitDone.Add(1)

	srv, err := startOauthHttpServer(httpServerExitDone, ":5992", "/oauthv2", fnCallback)
	if err != nil {
		return "", errors.Wrap(err, "couldn't get the auth code")
	}

	//open a browser to authenticate the user
	err = process.Start(authCodeUrl)
	if err != nil {
		return "", errors.Wrap(err, "couldn't get auth code")
	}

	fmt.Println("Please complete authentication process through your web browser...")
	httpServerExitDone.Wait()

	if authCode == "" {
		return authCode, fmt.Errorf("couldn't get the auth code (auth code was empty)")
	}

	return authCode, nil
}

func (o *OAuthClient) getToken(p getTokenParams) (*OAuthToken, error) {

	var token OAuthToken
	tokenPath := fmt.Sprintf("%s/token", p.OAuthEndpoint)

	scope := strings.Join(p.Scope, " ")

	data := url.Values{}
	data.Set("client_id", p.ClientId)
	data.Set("scope", scope)
	data.Set("code", string(p.AuthCode))
	data.Set("redirect_uri", p.RedirectUri)
	data.Set("grant_type", "authorization_code")

	var headers map[string]string = make(map[string]string)
	headers["Content-Type"] = "application/x-www-form-urlencoded"

	res, statusCode, err := o.restClient.Post(tokenPath, headers, strings.NewReader(data.Encode()))

	if statusCode != http.StatusOK {
		return nil, fmt.Errorf("couldn't get the oauth token")
	}

	err = json.Unmarshal(res, &token)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't deserialize the response")
	}

	return &token, nil
}

func startOauthHttpServer(wg *sync.WaitGroup, addr string, pattern string, callback http.HandlerFunc) (*http.Server, error) {
	srv := &http.Server{Addr: addr}
	http.HandleFunc(pattern, callback)

	var ch chan error = make(chan error)
	go func() error {
		defer wg.Done()

		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			ch <- errors.Wrap(err, "couldn't start the http server")
		}
		close(ch)
		return nil
	}()

	select {
	case v := <-ch:
		return nil, v
	default:
		return srv, nil
	}
}
