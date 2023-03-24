package cmd

import (
	"fmt"
	"langforge/python"
	"langforge/tui"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var integrationsCmd = &cobra.Command{
	Use:   "integrations",
	Short: "Edit the integrations for your langchain application",
	Long: `The integrations command allows you to view and edit the integrations
used in your langchain application.`,
	Run: func(cmd *cobra.Command, args []string) {
		editIntegrations()
	},
}

func init() {
	rootCmd.AddCommand(integrationsCmd)
}

func editIntegrations() {
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

	handler := python.NewPythonHandler(cwd)

	err = tui.EditAndUpdateIntegrations(handler, false, true)
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully updated integrations.")
}
