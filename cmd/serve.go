/*
Copyright © Markus Ecker <markus.ecker@gmail.com>
*/
package cmd

import (
	"fmt"
	"io"
	"langforge/python"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve [notebook.ipynb]",
	Short: "Serve a LangChain application",
	Long:  `The serve command serves a LangChain from a Jupyter notebook.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("notebook is missing")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		// Directly get the port value without checking for the flag's change
		port, err := cmd.Flags().GetInt("port")
		if err != nil {
			fmt.Printf("Error parsing port: %v\n", err)
			return
		}
		serveAppCmd(args[0], port)
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().Int("port", 2204, "port number to serve LangChain application")
}

func serveAppCmd(notebookPath string, port int) {

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

	// Add filename and --port arguments to the command
	cmd := exec.Command("python", "-", notebookPath, "--port", strconv.Itoa(port))
	stdin, err := cmd.StdinPipe()
	if err != nil {
		panic(err)
	}

	pythonScript, err := python.ServerPy()
	if err != nil {
		panic(err)
	}

	// Set Stdout and Stderr to stream the output
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Start()
	if err != nil {
		panic(err)
	}

	io.WriteString(stdin, strings.TrimSpace(string(pythonScript)))
	stdin.Close()

	err = cmd.Wait()
	if err != nil {
		panic(err)
	}
}
