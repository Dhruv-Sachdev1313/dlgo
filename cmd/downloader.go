package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"

	"github.com/spf13/cobra"
	//import downloader package
)

const chunkSize = 1024 * 1024 // 1 MB chunks

var downloaderCmd = &cobra.Command{
	Use:   "download",
	Short: "Download a file from the internet",
	Long:  `Download a file from the internet asyncronously.`,
	Run:   asyncDownload,
}

// func init() {
// 	// rootCmd.AddCommand(downloaderCmd)
// 	// downloaderCmd.Flags().StringP("url", "u", "", "URL to download file from")
// 	// downloaderCmd.Flags().IntP("worker_count", "w", 2, "Number of workers to use")
// 	// downloaderCmd.Flags().StringP("output_path", "o", "output", "Path to output file")
// }

func asyncDownload(cmd *cobra.Command, args []string) {
	fmt.Println("Downloading file...")
	url, _ := cmd.Flags().GetString("url")
	workerCount, _ := cmd.Flags().GetInt("worker_count")
	outputPath, _ := cmd.Flags().GetString("output_path")
	err := DownloadFileConcurrently(url, outputPath, workerCount)
	if err != nil {
		fmt.Println("Error:", err)
	}
	// go downloadFile(url, filepath)
	// fmt.Println("Download started...")
}

func DownloadFileConcurrently(url, outputPath string, numWorkers int) error {
	resp, err := http.Head(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP response error: %s", resp.Status)
	}

	fileSize := resp.ContentLength
	fmt.Printf("File size: %d bytes\n", fileSize)

	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	var wg sync.WaitGroup
	chunkSize := fileSize / int64(numWorkers)

	for i := 0; i < numWorkers; i++ {
		start := int64(i) * chunkSize
		end := start + chunkSize
		if i == numWorkers-1 {
			// Last worker may download slightly more if not divisible evenly
			end = fileSize
		}

		wg.Add(1)
		go func(start, end int64) {
			defer wg.Done()
			err := downloadChunk(url, file, start, end)
			if err != nil {
				fmt.Println("Error:", err)
			}
		}(start, end)
	}

	wg.Wait()

	fmt.Println("Download complete")
	return nil
}

func downloadChunk(url string, file *os.File, start, end int64) error {
	client := http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	// Specify the range of bytes to download
	req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", start, end))

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusPartialContent {
		return fmt.Errorf("HTTP response error: %s", resp.Status)
	}

	_, err = file.Seek(start, io.SeekStart)
	if err != nil {
		return err
	}

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}

	fmt.Printf("Downloaded chunk: %d-%d\n", start, end)
	return nil
}