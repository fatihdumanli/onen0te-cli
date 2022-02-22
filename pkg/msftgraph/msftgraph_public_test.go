package msftgraph_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetNotebooks(t *testing.T) {

	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(400)
			w.Write([]byte("hey"))
		}))

	defer server.Close()

	var c = server.Client()
	_ = c
	fmt.Println(server.URL)

	res, err := c.Get(server.URL)
	fmt.Println(res)
	fmt.Println(err)

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
//func readFromFile(f string) ([]byte, error) {
//	file, err := os.Open(f)
//	if err != nil {
//		return nil, err
//	}
//	bytes, err := ioutil.ReadAll(file)
//	if err != nil {
//		return nil, err
//	}
//
//	return bytes, nil
//}
