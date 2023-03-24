package python

import (
	"io/fs"
	"langforge/environment"

	"gopkg.in/yaml.v3"
)

var availableIntegrations []*environment.Integration

func init() {
	// Read the embedded YAML file
	data, err := fs.ReadFile(embeddedFS, "files/integrations.yaml")
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(data, &availableIntegrations)
	if err != nil {
		panic(err)
	}
}
