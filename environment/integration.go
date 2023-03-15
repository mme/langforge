package environment

import "langforge/tui"

type Integration struct {
	Name                string   `yaml:"name"`
	Title               string   `yaml:"title"`
	Selected            bool     `yaml:"selected"`
	Installed           bool     `yaml:"installed"`
	Packages            []string `yaml:"packages"`
	ApiKeys             []string `yaml:"apiKeys"`
	PreInstallCommands  []string `yaml:"preInstallCommands"`
	PostInstallCommands []string `yaml:"postInstallCommands"`
}

func (i *Integration) GetTitle() string {
	return i.Title
}

func (i *Integration) IsSelected() bool {
	return i.Selected
}

func (i *Integration) SetSelected(selected bool) {
	i.Selected = selected
}

func (i *Integration) Copy() *Integration {
	return &Integration{
		Name:                i.Name,
		Title:               i.Title,
		Selected:            i.Selected,
		Installed:           i.Installed,
		Packages:            i.Packages,
		ApiKeys:             i.ApiKeys,
		PreInstallCommands:  i.PreInstallCommands,
		PostInstallCommands: i.PostInstallCommands,
	}
}

func CopyIntegrations(integrations []*Integration) []*Integration {
	result := []*Integration{}
	for _, integration := range integrations {
		result = append(result, integration.Copy())
	}
	return result
}

func IntegrationsToSelectableItems(integrations []*Integration) []tui.SelectableItem {
	selectableItems := make([]tui.SelectableItem, len(integrations))
	for i, integration := range integrations {
		selectableItems[i] = integration
	}
	return selectableItems
}
