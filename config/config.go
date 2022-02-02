package config

import (
	"io"
	"os"
)

type AppOptions struct {
	Out io.Writer
	In  io.Reader
}

type MicrosoftGraphConfig struct {
	ClientId    string
	TenantId    string
	RedirectUrl string
}

var config = MicrosoftGraphConfig{
	ClientId:    "2124cbcc-943a-4a41-b8b2-efabbfc99b65",
	TenantId:    "31986ee9-8d0d-4a8e-8c6d-1d763b66d6c2",
	RedirectUrl: "http://localhost:5992/oauthv2",
}

func GetMicrosoftGraphConfig() MicrosoftGraphConfig {

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

	return config
}

func GetOptions() AppOptions {
	return AppOptions{
		In:  os.Stdin,
		Out: os.Stdout,
	}
}
