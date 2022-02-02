package survey

import (
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
	}

	var answer string
	if err := s.Ask([]*s.Question{noteContentQuestion}, &answer); err != nil {
		return "", err
	}
	return answer, nil
}

//TODO: a 'where' function might be useful here.
func AskNotebook(nlist []Notebook) (Notebook, error) {

	var notebookOptions []string
	var nindex int

	for _, n := range nlist {
		notebookOptions = append(notebookOptions, n.DisplayName)
	}

	var notebooksQuestion = &survey.Question{
		Name: "notebook",
		Prompt: &survey.Select{
			Message: "Select notebook\n",
			Options: notebookOptions,
		},
	}

	var answer string
	if err := s.Ask([]*s.Question{notebooksQuestion}, &answer); err != nil {
		panic(err)
	}

	for i, n := range nlist {
		if n.DisplayName == answer {
			nindex = i
		}
	}

	return nlist[nindex], nil
}

//TODO: a 'where' function might be useful here.
func AskSection(n Notebook) (Section, error) {

	var sections []string
	var sindex int

	for _, s := range n.Sections {
		sections = append(sections, s.Name)
	}

	var qsection = &survey.Question{
		Name: "qsection",
		Prompt: &survey.Select{
			Message: fmt.Sprintf("Select a section in %s\n", n.DisplayName),
			Options: sections,
		},
	}

	var answer string
	if err := s.Ask([]*s.Question{qsection}, &answer); err != nil {
		return n.Sections[0], err
	}

	for i, s := range n.Sections {
		if s.Name == answer {
			sindex = i
		}
	}

	return n.Sections[sindex], nil
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
