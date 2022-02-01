package survey

import (
	"fmt"

	s "github.com/AlecAivazis/survey/v2"
	"github.com/fatihdumanli/cnote/config"
	"github.com/fatihdumanli/cnote/pkg/onenote"
)

type AppOptions = config.AppOptions

func AskNoteContent(opts AppOptions) (string, error) {

	var answer string

	if err := s.Ask([]*s.Question{noteContentQuestion}, &answer); err != nil {
		return "", err
	}

	return answer, nil
}

//TODO: we might need to return a AskNotebook struct from here
func AskNotebook(opts AppOptions) (onenote.NotebookName, error) {
	var answer string

	if err := s.Ask([]*s.Question{notebooksQuestion}, &answer); err != nil {
		panic(err)
		return "", err
	}

	return onenote.NotebookName(answer), nil
}

func AskSection(opts AppOptions, n onenote.NotebookName) (string, error) {

	fmt.Printf("selected notebook is %s", n)

	var answer string
	var q = sectionQuestion(n)

	if err := s.Ask([]*s.Question{q}, &answer); err != nil {
		return "", err
	}

	return answer, nil
}
