package main

import (
	"crypto/sha1"
	"fmt"
	"log"
	"net/http"
	"strings"
	"text/template"

	utils "./utils"
	"github.com/gorilla/mux"
)

var bucketName = "urlShortener"
var db = utils.CreateDB(bucketName)

func shortenHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// Get hash value of url as shortened url
	hashValue := fmt.Sprintf("%x", sha1.Sum([]byte(vars["page"])))
	trimmedHashValue := strings.ReplaceAll(hashValue, " ", "")[:5]

	// Store and return response
	utils.WriteRecord(db, bucketName, trimmedHashValue, vars["page"])
	tmpl := template.Must(template.ParseFiles("static/shorten.html"))
	data := struct {
		URL  string
		Hash string
	}{
		URL:  vars["page"],
		Hash: r.Referer() + trimmedHashValue,
	}
	tmpl.Execute(w, data)
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	hashValue := r.URL.Path[len("/"):]

	// Retrieve url corresponding to url from bolt db and redirect
	redirectURL := utils.ReadRecord(db, bucketName, hashValue)
	http.Redirect(w, r, redirectURL, 302)
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/main.html")
}

func main() {
	defer db.Close()
	r := mux.NewRouter()
	r.HandleFunc("/shortener/{page}", shortenHandler)
	r.HandleFunc("/{hash:[a-z0-9]{5}}", redirectHandler)
	r.HandleFunc("/", mainHandler)
	log.Fatal(http.ListenAndServe(":8080", r))

}
