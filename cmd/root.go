package cmd

import (
	"fmt"
	"os"

	"github.com/fatihdumanli/cnote/config"
	"github.com/fatihdumanli/cnote/pkg/oauthv2"
	"github.com/fatihdumanli/cnote/pkg/onenote"
	"github.com/fatihdumanli/cnote/storage"
	"github.com/fatihdumanli/cnote/survey"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Long:                  "Take notes on your Onenote notebooks from terminal",
	Run:                   runRoot,
	Use:                   "cnote [command] [args] [flags]",
	DisableFlagsInUseLine: true,
}

func runRoot(c *cobra.Command, args []string) {

	var defaultOptions = config.GetOptions()
	var msGraphOptions = config.GetMicrosoftGraphConfig()

	//check token here.
	tokenStatus := storage.CheckToken()

	if tokenStatus == storage.DoesntExist {
		answer, err := survey.AskSetupAccount()
		if !answer || err != nil {
			os.Exit(1)
		}

		var p = oauthv2.OAuthParams{
			OAuthEndpoint: "https://login.microsoftonline.com/common/oauth2/v2.0",
			RedirectUri:   "http://localhost:5992/oauthv2",
			Scope:         []string{"offline_access", "Notes.ReadWrite.All", "Notes.Create", "Notes.Read", "Notes.ReadWrite"},
			ClientId:      msGraphOptions.ClientId,
		}

		authResult := onenote.Authorize(p, defaultOptions.Out)

		if authResult == onenote.Successful {
		}
	}

	noteContent, err := survey.AskNoteContent(defaultOptions)
	if err != nil {
		panic(err)
	}

	notebook, err := survey.AskNotebook(defaultOptions)
	section, err := survey.AskSection(defaultOptions, notebook)

	_ = notebook
	_ = noteContent
	_ = section
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stdout, err)
		os.Exit(1)
	}
}
