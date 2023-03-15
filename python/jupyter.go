package python

import (
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
	iPythonStartupDir := filepath.Join(dir, ".ipython", "profile_default", "startup")
	err := os.MkdirAll(iPythonStartupDir, 0755)
	if err != nil {
		return err
	}

	iPythonStartupFile := filepath.Join(iPythonStartupDir, "00-dotenv.py")
	iPythonStartupFileContents := []byte("from dotenv import load_dotenv\nload_dotenv()\n")

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
