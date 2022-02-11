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
	return pterm.Error.Sprintf("%s", s)
}

var Info Style = func(s string) string {
	pterm.Info.Scope = pterm.Scope{
		Style: pterm.NewStyle(pterm.FgDefault),
	}
	return pterm.Info.Sprintf("%s", s)
}

var Reminder Style = func(s string) string {
	pterm.Info.Prefix = pterm.Prefix{
		Text:  "💡 REMINDER",
		Style: pterm.NewStyle(pterm.BgYellow, pterm.FgBlack),
	}

	return pterm.Info.Sprintf(pterm.NewRGB(255, 217, 30).Sprintf(s))
}

var Warning Style = func(s string) string {
	return pterm.Warning.Sprintf("%s", s)
}
