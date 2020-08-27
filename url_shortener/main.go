package main

import (
	"log"
	"net/http"

	utils "./utils"
	"github.com/gorilla/mux"
)

var bucketName = "urlShortener"
var db = utils.CreateDB(bucketName)

func main() {
	defer db.Close()
	r := mux.NewRouter()

	// Define folder for static files
	s := http.StripPrefix("/static/", http.FileServer(http.Dir("./static/")))
	r.PathPrefix("/static/").Handler(s)

	// Define endpoints
	r.HandleFunc("/shortener/{page}", utils.ShortenHandler(db, bucketName))
	r.HandleFunc("/{hash:[a-z0-9]{5}}", utils.RedirectHandler(db, bucketName))
	r.HandleFunc("/", utils.MainHandler)

	log.Fatal(http.ListenAndServe(":8080", r))

}
