package utils

import (
	"crypto/sha1"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	bolt "go.etcd.io/bbolt"
)

// MainHandler serves mainpage
func MainHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "html/main.html")
}

// ShortenHandler shortens url of provided webpage
func ShortenHandler(db *bolt.DB, bucketName string) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		// Get hash value of url as shortened url
		hashValue := fmt.Sprintf("%x", sha1.Sum([]byte(vars["page"])))
		trimmedHashValue := strings.ReplaceAll(hashValue, " ", "")[:5]

		// Store and return response
		WriteRecord(db, bucketName, trimmedHashValue, vars["page"])
		tmpl := template.Must(template.ParseFiles("html/shorten.html"))
		data := struct {
			URL  string
			Hash string
		}{
			URL:  vars["page"],
			Hash: r.Referer() + trimmedHashValue,
		}
		tmpl.Execute(w, data)
	}

	return http.HandlerFunc(fn)
}

// RedirectHandler redirects to target page
func RedirectHandler(db *bolt.DB, bucketName string) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		hashValue := r.URL.Path[len("/"):]

		// Retrieve url corresponding to url from bolt db and redirect
		redirectURL := ReadRecord(db, bucketName, hashValue)
		http.Redirect(w, r, redirectURL, 302)
	}

	return http.HandlerFunc(fn)
}
