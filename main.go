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
		fmt.Printf("dir: %v: name: %s\n", info.IsDir(), path)
		if !info.IsDir() { // append only if we are a file since we only want to keep track of files
			fileNames = append(fileNames, strings.TrimPrefix(path, MEDIA_DIR))
		}
		return nil
	})
	if err != nil {
		logger.Fatal(err)
	}
	logger.Println("filenames:", fileNames)
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

	// serve list of files
	files, err := os.ReadDir(MEDIA_DIR)
	if err != nil {
		logger.Fatal("error reading directory: ", err)
	}
	var fileNames []string
	for _, file := range files {
		fileNames = append(fileNames, file.Name())
	}
	data := templates.DownloadData{UploadedFileNames: fileNames}
	templates.DownloadTemplate.Execute(w, data)
}

func main() {
	// serve landing page
	http.HandleFunc("/", index)

	// serve media files, handles byte range requests automatically :)
	http.Handle("/media/", http.FileServer(http.Dir(".")))

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
