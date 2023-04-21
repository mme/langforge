/*
Copyright © Markus Ecker <markus.ecker@gmail.com>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "langforge",
	Short: "LangForge: A toolkit for building LangChain applications",
	Long: `LangForge is a toolkit for building and deploying LangChain applications. 
	
It simplifies the process by handling dependencies, providing 
Jupyter notebooks for experimentation, and enabling you to 
interact with your chains via a REST API.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	defer recoverFromPanic()
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.langforge.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func recoverFromPanic() {
	if r := recover(); r != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", r)
	}
}
