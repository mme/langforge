package environment

type EnvironmentHandler interface {
	DetermineInstalledIntegrations() error
	NamesOfIntegrationsToInstall() []string
	NamesOfIntegrationsToUninstall() []string
	ExecuteChanges() error
	InstalledIntegrationsApiKeys() []string
	GetIntegrations() []*Integration
}
