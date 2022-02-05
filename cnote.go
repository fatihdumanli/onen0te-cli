package cnote

import (
	"os"

	"github.com/fatihdumanli/cnote/internal/authentication"
	"github.com/fatihdumanli/cnote/internal/storage"
	"github.com/fatihdumanli/cnote/internal/survey"
	"github.com/fatihdumanli/cnote/pkg/oauthv2"
	"github.com/fatihdumanli/cnote/pkg/onenote"
)

type cnote struct {
	storage storage.Storer
	auth    authentication.Authenticator
	api     onenote.Api
	//The nil value is important for the business logic
	//So we're using a ptr type rather than value type
	token *oauthv2.OAuthToken
}

var (
	root cnote
)

func GetNotebooks() []onenote.Notebook {
	checkTokenPresented()

	var notebooks, err = root.api.GetNotebooks(*root.token)
	if err != nil {
		panic(err)
	}
	return notebooks
}

func GetSections(n onenote.Notebook) []onenote.Section {
	checkTokenPresented()

	var sections, err = root.api.GetSections(*root.token, n)
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

func checkTokenPresented() {
	var opts = getOptions()

	if root.token == nil {
		answer, err := survey.AskSetupAccount()
		if !answer || err != nil {
			//TODO: Maybe we should prompt the user about the loss of the note that they've just taken.
			os.Exit(1)
		}

		authentication.AuthenticateUser(opts, root.storage)
	} else {
		//Check if the token has expired
		if root.token.IsExpired() {
			authentication.RefreshToken(opts, *root.token, root.storage)
		}
	}

}

//Grab the token from the local storage upon startup
func init() {
	api := onenote.NewApi()
	bitcask := &storage.Bitcask{}
	root = cnote{api: api, storage: bitcask}

	t, err := root.storage.Get(authentication.TOKEN_KEY)
	if err == nil {
		token, ok := t.(oauthv2.OAuthToken)
		if !ok {
			panic(authentication.TokenStorageError)
		}
		root.token = &token
	}
}
