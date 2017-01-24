package web

import (
	"io"
	"net/http"
	"os"
	"strings"
)

func uiHandler(w http.ResponseWriter, r *http.Request) {

	fileNameFromRequest := func() string {
		filePath := strings.TrimPrefix(r.URL.Path, "/public/")

		if filePath == "" {
			filePath = "index.html"
		}
		return filePath
	}

	publicDir := os.Getenv("PUBLIC_DIR")

	if len(publicDir) == 0 {
		w.WriteHeader(404)
		w.Write([]byte("Not Found (PUBLIC_DIR variable not set)."))
		return
	}

	//Remove trailing slash, if any
	publicDir = strings.TrimSuffix(publicDir, "/")

	fileName := fileNameFromRequest()
	filePath := publicDir + string(os.PathSeparator) + fileName

	file, err := os.Open(filePath)
	if err != nil {
		w.WriteHeader(404)
		w.Write([]byte("No such file: " + fileName))
		return
	}
	defer file.Close()

	io.Copy(w, file)
}
