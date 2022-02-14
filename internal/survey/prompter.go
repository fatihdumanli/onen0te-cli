package survey

import (
	"fmt"

	errors "github.com/pkg/errors"

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
		return "", errors.Wrap(err, "error while starting note content survey")
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
		return Notebook{}, errors.Wrap(err, "couldn't get notebook answer")
	}

	if n, ok := findNotebook(nlist, func(nx Notebook) bool {
		return nx.DisplayName == answer
	}); ok {
		//if notebook exist
		return n, nil
	} else {
		return n, fmt.Errorf("the notebook %s does not exist", answer)
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
		return slist[0], errors.Wrap(err, "couldn't start section survey")
	}

	if s, ok := findSection(slist, func(x Section) bool {
		return x.Name == answer
	}); ok {
		//section is found
		return s, nil
	} else {
		//section does not exist
		return s, fmt.Errorf("section %s could not be found", answer)
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

outer:
	//Prompt the user again if the alias already exist
	for {
		a, err := askSinglelineFreeText(promtMsg)
		if err != nil {
			return "", errors.Wrap(err, "couldn't ask alias")
		}

		for _, alias := range *aliaslist {
			//There's already an alias with the same name
			if alias.Short == a {
				var errorMsg = fmt.Sprintf("%s being used for another section", a)
				fmt.Println(style.Error(errorMsg))
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
		return answer, errors.Wrap(err, "couldn't ask confirmation question")
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
		return "", errors.Wrap(err, "couldn't ask single line text question")
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
