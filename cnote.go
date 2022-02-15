package cnote

import (
	"fmt"
	"log"
	"os"
	"time"

	errors "github.com/pkg/errors"

	"github.com/fatihdumanli/cnote/internal/authentication"
	"github.com/fatihdumanli/cnote/internal/storage"
	"github.com/fatihdumanli/cnote/internal/style"
	"github.com/fatihdumanli/cnote/internal/survey"
	"github.com/fatihdumanli/cnote/pkg/oauthv2"
	"github.com/fatihdumanli/cnote/pkg/onenote"
	"github.com/pterm/pterm"
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

var shouldRetry = func(statusCode onenote.HttpStatusCode) bool {
	var hashset map[int]bool
	hashset = make(map[int]bool, 0)
	hashset[503] = true
	hashset[504] = true
	return hashset[int(statusCode)]
}

//Get the list of notebooks belonging to the user logged in
func GetNotebooks() ([]onenote.Notebook, error) {
	checkTokenPresented()

	notebookSpinner, _ := pterm.DefaultSpinner.Start("Getting your notebooks...")
	var notebooks, statusCode, err = root.api.GetNotebooks(*root.token)

	if err != nil {
		if shouldRetry(statusCode) {
			//TODO: implement retry.
		}
		notebookSpinner.Fail(err.Error())
		return notebooks, errors.Wrap(err, "couldn't get the notebooks")
	}

	notebookSpinner.Success(pterm.FgDefault.Sprint("DONE"))
	return notebooks, nil
}

//Get the list of notebooks belonging to the user logged in
func GetSections(n onenote.Notebook) ([]onenote.Section, error) {
	checkTokenPresented()

	//TODO: We could wrap the code which the spinner would run while it's being executed
	sectionsSpinner, _ := pterm.DefaultSpinner.Start("Getting sections...")
	var sections, statusCode, err = root.api.GetSections(*root.token, n)

	if err != nil {
		//TODO: implement retry
		if shouldRetry(statusCode) {
		}
		sectionsSpinner.Fail(err.Error())
		return sections, errors.Wrap(err, "couldn't get the sections")
	}

	sectionsSpinner.Success(pterm.FgDefault.Sprint("DONE"))
	return sections, nil
}

//Save a note page using Onenote API
//Returns the link to the page.
func SaveNotePage(npage onenote.NotePage, remindAlias bool) (string, error) {
	checkTokenPresented()

	link, statusCode, err := root.api.SaveNote(*root.token, npage)
	if err != nil {
		if shouldRetry(statusCode) {
			//TODO: implement retry
		}
		return "", errors.Wrap(err, "couldn't save the note page")
	}

	var data = make([][]string, 2)
	data[0] = []string{"Status", "Notebook", "Section", "Title", "Note Size", "Elapsed", "Link", "SavedAt"}
	data[1] = []string{style.Success("OK"), npage.Section.Notebook.DisplayName, npage.Section.Name, npage.Title, "10 kB", "90 ms", "link FIXME", time.Now().Format(time.RFC3339)}
	pterm.DefaultTable.WithHasHeader().WithData(data).Render()
	fmt.Print("\n\n")

	aliases, err := GetAliases()
	if err != nil {
		return "", errors.Wrap(err, "couldn't get the alias list")
	}

	var answer string

	if !hasAlias(npage.Section, aliases) {
		answer, err := survey.AskAlias(npage.Section, aliases)
		if err != nil {
			return "", errors.Wrap(err, "couldn't ask the alias")
		}
		if answer != "" {
			err := SaveAlias(answer, *npage.Section.Notebook, npage.Section)
			if err != nil {
				return "", errors.Wrap(err, "couldn't save the alias")
			}
		}
	}

	//Print only if the alias didn't get created in this session.
	if answer == "" && remindAlias {
		printAliasReminder(npage.Section.Name)
	}

	return link, nil
}

func GetAliases() (*[]onenote.Alias, error) {

	var result []onenote.Alias
	keys, err := root.storage.GetKeys()
	if err != nil {
		return nil, errors.Wrap(err, "couldn't get the aliases")
	}

	var opts = getOptions()
	var hashset = make(map[string]bool, 0)
	for _, rk := range opts.ReservedKeys {
		hashset[rk] = true
	}

	if err != nil {
		return nil, errors.Wrap(err, "couldn't get the alias data")
	}

	for _, k := range *keys {
		if hashset[k] {
			continue
		}

		var a onenote.Alias
		root.storage.Get(k, &a)
		result = append(result, a)
	}

	return &result, nil
}

//Save the alias for a onenote section to use it later for quick save
func SaveAlias(name string, notebook onenote.Notebook, section onenote.Section) error {

	var isExist, err = GetAlias(name)
	if err != nil {
		return errors.Wrap(err, "couldn't check the alias if it already exists")
	}

	if isExist != nil {
		fmt.Println(style.Error(fmt.Sprintf("the alias %s already being used to identify the section %s", name, isExist.Section.Name)))
		return fmt.Errorf("the alias %s already exist", name)
	}

	err = root.storage.Set(name, onenote.Alias{
		Short:    name,
		Notebook: notebook,
		Section:  section})
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("couldn't save the alias %s to the local storage", name))
	}

	var msg = fmt.Sprintf("Alias '%s' has been saved.", name)
	fmt.Println(style.Success(msg))
	var infoMsg = "Now you can quickly add inline notes with the following command:"
	fmt.Println(style.Info(infoMsg))
	fmt.Println(fmt.Sprintf("$ cnote new -i \"inline text\" -a %s\n", name))
	return nil
}

//Get the details of given alias
//Returns nil if the alias does not found
func GetAlias(n string) (*onenote.Alias, error) {
	var alias onenote.Alias
	err := root.storage.Get(n, &alias)
	if err != nil {
		//TODO: Check if the error is KeyNotFound.
		return nil, nil
	}
	return &alias, nil
}

//Removes an alias
func RemoveAlias(a string) error {
	err := root.storage.Remove(a)
	if err != nil {
		var msg = fmt.Sprintf("The alias %s has not found.\n", a)
		fmt.Println(style.Error(msg))
		return errors.Wrap(err, "couldn't remove the alias")
	}

	var msg = fmt.Sprintf("The alias %s has been deleted.\n", a)
	fmt.Println(style.Success(msg))
	return nil
}

func hasAlias(section onenote.Section, aliasList *[]onenote.Alias) bool {
	for _, a := range *aliasList {
		if a.Section.ID == section.ID {
			return true
		}
	}
	return false
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

		token, err := authentication.AuthenticateUser(opts, root.storage)
		if err != nil {
			var errWrapped = errors.Wrap(err, "couldn't check the token if it's presented")
			log.Fatal(errWrapped)
		}
		root.token = &token
	} else {
		//Check if the token has expired
		if root.token.IsExpired() {
			token, err := authentication.RefreshToken(opts, *root.token, root.storage)
			if err != nil {
				var errWrapped = errors.Wrap(err, "couldn't check the token if it's presented")
				log.Fatal(errWrapped)
			}
			root.token = &token
		}
	}

}

//This function prints some alias instructions if the note has been created without using an alias
//Despite that the section the note was created in has an alias.
func printAliasReminder(section string) {
	var allAliases, _ = GetAliases() //we can ignore the err here
	for _, a := range *allAliases {
		if a.Section.Name == section {
			var msg = fmt.Sprintf("Existing alias for the section '%s' is '%s'", section, a.Short)
			fmt.Println(style.Reminder(msg))
			fmt.Printf("$ cnote new -i \"inline text\" -a %s\n", a.Short)
		}
	}
	fmt.Println()
}

//Grab the token from the local storage upon startup
func init() {
	api := onenote.NewApi()
	bitcask := &storage.Bitcask{}
	root = cnote{api: api, storage: bitcask}
	root.token = &oauthv2.OAuthToken{}

	err := root.storage.Get(authentication.TOKEN_KEY, root.token)
	if err != nil {
		//token does not exist on the local storage
		root.token = nil
	}
}
