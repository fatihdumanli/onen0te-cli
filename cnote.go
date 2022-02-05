package cnote

import (
	"github.com/fatihdumanli/cnote/internal/authentication"
	"github.com/fatihdumanli/cnote/internal/storage"
	"github.com/fatihdumanli/cnote/pkg/oauthv2"
	"github.com/fatihdumanli/cnote/pkg/onenote"
)

var authtoken oauthv2.OAuthToken

type cnote struct {
	storage storage.Storer
	auth    authentication.Authenticator
	api     onenote.Api
}

var (
	root = cnote{api: onenote.NewApi()}
)

func GetNotebooks() []onenote.Notebook {
	var token, _ = root.auth.GetToken()
	var notebooks, err = root.api.GetNotebooks(token)
	if err != nil {
		panic(err)
	}
	return notebooks
}

func GetSections(n onenote.Notebook) []onenote.Section {
	var token, _ = root.auth.GetToken()
	var sections, err = root.api.GetSections(token, n)
	if err != nil {
		panic(err)
	}
	return sections
}

func (cnote *cnote) SaveAlias(aname, nname, sname string) error {
	return nil
}

func (cnote *cnote) GetAlias(n string) onenote.Alias {
	return onenote.Alias{}
}

//Grab the token from the local storage upon startup
func init() {
	t, err := root.storage.Get(authentication.TOKEN_KEY)
	if err != nil {
		//Token is not found on the storage.
	}
	token, ok := t.(oauthv2.OAuthToken)
	if !ok {
		panic(authentication.TokenStorageError)
	}

	authtoken = token
}
