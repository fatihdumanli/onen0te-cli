package onenote

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	errors "github.com/pkg/errors"

	"github.com/fatihdumanli/cnote/pkg/oauthv2"
)

type HttpStatusCode int
type Api struct {
	GetNotebooks func(token oauthv2.OAuthToken) ([]Notebook, HttpStatusCode, error)
	GetSections  func(token oauthv2.OAuthToken, n Notebook) ([]Section, HttpStatusCode, error)
	//Saves the note and returns the link to the newly created note.
	SaveNote func(token oauthv2.OAuthToken, n NotePage) (string, HttpStatusCode, error)
}

func NewApi() Api {
	return Api{
		GetNotebooks: getNotebooks,
		GetSections:  getSections,
		SaveNote:     saveNote,
	}
}

func getNotebooks(token oauthv2.OAuthToken) ([]Notebook, HttpStatusCode, error) {

	var response struct {
		Notebooks []Notebook `json:"value"`
	}

	client := http.Client{}
	request, err := http.NewRequest(http.MethodGet, "https://graph.microsoft.com/v1.0/me/onenote/notebooks", nil)
	if err != nil {
		return nil, 000, errors.Wrap(err, "couldn't initialize the request while getting notebooks")
	}

	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))

	resp, err := client.Do(request)
	if err != nil {
		return nil, 000, errors.Wrap(err, "couldn't make the request while getting notebooks")
	}

	var statusCode HttpStatusCode = HttpStatusCode(resp.StatusCode)
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, statusCode, errors.Wrap(err, "couldn't read the response while getting notebooks")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, statusCode, fmt.Errorf("couldn't get the notebooks: %s", string(respBody))
	}

	err = json.Unmarshal(respBody, &response)
	if err != nil {
		return nil, statusCode, errors.Wrap(err, "couldn't deserialize response data while getting the notebooks")
	}

	return response.Notebooks, statusCode, nil
}

func getSections(t oauthv2.OAuthToken, n Notebook) ([]Section, HttpStatusCode, error) {
	var response struct {
		Sections []Section `json:"value"`
	}
	c := http.Client{}
	req, err := http.NewRequest(http.MethodGet, n.SectionsUrl, nil)
	if err != nil {
		return nil, 000, errors.Wrap(err, "couldn't initialize the request")
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", t.AccessToken))
	res, err := c.Do(req)
	if err != nil {
		return nil, 000, errors.Wrap(err, "couldn't make the request")
	}
	var statusCode HttpStatusCode = HttpStatusCode(res.StatusCode)
	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, statusCode, errors.Wrap(err, "couldn't read the response")
	}
	if res.StatusCode != http.StatusOK {
		return nil, statusCode, fmt.Errorf("couldn't load the sections: %s", string(bytes))
	}

	err = json.Unmarshal(bytes, &response)
	if err != nil {
		return nil, statusCode, errors.Wrap(err, "couldn't deserialize the response data")
	}
	//Set notebook ptr of each section in the response
	for i := 0; i < len(response.Sections); i++ {
		response.Sections[i].Notebook = &n
	}

	return response.Sections, statusCode, nil
}

func saveNote(t oauthv2.OAuthToken, n NotePage) (string, HttpStatusCode, error) {
	c := http.Client{}

	var body = `<html>
<head>
<title>` + n.Title + `</title>
</head>
<body>
<p>` + n.Content + `</p>
</bod>
</html>`

	url := fmt.Sprintf("https://graph.microsoft.com/v1.0/me/onenote/sections/%s/pages", n.Section.ID)
	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(body))
	if err != nil {
		return "", 000, errors.Wrap(err, "couldn't initialize the request")
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", t.AccessToken))
	req.Header.Add("Content-Type", "application/xhtml+xml")

	res, err := c.Do(req)
	if err != nil {
		return "", 000, errors.Wrap(err, "couldn't make the request")
	}

	var statusCode HttpStatusCode = HttpStatusCode(res.StatusCode)
	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return "", statusCode, errors.Wrap(err, "couldn't read the response")
	}
	if res.StatusCode != http.StatusCreated {
		return "", statusCode, fmt.Errorf("couldn't save the note: %s", string(bytes))
	}

	var response struct {
		Links struct {
			OneNoteWebUrl struct {
				Href string `json:"href"`
			} `json:"oneNoteWebUrl"`
		} `json:"links"`
	}

	err = json.Unmarshal(bytes, &response)
	if err != nil {
		return "", statusCode, errors.Wrap(err, "couldn't deserialize the response data")
	}

	return response.Links.OneNoteWebUrl.Href, statusCode, nil
}

func readResBody(r io.Reader) *string {
	var bytes, err = ioutil.ReadAll(r)
	if err != nil {
	}
	var strBody = string(bytes)
	return &strBody
}
