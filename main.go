package main

import (
	"log"
	"net/http"
)

func serveContent(fs http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		log.Print("request path: ", r.URL.Path)
		log.Print("request headers: ", r.Header)
		// add mimetype="application/dash+xml" for .mpd files
		if len(r.URL.Path) > 4 && r.URL.Path[len(r.URL.Path)-4:] == ".mpd" {
			w.Header().Add("mimetype", "application/dash+xml")
		}

		fs.ServeHTTP(w, r)
	}
}

func main() {
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", serveContent(fs))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
