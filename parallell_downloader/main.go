package main

import (

	tm "github.com/buger/goterm"

	cli "lets_go/parallell_downloader/utils/cli"
	download "lets_go/parallell_downloader/utils/download"
	progressbar "lets_go/parallell_downloader/utils/progressbar"
)

func main() {
	urls := cli.GetURLs()
	c := make(chan download.Statistics)
	n := len(urls)
	overview := make([]float32, n)
	fileNames := make([]string, n)

	for i, url := range urls {
		go download.Downloader(url, i, c, fileNames)
	}
	
	for {
		select {
		case update := <-c:
			// Update progressbar with newly received statistics
			overview[update.Item] = update.Percentage
			progressbar.PrintStatus(overview, fileNames)

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
