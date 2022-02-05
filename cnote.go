package cnote

import (
	"github.com/fatihdumanli/cnote/internal/authentication"
	"github.com/fatihdumanli/cnote/internal/storage"
	"github.com/fatihdumanli/cnote/pkg/oauthv2"
	"github.com/fatihdumanli/cnote/pkg/onenote"
)

type cnote struct {
	storage storage.Storer
	auth    authentication.Authenticator
	api     onenote.Api
	token   oauthv2.OAuthToken
}

var (
	root cnote
)

func GetNotebooks() []onenote.Notebook {
	var notebooks, err = root.api.GetNotebooks(root.token)
	if err != nil {
		panic(err)
	}
	return notebooks
}

func GetSections(n onenote.Notebook) []onenote.Section {
	var sections, err = root.api.GetSections(root.token, n)
	if err != nil {
		panic(err)
	}
	return sections
}

func SaveAlias(aname, nname, sname string) error {
	return nil
}

func GetAlias(n string) onenote.Alias {
	return onenote.Alias{}
}

//Grab the token from the local storage upon startup
func init() {
	api := onenote.NewApi()
	bitcask := &storage.Bitcask{}
	root = cnote{api: api, storage: bitcask}

	t, err := root.storage.Get(authentication.TOKEN_KEY)
	if err != nil {
		//Token is not found on the storage.
	}
	token, ok := t.(oauthv2.OAuthToken)
	if !ok {
		panic(authentication.TokenStorageError)
	}

	root.token = token
}
