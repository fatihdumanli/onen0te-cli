package survey

import (
	"errors"
	"fmt"
	"log"

	"github.com/AlecAivazis/survey/v2"
	s "github.com/AlecAivazis/survey/v2"
	"github.com/fatihdumanli/cnote/internal/config"
	"github.com/fatihdumanli/cnote/pkg/onenote"
)

//Friendly names
type AppOptions = config.AppOptions
type Notebook = onenote.Notebook
type Section = onenote.Section
type SectionName = onenote.SectionName
type NotebookName = onenote.NotebookName

func AskNoteContent() (string, error) {

	var noteContentQuestion = &survey.Question{
		Name: "notecontent",
		Prompt: &survey.Multiline{
			Message: "Enter your note content.\n",
		},
		Validate: survey.Required,
	}

	var answer string
	if err := s.Ask([]*s.Question{noteContentQuestion}, &answer); err != nil {
		return "", err
	}
	return answer, nil
}

//Ask for the notebook to save the note.
//TODO: Add error handling and output if an error occured.
func AskNotebook(nlist []Notebook) (Notebook, error) {

	var notebookOptions []string
	for _, n := range nlist {
		notebookOptions = append(notebookOptions, n.DisplayName)
	}

	var notebooksQuestion = &survey.Question{
		Name: "notebook",
		Prompt: &survey.Select{
			Message: "Select notebook\n",
			Options: notebookOptions,
		},
		Validate: survey.Required,
	}

	var answer string
	if err := s.Ask([]*s.Question{notebooksQuestion}, &answer); err != nil {
		panic(err)
	}

	if n, ok := findNotebook(nlist, func(nx Notebook) bool {
		return nx.DisplayName == answer
	}); ok {
		return n, nil
	} else {
		return n, errors.New("An error has occured")
	}

}

//Ask for the section in which the note(onenote page) will be created.
func AskSection(n Notebook, slist []Section) (Section, error) {

	var sections []string
	for _, s := range slist {
		sections = append(sections, s.Name)
	}

	var qsection = &survey.Question{
		Name: "qsection",
		Prompt: &survey.Select{
			Message:  fmt.Sprintf("Select a section in %s\n", n.DisplayName),
			Options:  sections,
			PageSize: 100,
		},
		Validate: survey.Required,
	}

	var answer string
	if err := s.Ask([]*s.Question{qsection}, &answer); err != nil {
		return slist[0], err
	}

	if s, ok := findSection(slist, func(x Section) bool {
		return x.Name == answer
	}); ok {
		return s, nil
	} else {
		return s, errors.New("An error has occured")
	}
}

//In case there's no account have been set yet, prompt the user to ask whether create one at that moment.
func AskSetupAccount() (bool, error) {
	return askConfirmation("You haven't setup a Onenote account yet, would you like to setup one now?")
}

//Promps the user to get confirmation on creating alias to given notebook&section combination.
//Returns the answer and the error if any
//Ask alias to make it easy to create a note within the section
func AskAlias(section onenote.Section, aliaslist *[]onenote.Alias) (string, error) {
	promtMsg := fmt.Sprintf("Enter an alias for %s (Press <Enter> to skip.)", section.Name)
	//Ask for an alias as long as there's already an existing one
outer:
	for {
		a, err := askSinglelineFreeText(promtMsg)
		if err != nil {
			log.Fatal(err)
			return "", err
		}

		fmt.Println("checking if tthere's an already an alias", a)
		for _, alias := range *aliaslist {
			//There's already an alias with the same name
			if alias.Short == a {
				continue outer
			}
		}
		return a, nil
	}

}

//Asks title if the user is about to create an untitled note with the command
//"cnote new <input-file>" (without the -title flag.)
func AskForgottenTitle() (string, error) {
	var msg = "Did you forget to pass the title flag? Set it now: (Press <Enter> to skip.)"
	return askSinglelineFreeText(msg)
}

//Ask title when creating the note without any flags.
func AskTitle() (string, error) {
	var msg = "Enter title: (Press <Enter> to skip.)"
	return askSinglelineFreeText(msg)
}

//Prepares a confirmation msg (y/N)
func askConfirmation(msg string) (bool, error) {
	var q = &survey.Question{
		Name: "q",
		Prompt: &survey.Confirm{
			Message: msg,
		},
	}
	var answer bool
	if err := s.Ask([]*s.Question{q}, &answer); err != nil {
		return answer, err
	}
	return answer, nil
}

//Prepares a question that can be answered with a free text.
func askSinglelineFreeText(msg string) (string, error) {
	var answer string
	var q = &survey.Question{
		Name: "qtitle",
		Prompt: &survey.Input{
			Message: msg,
		},
	}
	if err := s.Ask([]*s.Question{q}, &answer); err != nil {
		return "", err
	}
	return answer, nil
}

//Iterates over the sections and returns the one satisfies the given condition
func findSection(arr []Section, f func(s Section) bool) (Section, bool) {
	for _, x := range arr {
		if f(x) {
			return x, true
		}
	}
	return Section{}, false
}

//Iterates over the notebooks and returns the one satisfies the given condition
func findNotebook(arr []Notebook, f func(n Notebook) bool) (Notebook, bool) {
	for _, x := range arr {
		if f(x) {
			return x, true
		}
	}
	return Notebook{}, false
}
