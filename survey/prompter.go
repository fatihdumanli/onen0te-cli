package survey

import (
	"errors"
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	s "github.com/AlecAivazis/survey/v2"
	"github.com/fatihdumanli/cnote/config"
	"github.com/fatihdumanli/cnote/pkg/onenote"
)

//Shortcuts
type AppOptions = config.AppOptions
type Notebook = onenote.Notebook
type Section = onenote.Section

func AskNoteContent(opts AppOptions) (string, error) {

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

func AskSection(n Notebook) (Section, error) {
	var sections []string

	for _, s := range n.Sections {
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
		return n.Sections[0], err
	}

	if s, ok := findSection(n.Sections, func(x Section) bool {
		return x.Name == answer
	}); ok {
		return s, nil
	} else {
		return s, errors.New("An error has occured")
	}
}

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

func findSection(arr []Section, f func(x Section) bool) (Section, bool) {
	for _, s := range arr {
		if f(s) {
			return s, true
		}
	}
	return Section{}, false
}

func findNotebook(arr []Notebook, f func(x Notebook) bool) (Notebook, bool) {
	for _, s := range arr {
		if f(s) {
			return s, true
		}
	}
	return Notebook{}, false
}
