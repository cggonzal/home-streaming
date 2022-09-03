/* NOTE: in order to add a new template there are 3 steps
   1. Add the new template name to the list of variables
   2. Create a new struct for the data that gets passed into the template
   3. Add the file(s) that get parsed to the init templates function using ParseFiles() as shown below
*/
package templates

import (
	"homeStreaming/customLogger"
	"html/template"
)

// 1.
// template variable names
var (
	IndexTemplate    *template.Template
	StreamTemplate   *template.Template
	DownloadTemplate *template.Template
)

// 2.
// template data definitions. One per template.
type IndexData struct {
	UploadedFileNames []string
}

type StreamData struct {
	VideoName string
}

type DownloadData struct {
	VideoName string
}

// 3.
// initialize templates. Store them in global variables so that files don't have to be parsed on every request
func InitTemplates() {
	logger := customLogger.GetLogger()

	var err error

	IndexTemplate, err = template.ParseFiles("templates/index.html")
	if err != nil {
		logger.Fatal(err)
	}

	StreamTemplate, err = template.ParseFiles("templates/stream.html")
	if err != nil {
		logger.Fatal(err)
	}

	DownloadTemplate, err = template.ParseFiles("templates/download.html")
	if err != nil {
		logger.Fatal(err)
	}
}
