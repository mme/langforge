package python

import "embed"

//go:embed files/integrations.yaml
//go:embed files/startup/00-dotenv.py
//go:embed files/startup/10-extension-support.py
var embeddedFS embed.FS
