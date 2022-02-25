package msftgraph_test

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/fatihdumanli/onenote/pkg/msftgraph"
	"github.com/fatihdumanli/onenote/pkg/oauthv2"
	"github.com/fatihdumanli/onenote/pkg/rest"
	"github.com/google/go-cmp/cmp"
)

//Ergonomic alias
type AuthToken = oauthv2.OAuthToken
type Notebook = msftgraph.Notebook
type Section = msftgraph.Section
type HttpStatusCode = rest.HttpStatusCode
type NotePage = msftgraph.NotePage

//Assure that we're sending the following data to the remote server.
type apiDebugInfo struct {
	statusCode   int
	requestBody  io.Reader
	headers      map[string]string
	responseBody []byte
}

var token = AuthToken{
	AccessToken: "some-secret-stuff",
}

var notebooks = []msftgraph.Notebook{
	{
		"a-id", "Notebook A", "http://link-to-sections-of-notebook-a",
	},
	{
		"b-id", "Notebook B", "http://link-to-sections-of-notebook-b",
	},
	{
		"c-id", "Notebook C", "http://link-to-sections-of-notebook-c",
	},
}

//TODO: Centralize this function.
func launchTestHttpServer(io apiDebugInfo, t *testing.T) *httptest.Server {
	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			//Return 400 if the anticipated body is different than the one has sent to the server
			if io.requestBody != nil {
				if diff := cmp.Diff(io.requestBody, r.Body); diff != "" {
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

			//Write the anticipated status code if everything goes well
			w.WriteHeader(io.statusCode)
			w.Write(io.responseBody)

		}))
	return server

}

//Tests the function msftgraph.GetNotebooks()
func TestGetNotebooks(t *testing.T) {

	var io apiDebugInfo = apiDebugInfo{requestBody: nil, statusCode: 200}
	var headers map[string]string = make(map[string]string)
	headers["Authorization"] = "Bearer " + token.AccessToken
	io.headers = headers
	bytes, err := readFromFile("testdata/getnotebooks-200.json")
	if err != nil {
		t.Error(err.Error())
	}
	io.responseBody = bytes

	server := launchTestHttpServer(io, t)
	defer server.Close()

	var api = msftgraph.NewApi(&rest.RestClient{}, server.URL)
	data := []struct {
		name       string
		token      AuthToken
		apiio      apiDebugInfo
		statusCode int
		notebooks  []Notebook
		errMsg     string
	}{
		{"getnotebooks-200", token, io, io.statusCode, notebooks, ""},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {

			notebooks, statusCode, err := api.GetNotebooks(d.token)

			if int(statusCode) != d.statusCode {
				t.Errorf("expected status code %d, got %d", d.statusCode, statusCode)
			}

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

func Test_GetSections(t *testing.T) {

	var io apiDebugInfo = apiDebugInfo{requestBody: nil, statusCode: 200}
	var headers map[string]string = make(map[string]string)
	headers["Authorization"] = "Bearer " + token.AccessToken
	io.headers = headers
	bytes, err := readFromFile("testdata/getsections-200.json")
	if err != nil {
		t.Error(err.Error())
	}
	io.responseBody = bytes

	server := launchTestHttpServer(io, t)
	defer server.Close()
	var api = msftgraph.NewApi(&rest.RestClient{}, server.URL)

	var notebook = Notebook{SectionsUrl: server.URL}
	var sections = []msftgraph.Section{
		{
			"Section A1", "a1", &notebook,
		},
		{
			"Section A2", "a2", &notebook,
		},
	}

	data := []struct {
		name       string
		token      AuthToken
		notebook   Notebook
		apiio      apiDebugInfo
		statusCode int
		sections   []Section
		errMsg     string
	}{
		{"getsections-200", token, Notebook{SectionsUrl: server.URL}, io, io.statusCode, sections, ""},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {

			sections, statusCode, err := api.GetSections(d.token, d.notebook)

			if int(statusCode) != d.statusCode {
				t.Errorf("expected status code %d, got %d", d.statusCode, statusCode)
			}

			if diff := cmp.Diff(sections, d.sections); diff != "" {
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

func Test_SaveNote(t *testing.T) {
	var io apiDebugInfo = apiDebugInfo{requestBody: nil, statusCode: 201}
	var headers map[string]string = make(map[string]string)
	headers["Authorization"] = "Bearer " + token.AccessToken
	io.headers = headers
	bytes, err := readFromFile("testdata/savenote-200.json")
	if err != nil {
		t.Error(err.Error())
	}
	io.responseBody = bytes

	server := launchTestHttpServer(io, t)
	defer server.Close()
	var api = msftgraph.NewApi(&rest.RestClient{}, server.URL)

	var notepage = msftgraph.NewNotePage(msftgraph.Section{Name: "pack-rat"}, "title", "some content")
	data := []struct {
		name       string
		token      AuthToken
		notepage   NotePage
		apiio      apiDebugInfo
		statusCode int
		link       string
		errMsg     string
	}{
		//NOTE: The expected note link is defined in /testdata/savenote-200.json
		{"savenote-200", token, *notepage, io, io.statusCode, "http://new-note", ""},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			link, statusCode, err := api.SaveNote(d.token, d.notepage)

			if int(statusCode) != d.statusCode {
				t.Errorf("expected status code %d, got %d", d.statusCode, statusCode)
			}

			if link != d.link {
				t.Errorf("expected url %s, got %s", d.link, link)
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

//type RestStub struct {
//	get  func(url string, headers map[string]string) ([]byte, rest.HttpStatusCode, error)
//	post func(url string, headers map[string]string, body io.Reader) ([]byte, rest.HttpStatusCode, error)
//}
//
//func (stub *RestStub) Get(url string, headers map[string]string) ([]byte, rest.HttpStatusCode, error) {
//	return stub.get(url, headers)
//}
//
//func (stub *RestStub) Post(url string, headers map[string]string, body io.Reader) ([]byte, rest.HttpStatusCode, error) {
//	return stub.Post(url, headers, body)
//}
//
////Tests msftgraph.GetNotebooks()
//func TestGetNotebooks(t *testing.T) {
//
//	data := []struct {
//		name       string
//		get        func(url string, headers map[string]string) ([]byte, rest.HttpStatusCode, error)
//		token      oauthv2.OAuthToken
//		statusCode int
//		notebooks  []msftgraph.Notebook
//		errMsg     string
//	}{
//		{"getnotebooks-successful",
//			func(url string, headers map[string]string) ([]byte, rest.HttpStatusCode, error) {
//				var body, err = readFromFile("testdata/getnotebooks-200.json")
//				if err != nil {
//					return nil, 000, errors.New("check the json file path")
//				}
//				return body, http.StatusOK, nil
//			},
//			oauthv2.OAuthToken{},
//			200,
//			[]msftgraph.Notebook{
//				{
//					"a-id", "Notebook A", "http://link-to-sections-of-notebook-a",
//				},
//				{
//					"b-id", "Notebook B", "http://link-to-sections-of-notebook-b",
//				},
//				{
//					"c-id", "Notebook C", "http://link-to-sections-of-notebook-c",
//				},
//			},
//			""},
//
//		{"getnotebooks-4xx",
//			func(url string, headers map[string]string) ([]byte, rest.HttpStatusCode, error) {
//				return nil, 400, nil
//			},
//			oauthv2.OAuthToken{},
//			200, nil, ""},
//	}
//
//	for _, d := range data {
//		var api = msftgraph.NewApi(
//			&RestStub{get: d.get})
//
//		t.Run(d.name, func(t *testing.T) {
//
//			notebooks, _, err := api.GetNotebooks(d.token)
//			if diff := cmp.Diff(notebooks, d.notebooks); diff != "" {
//				t.Error(diff)
//			}
//
//			var errMsg string
//			if err != nil {
//				errMsg = err.Error()
//			}
//
//			if errMsg != d.errMsg {
//				t.Errorf("expected error message `%s`, got `%s`", d.errMsg, errMsg)
//			}
//		})
//
//	}
//
//	_ = data
//}
//
////Tests the function msftgraph.GetSections()
//func TestGetSections(t *testing.T) {
//	data := []struct {
//		name       string
//		notebook   msftgraph.Notebook
//		get        func(url string, headers map[string]string) ([]byte, rest.HttpStatusCode, error)
//		token      oauthv2.OAuthToken
//		statusCode int
//		sections   []msftgraph.Section
//		errMsg     string
//	}{
//		{"getsections-200", msftgraph.Notebook{DisplayName: "Notebook A", SectionsUrl: "http://notebook-a/sections"},
//			func(url string, headers map[string]string) ([]byte, rest.HttpStatusCode, error) {
//				var body, err = readFromFile("testdata/getsections-200.json")
//				if err != nil {
//					return nil, 000, errors.New("check the json file path")
//				}
//				return body, http.StatusOK, nil
//			},
//			oauthv2.OAuthToken{}, 200, []msftgraph.Section{
//				{
//					"Section A1", "a1", nil,
//				},
//				{
//					"Section A2", "a2", nil,
//				},
//			},
//			""},
//		{"getsections-4xx", msftgraph.Notebook{}, func(url string, headers map[string]string) ([]byte, rest.HttpStatusCode, error) {
//			return nil, 400, nil
//		}, oauthv2.OAuthToken{}, 400, nil, "couldn't get the sections from the server"},
//	}
//
//	sectionComparer := cmp.Comparer(func(x, y msftgraph.Section) bool {
//		return x.Name == y.Name && x.ID == y.ID
//	})
//
//	for _, d := range data {
//		var api = msftgraph.NewApi(&RestStub{
//			get: d.get,
//		})
//		t.Run(d.name, func(t *testing.T) {
//			sections, _, err := api.GetSections(d.token, d.notebook)
//			if diff := cmp.Diff(sections, d.sections, sectionComparer); diff != "" {
//				t.Error(diff)
//			}
//
//			var errMsg string
//			if err != nil {
//				errMsg = err.Error()
//			}
//
//			if errMsg != d.errMsg {
//				t.Errorf("expected error message `%s`, got `%s`", d.errMsg, errMsg)
//			}
//		})
//
//	}
//
//}
//
////Tests the function msftgraph.SaveNote()
//func TestSaveNote(t *testing.T) {
//	data := []struct {
//		name       string
//		notepage   msftgraph.NotePage
//		post       func(url string, headers map[string]string, body io.Reader) ([]byte, rest.HttpStatusCode, error)
//		token      oauthv2.OAuthToken
//		statusCode int
//		link       string
//		errMsg     string
//	}{
//		{"savenote-201", msftgraph.NotePage{},
//			func(url string, headers map[string]string, body io.Reader) ([]byte, rest.HttpStatusCode, error) {
//				var b, err = readFromFile("testdata/savenote-200.json")
//				if err != nil {
//					return nil, 000, errors.New("check the json file path")
//				}
//				return b, http.StatusCreated, nil
//			},
//			oauthv2.OAuthToken{}, 201, "http://new-note", ""},
//		{"savenote-4xx", msftgraph.NotePage{},
//			func(url string, headers map[string]string, body io.Reader) ([]byte, rest.HttpStatusCode, error) {
//				return nil, 400, nil
//			},
//			oauthv2.OAuthToken{}, 400, "", "couldn't save the note"},
//	}
//
//	for _, d := range data {
//		var api = msftgraph.NewApi(&RestStub{post: d.post})
//		t.Run(d.name, func(t *testing.T) {
//			link, _, err := api.SaveNote(d.token, d.notepage)
//
//			if link != d.link {
//				t.Errorf("expected link %s, got %s", d.link, link)
//			}
//
//			var errMsg string
//			if err != nil {
//				errMsg = err.Error()
//			}
//
//			if errMsg != d.errMsg {
//				t.Errorf("expected error message `%s`, got `%s`", d.errMsg, errMsg)
//			}
//		})
//	}
//
//}
//
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
