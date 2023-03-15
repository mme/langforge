package environment

type EnvironmentHandler interface {
	DetermineInstalledIntegrations() error
	NamesOfIntegrationsToInstall() []string
	NamesOfIntegrationsToUninstall() []string
	ExecuteChanges(dir string) error
	InstalledIntegrationsApiKeys() []string
	GetIntegrations() []*Integration
}
