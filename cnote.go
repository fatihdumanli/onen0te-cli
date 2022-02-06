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

//Get the list of notebooks belonging to the user logged in
func GetNotebooks() []onenote.Notebook {
	checkTokenPresented()
	var notebooks, err = root.api.GetNotebooks(*root.token)
	if err != nil {
		panic(err)
	}
	return notebooks
}

//Get the list of notebooks belonging to the user logged in
func GetSections(n onenote.Notebook) []onenote.Section {
	checkTokenPresented()

	var sections, err = root.api.GetSections(*root.token, n)
	if err != nil {
		panic(err)
	}
	return sections
}

//Save the alias for a onenote section to use it later for quick save
func SaveAlias(aname, nname, sname string) error {
	err := root.storage.Set(aname, onenote.Alias{
		Notebook: onenote.NotebookName(nname),
		Section:  onenote.SectionName(sname),
	})
	if err != nil {
		log.Fatalf("An error has occured while saving the alias to the local storage %s", err.Error())
		return err
	}
	return nil
}

//Get the details of given alias
func GetAlias(n string) *onenote.Alias {
	var alias onenote.Alias
	err := root.storage.Get(n, &alias)
	if err != nil {
		log.Fatalf("An error has occured while trying to get the alias data %s", err.Error())
		return nil
	}

	return &alias
}

//Save a note page using Onenote API
func SaveNotePage(npage onenote.NotePage) error {
	checkTokenPresented()
	err := root.api.SaveNote(*root.token, npage)
	if err != nil {
		return err
	}
	return nil
}

//This function gets called prior to each API operation to make sure that we're not going to deal with any stale token.
//Check if the OAuth token has presented on the local storage.
//Prompt the user if token doesn't exist on the local storage
//Refresh it if the token has expired
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

	err := root.storage.Get(authentication.TOKEN_KEY, root.token)
	if err != nil {
		root.token = nil
	}
}
