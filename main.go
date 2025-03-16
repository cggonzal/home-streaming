package main

import (
	"fmt"
	"homeStreaming/customLogger"
	"homeStreaming/templates"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var MEDIA_DIR string = os.Getenv("MEDIA_DIR")

func index(w http.ResponseWriter, r *http.Request) {
	logger := customLogger.GetLogger()

	// walk through media directory and find path to all files
	var fileNames []string
	err := filepath.Walk(MEDIA_DIR, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			logger.Fatal(err)
		}

		// only allow selection to stream mp4 files that do not start with a "."
		_, fileName := filepath.Split(path)
		if !info.IsDir() && !strings.HasPrefix(fileName, ".") && strings.HasSuffix(fileName, ".mp4") {
			fileNames = append(fileNames, strings.TrimPrefix(path, MEDIA_DIR))
		}

		return nil
	})
	if err != nil {
		logger.Fatal(err)
	}

	// execute template
	data := templates.IndexData{UploadedFileNames: fileNames}
	templates.IndexTemplate.Execute(w, data)
}

func upload(w http.ResponseWriter, r *http.Request) {
	logger := customLogger.GetLogger()

	if r.Method != http.MethodPost {
		http.ServeFile(w, r, "./static/upload.html")
		return
	}

	// get uploaded file
	file, header, err := r.FormFile("filename")
	if err != nil {
		logger.Fatal("error getting file...", err)
	}

	data := make([]byte, header.Size)
	_, err = file.Read(data)
	if err != nil {
		logger.Fatal("error reading file: ", err)
	}

	err = os.WriteFile(MEDIA_DIR+header.Filename, data, 0666)
	if err != nil {
		logger.Fatal("Error saving file:", err)
	}

	// redirect to landing page
	http.Redirect(w, r, "/", http.StatusFound)
}

func stream(w http.ResponseWriter, r *http.Request) {
	data := templates.StreamData{VideoName: r.URL.Path[len("/stream/"):]}
	templates.StreamTemplate.Execute(w, data)
}

func download(w http.ResponseWriter, r *http.Request) {
	logger := customLogger.GetLogger()

	// if asking for file, serve file
	if r.URL.Path != "/download/" {
		filename := r.URL.Path[len("/download/"):]
		w.Header().Set("Content-Disposition", "attachment; filename="+filename)
		w.Header().Set("Content-Type", "application/octet-stream")
		http.ServeFile(w, r, MEDIA_DIR+filename)
		return
	}

	// walk through media directory and find path to all files
	var fileNames []string
	err := filepath.Walk(MEDIA_DIR, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			logger.Fatal(err)
		}

		// don't append directories or files that start with a "." since those are hidden files with useless meta data
		_, fileName := filepath.Split(path)
		if !info.IsDir() && !strings.HasPrefix(fileName, ".") {
			fileNames = append(fileNames, strings.TrimPrefix(path, MEDIA_DIR))
		}
		return nil
	})
	if err != nil {
		logger.Fatal(err)
	}

	// serve list of files
	data := templates.DownloadData{UploadedFileNames: fileNames}
	templates.DownloadTemplate.Execute(w, data)
}

func main() {
	// serve landing page
	http.HandleFunc("/", index)

	// serve static files
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	// serve media files. Translate path so that /media/ requests are routed to serve files from $MEDIA_DIR.
	// See here for godoc example: https://pkg.go.dev/net/http#example-FileServer-StripPrefix
	// Note: http.FileServer handles byte range requests automatically :)
	http.Handle("/media/", http.StripPrefix("/media/", http.FileServer(http.Dir(MEDIA_DIR))))

	// upload file
	http.HandleFunc("/upload/", upload)

	// stream video
	http.HandleFunc("/stream/", stream)

	// download video
	http.HandleFunc("/download/", download)

	// initialize Logger, this has to come before all other initializations since they use the logger
	customLogger.InitLogger()

	// initialize templates
	templates.InitTemplates()

	// check if $PORT environment variable is set
	logger := customLogger.GetLogger()
	if os.Getenv("PORT") == "" {
		logger.Fatal("ERROR... No $PORT environment variable set... Exiting...")
	}

	// check if $MEDIA_DIR environment variable is set
	if os.Getenv("MEDIA_DIR") == "" {
		logger.Fatal("ERROR... No $MEDIA_DIR environment variable set... Exiting...")
	}

	// append a "/" to the $MEDIA_DIR environment variable if it is not set
	if !strings.HasSuffix(MEDIA_DIR, "/") {
		MEDIA_DIR += "/"
	}

	// start the server on given $PORT
	PORT := fmt.Sprintf(":%s", os.Getenv("PORT"))
	logger.Print("starting app on port ", PORT)
	logger.Fatal(http.ListenAndServe(PORT, nil))
}
