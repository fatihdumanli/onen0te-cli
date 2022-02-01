package survey

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatihdumanli/cnote/pkg/onenote"
)

//TODO: add validations
var getNotebookOptions = func() []string {
	notebooks, err := onenote.GetNotebooks()
	//TODO: handle
	if err != nil {
	}
	var result []string

	for _, n := range notebooks {
		result = append(result, n.Name)
	}
	return result
}

var getSectionOptions = func(n onenote.NotebookName) []string {
	sections, err := onenote.GetSections(n)

	//TODO: handle
	if err != nil {
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

//add title question
