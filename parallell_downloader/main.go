package main

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"time"
)

type response struct {
	status  bool
	endtime time.Time
}

func download(url string, c chan response) {
	resp, _ := http.Get(url)
	defer resp.Body.Close()

	name := fmt.Sprintf("test%d.tar.gz", rand.Intn(10000))
	f, _ := os.Create(name)
	defer f.Close()

	io.Copy(f, resp.Body)

	c <- response{status: true, endtime: time.Now()}
}

func main() {
	url := "http://ipv4.download.thinkbroadband.com/5MB.zip"
	n := 3
	c := make(chan response)
	for i := 0; i < n; i++ {
		go download(url, c)
	}

	for i := 0; i < n; i++ {
		res := <-c
		fmt.Printf("%t, %s\n", res.status, res.endtime)
	}
}
