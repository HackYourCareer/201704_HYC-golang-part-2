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

	reader, err := h.engine.GetResult(id)

	if err != nil {
		handleError(err, w, r)
		return
	}

	io.Copy(w, reader)
}
