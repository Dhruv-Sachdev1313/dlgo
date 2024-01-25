package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var downloaderCmd = &cobra.Command{
	Use:   "download",
	Short: "Download a file from the internet",
	Long:  `Download a file from the internet asyncronously.`,
	Run:   asyncDownload,
}

func init() {
	rootCmd.AddCommand(downloaderCmd)
	downloaderCmd.Flags().StringP("url", "u", "", "URL to download file from")
	downloaderCmd.Flags().IntP("worker_count", "w", 1, "Number of workers to use")
}

func asyncDownload(cmd *cobra.Command, args []string) {
	fmt.Println("Downloading file...")
	url, _ := cmd.Flags().GetString("url")
	workerCount, _ := cmd.Flags().GetInt("worker_count")
	// go downloadFile(url, filepath)
	fmt.Println("Download started...")
}
