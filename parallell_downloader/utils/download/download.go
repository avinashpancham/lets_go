package download

import (
		"fmt"
		"io"
		"math/rand"
		"net/http"
		"os"
		"strconv"
	)

// Statistics is a container for download progress statistics
type Statistics struct {
	Item       int
	Percentage float32
}

type tracker struct {
	channel  chan Statistics
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
	case tr.channel <- Statistics{Item: tr.item, Percentage: percentage}:
	}

	return n, nil
}

// Downloader downloads the content from url
func Downloader(url string, item int, c chan Statistics) {
	resp, _ := http.Get(url)
	defer resp.Body.Close()
	size, _ := strconv.Atoi(resp.Header.Get("Content-Length"))

	name := fmt.Sprintf("test%d.tar.gz", rand.Intn(10000))
	f, _ := os.Create(name)
	defer f.Close()

	tr := &tracker{size: size, item: item, channel: c}
	io.Copy(f, io.TeeReader(resp.Body, tr))
}