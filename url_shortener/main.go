package main

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	utils "./utils"
)

var bucketName = "urlShortener"
var db = utils.CreateDB(bucketName)

func shortenHandler(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Path[len("/shortener/"):]

	// Get hash value of url as shortened url
	hashValue := fmt.Sprintf("%x", sha1.Sum([]byte(url)))
	trimmedHashValue := strings.ReplaceAll(hashValue, " ", "")[:5]

	// Store and return response
	mappedReturn, _ := json.Marshal(map[string]string{url: trimmedHashValue})
	w.Write([]byte(string(mappedReturn)))
	utils.WriteRecord(db, bucketName, trimmedHashValue, url)
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	hashValue := r.URL.Path[len("/"):]

	// Retrieve url corresponding to url from bolt db and redirect
	redirectURL := utils.ReadRecord(db, bucketName, hashValue)
	http.Redirect(w, r, redirectURL, 302)
}

func main() {
	defer db.Close()

	http.HandleFunc("/shortener/", shortenHandler)
	http.HandleFunc("/", redirectHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))

}
