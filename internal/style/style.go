package style

import "github.com/pterm/pterm"

type Style func(s string) string

var Success Style = func(s string) string {
	return pterm.Success.Sprintf("%s", s)
}
var Section Style = func(s string) string {
	return pterm.FgMagenta.Sprintf("%s", s)
}

var Alias Style = func(s string) string {
	return pterm.NewRGB(255, 165, 96).Sprintf(s)
}

var Error Style = func(s string) string {
	return pterm.Error.Sprintf("&s", s)
}
