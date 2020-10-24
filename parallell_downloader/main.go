package main

import (

	tm "github.com/buger/goterm"

	download "lets_go/parallell_downloader/utils/download"
	progressbar "lets_go/parallell_downloader/utils/progressbar"
)

func main() {
	url := "http://ipv4.download.thinkbroadband.com/5MB.zip"
	n := 2
	c := make(chan download.Statistics)
	overview := make([]float32, n)

	for i := 0; i < n; i++ {
		go download.Downloader(url, i, c)
	}
	
	for {
		select {
		case update := <-c:
			// Update progressbar with newly received statistics
			overview[update.Item] = update.Percentage
			progressbar.PrintStatus(overview)

		}

		if progressbar.DownloadCompleted(overview) {
			// Stop listening once download is complete
			break
		} else {
			// Move cursor back to top to easily update progressbar
			tm.MoveCursorUp(n)
			tm.Flush()
		}
	}
}
