package web

import (
	"net/http"
	"os"
	"io"
)

func uiHandler(w http.ResponseWriter, r *http.Request) {
	indexFilePath := os.Getenv("INDEX_FILE_PATH")

	if len(indexFilePath) == 0 {
		w.WriteHeader(404)
		w.Write([]byte("Not Found (INDEX_FILE_PATH variable not set)."))
	} else {
		file, err := os.Open(indexFilePath)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
		}
		defer file.Close()
		io.Copy(w, file)
	}
}