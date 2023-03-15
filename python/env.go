package python

import (
	"fmt"
	"langforge/system"
	"os"
	"os/exec"
	"path/filepath"
)

// ActivateEnvironment checks if the specified environment exists, and activates it if it does.
// If a directory is provided, the environment is looked for within that directory. Otherwise,
// it is looked for in the default environment directory.
func ActivateEnvironment(envName string, envDir ...string) error {
	var activateScript string
	if len(envDir) > 0 {
		activateScript = fmt.Sprintf("%s/%s/bin/activate", envDir[0], envName)
	} else {
		activateScript = fmt.Sprintf("%s/bin/activate", envName)
	}

	// Check if the environment exists
	if _, err := os.Stat(activateScript); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("environment %q not found", envName)
		}
		return fmt.Errorf("failed to stat environment %q: %v", envName, err)
	}

	// Activate the environment
	if err := system.ShellSource(activateScript); err != nil {
		return fmt.Errorf("failed to activate environment %q: %v", envName, err)
	}

	return nil
}

// PythonCreateVirtualEnv creates a new Python virtual environment using the `venv` module.
// It takes the name of the environment and an optional directory containing the
// virtual environment as arguments, and returns an error if the environment creation fails.
func CreateVirtualEnv(envName string, envDir ...string) error {
	// Determine the path where the virtual environment will be created
	envPath := envName
	if len(envDir) > 0 {
		envPath = filepath.Join(envDir[0], envName)
	}
	envAbsPath, err := filepath.Abs(envPath)
	if err != nil {
		return err
	}

	// Find the path to the Python interpreter
	pythonPath, err := system.FindPython()
	if err != nil {
		return err
	}

	// Create the virtual environment using the venv module
	cmd := exec.Command(pythonPath, "-m", "venv", "--clear", "--symlinks", envAbsPath)

	// Set the command's Stdout and Stderr fields to os.Stdout and os.Stderr respectively
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Start the command
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func WriteRequirementsTxt(path string) error {
	pipPath, err := system.FindPip()
	if err != nil {
		return err
	}

	cmd := exec.Command(pipPath, "freeze", "--local")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to get the output of pip freeze: %v", err)
	}

	err = os.WriteFile(path, output, 0644)
	if err != nil {
		return fmt.Errorf("failed to write requirements.txt: %v", err)
	}

	return nil
}
