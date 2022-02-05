package oauthv2

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"time"
)

/*
NOTE
Go uses parameters of pointer type to indicate that a parameter might be modified by the function.
The same rules apply for method receivers, too.

--> If your method modifies the receiver, you must use a pointer receiver
--> If your method needs to handle nil instances, then it must use a pointer receiver
--> If your method doesn't modify the receiver you can use a value receiver.


notice that if we had used a value receiver, there wouldn't be a way to mutate the receiver...
however, we don't need  to modify the receiver within this method.
*/
func (t *OAuthToken) IsExpired() bool {
	return time.Now().After(t.ExpiresAt)
}

func RefreshToken(p OAuthParams, r string) (OAuthToken, error) {
	var newToken OAuthToken

	var data = url.Values{}
	data.Set("client_id", p.ClientId)
	data.Set("scope", strings.Join(p.Scope, " "))
	data.Set("refresh_token", r)
	data.Set("redirect_uri", p.RedirectUri)
	data.Set("grant_type", "refresh_token")

	c := http.Client{}
	req, err := http.NewRequest(http.MethodPost, p.RefreshTokenEndpoint, strings.NewReader(data.Encode()))
	if err != nil {
		log.Fatal(err)
		return newToken, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := c.Do(req)
	if err != nil {
		log.Fatal(err)
		return newToken, err
	}

	if res.StatusCode != http.StatusOK {
		return newToken, errors.New("Refreshtoken didn't return a 200")
	}

	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
		return newToken, err
	}
	defer res.Body.Close()
	err = json.Unmarshal(bytes, &newToken)
	if err != nil {
		log.Fatal(err)
		return newToken, err

	}
	//set expiredAt proprety so that we can check if the token has expired
	t := time.Duration(newToken.ExpiresIn)
	newToken.ExpiresAt = time.Now().Add(time.Second * t)

	return newToken, nil
}

/*
NOTE
Rather than returning a pointer set to nil,
Use comma ok idiom
return a boolean and a value type
*/
func Authorize(p OAuthParams, out io.Writer) (OAuthToken, error) {

	var token OAuthToken

	authCode, err := getAuthCode(p, out)
	if err != nil {
		return token, err
	}

	var getTokenParams = getTokenParams{
		OAuthParams: p,
		AuthCode:    authCode,
	}

	token, err = getToken(getTokenParams)

	//set expiredAt proprety so that we can check if the token has expired
	t := time.Duration(token.ExpiresIn)
	token.ExpiresAt = time.Now().Add(time.Second * t)

	if err != nil {
		return token, err
	}

	return token, nil
}

func getAuthCode(p OAuthParams, out io.Writer) (AuthorizationCode, error) {

	authCodeUrl := fmt.Sprintf("%s/authorize?client_id=%s&response_type=code&redirect_uri=%s&response_mode=query&scope=%s&state=%s", p.OAuthEndpoint, p.ClientId, p.RedirectUri, strings.Join(p.Scope, " "), p.State)

	var authCode AuthorizationCode

	//closure
	//notice that we've passed a closure and we've utilized a local variable, awesome!
	var fnCallback = func(w http.ResponseWriter, r *http.Request) {
		values := r.URL.Query()
		authCode = AuthorizationCode(values.Get("code"))
		fmt.Fprintln(w, "You can close this window now.")
	}

	httpServerExitDone := &sync.WaitGroup{}
	httpServerExitDone.Add(1)
	srv := startOauthHttpServer(httpServerExitDone, ":5992", "/oauthv2", fnCallback)

	//open a browser to authenticate the user
	openWebBrowser(authCodeUrl)

	fmt.Fprintln(out, "Please complete authentication process through your web browser...")
	time.Sleep(10 * time.Second)

	if err := srv.Shutdown(context.TODO()); err != nil {
		panic(err)
	}

	httpServerExitDone.Wait()

	if authCode == "" {
		return authCode, CannotGetAuthCode
	}

	return authCode, nil
}

func getToken(p getTokenParams) (OAuthToken, error) {

	var token OAuthToken
	tokenPath := fmt.Sprintf("%s/token", p.OAuthEndpoint)

	scope := strings.Join(p.Scope, " ")

	data := url.Values{}
	data.Set("client_id", p.ClientId)
	data.Set("scope", scope)
	data.Set("code", string(p.AuthCode))
	data.Set("redirect_uri", p.RedirectUri)
	data.Set("grant_type", "authorization_code")

	client := http.Client{}

	request, _ := http.NewRequest(http.MethodPost, tokenPath, strings.NewReader(data.Encode()))
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(request)

	if err != nil {
		return token, err
	}

	if resp.StatusCode != http.StatusOK {
		return token, FailedToGetToken
	}

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	err = json.Unmarshal(bytes, &token)
	if err != nil {
		return token, err
	}

	return token, nil
}

func openWebBrowser(url string) {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		log.Fatal("your OS is not supported.")
	}

	if err != nil {
		log.Fatal(err)
	}
}

func startOauthHttpServer(wg *sync.WaitGroup, addr string, pattern string, callback http.HandlerFunc) *http.Server {
	srv := &http.Server{Addr: addr}
	http.HandleFunc(pattern, callback)
	go func() {
		defer wg.Done()

		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe(): %v", err)
		}
	}()
	return srv
}
