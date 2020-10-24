package main

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/buger/goterm"
	tm "github.com/buger/goterm"
)

type itemUpdate struct {
	item       int
	percentage float32
}

type progressBar struct {
	total    int
	item     int
	progress int
	channel  chan itemUpdate
}

func (pb *progressBar) Write(p []byte) (int, error) {
	n := len(p)
	pb.progress += int(n)
	update := itemUpdate{item: pb.item, percentage: 100 * float32(pb.progress) / float32(pb.total)}
	select {
	case pb.channel <- update:
	}

	return n, nil
}

func download(url string, item int, c chan itemUpdate) {
	resp, _ := http.Get(url)
	defer resp.Body.Close()
	size, _ := strconv.Atoi(resp.Header.Get("Content-Length"))

	name := fmt.Sprintf("test%d.tar.gz", rand.Intn(10000))
	f, _ := os.Create(name)
	defer f.Close()

	pb := &progressBar{total: size, item: item, channel: c}
	io.Copy(f, io.TeeReader(resp.Body, pb))
}

func sum(s []float32) int {
	var sum float32
	for _, v := range s {
		sum += v
	}
	return int(sum)
}

func printProgressBar(arr []float32) {

	for i, percentage := range arr {

		tm.Printf("item %d: %s %5.2f\n", i, pgb(percentage), percentage)
	}
	tm.Flush()

}

func pgb(percentage float32) string {
	var elementWidth float32 = 2.5
	numEqual := int(percentage / elementWidth)

	pg := strings.Repeat("=", numEqual) + ">" + strings.Repeat(" ", int(100/elementWidth))
	return pg[:int(100/elementWidth)]
}

func main() {
	url := "http://ipv4.download.thinkbroadband.com/5MB.zip"
	n := 2
	c := make(chan itemUpdate)

	for i := 0; i < n; i++ {
		go download(url, i, c)
	}

	done := true
	tracker := make([]float32, n)
	for done == true {
		select {
		case res := <-c:
			tracker[res.item] = res.percentage
			printProgressBar(tracker)

			if sum(tracker) == n*100 {
				done = false
			} else {
				goterm.MoveCursorUp(n)
				tm.Flush()

			}

		}
	}
}
