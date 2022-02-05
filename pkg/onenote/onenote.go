package onenote

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/fatihdumanli/cnote/pkg/oauthv2"
)

type Api struct {
	GetNotebooks func(token oauthv2.OAuthToken) ([]Notebook, error)
	GetSections  func(token oauthv2.OAuthToken, n Notebook) ([]Section, error)
}

func NewApi() Api {
	return Api{
		GetNotebooks: getNotebooks,
		GetSections:  getSections,
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

	return response.Sections, nil
}
