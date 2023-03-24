package python

import (
	"bytes"
	"fmt"
	"langforge/system"
	"os"
	"os/exec"
	"strings"
)

// PythonPackage represents a Python package with its name and version.
type PythonPackage struct {
	Name    string
	Version string
}

// GetInstalledPackages retrieves a list of currently installed Python packages
// with their name and version. It returns an error if it fails to locate pip or
// execute the pip command.
func GetInstalledPackages() ([]PythonPackage, error) {
	pipPath, err := system.FindPip()
	if err != nil {
		return nil, fmt.Errorf("failed to locate pip: %v", err)
	}

	// Build the pip command.
	cmd := exec.Command(pipPath, "list", "--format=freeze")

	// Run the command and capture its output.
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("failed to list installed packages: %v", err)
	}

	// Parse the output and build a list of PythonPackage objects.
	var packages []PythonPackage
	for _, line := range bytes.Split(stdout.Bytes(), []byte{'\n'}) {
		if len(line) == 0 {
			continue
		}
		parts := bytes.SplitN(line, []byte("=="), 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid package specification: %s", line)
		}
		packages = append(packages, PythonPackage{
			Name:    string(parts[0]),
			Version: string(parts[1]),
		})
	}

	return packages, nil
}

// managePackages is a helper function to handle common tasks for installing
// and uninstalling Python packages. It takes a list of packages and an action
// ("install" or "uninstall") as arguments. It returns an error if it fails to
// locate the Python interpreter or execute the pip command.
func managePackages(packages []string, action string) error {
	if len(packages) == 0 {
		return nil
	}

	// Get unique packages
	uniquePackages := make(map[string]bool)
	for _, pkg := range packages {
		uniquePackages[pkg] = true
	}
	packages = []string{}
	for pkg := range uniquePackages {
		packages = append(packages, pkg)
	}

	// Find the path to the Python interpreter
	pythonPath, err := system.FindPython()
	if err != nil {
		return err
	}

	// Manage the packages using pip
	args := append([]string{"-m", "pip"}, strings.Split(action, " ")...)
	args = append(args, packages...)
	args = append(args, "--disable-pip-version-check")
	cmd := exec.Command(pythonPath, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()

	return err
}

// InstallPackages installs the specified Python packages. It returns an error
// if it fails to locate the Python interpreter, execute the pip command or
// manage packages.
func InstallPackages(packages []string) error {
	return managePackages(packages, "install")
}

// UninstallPackages uninstalls the specified Python packages. It returns an error
// if it fails to locate the Python interpreter, execute the pip command or
// manage packages.
func UninstallPackages(packages []string) error {
	return managePackages(packages, "uninstall -y")
}
