/*
Copyright © Markus Ecker <markus.ecker@gmail.com>
*/
package cmd

import (
	"fmt"
	"langforge/environment"
	"langforge/python"
	"langforge/system"
	"langforge/tui"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create [app-name]",
	Short: "Create a new langchain application",
	Long: `The create command generates a new langchain application with LangForge. 
	
It sets up a virtual environment, installs dependencies, 
and generates a Jupyter notebook for experimentation. 
Simply run "langforge create myapp" to create a new 
directory with all the necessary files to get started.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("app name is missing")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		createAppCmd(args[0])
	},
}

func init() {
	rootCmd.AddCommand(createCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func createAppCmd(appName string) {

	handler := python.NewPythonHandler()

	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory:", err)
		return
	}

	dir := filepath.Join(currentDir, appName)

	// Check if a file with the specified app name already exists
	if _, err := os.Stat(dir); err == nil {
		panic(fmt.Errorf("file with name '%s' already exists", dir))
	}

	tui.DisplayBanner()

	// Check if a virtual environment should be created
	shouldCreateEnvironment, err := tui.PromptYesNo("Create a virtual environment for your 🦜️🔗LangChain app?", true)
	if err != nil {
		panic(err)
	}

	tui.EmptyLine()
	if shouldCreateEnvironment {
		fmt.Println("Yes, create a virtual environment.")
	} else {
		fmt.Println("No, install dependencies in this environment.")
	}
	tui.EmptyLine()

	if !shouldCreateEnvironment {
		err := handler.DetermineInstalledIntegrations()
		if err != nil {
			panic(err)
		}
	}

	if len(handler.NamesOfIntegrationsToInstall()) != 0 {
		fmt.Println("The following integrations will be installed:", handler.NamesOfIntegrationsToInstall())
	} else {
		fmt.Println("No additional integrations will be installed.")
	}
	tui.EmptyLine()

	if len(handler.NamesOfIntegrationsToUninstall()) != 0 {
		fmt.Println("The following integrations will be uninstalled:", handler.NamesOfIntegrationsToUninstall())
		tui.EmptyLine()
	}

	shouldEditIntegrations, err := tui.PromptYesNo("Would you like to edit the list of integrations?", false)
	if err != nil {
		panic(err)
	}
	tui.EmptyLine()

	if shouldEditIntegrations {
		// Ask the user which integrations to install
		_, err := tui.EditMultiSelect("Select the integrations you want to install:", environment.IntegrationsToSelectableItems(handler.GetIntegrations()))
		if err != nil {
			panic(err)
		}
	}

	// Create the app directory
	if err := os.Mkdir(appName, 0755); err != nil {
		panic(err)
	}

	if shouldCreateEnvironment {
		fmt.Println("Creating virtual environment...")
		tui.EmptyLine()

		// Create the virtual environment
		if err := python.CreateVirtualEnv(".venv", appName); err != nil {
			panic(err)
		}

		// activate the virtual environment
		if err := python.ActivateEnvironment(".venv", appName); err != nil {
			panic(err)
		}
	}

	// Install dependencies
	fmt.Println("Installing dependencies...")
	tui.EmptyLine()

	err = handler.ExecuteChanges(dir)
	if err != nil {
		panic(err)
	}

	// Ensure the environment has all required keys in the .env file
	dotEnvPath := filepath.Join(appName, ".env")
	apiKeys := handler.InstalledIntegrationsApiKeys()
	err = system.EnsureEnv(dotEnvPath, apiKeys)
	if err != nil {
		panic(err)
	}

	if len(apiKeys) > 0 {
		unsetKeys, err := system.UnsetAPIKeys(dotEnvPath, apiKeys)
		if err != nil {
			panic(err)
		}

		unsetKeysMap := make(map[string]bool)
		for _, key := range unsetKeys {
			unsetKeysMap[key] = true
		}

		tui.EmptyLine()
		nAPIKeys := len(apiKeys)
		if nAPIKeys == 1 {
			fmt.Println(tui.Bold("We found 1 API key associated with your installed integration:"))
		} else {
			fmt.Printf(tui.Bold("We found %d API keys associated with your installed integrations:\n"), nAPIKeys)
		}
		for _, key := range apiKeys {
			if unsetKeysMap[key] {
				fmt.Printf("- %s\n", key)
			} else {
				fmt.Printf("- %s (set)\n", key)
			}
		}

		tui.EmptyLine()
		shouldEditApiKeys, err := tui.PromptYesNo("Would you like to edit your API keys now?", true)
		if err != nil {
			panic(err)
		}

		if shouldEditApiKeys {
			tui.EmptyLine()
			// Load the environment from the .env file
			env, err := system.ReadEnv(dotEnvPath)
			if err != nil {
				panic(err)
			}

			// Edit the environment (only API keys)
			editedEnv, err := tui.EditApiKeys(apiKeys, env)
			if err != nil {
				panic(err)
			}

			// Save the edited environment back to the .env file
			err = system.WriteEnv(dotEnvPath, editedEnv)
			if err != nil {
				panic(err)
			}
		} else {
			tui.EmptyLine()
		}
	}

	fmt.Printf("Successfully created 🦜️🔗LangChain application '%s'.\n", appName)
}
