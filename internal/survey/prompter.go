package survey

import (
	"errors"
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	s "github.com/AlecAivazis/survey/v2"
	"github.com/fatihdumanli/cnote/internal/config"
	"github.com/fatihdumanli/cnote/internal/style"
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

	var setupQuestion = &survey.Question{
		Name: "setup",
		Prompt: &survey.Confirm{
			Message: "You haven't setup a Onenote account yet, would you like to setup one now?",
		},
	}

	var answer bool
	if err := s.Ask([]*s.Question{setupQuestion}, &answer); err != nil {
		return answer, err
	}
	return answer, nil
}

//Promps the user to get confirmation on creating alias to given notebook&section combination.
//Returns the answer and the error if any
func AskAlias(n NotebookName, sn SectionName) (string, error) {
	var answer string

	promtMsg := fmt.Sprintf("Enter an alias for %s (Press <Enter> to skip.)", style.Section(string(sn)))

	var aliasQuestion = &survey.Question{
		Name: "salias",
		Prompt: &survey.Input{
			Message: promtMsg,
		},
	}

	if err := s.Ask([]*s.Question{aliasQuestion}, &answer); err != nil {
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
