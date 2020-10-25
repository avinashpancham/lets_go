package progressbar

import (
	"strings"
	tm "github.com/buger/goterm"
)

// PrintStatus prints the status of all downloads to the terminals
func PrintStatus(arr []float32, names []string) {

	for index, percentage := range arr {
		tm.Printf("%s \t[%s] %5.2f%%\n",names[index] , printProgressBar(percentage), percentage)
	}
	tm.Flush()

}

func printProgressBar(arr float32) string {
	var elementWidth float32 = 2.5 
	progressBarString := strings.Repeat("=", int(arr / elementWidth)) + ">" + strings.Repeat(" ", int(100/elementWidth))

	return progressBarString[:int(100/elementWidth)]
}


// DownloadCompleted checks whether the download is completed
func DownloadCompleted(arr []float32) bool {
	var sum float32
	for _, value := range arr {
		sum += value
	}

	return int(sum) == len(arr)*100
}
