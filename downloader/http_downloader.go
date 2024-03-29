package downloader

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
)

func HTTPDownloadFileConcurrently(url, outputPath string, numWorkers int) error {

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
		return fmt.Errorf("HTTP response code error: %d", resp.StatusCode)
	}

	fileSize := resp.ContentLength

	fetchContentSize := fileSize / int64(numWorkers)

	var mutex sync.Mutex
	var wg sync.WaitGroup
	progress := make(chan int64)

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
		wg.Add(1)
		go func(start, end int64) {
			defer wg.Done()
			err := downloadChunk(url, file, start, end, &mutex, progress)
			if err != nil {
				fmt.Println("Error:", err)
			}
		}(start, end)
	}
	go func() {
		var downloadedBytes int64
		for n := range progress {
			downloadedBytes += n
			fmt.Printf("\rDownloading.. %d%% ", downloadedBytes*100/fileSize)
		}
	}()
	wg.Wait()
	close(progress)
	return nil
}

func downloadChunk(url string, file *os.File, start int64, end int64, mutex *sync.Mutex, progress chan<- int64) error {
	client := http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	rangeHeader := fmt.Sprintf("bytes=%d-%d", start, end)
	req.Header.Set("Range", rangeHeader)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusPartialContent {
		return fmt.Errorf("server does not support partial content")
	}
	defer resp.Body.Close()

	mutex.Lock()

	_, err = file.Seek(start, io.SeekStart)
	if err != nil {
		return err
	}

	n, err := io.Copy(file, resp.Body)
	if err != nil {
		return err
	}
	progress <- n
	mutex.Unlock()

	return nil

}
