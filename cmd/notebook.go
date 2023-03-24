package cmd

import (
	"bytes"
	"fmt"
	"langforge/python"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

var notebookCmd = &cobra.Command{
	Use:   "notebook",
	Short: "Starts a Jupyter Notebook server and opens a browser",
	Run: func(cmd *cobra.Command, args []string) {
		startNotebookCmd()
	},
}

func init() {
	rootCmd.AddCommand(notebookCmd)
}

func startNotebookCmd() {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory:", err)
		return
	}

	venvDir := filepath.Join(cwd, ".venv")
	if _, err := os.Stat(venvDir); err == nil {
		// Activate the virtual environment
		err = python.ActivateEnvironment(venvDir)
		if err != nil {
			fmt.Println("Error activating virtual environment:", err)
			return
		}
	} else {
		fmt.Println("No virtual environment found. Continuing in the current environment.")
	}

	err = python.SetJupyterEnvironmentVariables(cwd)
	if err != nil {
		panic(err)
	}

	// Start the Jupyter Notebook server
	notebookCmd := exec.Command("jupyter", "lab", "--LabApp.default_url='/doc'")
	var notebookOutput bytes.Buffer
	notebookCmd.Stdout = &notebookOutput
	notebookCmd.Stderr = os.Stderr

	err = notebookCmd.Run()
	if err != nil {
		panic(err)
	}
}