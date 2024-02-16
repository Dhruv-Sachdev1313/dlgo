package downloader

import (
	"fmt"
	"net/http"
	"os"
)

func DownloadFileConcurrently(url, outputPath string, numWorkers int) error {

	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	resp, err := http.Head(url)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP response code error: %s", resp.StatusCode)
	}

	fileSize := resp.ContentLength

	fetchContentSize := fileSize / int64(numWorkers)

	var start, end int64

	for i := 0; i < numWorkers; i++ {
		if i == 0 {
			start = fetchContentSize * int64(i)
		} else {
			start = end + 1
		}
		end = start + fetchContentSize
		if i == numWorkers-1 {
			end = fileSize
		}
		go func(start, end int64) {
			err := downloadChunk(url, file, start, end)
			if err != nil {
				fmt.Println("Error:", err)
			}
		}(start, end)
	}
	return nil
}

func downloadChunk(url string, file *os.File, start int64, end int64) error {
	client := http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	rangeHeader := fmt.Sprintf("bytes=%d-%d", start, end)
	req.Header.Set("Range", rangeHeader)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusPartialContent {

	}
	defer resp.Body.Close()

	return nil

}
