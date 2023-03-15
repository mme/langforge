package tui

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/pterm/pterm"
)

func PromptYesNo(message string, defaultValue bool) (bool, error) {
	confirm := pterm.InteractiveConfirmPrinter{
		DefaultValue: defaultValue,
		DefaultText:  message,
		TextStyle:    &pterm.Style{pterm.Bold, pterm.BgDefault},
		ConfirmText:  "Yes",
		ConfirmStyle: &pterm.ThemeDefault.DefaultText,
		RejectText:   "No",
		RejectStyle:  &pterm.ThemeDefault.DefaultText,
		SuffixStyle:  &pterm.ThemeDefault.DefaultText,
	}

	return confirm.Show()
}

func PromptPassword(message string) (string, error) {
	var password string
	prompt := &survey.Password{
		Message: message,
	}
	err := survey.AskOne(prompt, &password)
	if err != nil {
		return "", err
	}
	return password, nil
}
