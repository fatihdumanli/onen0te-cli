package msftgraph_test

import (
	"io"
	"testing"

	"github.com/fatihdumanli/onenote/pkg/msftgraph"
	"github.com/fatihdumanli/onenote/pkg/oauthv2"
	"github.com/fatihdumanli/onenote/pkg/rest"
)

type RestStub struct {
	rest.Requester
}

func (r RestStub) Get(url string, headers map[string]string) ([]byte, rest.HttpStatusCode, error) {
	return nil, 000, nil
}

func (r RestStub) Post(url string, headers map[string]string, body io.Reader) ([]byte, rest.HttpStatusCode, error) {
	return nil, 000, nil
}

func Test_GetNotebooks(t *testing.T) {

	var api = msftgraph.NewApi(RestStub{})
	notebooks, statusCode, err := api.GetNotebooks(oauthv2.OAuthToken{})

	t.Log("the status code is", statusCode)

	_ = notebooks
	_ = statusCode
	_ = err
}
