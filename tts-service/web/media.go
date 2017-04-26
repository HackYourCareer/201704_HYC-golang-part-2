package web

import (
	"io"
	"net/http"
)

func onGetMediaRequest(h mediaHandling, w http.ResponseWriter, r *http.Request) {

	//Invoke service
	id, err := getId(h.pathPrefix, r)
	if err != nil {
		handleError(err, w, r)
		return
	}

	reader, err := h.engine.Result(id)

	if err != nil {
		handleError(err, w, r)
		return
	}

	io.Copy(w, reader)
	reader.Close() // media file has to be closed so it can be deleted on windows
}
