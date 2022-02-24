package oauthv2

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/fatihdumanli/onenote/pkg/rest"
	"github.com/google/go-cmp/cmp"
)

//Assure that we're sending the following data to the remote server.
type apiDebugInfo struct {
	statusCode   int
	requestBody  io.Reader
	headers      map[string]string
	responseBody []byte
}

func launchTestHttpServer(io apiDebugInfo, t *testing.T) *httptest.Server {
	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			//Return 400 if the expected body is different than the one has sent to the server
			if io.requestBody != nil {
				bytes, _ := ioutil.ReadAll(r.Body)
				expectedBytes, _ := ioutil.ReadAll(io.requestBody)

				if diff := cmp.Diff(bytes, expectedBytes); diff != "" {
					t.Error(diff)
					w.WriteHeader(400)
					return
				}
			}

			for k, v := range io.headers {
				var headerValue = r.Header.Get(k)

				if headerValue == "" {
					t.Error(fmt.Sprintf("the header %s was not sent to the server", k))
				}

				if headerValue != v {
					t.Error(fmt.Sprintf("exptected value for the header %s was %s, got %s", k, v, headerValue))
					w.WriteHeader(400)
					return
				}
			}

			//Write the expected status code if everything goes well
			w.WriteHeader(io.statusCode)
			w.Write(io.responseBody)
		}))
	return server
}

var oauthParams = OAuthParams{
	ClientId:      "client-id",
	RedirectUri:   "redirect-uri",
	Scope:         []string{"read", "write", "execute"},
	OAuthEndpoint: "http://some-oauth-endpoint",
	State:         "some-state",
}

var scope = strings.Join(oauthParams.Scope, " ")

var tokenParams = getTokenParams{
	OAuthParams: oauthParams,
	AuthCode:    "auth-code",
}

var token = OAuthToken{
	AccessToken:  "some-secret-stuff",
	RefreshToken: "refresh-your-token-with-this",
	ExpiresIn:    3600,
}

func TestOauth_GetToken(t *testing.T) {
	var io apiDebugInfo = apiDebugInfo{statusCode: 200}

	data := url.Values{}
	data.Set("client_id", oauthParams.ClientId)
	data.Set("scope", scope)
	data.Set("code", string(tokenParams.AuthCode))
	data.Set("redirect_uri", oauthParams.RedirectUri)
	data.Set("grant_type", "authorization_code")
	io.requestBody = strings.NewReader(data.Encode())

	var headers map[string]string = make(map[string]string)
	headers["Content-Type"] = "application/x-www-form-urlencoded"
	io.headers = headers

	var tokenInjson, _ = json.Marshal(token)
	io.responseBody = []byte(tokenInjson)

	server := launchTestHttpServer(io, t)
	defer server.Close()
	tokenParams.OAuthEndpoint = server.URL

	var oauthClient = NewOAuthClient(&rest.RestClient{})

	//name: Name for the test case
	//apiio: Expected Api Input / Output
	//tokenParams: Token parameters that'll be passed in to getTokenParams() function
	//token: Expected token result
	//errMsg: Expected error msg
	testCases := []struct {
		name        string
		apiio       apiDebugInfo
		tokenParams getTokenParams
		token       OAuthToken
		errMsg      string
	}{
		{"gettoken-200", io, tokenParams, token, ""},
	}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			res, err := oauthClient.getToken(c.tokenParams)

			if diff := cmp.Diff(*res, c.token); diff != "" {
				t.Error(diff)
			}

			//compare err and c.errMsg
			var errMsg string
			if err != nil {
				errMsg = err.Error()
			}

			if errMsg != c.errMsg {
				t.Errorf("expected error message `%s`, got `%s`", c.errMsg, errMsg)
			}

		})
	}

}
