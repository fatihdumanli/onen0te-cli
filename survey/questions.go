package survey

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatihdumanli/cnote/pkg/onenote"
	"github.com/fatihdumanli/cnote/storage"
)

//TODO: add validations
var getNotebookOptions = func() []string {

	//TODO: Note that storage.CheckToken() gets called multiple times, code duplication and its an expensive operation...
	t, st := storage.CheckToken()
	if st != storage.Valid {
	}

	notebooks, err := onenote.GetNotebooks(t)
	if err != nil {
		panic(err)
	}
	var result []string

	for _, n := range notebooks {
		result = append(result, n.DisplayName)
	}
	return result
}

var getSectionOptions = func(n onenote.NotebookName) []string {
	sections, err := onenote.GetSections(n)

	if err != nil {
		panic(err)
	}

	var result []string

	for _, n := range sections {
		result = append(result, n.Name)
	}
	return result
}

var noteContentQuestion = &survey.Question{
	Name: "notecontent",
	Prompt: &survey.Multiline{
		Message: "Enter your note content.\n",
	},
}

var notebooksQuestion = &survey.Question{
	Name: "notebook",
	Prompt: &survey.Select{
		Message: "Select notebook\n",
		Options: getNotebookOptions(),
	},
}

var sectionQuestion = func(n onenote.NotebookName) *survey.Question {

	var opts = getSectionOptions(n)

	return &survey.Question{
		Name: "qsection",
		Prompt: &survey.Select{
			Message: fmt.Sprintf("Select a section in %s\n", n),
			Options: opts,
		},
	}

}

var setupQuestion = &survey.Question{
	Name: "setup",
	Prompt: &survey.Confirm{
		Message: "You haven't setup a Onenote account yet, would you like to setup one now?",
	},
}

//add title question
