package cmd

import (
	"fmt"
	"langforge/python"
	"langforge/system"
	"langforge/tui"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// keysCmd represents the keys command
var keysCmd = &cobra.Command{
	Use:   "keys",
	Short: "Edit API keys for a LangChain application",
	Long:  `The keys command allows you to edit the API keys for a LangChain application.`,
	Run: func(cmd *cobra.Command, args []string) {
		editApiKeysCmd()
	},
}

func init() {
	rootCmd.AddCommand(keysCmd)
}

func editApiKeysCmd() {

	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory:", err)
		return
	}

	handler := python.NewPythonHandler(currentDir)

	venvDir := filepath.Join(currentDir, ".venv")
	if _, err := os.Stat(venvDir); err == nil {
		// Activate the virtual environment
		err = python.ActivateEnvironment(venvDir)
		if err != nil {
			fmt.Println("Error activating virtual environment:", err)
			return
		}
	}

	err = handler.DetermineInstalledIntegrations()
	if err != nil {
		panic(err)
	}

	// Ensure the environment has all required keys in the .env file
	dotEnvPath := filepath.Join(currentDir, ".env")

	apiKeys := handler.InstalledIntegrationsApiKeys()
	err = system.EnsureEnv(dotEnvPath, apiKeys)
	if err != nil {
		panic(err)
	}

	// Load the environment from the .env file
	env, err := system.ReadEnv(dotEnvPath)
	if err != nil {
		panic(err)
	}

	editedEnv, err := tui.EditApiKeys(apiKeys, env)
	if err != nil {
		panic(err)
	}

	env = editedEnv

	// Save the edited environment back to the .env file
	err = system.WriteEnv(dotEnvPath, env)
	if err != nil {
		panic(err)
	}
}
