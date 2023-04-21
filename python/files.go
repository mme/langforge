package python

import (
	"embed"
	"io/fs"
)

//go:embed files/integrations.yaml
//go:embed files/startup/00-dotenv.py
//go:embed files/startup/10-extension-support.py
//go:embed files/startup/20-utilities.py
//go:embed files/server.py
//go:embed files/langforge-0.1.0-py3-none-any.whl
var embeddedFS embed.FS

func ServerPy() ([]byte, error) {
	return fs.ReadFile(embeddedFS, "files/server.py")
}
