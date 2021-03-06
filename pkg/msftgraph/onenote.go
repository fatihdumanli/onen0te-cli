package msftgraph

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	errors "github.com/pkg/errors"

	"github.com/fatihdumanli/onen0te-cli/pkg/oauthv2"
	"github.com/fatihdumanli/onen0te-cli/pkg/rest"
)

type HttpStatusCode = rest.HttpStatusCode

type Api struct {
	msftgraphURL string
	restClient   rest.Requester
}

func NewApi(r rest.Requester, msftgraphApiUrl string) Api {
	return Api{
		msftgraphURL: msftgraphApiUrl,
		restClient:   r,
	}
}

//https://docs.microsoft.com/en-us/graph/api/onenote-list-notebooks?view=graph-rest-1.0&tabs=http
func (a *Api) GetNotebooks(token oauthv2.OAuthToken) ([]Notebook, HttpStatusCode, error) {

	var response struct {
		Notebooks []Notebook `json:"value"`
	}

	var headers = make(map[string]string, 0)
	headers["Authorization"] = fmt.Sprintf("Bearer %s", token.AccessToken)

	res, statusCode, err := a.restClient.Get(a.msftgraphURL+"/me/onenote/notebooks", headers)
	if statusCode != http.StatusOK {
		return nil, statusCode, fmt.Errorf("couldn't get the notebooks from the server %s", string(res))
	}

	err = json.Unmarshal(res, &response)
	if err != nil {
		return nil, statusCode, errors.Wrap(err, "couldn't deserialize response data while getting the notebooks")
	}
	return response.Notebooks, statusCode, nil
}

//https://docs.microsoft.com/en-us/graph/api/onenote-list-sections?view=graph-rest-1.0&tabs=http
func (a *Api) GetSections(token oauthv2.OAuthToken, n Notebook) ([]Section, HttpStatusCode, error) {
	var response struct {
		Sections []Section `json:"value"`
	}

	var headers = make(map[string]string)
	headers["Authorization"] = fmt.Sprintf("Bearer %s", token.AccessToken)

	res, statusCode, err := a.restClient.Get(n.SectionsUrl, headers)

	if statusCode != http.StatusOK {
		return nil, statusCode, fmt.Errorf("couldn't get the sections from the server")
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

//https://docs.microsoft.com/en-us/graph/api/onenote-list-pages?view=graph-rest-1.0
func (a *Api) GetPages(token oauthv2.OAuthToken, section Section) ([]NotePage, HttpStatusCode, error) {

	var response struct {
		NotePages []NotePage `json:"value"`
	}

	var headers = make(map[string]string)
	headers["Authorization"] = fmt.Sprintf("Bearer %s", token.AccessToken)

	var url = fmt.Sprintf("%s/me/onenote/sections/%s/pages", a.msftgraphURL, section.ID)
	res, statusCode, err := a.restClient.Get(url, headers)
	if statusCode != http.StatusOK {
		return nil, statusCode, errors.Wrap(err, "couldn't get the note pages")
	}

	err = json.Unmarshal(res, &response)
	if err != nil {
		return nil, statusCode, errors.Wrap(err, "couldn't deserialize the response data")
	}

	return response.NotePages, statusCode, nil
}

//Searchs in all notes
//https://docs.microsoft.com/en-us/graph/search-query-parameter
func (a *Api) SearchPage(token oauthv2.OAuthToken, q string) ([]NotePage, HttpStatusCode, error) {
	var response struct {
		NotePages []NotePage `json:"value"`
	}

	var queryParam = fmt.Sprintf(`"%s"`, q)
	queryParam = url.QueryEscape(queryParam)

	url := fmt.Sprintf(`%s/me/onenote/pages?$search=%s`, a.msftgraphURL, queryParam)

	var headers = make(map[string]string)
	headers["Authorization"] = fmt.Sprintf("Bearer %s", token.AccessToken)

	res, statusCode, err := a.restClient.Get(url, headers)
	if err != nil {
		return nil, statusCode, errors.Wrap(err, "couldn't perform the search\n")
	}

	if statusCode != http.StatusOK {
		return nil, statusCode, fmt.Errorf("couldn't perform the search: %s\n", string(res))
	}

	err = json.Unmarshal(res, &response)
	if err != nil {
		return nil, statusCode, errors.Wrap(err, "couldn't deserialize the response data")
	}

	return response.NotePages, statusCode, nil
}

//Returns the page content as string
func (a *Api) GetContent(token oauthv2.OAuthToken, notepage NotePage) ([]byte, HttpStatusCode, error) {
	var headers = make(map[string]string)
	headers["Authorization"] = fmt.Sprintf("Bearer %s", token.AccessToken)

	res, statusCode, err := a.restClient.Get(notepage.ContentUrl, headers)
	if err != nil {
		return nil, statusCode, errors.Wrap(err, "couldn't get the note content\n")
	}

	if statusCode != http.StatusOK {
		return nil, statusCode, fmt.Errorf("couldn't get the note content %d - %s", statusCode, string(res))
	}

	return res, statusCode, nil
}

//Gets the section data by given id
//https://docs.microsoft.com/en-us/graph/api/section-get?view=graph-rest-1.0&tabs=go
func (a *Api) GetSection(token oauthv2.OAuthToken, id string) (Section, HttpStatusCode, error) {

	var headers = make(map[string]string)
	headers["Authorization"] = fmt.Sprintf("Bearer %s", token.AccessToken)

	url := fmt.Sprintf("%s/me/onenote/sections/%s", a.msftgraphURL, id)
	res, statusCode, err := a.restClient.Get(url, headers)
	if err != nil {
		return Section{}, statusCode, errors.Wrap(err, "couldn't get the section data\n")
	}

	if statusCode != http.StatusOK {
		return Section{}, statusCode, fmt.Errorf("couldn't get the section data %d - %s", statusCode, string(res))
	}

	var response Section
	err = json.Unmarshal(res, &response)
	if err != nil {
		return Section{}, statusCode, errors.Wrap(err, "couldn't deserialize the response data")
	}

	return response, statusCode, nil
}

//https://docs.microsoft.com/en-us/graph/api/onenote-post-pages?view=graph-rest-1.0
func (a *Api) SaveNote(t oauthv2.OAuthToken, n NotePage) (string, HttpStatusCode, error) {
	url := fmt.Sprintf("%s/me/onenote/sections/%s/pages", a.msftgraphURL, n.Section.ID)
	body := getNoteTemplate(n.Title, n.Content)

	var headers map[string]string = make(map[string]string, 0)
	headers["Authorization"] = fmt.Sprintf("Bearer %s", t.AccessToken)
	headers["Content-Type"] = "application/xhtml+xml"

	res, statusCode, err := a.restClient.Post(url, headers, strings.NewReader(body))
	if statusCode != http.StatusCreated {
		return "", statusCode, fmt.Errorf("couldn't save the note")
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
