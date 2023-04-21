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

var jupyterLabCmd = &cobra.Command{
	Use:   "lab",
	Short: " launches JupyterLab for instant coding in your virtual environment.",
	Long: `The lab command starts a JupyterLab server within your virtual environment and 
launches a browser window, enabling you to begin coding immediately in an 
interactive workspace.`,
	Run: func(cmd *cobra.Command, args []string) {
		startJupyterLabCmd()
	},
}

func init() {
	rootCmd.AddCommand(jupyterLabCmd)
}

func startJupyterLabCmd() {
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
	labCmd := exec.Command("jupyter", "lab")
	var labOutput bytes.Buffer
	labCmd.Stdout = &labOutput
	labCmd.Stderr = os.Stderr

	err = labCmd.Run()
	if err != nil {
		panic(err)
	}
}
