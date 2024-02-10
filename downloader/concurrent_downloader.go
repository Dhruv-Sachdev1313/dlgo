package downloader

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
)

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
