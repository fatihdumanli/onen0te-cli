package onenote

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/fatihdumanli/cnote/pkg/oauthv2"
)

func GetNotebooks(token oauthv2.OAuthToken) ([]Notebook, error) {

	var response GetNotebooksResponse

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
		return []Notebook{}, err
	}

	err = json.Unmarshal(bytes, &response)
	if err != nil {
		return []Notebook{}, err
	}

	return response.Notebooks, nil
}

func GetSections(n NotebookName) ([]Section, error) {
	return getDummySections(), nil
}

func getDummySections() []Section {
	return []Section{
		{"Quick notes"},
		{"Go"},
		{"Projects"},
		{"Todos"},
	}
}
