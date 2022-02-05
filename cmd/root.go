package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/fatihdumanli/cnote"
	"github.com/fatihdumanli/cnote/internal/authentication"
	"github.com/fatihdumanli/cnote/internal/storage"
	"github.com/fatihdumanli/cnote/internal/survey"
	"github.com/fatihdumanli/cnote/pkg/oauthv2"
	"github.com/fatihdumanli/cnote/pkg/onenote"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Long: "Take notes on your Onenote notebooks from terminal",
	Run:  runRoot,
	Use:  "cnote [command] [args] [flags]",
}

type Notebook onenote.Notebook
type Section onenote.Section
type NotebookName = onenote.NotebookName
type SectionName = onenote.SectionName

//The function gets executed once the application starts without any commands/arguments.
func runRoot(c *cobra.Command, args []string) {

	var cn = cnote.Cnote{}

	t, err := getValidAccount()
	if err != nil {
		panic(err)
	}

	_, err = survey.AskNoteContent()
	if err != nil {
		panic(err)
	}

	notebooks, err := cn.GetNotebooks(t)
	fmt.Println("Getting your notebooks... This might take a while...")

	if err != nil {
		panic(err)
	}

	n, err := survey.AskNotebook(notebooks)
	sections, err := cn.GetSections(t, n)
	if err != nil {
		panic(err)
	}

	//iterating over the sections and adding them to the notebook struct
	for _, s := range sections {
		n.Sections = append(n.Sections, s)
	}

	section, err := survey.AskSection(n)

	var defaultOptions = cnote.GetOptions()
	fmt.Fprintf(defaultOptions.Out, "Your note has saved to the notebook %s and the section %s",
		n.DisplayName, section.Name)

	a, err := survey.AskAlias(NotebookName(n.DisplayName), SectionName(section.Name))
	if err != nil {
		panic(err)
	}

	if a != "" {
		//TODO: save the alias
		fmt.Println("we need to save the alias.")
	}

}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stdout, err)
		os.Exit(1)
	}
}

func getValidAccount() (oauthv2.OAuthToken, error) {
	//TODO: I feel like we shouldn't expose GetOptions() out of the cnote packge.
	var defaultOptions = cnote.GetOptions()
	var oauthParams = cnote.GetOAuthParams()

	t, st := storage.CheckToken()
	if st == authentication.DoesntExist {
		setupAccount(oauthParams, defaultOptions.Out)
	} else if st == authentication.Expired {
		//Need to refresh the token
		refreshToken(oauthParams, t)
	}

	return t, nil

}

//Refresh the token and save the new one on local storage.
func refreshToken(oauthParams oauthv2.OAuthParams, t oauthv2.OAuthToken) (oauthv2.OAuthToken, error) {

	newToken, err := oauthv2.RefreshToken(oauthParams, t.RefreshToken)
	if err != nil {
		panic(err)
	} else {
		//Save the token on local storage
		err = storage.StoreToken(newToken)
		if err != nil {
			return t, nil
		}
	}
	return newToken, nil
}

//Setup a onenote account for the very first time.
func setupAccount(oauthParams oauthv2.OAuthParams, out io.Writer) {
	answer, err := survey.AskSetupAccount()
	if !answer || err != nil {
		os.Exit(1)
	}

	//If the user confirms to setup an account now we trigger the authentication process.
	token, err := oauthv2.Authorize(oauthParams, out)
	if err != nil {
		log.Fatalf("An error occured while trying to authenticate you. %s", err.Error())
	}

	//Save the token on local storage
	err = storage.StoreToken(token)
	if err != nil {
		log.Fatalf("An error occured while trying to save the token. %s", err.Error())
	}

}
