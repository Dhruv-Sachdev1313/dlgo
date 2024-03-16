package cmd

import (
	"fmt"

	downloader "github.com/Dhruv-Sachdev1313/dlgo/downloader"
	"github.com/spf13/cobra"
)

var downloaderCmd = &cobra.Command{
	Use:   "download",
	Short: "Download a file from the internet",
	Long: `Download a file from the internet asyncronously.
It does this by splitting the file and downloading it concurrently using multiple workers.`,
	Run: asyncDownload,
}

func init() {
	rootCmd.AddCommand(downloaderCmd)
	downloaderCmd.Flags().StringP("url", "u", "", "URL to download file from")
	downloaderCmd.Flags().IntP("worker_count", "w", 2, "Number of workers to use")
	downloaderCmd.Flags().StringP("output_path", "o", "output", "Path to output file")
}

func asyncDownload(cmd *cobra.Command, args []string) {
	url, _ := cmd.Flags().GetString("url")
	workerCount, _ := cmd.Flags().GetInt("worker_count")
	outputPath, _ := cmd.Flags().GetString("output_path")
	if url == "" {
		fmt.Println("Error: URL is required. Use -u flag to specify the URL")
		return
	}
	fmt.Println("Downloading file...")
	err := downloader.HTTPDownloadFileConcurrently(url, outputPath, workerCount)
	if err != nil {
		fmt.Println("Error:", err)
	}
}
