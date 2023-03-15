package tui

import (
	"fmt"
	"os"
	"strings"

	"atomicgo.dev/keyboard/keys"
	"github.com/pterm/pterm"
)

// run this on program start:
func init() {
	pterm.ThemeDefault.SecondaryStyle = pterm.Style{pterm.FgGray, pterm.BgDefault}
}

type SelectableItem interface {
	GetTitle() string
	IsSelected() bool
	SetSelected(bool)
}

func EditMultiSelect(promptMessage string, items []SelectableItem) ([]SelectableItem, error) {
	var options []string
	var defaultOptions []string
	for _, item := range items {
		options = append(options, item.GetTitle())
		if item.IsSelected() {
			defaultOptions = append(defaultOptions, item.GetTitle())
		}
	}

	multiSelect := pterm.InteractiveMultiselectPrinter{
		TextStyle:      &pterm.Style{pterm.Bold, pterm.BgDefault},
		DefaultText:    promptMessage,
		Options:        options,
		OptionStyle:    &pterm.ThemeDefault.DefaultText,
		DefaultOptions: defaultOptions,
		MaxHeight:      5,
		Selector:       ">",
		SelectorStyle:  &pterm.ThemeDefault.DefaultText,
		Filter:         true,
		KeySelect:      keys.Enter,
		KeyConfirm:     keys.Tab,
		Checkmark: &pterm.Checkmark{
			Checked:   pterm.DefaultInteractiveMultiselect.Checkmark.Checked,
			Unchecked: " ",
		},
	}

	// Create a MultiSelect prompt object
	selectedOptions, err := multiSelect.Show()
	if err != nil {
		return nil, err
	}

	// Update the Selected field of the items according to the user's selection
	for _, item := range items {
		selected := false
		for _, title := range selectedOptions {
			if item.GetTitle() == title {
				selected = true
				break
			}
		}
		item.SetSelected(selected)
	}

	return items, nil
}

func EditSelect(message string, options []string, showDone bool) (int, error) {
	doneOption := Bold("✓ Done")
	if showDone {
		options = append([]string{doneOption}, options...)
	}

	selectPrompt := pterm.InteractiveSelectPrinter{
		TextStyle:     &pterm.Style{pterm.Bold, pterm.BgDefault},
		DefaultText:   message,
		Options:       options,
		OptionStyle:   &pterm.ThemeDefault.DefaultText,
		DefaultOption: "",
		MaxHeight:     5,
		Selector:      ">",
		SelectorStyle: &pterm.ThemeDefault.DefaultText,
	}

	choice, err := selectPrompt.Show()
	if err != nil {
		return -1, err
	}

	if showDone && choice == doneOption {
		return -1, nil
	}

	// Find the index of the selected option
	for i, option := range options {
		if choice == option {
			if showDone {
				return i - 1, nil // Subtract 1 to account for the "Done" option
			}
			return i, nil
		}
	}

	// This should never happen, but just in case
	return -1, fmt.Errorf("selected option not found in the options list")
}

func EditApiKeys(apiKeys []string, env map[string]string) (map[string]string, error) {
	// Ensure all apiKeys are present in env
	for _, apiKey := range apiKeys {
		if _, ok := env[apiKey]; !ok {
			// If the environment variable is set in the system, use its value; otherwise, use the empty string
			envValue, ok := os.LookupEnv(apiKey)
			if !ok {
				envValue = ""
			}
			env[apiKey] = envValue
		}
	}

	const setIndicator = " (set)"

	for {
		options := make([]string, len(apiKeys))
		copy(options, apiKeys)

		for i, apiKey := range apiKeys {
			if env[apiKey] != "" {
				options[i] = fmt.Sprintf("%s%s", apiKey, setIndicator)
			} else {
				options[i] = apiKey
			}
		}

		choiceIndex, err := EditSelect("Choose an API key to edit or select Done to finish", options, true)
		if err != nil {
			return nil, err
		}

		if choiceIndex == -1 {
			break
		}

		selectedKey := strings.TrimSuffix(strings.TrimPrefix(options[choiceIndex], setIndicator), " ")
		fmt.Printf("Editing API key: %s\n", selectedKey)

		newValue, err := PromptPassword("Enter the new value:")
		if err != nil {
			return nil, err
		}
		EmptyLine()

		env[selectedKey] = newValue
	}

	return env, nil
}
