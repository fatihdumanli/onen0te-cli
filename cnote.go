package cnote

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

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
	AliasAlreadyExists = errors.New("alias already exists")
)

var (
	root cnote
)

//Get the list of notebooks belonging to the user logged in
func GetNotebooks() ([]onenote.Notebook, bool) {
	checkTokenPresented()

	notebookSpinner, _ := pterm.DefaultSpinner.Start("Getting your notebooks...")
	var notebooks, err = root.api.GetNotebooks(*root.token)
	if err != nil {
		notebookSpinner.Fail(err.Error())
		return notebooks, false
	}
	//TODO: What if it fails, consider use retry.
	notebookSpinner.Success()

	return notebooks, true
}

//Get the list of notebooks belonging to the user logged in
func GetSections(n onenote.Notebook) ([]onenote.Section, bool) {
	checkTokenPresented()

	sectionsSpinner, _ := pterm.DefaultSpinner.Start("Getting sections...")
	var sections, err = root.api.GetSections(*root.token, n)
	if err != nil {
		sectionsSpinner.Fail(err.Error())
		return sections, false
	}
	//TODO: What if it fails, consider use retry.
	sectionsSpinner.Success()
	return sections, true
}

//Save a note page using Onenote API
//Returns the link to the page.
func SaveNotePage(npage onenote.NotePage, remindAlias bool) (string, error) {
	checkTokenPresented()

	link, err := root.api.SaveNote(*root.token, npage)
	if err != nil {
		log.Fatal("couldn't save the note.")
		return "", err
	}

	var data = make([][]string, 2)
	data[0] = []string{"Status", "Notebook", "Section", "Title", "Note Size", "Elapsed", "Link", "SavedAt"}
	data[1] = []string{style.Success("OK"), npage.Section.Notebook.DisplayName, npage.Section.Name, npage.Title, "10 kB", "90 ms", "link FIXME", time.Now().Format(time.RFC3339)}
	pterm.DefaultTable.WithHasHeader().WithData(data).Render()
	fmt.Print("\n\n")

	ok := askAlias(npage.Section)
	//Print only if the alias didn't get created in this session.
	if !ok && remindAlias {
		printAliasReminder(npage.Section.Name)
	}

	return link, nil
}

func GetAliases() *[]onenote.Alias {

	var result []onenote.Alias
	keys, err := root.storage.GetKeys()

	var opts = getOptions()
	var hashset = make(map[string]bool, 0)
	for _, rk := range opts.ReservedKeys {
		hashset[rk] = true
	}

	if err != nil {
		log.Fatalf("An error has occured while trying to get the alias data %s", err.Error())
		return nil
	}

	for _, k := range *keys {
		if hashset[k] {
			continue
		}

		var a onenote.Alias
		root.storage.Get(k, &a)
		result = append(result, a)
	}

	return &result
}

//Save the alias for a onenote section to use it later for quick save
func saveAlias(name string, notebook onenote.Notebook, section onenote.Section) error {

	var isExist = GetAlias(name)
	if isExist != nil {
		fmt.Println(style.Error(fmt.Sprintf("the alias %s already being used to identify the section %s", name, isExist.Section.Name)))
		return AliasAlreadyExists
	}

	err := root.storage.Set(name, onenote.Alias{
		Short:    name,
		Notebook: notebook,
		Section:  section})
	if err != nil {
		log.Fatalf("An error has occured while saving the alias to the local storage %s", err.Error())
		return err
	}

	var msg = fmt.Sprintf("Alias '%s' has been saved.\n", name)
	fmt.Println(style.Success(msg))
	var infoMsg = "Now you can quickly add notes with the following command:"
	fmt.Println(style.Info(infoMsg))
	fmt.Println(fmt.Sprintf("$ cnote new <path-to-input> -a %s\n", name))
	return nil
}

//Get the details of given alias
//Returns nil if the alias does not found
func GetAlias(n string) *onenote.Alias {
	var alias onenote.Alias
	err := root.storage.Get(n, &alias)
	if err != nil {
		return nil
	}
	return &alias
}

//Removes an alias
func RemoveAlias(a string) error {
	err := root.storage.Remove(a)
	if err != nil {
		var msg = fmt.Sprintf("The alias %s has not found.\n", a)
		fmt.Println(style.Error(msg))
		return err
	}

	var msg = fmt.Sprintf("The alias %s has been deleted.\n", a)
	fmt.Println(style.Success(msg))
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

//Ask alias to make it easy to create a note within the section
//Returns TRUE if the user has craeted an alias through this function
//If the user has just created the alias, we're not going to display alias reminder.
func askAlias(s onenote.Section) bool {
	allAliases := GetAliases()
	//Exit the funciton if there's already an alias for the section.
	for _, a := range *allAliases {
		if a.Section.ID == s.ID {
			return false
		}
	}

	for {
		a, _ := survey.AskAlias(onenote.NotebookName(s.Notebook.DisplayName), onenote.SectionName(s.Name))
		if a != "" {
			//User answered with an alias
			err := saveAlias(a, *s.Notebook, s)
			if err != nil {
				//if alias already exists
				if err == AliasAlreadyExists {
					continue
				}

				log.Fatal(style.Error("Couldn't save the alias."))
				return false
			}
			return true
		} else {
			break
		}
	}

	return false
}

//This function prints some alias instructions if the note has been created without using an alias
//Despite that the section the note was created in has an alias.
func printAliasReminder(section string) {
	var allAliases = GetAliases()
	for _, a := range *allAliases {
		if a.Section.Name == section {
			var msg = fmt.Sprintf("Existing alias for the section '%s' is '%s'", section, a.Short)
			fmt.Println(style.Reminder(msg))
			fmt.Printf("$ cnote new <input-file-path> -a %s will do the trick.\n", a.Short)
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
		root.token = nil
	}
}
