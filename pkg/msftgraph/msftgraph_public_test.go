package msftgraph_test

import (
	"errors"
	"io"
	"io/ioutil"
	"os"
	"testing"

	"github.com/fatihdumanli/onenote/pkg/msftgraph"
	"github.com/fatihdumanli/onenote/pkg/oauthv2"
	"github.com/fatihdumanli/onenote/pkg/rest"
	"github.com/google/go-cmp/cmp"
)

type SuccessfulRestStub struct {
	rest.Requester
}

func (r SuccessfulRestStub) Get(url string, headers map[string]string) ([]byte, rest.HttpStatusCode, error) {

	//Get notebooks
	if url == "https://graph.microsoft.com/v1.0/me/onenote/notebooks" {
		//Return a successful json
		var body, err = readFromFile("testdata/getnotebook-success.json")
		if err != nil {
			return nil, 000, errors.New("check the json file path")
		}
		return body, 200, nil
	}

	//TODO Get section (need regex here)
	if url == "" {
	}

	//TODO Save note (need regex here)
	if url == "https://graph.microsoft.com/v1.0/me/onenote/sections/%s/pages" {
	}

	return nil, 000, nil
}

func (r SuccessfulRestStub) Post(url string, headers map[string]string, body io.Reader) ([]byte, rest.HttpStatusCode, error) {
	return nil, 000, nil
}

type BadRequestRestStub struct {
	rest.Requester
}

func (r BadRequestRestStub) Get(url string, headers map[string]string) ([]byte, rest.HttpStatusCode, error) {
	return nil, 400, nil
}

func (r BadRequestRestStub) Post(url string, headers map[string]string, body io.Reader) ([]byte, rest.HttpStatusCode, error) {
	return nil, 400, nil
}

func Test_GetNotebooks(t *testing.T) {

	data := []struct {
		name     string
		restStub rest.Requester
		oauthv2.OAuthToken
		statusCode int
		notebooks  []msftgraph.Notebook
		errMsg     string
	}{
		{"getnotebooks-successful", SuccessfulRestStub{}, oauthv2.OAuthToken{},
			200, []msftgraph.Notebook{
				{
					"a-id", "Notebook A", "http://link-to-sections-of-notebook-a",
				},
				{
					"b-id", "Notebook B", "http://link-to-sections-of-notebook-b",
				},
				{
					"c-id", "Notebook C", "http://link-to-sections-of-notebook-c",
				},
			}, ""},
		{"getnotebooks-4xx", BadRequestRestStub{}, oauthv2.OAuthToken{}, 400, nil, "couldn't get the notebooks from the server"},
	}

	for _, d := range data {

		var api = msftgraph.NewApi(d.restStub)
		t.Run(d.name, func(t *testing.T) {

			notebooks, _, err := api.GetNotebooks(d.OAuthToken)
			if diff := cmp.Diff(notebooks, d.notebooks); diff != "" {
				t.Error(diff)
			}

			var errMsg string
			if err != nil {
				errMsg = err.Error()
			}

			if errMsg != d.errMsg {
				t.Errorf("expected error message `%s`, got `%s`", d.errMsg, errMsg)
			}
		})
	}

}

func readFromFile(f string) ([]byte, error) {
	file, err := os.Open(f)
	if err != nil {
		return nil, err
	}
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}
