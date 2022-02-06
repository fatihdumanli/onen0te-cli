package cnote

import (
	"log"
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

func SaveNotePage(npage onenote.NotePage) error {
	checkTokenPresented()
	err := root.api.SaveNote(*root.token, npage)
	if err != nil {
		return err
	}
	return nil
}

func checkTokenPresented() {
	var opts = getOptions()

	if root.token == nil {
		answer, err := survey.AskSetupAccount()
		if !answer || err != nil {
			//TODO: Maybe we should prompt the user about the loss of the note that they've just taken.
			os.Exit(1)
		}

		token, _ := authentication.AuthenticateUser(opts, root.storage)
		root.token = &token
	} else {
		//Check if the token has expired
		if root.token.IsExpired() {
			token, err := authentication.RefreshToken(opts, *root.token, root.storage)
			if err != nil {
				log.Fatal("An error has occured while trying to refresh OAuth token")
				panic(err)
			}
			root.token = &token
		}
	}

}

//Grab the token from the local storage upon startup
func init() {

	api := onenote.NewApi()
	bitcask := &storage.Bitcask{}
	root = cnote{api: api, storage: bitcask}
	root.token = &oauthv2.OAuthToken{}
	//root.storage.Remove(authentication.TOKEN_KEY)
	err := root.storage.Get(authentication.TOKEN_KEY, root.token)
	if err != nil {
		root.token = nil
	}
}
