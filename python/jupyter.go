package python

import (
	"io/fs"
	"os"
	"path/filepath"
)

func EnableJupyterLabExtensions(dir string) error {
	configDir := filepath.Join(dir, ".jupyter", "@jupyterlab", "extensionmanager-extension")
	err := os.MkdirAll(configDir, 0755)
	if err != nil {
		return err
	}
	configFile := filepath.Join(configDir, "plugin.jupyterlab-settings")

	// check if file exists
	if _, err = os.Stat(configFile); err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	} else {
		return nil
	}

	configFileContents := []byte(`{"disclaimed": true}`)

	err = os.WriteFile(configFile, configFileContents, 0644)
	if err != nil {
		return err
	}
	return nil
}

func WriteIPythonStartupScripts(dir string) error {
	err := writeIPythonStartupScript(dir, "00-dotenv.py")
	if err != nil {
		return err
	}

	err = writeIPythonStartupScript(dir, "10-extension-support.py")
	if err != nil {
		return err
	}

	return nil
}

func writeIPythonStartupScript(dir string, filename string) error {
	iPythonStartupDir := filepath.Join(dir, ".ipython", "profile_default", "startup")
	err := os.MkdirAll(iPythonStartupDir, 0755)
	if err != nil {
		return err
	}

	iPythonStartupFile := filepath.Join(iPythonStartupDir, filename)
	iPythonStartupFileContents, err := fs.ReadFile(embeddedFS, "files/startup/"+filename)
	if err != nil {
		panic(err)
	}

	err = os.WriteFile(iPythonStartupFile, iPythonStartupFileContents, 0644)
	if err != nil {
		return err
	}

	return nil
}

func SetJupyterEnvironmentVariables(dir string) error {
	iPythonDir := filepath.Join(dir, ".ipython")
	err := os.Setenv("IPYTHONDIR", iPythonDir)
	if err != nil {
		return err
	}

	jupyterlabSettingsDir := filepath.Join(dir, ".jupyter")
	err = os.Setenv("JUPYTERLAB_SETTINGS_DIR", jupyterlabSettingsDir)
	if err != nil {
		return err
	}

	jupyterConfigDir := filepath.Join(dir, ".jupyter")
	err = os.Setenv("JUPYTER_CONFIG_DIR", jupyterConfigDir)
	if err != nil {
		return err
	}
	return nil
}
