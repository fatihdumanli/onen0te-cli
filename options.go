package onenote

import (
	"os"

	"github.com/fatihdumanli/onenote/internal/authentication"
	"github.com/fatihdumanli/onenote/internal/config"
	"github.com/fatihdumanli/onenote/pkg/oauthv2"
)

var c = config.MicrosoftGraphConfig{
	ClientId:    "2124cbcc-943a-4a41-b8b2-efabbfc99b65",
	TenantId:    "31986ee9-8d0d-4a8e-8c6d-1d763b66d6c2",
	RedirectUrl: "http://localhost:5992/oauthv2",
}

//TODO: Define a standard for colors here.
//For aliases: CyanBg-RedFg
//For sections: Red
//...
func getMicrosoftGraphConfig() config.MicrosoftGraphConfig {

	//NOTE
	//if we instantiate the config struct here,
	//that means we're instantiatng a new struct each time this func gets called.
	//and this is not good.
	//return MicrosoftGraphConfig{
	//	ClientId:    "2124cbcc-943a-4a41-b8b2-efabbfc99b65",
	//	TenantId:    "31986ee9-8d0d-4a8e-8c6d-1d763b66d6c2",
	//	RedirectUrl: "http://localhost:5992/oauthv2",
	//}

	//and if we return a pointer of MicrosoftGraphConfig
	//it's dangerous bc we could end up with a mutated ms graph config
	//which could lead the app a subtle bug
	//return &config

	return c
}

//TODO: Notice that this method gets called everywhere in the app
//We might need to come up with a DI trick.
func getOptions() config.AppOptions {
	var msGraphOptions = getMicrosoftGraphConfig()

	var oauthParams = oauthv2.OAuthParams{
		OAuthEndpoint:        "https://login.microsoftonline.com/common/oauth2/v2.0",
		RedirectUri:          "http://localhost:5992/oauthv2",
		Scope:                []string{"offline_access", "Notes.ReadWrite.All", "Notes.Create", "Notes.Read", "Notes.ReadWrite"},
		ClientId:             msGraphOptions.ClientId,
		RefreshTokenEndpoint: "https://login.microsoftonline.com/common/oauth2/v2.0/token",
	}

	return config.AppOptions{
		In:           os.Stdin,
		Out:          os.Stdout,
		OAuthParams:  oauthParams,
		ReservedKeys: []string{authentication.TOKEN_KEY},
	}
}
