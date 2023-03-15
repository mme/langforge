package python

import (
	"langforge/environment"
	"langforge/system"
	"path/filepath"
)

type PythonHandler struct {
	integrations []*environment.Integration
}

func NewPythonHandler() environment.EnvironmentHandler {
	return &PythonHandler{
		integrations: environment.CopyIntegrations(availableIntegrations),
	}
}

func (h *PythonHandler) DetermineInstalledIntegrations() error {
	packages, err := GetInstalledPackages()
	if err != nil {
		return err
	}

	// Create a map of installed package names for fast lookups
	packagesMap := make(map[string]bool)
	for _, p := range packages {
		packagesMap[p.Name] = true
	}

	// Check each module to see if its packages are installed
	for _, m := range h.integrations {
		m.Installed = true

		for _, packageName := range m.Packages {
			if !packagesMap[packageName] {
				m.Installed = false
				break
			}
		}
	}

	return nil
}

func (h *PythonHandler) ExecuteChanges(dir string) error {

	err := SetJupyterEnvironmentVariables(dir)
	if err != nil {
		return err
	}

	install := []*environment.Integration{}
	uninstall := []*environment.Integration{}

	for _, integration := range h.integrations {
		if integration.Selected && !integration.Installed {
			install = append(install, integration)
		} else if !integration.Selected && integration.Installed {
			uninstall = append(uninstall, integration)
		}
	}

	pre := []string{}
	packages := []string{}
	post := []string{}
	uninstallPackages := []string{}

	for _, integration := range install {
		pre = append(pre, integration.PreInstallCommands...)
		packages = append(packages, integration.Packages...)
		post = append(post, integration.PostInstallCommands...)
	}

	for _, integration := range uninstall {
		uninstallPackages = append(uninstallPackages, integration.Packages...)
	}

	err = UninstallPackages(uninstallPackages)
	if err != nil {
		return err
	}

	err = system.ExecuteCommands(pre, dir)
	if err != nil {
		return err
	}

	err = InstallPackages(packages)
	if err != nil {
		return err
	}

	err = system.ExecuteCommands(post, dir)
	if err != nil {
		return err
	}

	for _, integration := range h.integrations {
		integration.Installed = integration.Selected
	}

	err = WriteRequirementsTxt(filepath.Join(dir, "requirements.txt"))
	if err != nil {
		panic(err)
	}

	isJupyterLabInstalled := false
	for _, integration := range h.integrations {
		if integration.Name == "jupyterlab" {
			isJupyterLabInstalled = integration.Installed
			break
		}
	}

	if isJupyterLabInstalled {
		err = EnableJupyterLabExtensions(dir)
		if err != nil {
			panic(err)
		}

		err = WriteIPythonStartupScripts(dir)
		if err != nil {
			panic(err)
		}
	}

	return nil
}

func (h *PythonHandler) InstalledIntegrationsApiKeys() []string {
	apiKeys := []string{}

	for _, integration := range h.integrations {
		if integration.Installed {
			apiKeys = append(apiKeys, integration.ApiKeys...)
		}
	}

	return apiKeys
}

func (h *PythonHandler) NamesOfIntegrationsToInstall() []string {
	names := []string{}

	for _, integration := range h.integrations {
		if integration.Selected && !integration.Installed {
			names = append(names, integration.Title)
		}
	}

	return names
}

func (h *PythonHandler) NamesOfIntegrationsToUninstall() []string {
	names := []string{}

	for _, integration := range h.integrations {
		if !integration.Selected && integration.Installed {
			names = append(names, integration.Title)
		}
	}

	return names
}

func (h *PythonHandler) GetIntegrations() []*environment.Integration {
	return h.integrations
}
