package msftgraph

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	errors "github.com/pkg/errors"

	"github.com/fatihdumanli/onenote/pkg/oauthv2"
	"github.com/fatihdumanli/onenote/pkg/rest"
)

type HttpStatusCode = rest.HttpStatusCode

type Api struct {
	restClient rest.Requester
}

func NewApi() Api {
	return Api{
		restClient: &rest.RestClient{},
	}
}

func (a *Api) GetNotebooks(token oauthv2.OAuthToken) ([]Notebook, HttpStatusCode, error) {

	var response struct {
		Notebooks []Notebook `json:"value"`
	}

	var headers = make(map[string]string, 0)
	headers["Authorization"] = fmt.Sprintf("Bearer %s", token.AccessToken)

	res, statusCode, err := a.restClient.Get("https://graph.microsoft.com/v1.0/me/onenote/notebooks", headers)
	if statusCode != http.StatusOK {
		return nil, statusCode, fmt.Errorf("the statusCode doesn't indicate a successful operation: %d", statusCode)
	}

	err = json.Unmarshal(res, &response)
	if err != nil {
		return nil, statusCode, errors.Wrap(err, "couldn't deserialize response data while getting the notebooks")
	}

	return response.Notebooks, statusCode, nil
}

func (a *Api) GetSections(token oauthv2.OAuthToken, n Notebook) ([]Section, HttpStatusCode, error) {
	var response struct {
		Sections []Section `json:"value"`
	}

	var headers = make(map[string]string, 0)
	headers["Authorization"] = fmt.Sprintf("Bearer %s", token.AccessToken)

	res, statusCode, err := a.restClient.Get(n.SectionsUrl, headers)
	if statusCode != http.StatusOK {
		return nil, statusCode, fmt.Errorf("couldn't load the sections: %s", string(res))
	}

	err = json.Unmarshal(res, &response)
	if err != nil {
		return nil, statusCode, errors.Wrap(err, "couldn't deserialize the response data")
	}

	//Set notebook ptr of each section in the response
	for i := 0; i < len(response.Sections); i++ {
		response.Sections[i].Notebook = &n
	}

	return response.Sections, statusCode, nil
}

//TODO:Complete
//func getPages(token oauthv2.OAuthToken) (*Notebook, HttpStatusCode, error) {
//
//	var headers = make(map[string]string, 0)
//	headers["Authorization"] = fmt.Sprintf("Bearer %s", token.AccessToken)
//	res, statusCode, err := makeHttpRequest("https://graph.microsoft.com/v1.0/me/onenote/pages", http.MethodGet, nil, headers)
//	_ = res
//	_ = statusCode
//	_ = err
//
//	return nil, 200, nil
//}

func (a *Api) SaveNote(t oauthv2.OAuthToken, n NotePage) (string, HttpStatusCode, error) {
	url := fmt.Sprintf("https://graph.microsoft.com/v1.0/me/onenote/sections/%s/pages", n.Section.ID)
	body := getNoteTemplate(n.Title, n.Content)

	var headers map[string]string
	headers["Authorization"] = fmt.Sprintf("Bearer %s", t.AccessToken)
	headers["Content-Type"] = "application/xhtml+xml"

	res, statusCode, err := a.restClient.Post(url, headers, strings.NewReader(body))

	if statusCode != http.StatusCreated {
		return "", statusCode, fmt.Errorf("couldn't save the note: %s", string(res))
	}

	var response struct {
		Links struct {
			OneNoteWebUrl struct {
				Href string `json:"href"`
			} `json:"oneNoteWebUrl"`
		} `json:"links"`
	}

	err = json.Unmarshal(res, &response)
	if err != nil {
		return "", statusCode, errors.Wrap(err, "couldn't deserialize the response data")
	}

	return response.Links.OneNoteWebUrl.Href, statusCode, nil
}

func getNoteTemplate(title, content string) string {

	var body = `<html>
			<head>
				<title>` + title + `</title>
			</head>
			<body>
				<p>` + content + `</p>
			</body>
		</html>`

	return body
}
