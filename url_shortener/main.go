package main

import (
	"log"
	"net/http"

	db_utils "lets_go/url_shortener/utils/db"
	view_utils "lets_go/url_shortener/utils/views"

	"github.com/gorilla/mux"
)

var bucketName = "urlShortener"
var db = db_utils.CreateDB(bucketName)

func main() {
	defer db.Close()
	r := mux.NewRouter()

	// Define folder for static files
	s := "/static/"
	r.PathPrefix(s).Handler(http.StripPrefix(s, http.FileServer(http.Dir("."+s))))

	// Define app endpoints
	r.HandleFunc("/shortener/{page}", view_utils.ShortenHandler(db, bucketName))
	r.HandleFunc("/{hash:[a-z0-9]{5}}", view_utils.RedirectHandler(db, bucketName))
	r.HandleFunc("/", view_utils.MainHandler)

	log.Panic(http.ListenAndServe(":8080", r))

}
