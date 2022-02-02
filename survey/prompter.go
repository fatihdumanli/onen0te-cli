package survey

import (
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

//TODO: we might need to return a Notebook struct from here
func AskNotebook(opts AppOptions) (onenote.NotebookName, error) {
	var answer string

	if err := s.Ask([]*s.Question{notebooksQuestion}, &answer); err != nil {
		panic(err)
		return "", err
	}

	return onenote.NotebookName(answer), nil
}

func AskSection(opts AppOptions, n onenote.NotebookName) (string, error) {

	var answer string
	var q = sectionQuestion(n)

	if err := s.Ask([]*s.Question{q}, &answer); err != nil {
		return "", err
	}

	return answer, nil
}

func AskSetupAccount() (bool, error) {
	var answer bool
	if err := s.Ask([]*s.Question{setupQuestion}, &answer); err != nil {
		return answer, err
	}

	return answer, nil
}
