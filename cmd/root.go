package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/fatihdumanli/cnote/config"
	"github.com/fatihdumanli/cnote/pkg/oauthv2"
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
	var oauthParams = getOAuthParams()

	//check token here.
	t, st := storage.CheckToken()

	fmt.Fprintf(defaultOptions.Out, "token status is %d", st)

	if st == storage.DoesntExist {
		answer, err := survey.AskSetupAccount()
		if !answer || err != nil {
			os.Exit(1)
		}

		token, err := oauthv2.Authorize(oauthParams, defaultOptions.Out)
		if err != nil {
			log.Fatalf("An error occured while trying to authenticate you. %s", err.Error())
		}

		//save the token on local storage
		err = storage.StoreToken(token)
		if err != nil {
			log.Fatalf("An error occured while trying to save the token. %s", err.Error())
		}
	} else if st == storage.Expired {
		//need to refresh the token
		newToken, err := oauthv2.RefreshToken(oauthParams, t.RefreshToken)
		if err != nil {
			panic(err)
		} else {
			//save the token on local storage
			err = storage.StoreToken(newToken)
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

func getOAuthParams() oauthv2.OAuthParams {
	var msGraphOptions = config.GetMicrosoftGraphConfig()
	var oauthParams = oauthv2.OAuthParams{
		OAuthEndpoint:        "https://login.microsoftonline.com/common/oauth2/v2.0",
		RedirectUri:          "http://localhost:5992/oauthv2",
		Scope:                []string{"offline_access", "Notes.ReadWrite.All", "Notes.Create", "Notes.Read", "Notes.ReadWrite"},
		ClientId:             msGraphOptions.ClientId,
		RefreshTokenEndpoint: "https://login.microsoftonline.com/common/oauth2/v2.0/token",
	}
	return oauthParams
}
