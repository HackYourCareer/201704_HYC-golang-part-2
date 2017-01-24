package web

import (
	"encoding/json"
	"net/http"
)

func onGetByIdRequest(h getHandling, w http.ResponseWriter, r *http.Request) {

	//Invoke service
	id, err := getId(h.pathPrefix, r)
	if err != nil {
		handleError(err, w, r)
		return
	}

	result, serviceErr := h.service.Get(id)

	if serviceErr != nil {
		handleError(convertError(serviceErr), w, r)
	} else {
		addJsonHeader(w)
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(toResultDTO(result, h.mediaUrl))
	}
}
