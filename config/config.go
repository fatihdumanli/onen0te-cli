package config

import (
	"io"
	"os"
)

type AppOptions struct {
	out io.Writer
	in  io.Reader
}

type MicrosoftGraphConfig struct {
	ClientId    string
	TenantId    string
	RedirectUrl string
}

func GetMicrosoftGraphConfig() MicrosoftGraphConfig {
	return MicrosoftGraphConfig{
		ClientId:    "2124cbcc-943a-4a41-b8b2-efabbfc99b65",
		TenantId:    "31986ee9-8d0d-4a8e-8c6d-1d763b66d6c2",
		RedirectUrl: "http://localhost:5992/oauthv2",
	}
}

func GetOptions() AppOptions {
	return AppOptions{
		in:  os.Stdin,
		out: os.Stdout,
	}
}
