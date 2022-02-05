package cnote

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/fatihdumanli/cnote/internal/storage"
	"github.com/fatihdumanli/cnote/pkg/oauthv2"
	"github.com/fatihdumanli/cnote/pkg/onenote"
)

type Cnote struct {
	storage storage.Storer
	//	auth    authentication.Authenticator
}

func (cnote *Cnote) GetNotebooks(token oauthv2.OAuthToken) ([]onenote.Notebook, error) {

	var response struct {
		Notebooks []onenote.Notebook `json:"value"`
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

func (cnote *Cnote) GetSections(t oauthv2.OAuthToken, n onenote.Notebook) ([]onenote.Section, error) {
	var response struct {
		Sections []onenote.Section `json:"value"`
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
func (cnote *Cnote) SaveAlias(aname, nname, sname string) error {
	return nil
}

func (cnote *Cnote) GetAlias(n string) onenote.Alias {
	return onenote.Alias{}
}
