package msftgraph_test

import (
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
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
		var body, err = readFromFile("testdata/getnotebooks-200.json")
		if err != nil {
			return nil, 000, errors.New("check the json file path")
		}
		return body, http.StatusOK, nil
	}

	//Get sections (Notebook A)
	if strings.HasSuffix(url, "/sections") {
		var body, err = readFromFile("testdata/getsections-200.json")
		if err != nil {
			return nil, 000, errors.New("check the json file path")
		}
		return body, http.StatusOK, nil
	}

	return nil, 000, nil
}

func (r SuccessfulRestStub) Post(url string, headers map[string]string, body io.Reader) ([]byte, rest.HttpStatusCode, error) {

	//Save note
	if strings.HasSuffix(url, "/pages") {
		var body, err = readFromFile("testdata/savenote-200.json")
		if err != nil {
			return nil, 000, errors.New("check the json file path")
		}
		return body, http.StatusCreated, nil
	}

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

//Tests msftgraph.GetNotebooks()
func TestGetNotebooks(t *testing.T) {
	data := []struct {
		name       string
		restStub   rest.Requester
		token      oauthv2.OAuthToken
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

			notebooks, _, err := api.GetNotebooks(d.token)
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

//Tests the function msftgraph.GetSections()
func TestGetSections(t *testing.T) {
	data := []struct {
		name       string
		notebook   msftgraph.Notebook
		restStub   rest.Requester
		token      oauthv2.OAuthToken
		statusCode int
		sections   []msftgraph.Section
		errMsg     string
	}{
		{"getsections-200", msftgraph.Notebook{DisplayName: "Notebook A", SectionsUrl: "http://notebook-a/sections"}, SuccessfulRestStub{}, oauthv2.OAuthToken{},
			200, []msftgraph.Section{
				{
					"Section A1", "a1", nil,
				},
				{
					"Section A2", "a2", nil,
				},
			}, "",
		},
		{"getsections-4xx", msftgraph.Notebook{}, BadRequestRestStub{}, oauthv2.OAuthToken{}, 400, nil, "couldn't get the sections from the server"},
	}

	sectionComparer := cmp.Comparer(func(x, y msftgraph.Section) bool {
		return x.Name == y.Name && x.ID == y.ID
	})

	for _, d := range data {
		var api = msftgraph.NewApi(d.restStub)
		t.Run(d.name, func(t *testing.T) {
			sections, _, err := api.GetSections(d.token, d.notebook)
			if diff := cmp.Diff(sections, d.sections, sectionComparer); diff != "" {
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

//Tests the function msftgraph.SaveNote()
func TestSaveNote(t *testing.T) {

	data := []struct {
		name       string
		notepage   msftgraph.NotePage
		token      oauthv2.OAuthToken
		restStub   rest.Requester
		statusCode int
		link       string
		errMsg     string
	}{
		{"savenote-200", msftgraph.NotePage{}, oauthv2.OAuthToken{}, SuccessfulRestStub{}, 201, "http://new-note", ""},
		{"savenote-4xx", msftgraph.NotePage{}, oauthv2.OAuthToken{}, BadRequestRestStub{}, 400, "", "couldn't save the note"},
	}

	for _, d := range data {
		var api = msftgraph.NewApi(d.restStub)
		t.Run(d.name, func(t *testing.T) {
			link, _, err := api.SaveNote(d.token, d.notepage)

			if link != d.link {
				t.Errorf("expected link %s, got %s", d.link, link)
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
