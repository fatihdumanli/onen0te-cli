package onenote

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/fatihdumanli/cnote/pkg/oauthv2"
)

type Api struct {
	GetNotebooks func(token oauthv2.OAuthToken) ([]Notebook, error)
	GetSections  func(token oauthv2.OAuthToken, n Notebook) ([]Section, error)
	//Saves the note and returns the link to the newly created note.
	SaveNote func(token oauthv2.OAuthToken, n NotePage) (string, error)
}

func NewApi() Api {
	return Api{
		GetNotebooks: getNotebooks,
		GetSections:  getSections,
		SaveNote:     saveNote,
	}
}

func getNotebooks(token oauthv2.OAuthToken) ([]Notebook, error) {

	var response struct {
		Notebooks []Notebook `json:"value"`
	}

	client := http.Client{}
	request, err := http.NewRequest(http.MethodGet, "https://graph.microsoft.com/v1.0/me/onenote/notebooks", nil)

	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))
	resp, err := client.Do(request)

	if err != nil {
		panic(err)
	}

	if resp.StatusCode != http.StatusOK {
		errStr := fmt.Sprintf("An error has occured while fetching your notebooks...")
		return response.Notebooks, errors.New(errStr)
	}

	bytes, err := io.ReadAll(resp.Body)

	if err != nil {
		return response.Notebooks, err
	}

	err = json.Unmarshal(bytes, &response)
	if err != nil {
		return response.Notebooks, err
	}

	return response.Notebooks, nil
}

func getSections(t oauthv2.OAuthToken, n Notebook) ([]Section, error) {
	var response struct {
		Sections []Section `json:"value"`
	}

	//Ditch the nil pointer.
	//TODO: We need to find out how to set section.Notebook
	for _, s := range response.Sections {
		s.Notebook = &Notebook{}
	}

	c := http.Client{}
	req, err := http.NewRequest(http.MethodGet, n.SectionsUrl, nil)
	if err != nil {
		return response.Sections, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", t.AccessToken))
	res, err := c.Do(req)
	if err != nil {
		return response.Sections, err
	}

	if res.StatusCode != http.StatusOK {
		errStr := fmt.Sprintf("An error has occured while fetching sections..")
		return response.Sections, errors.New(errStr)
	}

	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return response.Sections, err
	}

	err = json.Unmarshal(bytes, &response)
	if err != nil {
		return response.Sections, err
	}

	//Set notebook ptr of each section in the response
	for _, s := range response.Sections {
		s.Notebook = &n
	}

	return response.Sections, nil
}

func saveNote(t oauthv2.OAuthToken, n NotePage) (string, error) {
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
		log.Fatal(err)
		return "", err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", t.AccessToken))
	req.Header.Add("Content-Type", "application/xhtml+xml")

	res, err := c.Do(req)
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	resBody := res.Body
	bytes, _ := io.ReadAll(resBody)
	if res.StatusCode != http.StatusCreated {
		log.Fatal(string(bytes))
		return "", errors.New("Couldn't save the note.")
	}

	var response struct {
		Links struct {
			OneNoteWebUrl struct {
				Href string `json:"href"`
			} `json:"oneNoteWebUrl"`
		} `json:"links"`
	}

	json.Unmarshal(bytes, &response)

	return response.Links.OneNoteWebUrl.Href, nil
}
