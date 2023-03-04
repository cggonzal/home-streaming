package main

import (
	"fmt"
	"homeStreaming/customLogger"
	"homeStreaming/templates"
	"net/http"
	"os"
)

func index(w http.ResponseWriter, r *http.Request) {
	logger := customLogger.GetLogger()
	saveDir := "./media/"

	// serve list of files
	files, err := os.ReadDir(saveDir)
	if err != nil {
		logger.Fatal("error reading directory: ", err)
	}
	var fileNames []string
	for _, file := range files {
		fileNames = append(fileNames, file.Name())
	}
	data := templates.IndexData{UploadedFileNames: fileNames}
	templates.IndexTemplate.Execute(w, data)
}

func upload(w http.ResponseWriter, r *http.Request) {
	logger := customLogger.GetLogger()

	// directory where uploaded files are saved to and read from
	saveDir := "./media/"

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

	err = os.WriteFile(saveDir+header.Filename, data, 0666)
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
	saveDir := "./media/"

	// if asking for file, serve file
	if r.URL.Path != "/download/" {
		filename := r.URL.Path[len("/download/"):]
		w.Header().Set("Content-Disposition", "attachment; filename="+filename)
		w.Header().Set("Content-Type", "application/octet-stream")
		http.ServeFile(w, r, saveDir+filename)
		return
	}

	// serve list of files
	files, err := os.ReadDir(saveDir)
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

	// serve media files
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

	// start the server on given $PORT
	PORT := fmt.Sprintf(":%s", os.Getenv("PORT"))
	logger.Print("starting app on port ", PORT)
	logger.Fatal(http.ListenAndServe(PORT, nil))
}
