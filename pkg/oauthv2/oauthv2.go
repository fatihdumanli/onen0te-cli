package oauthv2

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"time"
)

type OAuthParams struct {
	ClientId      string
	RedirectUri   string
	Scope         []string
	OAuthEndpoint string
	State         string
}

type getTokenParams struct {
	OAuthParams
	AuthCode string
}

type OAuthToken struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	ExpiresIn    int    `json:"expires_in"`
}

func Authorize(p OAuthParams) error {
	authCode, err := getAuthCode(p)
	if err != nil {
		return err
	}

	var getTokenParams = getTokenParams{
		OAuthParams: p,
		AuthCode:    authCode,
	}

	token, err := getToken(getTokenParams)
	_ = token

	return nil
}

func getAuthCode(p OAuthParams) (string, error) {

	authCodeUrl := fmt.Sprintf("%s/authorize?client_id=%s&response_type=code&redirect_uri=%s&response_mode=query&scope=%s&state=%s", p.OAuthEndpoint, p.ClientId, p.RedirectUri, strings.Join(p.Scope, " "), p.State)

	var authCode string

	//closure
	//notice that we've passed a closure and we've utilized a local variable, awesome!
	var fnCallback = func(w http.ResponseWriter, r *http.Request) {
		values := r.URL.Query()
		authCode = values.Get("code")
		fmt.Fprintln(w, "You can close this window now.")
	}

	httpServerExitDone := &sync.WaitGroup{}
	httpServerExitDone.Add(1)
	srv := startOauthHttpServer(httpServerExitDone, ":5992", "/oauthv2", fnCallback)

	//open a browser to authenticate the user
	openWebBrowser(authCodeUrl)

	fmt.Println("Please complete authentication process through your web browser...")
	time.Sleep(20 * time.Second)

	if err := srv.Shutdown(context.TODO()); err != nil {
		panic(err)
	}

	httpServerExitDone.Wait()

	return authCode, nil
}

func getToken(p getTokenParams) (OAuthToken, error) {
	tokenPath := fmt.Sprintf("%s/token", p.OAuthEndpoint)

	scope := strings.Join(p.Scope, " ")

	data := url.Values{}
	data.Set("client_id", p.ClientId)
	data.Set("scope", scope)
	data.Set("code", p.AuthCode)
	data.Set("redirect_uri", p.RedirectUri)
	data.Set("grant_type", "authorization_code")

	client := http.Client{}

	request, _ := http.NewRequest(http.MethodPost, tokenPath, strings.NewReader(data.Encode()))
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, _ := client.Do(request)

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var token OAuthToken
	json.Unmarshal(bytes, &token)

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
