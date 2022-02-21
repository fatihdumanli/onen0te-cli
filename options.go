package onenote

import (
	"os"

	"github.com/fatihdumanli/onenote/internal/authentication"
	"github.com/fatihdumanli/onenote/internal/config"
	"github.com/fatihdumanli/onenote/pkg/oauthv2"
	_ "github.com/joho/godotenv/autoload"
)

func getOptions() config.AppOptions {

	var oauthEndpoint = os.Getenv("OAUTH_ENDPOINT")
	var redirectUri = os.Getenv("REDIRECT_URL")
	var clientId = os.Getenv("CLIENT_ID")
	var tenantId = os.Getenv("TENANT_ID")

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
