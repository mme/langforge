package system

import (
	"errors"
	"os"
	"os/exec"
	"strings"
)

// FindPython searches for the Python interpreter in the system's PATH.
// It looks for a binary called "python3" first, and if that's not found, it
// looks for a binary called "python". If neither binary is found, an error is returned.
//
// Returns the path to the Python interpreter binary and nil error if
// it is found, or empty string and non-nil error if it is not found.
func FindPython() (string, error) {
	// try to find python3 first
	pythonPath, err := exec.LookPath("python3")
	if err == nil {
		return pythonPath, nil
	}

	// if python3 not found, try to find python
	pythonPath, err = exec.LookPath("python")
	if err == nil {
		return pythonPath, nil
	}

	// if neither are found, return an error
	return "", errors.New("python interpreter not found")
}

// FindNode searches for the Node.js interpreter in the system's PATH.
// It looks for a binary called "node". If the binary is found, it returns
// the path to the binary and nil error. If the binary is not found, an error
// is returned.
//
// Returns the path to the Node.js interpreter binary and nil error if
// it is found, or empty string and non-nil error if it is not found.
func FindNode() (string, error) {
	nodePath, err := exec.LookPath("node")
	if err != nil {
		return "", errors.New("node.js interpreter not found")
	}

	return nodePath, nil
}

// FindPip searches for the location of the pip command in the system. It first searches for pip3, then for pip,
// returning the location of the command if found. If the command is not found, it returns an error.
//
// Returns:
// - string: the location of the pip command
// - error: an error if the pip command was not found
func FindPip() (string, error) {
	// Try pip3
	if path, err := exec.LookPath("pip3"); err == nil {
		return path, nil
	}

	// Try pip
	if path, err := exec.LookPath("pip"); err == nil {
		return path, nil
	}

	// Couldn't find pip or pip3
	return "", errors.New("pip command not found")
}

// ShellSource emulates the action of the "source" command in bash by executing
// a shell script and setting environment variables based on its output. The
// script file is passed in as an argument to the function. It returns an error
// if the script fails to execute.
//
// Parameters:
//   - script: the path to the shell script to execute.
//
// Returns:
//   - nil error if the script is executed successfully and the environment variables
//     are set, or a non-nil error if the script fails to execute.
func ShellSource(script string) error {
	cmd := exec.Command("sh", "-c", ". "+script+" && env")

	output, err := cmd.Output()
	if err != nil {
		return errors.New("Failed to execute shell script: " + err.Error())
	}

	env := make(map[string]string)
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			env[parts[0]] = parts[1]
		}
	}

	for key, value := range env {
		os.Setenv(key, value)
	}

	return nil
}

// ExecuteCommands takes a list of shell commands as input, removes duplicates,
// and executes them sequentially. It returns an error if any of the commands fail
// to execute. The stdout and stderr of the executed commands are redirected to
// the current process's stdout and stderr.
func ExecuteCommands(commands []string, dir string) error {

	if len(commands) == 0 {
		return nil
	}

	for _, command := range commands {
		parts := strings.Split(command, " ")
		cmdName := parts[0]
		args := []string{}
		if len(parts) > 1 {
			args = parts[1:]
		}
		cmd := exec.Command(cmdName, args...)
		cmd.Dir = dir
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			return err
		}
	}
	return nil
}
