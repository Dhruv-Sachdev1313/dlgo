/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "dlgo",
	Short: "A download accelerator written in Go",
	Long: `
A download accelerator written in Go.
This Splits file and downloads it concurrently using multiple workers.
This currently supports only HTTP/HTTPS URLs.

Example usage:
dlgo download -u <URL> -w <Number of workers> -o <Output file path>`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.dlgo.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	// rootCmd.AddCommand(downloaderCmd)
	// downloaderCmd.Flags().StringP("url", "u", "", "URL to download file from")
	// downloaderCmd.Flags().IntP("worker_count", "w", 2, "Number of workers to use")
	// downloaderCmd.Flags().StringP("output_path", "o", "output", "Path to output file")
}
