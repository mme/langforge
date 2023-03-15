package python

import (
	"embed"
	"io/fs"
	"langforge/environment"

	"gopkg.in/yaml.v3"
)

//go:embed integrations.yaml
var embeddedFS embed.FS
var availableIntegrations []*environment.Integration

func init() {
	// Read the embedded YAML file
	data, err := fs.ReadFile(embeddedFS, "integrations.yaml")
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(data, &availableIntegrations)
	if err != nil {
		panic(err)
	}
}
