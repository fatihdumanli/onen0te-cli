package onenote

import (
	"fmt"
	"log"
	"os"

	errors "github.com/pkg/errors"

	"github.com/fatihdumanli/onen0te-cli/internal/authentication"
	"github.com/fatihdumanli/onen0te-cli/internal/storage"
	"github.com/fatihdumanli/onen0te-cli/internal/style"
	"github.com/fatihdumanli/onen0te-cli/internal/survey"
	"github.com/fatihdumanli/onen0te-cli/internal/util/file"
	"github.com/fatihdumanli/onen0te-cli/pkg/msftgraph"
	"github.com/fatihdumanli/onen0te-cli/pkg/oauthv2"
	"github.com/fatihdumanli/onen0te-cli/pkg/rest"
	"github.com/pterm/pterm"
)

type onenote struct {
	storage     storage.Storer
	auth        authentication.Authenticator
	api         msftgraph.Api
	token       *oauthv2.OAuthToken
	oauthClient *oauthv2.OAuthClient
}

var (
	root onenote
)

var shouldRetry = func(statusCode msftgraph.HttpStatusCode) bool {
	var hashset map[int]bool
	hashset = make(map[int]bool, 0)
	hashset[503] = true
	hashset[504] = true
	return hashset[int(statusCode)]
}

//Get the list of notebooks belonging to the user logged in
func GetNotebooks() ([]msftgraph.Notebook, error) {
	checkTokenPresented()

	notebookSpinner, _ := pterm.DefaultSpinner.Start("Getting your notebooks...")
	notebookSpinner.RemoveWhenDone = true
	var notebooks, statusCode, err = root.api.GetNotebooks(*root.token)

	if err != nil {
		if shouldRetry(statusCode) {
			//TODO: implement retry.
		}
		notebookSpinner.Fail(err.Error())
		return notebooks, errors.Wrap(err, "couldn't get the notebooks\n")
	}

	notebookSpinner.Success(pterm.FgDefault.Sprint("DONE"))
	return notebooks, nil
}

//Get the list of notebooks belonging to the user logged in
func GetSections(n msftgraph.Notebook) ([]msftgraph.Section, error) {
	checkTokenPresented()

	sectionsSpinner, _ := pterm.DefaultSpinner.Start("Getting sections...")
	sectionsSpinner.RemoveWhenDone = true
	var sections, statusCode, err = root.api.GetSections(*root.token, n)

	if err != nil {
		//TODO: implement retry
		if shouldRetry(statusCode) {
		}
		sectionsSpinner.Fail(err.Error())
		return sections, errors.Wrap(err, "couldn't get the sections\n")
	}

	sectionsSpinner.Success(pterm.FgDefault.Sprint("DONE"))
	return sections, nil
}

//Gets the section by id
func GetSection(id string) (msftgraph.Section, error) {
	checkTokenPresented()

	var section, statusCode, err = root.api.GetSection(*root.token, id)

	if err != nil {
		//TODO: implement retry
		if shouldRetry(statusCode) {
		}
		return section, errors.Wrap(err, "couldn't get the section data\n")
	}
	return section, nil
}

//Get the lisft of pages of given section
func GetPages(section msftgraph.Section) ([]msftgraph.NotePage, error) {
	checkTokenPresented()

	pagesSpinner, _ := pterm.DefaultSpinner.Start("Getting note pages...")
	pagesSpinner.RemoveWhenDone = true
	var pages, statusCode, err = root.api.GetPages(*root.token, section)

	if err != nil {
		//TODO: Implement retry
		if shouldRetry(statusCode) {
		}
		pagesSpinner.Fail(err.Error())
		return pages, errors.Wrap(err, "couldn't get the note pages\n")
	}
	pagesSpinner.Success(pterm.FgDefault.Sprint("DONE"))

	return pages, nil
}

//Perfrom search in all your notebooks
func Search(phrase string) ([]msftgraph.NotePage, error) {
	checkTokenPresented()

	spinner, _ := pterm.DefaultSpinner.Start("Searching...")
	spinner.RemoveWhenDone = true
	searchResult, statusCode, err := root.api.SearchPage(*root.token, phrase)

	if err != nil {
		if shouldRetry(statusCode) {
		}
		spinner.Fail(err.Error())
		return nil, errors.Wrap(err, "failed to perform search\n")
	}
	spinner.Success(pterm.FgDefault.Sprint("DONE"))

	return searchResult, nil

}
func GetPageContent(notepage msftgraph.NotePage) ([]byte, error) {
	checkTokenPresented()

	spinner, _ := pterm.DefaultSpinner.Start("Getting the note content...")
	spinner.RemoveWhenDone = true
	content, statusCode, err := root.api.GetContent(*root.token, notepage)

	if err != nil {
		//TODO: Implement retry
		if shouldRetry(statusCode) {
		}
		spinner.Fail(err.Error())
		return content, errors.Wrap(err, "couldn't get the note content\n")
	}
	spinner.Success(pterm.FgDefault.Sprint("DONE"))

	return content, nil
}

//Save a note page using Onenote API
//Returns the link to the page.
func SaveNotePage(npage msftgraph.NotePage, remindAlias bool) (string, error) {
	checkTokenPresented()

	spinner, _ := pterm.DefaultSpinner.Start("Saving your note...")
	spinner.RemoveWhenDone = true
	link, statusCode, err := root.api.SaveNote(*root.token, npage)
	if err != nil {
		if shouldRetry(statusCode) {
			//TODO: implement retry
		}
		return "", errors.Wrap(err, "couldn't save the note page")
	}

	spinner.Success(pterm.FgDefault.Sprint("DONE"))
	printSuccessOutput(link, npage)

	aliases, err := GetAliases()
	if err != nil {
		return "", errors.Wrap(err, "couldn't get the alias list\n")
	}

	var hasAlias = hasAlias(npage.Section, aliases)
	if !hasAlias {
		answer, err := survey.AskAlias(npage.Section, aliases)
		if err != nil {
			return "", errors.Wrap(err, "couldn't ask the alias\n")
		}
		if answer != "" {
			err := SaveAlias(answer, *npage.Section.Notebook, npage.Section)
			if err != nil {
				return "", errors.Wrap(err, "couldn't save the alias\n")
			}
		}
	}

	//Print only if the alias didn't get created in this session.
	if hasAlias && remindAlias {
		printAliasReminder(npage.Section.Name)
	}

	return link, nil
}

func GetAliases() (*[]msftgraph.Alias, error) {

	var result []msftgraph.Alias
	keys, err := root.storage.GetKeys()
	if err != nil {
		return nil, errors.Wrap(err, "couldn't get the aliases\n")
	}

	var opts = getOptions()
	var hashset = make(map[string]bool, 0)
	for _, rk := range opts.ReservedKeys {
		hashset[rk] = true
	}

	if err != nil {
		return nil, errors.Wrap(err, "couldn't get the alias data\n")
	}

	for _, k := range *keys {
		if hashset[k] {
			continue
		}

		var a msftgraph.Alias
		root.storage.Get(k, &a)
		result = append(result, a)
	}

	return &result, nil
}

//Save the alias for a onenote section to use it later for quick save
func SaveAlias(name string, notebook msftgraph.Notebook, section msftgraph.Section) error {

	var isExist, err = GetAlias(name)
	if err != nil {
		return errors.Wrap(err, "couldn't check the alias if it already exists\n")
	}

	if isExist != nil {
		fmt.Println(style.Error(fmt.Sprintf("the alias %s already being used to identify the section %s", name, isExist.Section.Name)))
		return fmt.Errorf("the alias %s already exist", name)
	}

	err = root.storage.Set(name, msftgraph.Alias{
		Short:    name,
		Notebook: notebook,
		Section:  section})
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("couldn't save the alias %s to the local storage\n", name))
	}

	var msg = fmt.Sprintf("Alias '%s' has been saved.", name)
	fmt.Println(style.Success(msg))
	var infoMsg = "Now you can quickly add inline notes with the following command:"
	fmt.Println(style.Info(infoMsg))
	fmt.Println(fmt.Sprintf("$ nnote new -i \"inline text\" -a %s\n", name))
	return nil
}

//Get the details of given alias
//Returns nil if the alias does not found
func GetAlias(n string) (*msftgraph.Alias, error) {
	var alias msftgraph.Alias
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
		return errors.Wrap(err, "couldn't remove the alias\n")
	}

	var msg = fmt.Sprintf("The alias %s has been deleted.\n", a)
	fmt.Println(style.Success(msg))
	return nil
}

func hasAlias(section msftgraph.Section, aliasList *[]msftgraph.Alias) bool {
	for _, a := range *aliasList {
		if a.Section.ID == section.ID {
			return true
		}
	}
	return false
}

//Prints the outbut as a table upon saving a note
func printSuccessOutput(link string, npage msftgraph.NotePage) {

	var size = file.HumanizeSize(len([]byte(npage.Content)))

	var data = make([][]string, 2)
	data[0] = []string{"Status", "Notebook", "Section", "Title", "Size"}
	data[1] = []string{style.Success("OK"), npage.Section.Notebook.DisplayName, npage.Section.Name, npage.Title, size}
	pterm.DefaultTable.WithHasHeader().WithData(data).Render()
	fmt.Print("\n")
	fmt.Println(link)
	fmt.Print("\n")
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

		token, err := authentication.AuthenticateUser(root.oauthClient, opts, root.storage)
		if err != nil {
			var errWrapped = errors.Wrap(err, "couldn't check the token if it's presented\n")
			log.Fatal(errWrapped)
		}
		root.token = &token
	} else {
		//Check if the token has expired
		if root.token.IsExpired() {
			token, err := authentication.RefreshToken(root.oauthClient, opts, *root.token, root.storage)
			if err != nil {
				var errWrapped = errors.Wrap(err, "couldn't check the token if it's presented\n")
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
			fmt.Printf("$ nnote new -i \"inline text\" -a %s\n", a.Short)
		}
	}
	fmt.Println()
}

//Grab the token from the local storage upon startup
func init() {
	api := msftgraph.NewApi(&rest.RestClient{}, "https://graph.microsoft.com/v1.0")
	bitcask := &storage.Bitcask{}
	root = onenote{api: api, storage: bitcask}
	root.token = &oauthv2.OAuthToken{}
	root.oauthClient = oauthv2.NewOAuthClient(&rest.RestClient{})

	err := root.storage.Get(authentication.TOKEN_KEY, root.token)
	if err != nil {
		//token does not exist on the local storage
		root.token = nil
	}
}
