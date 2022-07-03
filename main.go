package main

import (
	"log"
	"net/http"
)

func serveContent(fs http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// add mimetype="application/dash+xml" for .mpd files
		if len(r.URL.Path) > 4 && r.URL.Path[len(r.URL.Path)-4:] == ".mpd" {
			w.Header().Add("mimetype", "application/dash+xml")
		}
		log.Print("url path is:", r.URL.Path)
		fs.ServeHTTP(w, r)
	}
}

func main() {
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", serveContent(fs))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
