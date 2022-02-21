package config

import (
	"github.com/fatihdumanli/onenote/pkg/oauthv2"
)

type AppOptions struct {
	//OAuth parameters that contains essential data for authentication operations.
	OAuthParams oauthv2.OAuthParams
	//Reserved keys for the app to run properly. These keys may not be used for any alias.
	ReservedKeys []string
}

type MicrosoftGraphConfig struct {
	ClientId    string
	TenantId    string
	RedirectUrl string
}
