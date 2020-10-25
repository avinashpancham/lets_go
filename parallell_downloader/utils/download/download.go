package download

import (
		"fmt"
		"io"
		"log"
		"net/http"
		"os"
		"strings"
	)

// Statistics is a container for download progress statistics
type Statistics struct {
	Item       int
	FileName	string
	Percentage float32
}

type tracker struct {
	channel  chan Statistics
	fileName string
	item     int
	progress int
	size    int
}

func (tr *tracker) Write(p []byte) (int, error) {
	// Increment downloadedd data
	n := len(p)
	tr.progress += int(n)
	percentage := 100 * float32(tr.progress) / float32(tr.size)

	// Forward all data to a central location
	select {
	case tr.channel <- Statistics{Item: tr.item, FileName: tr.fileName, Percentage: percentage}:
	}

	return n, nil
}

func getFileName(responseHeader http.Header, url string) string {
	headerName := responseHeader.Get("Content-Disposition")

	var fileName string

	// Check if filename is present in the header else make it from url
	if headerName == "" {		
		splitURL := strings.Split(url,"/")
		fileName= splitURL[len(splitURL)-1] 
	} else {
		fileName = headerName 
	}

	// Check if filename already exists
	if _, err := os.Stat(fileName); err == nil {
		// If filename exists try saving with suffix (i)
		i := 1
		fileNameComponents := strings.Split(fileName,".")	
		
		for {
			fileName := fmt.Sprintf(fileNameComponents[0] + "(%d)." + fileNameComponents[1], i)
			if _, err := os.Stat(fileName); os.IsNotExist(err) {
				// Check whether the form with suffix (i) exists, else increment i
				return fileName
			}
			i++
		}
	} else if !os.IsNotExist(err) {
		// Catch errors due to permission etc
		log.Fatal(err)	  
	} 
	return fileName
}


// Downloader downloads the content from url
func Downloader(url string, item int, c chan Statistics, fileNames []string) {
	resp, _ := http.Get(url)
	defer resp.Body.Close()

	// Get file info
	size := int(resp.ContentLength)
	fileName := getFileName(resp.Header, url)
	fileNames[item] = fileName

	f, _ := os.Create(fileName)
	defer f.Close()

	tr := &tracker{size: size, item: item, fileName: fileName, channel: c}
	io.Copy(f, io.TeeReader(resp.Body, tr))
}