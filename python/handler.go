package python

import (
	"langforge/environment"
	"langforge/system"
	"path/filepath"
)

type PythonHandler struct {
	integrations []*environment.Integration
	dir          string
}

func NewPythonHandler(dir string) environment.EnvironmentHandler {
	return &PythonHandler{
		integrations: environment.CopyIntegrations(availableIntegrations),
		dir:          dir,
	}
}

func (h *PythonHandler) DetermineInstalledIntegrations() error {
	packages, err := GetInstalledPackages()
	if err != nil {
		return err
	}

	env, err := system.GetEnv(h.dir)
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

		if len(m.Packages) == 0 || m.Name == "gooseai" {
			// If the module has no packages, check if the API key is set
			for _, apiKey := range m.ApiKeys {
				if _, ok := env[apiKey]; !ok {
					m.Installed = false
					break
				}
			}
		} else {
			for _, packageName := range m.Packages {
				if !packagesMap[packageName] {
					m.Installed = false
					break
				}
			}

			if m.Installed {
				m.Selected = true
			}
		}

	}

	return nil
}

func (h *PythonHandler) ExecuteChanges() error {

	err := SetJupyterEnvironmentVariables(h.dir)
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
	removeApiKeys := []string{}

	for _, integration := range install {
		pre = append(pre, integration.PreInstallCommands...)
		packages = append(packages, integration.Packages...)
		post = append(post, integration.PostInstallCommands...)
	}

	for _, integration := range uninstall {
		uninstallPackages = append(uninstallPackages, integration.Packages...)
		removeApiKeys = append(removeApiKeys, integration.ApiKeys...)
	}

	err = UninstallPackages(uninstallPackages)
	if err != nil {
		return err
	}

	err = system.ExecuteCommands(pre, h.dir)
	if err != nil {
		return err
	}

	err = InstallPackages(packages)
	if err != nil {
		return err
	}

	err = system.ExecuteCommands(post, h.dir)
	if err != nil {
		return err
	}

	for _, integration := range h.integrations {
		integration.Installed = integration.Selected
	}

	err = WriteRequirementsTxt(filepath.Join(h.dir, "requirements.txt"))
	if err != nil {
		panic(err)
	}

	wasJupyterLabInstalled := false
	for _, integration := range install {
		if integration.Name == "jupyterlab" {
			wasJupyterLabInstalled = true
			break
		}
	}

	if len(removeApiKeys) > 0 {
		env, err := system.GetEnv(h.dir)
		if err != nil {
			return err
		}

		for _, apiKey := range removeApiKeys {
			delete(env, apiKey)
		}

		err = system.WriteEnv(h.dir, env)
		if err != nil {
			return err
		}
	}

	if wasJupyterLabInstalled {
		err = EnableJupyterLabExtensions(h.dir)
		if err != nil {
			panic(err)
		}

		err = InstallLangforgeJupyterExtension(h.dir)
		if err != nil {
			panic(err)
		}

		err = WriteIPythonStartupScripts(h.dir)
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
