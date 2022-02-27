package onenote

import (
	"github.com/fatihdumanli/onen0te-cli/internal/authentication"
	"github.com/fatihdumanli/onen0te-cli/internal/config"
	"github.com/fatihdumanli/onen0te-cli/pkg/oauthv2"
)

func getOptions() config.AppOptions {

	var oauthEndpoint = "https://login.microsoftonline.com/common/oauth2/v2.0"
	var redirectUri = "http://localhost:5992/oauthv2"
	var clientId = "2124cbcc-943a-4a41-b8b2-efabbfc99b65"
	var tenantId = "31986ee9-8d0d-4a8e-8c6d-1d763b66d6c2"

	var msGraphOptions = config.MicrosoftGraphConfig{
		ClientId:    clientId,
		TenantId:    tenantId,
		RedirectUrl: redirectUri,
	}

	var oauthParams = oauthv2.OAuthParams{
		OAuthEndpoint: oauthEndpoint,
		RedirectUri:   redirectUri,
		Scope:         []string{"offline_access", "Notes.ReadWrite.All", "Notes.Create", "Notes.Read", "Notes.ReadWrite"},
		ClientId:      msGraphOptions.ClientId,
	}

	return config.AppOptions{
		OAuthParams:  oauthParams,
		ReservedKeys: []string{authentication.TOKEN_KEY},
	}
}
