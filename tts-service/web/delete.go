package web

import (
	"net/http"
)

func onDeleteByIdRequest(h getOrDeleteHandling, w http.ResponseWriter, r *http.Request) {

	//Invoke service
	id, err := getId(h.pathPrefix, r)
	if err != nil {
		handleError(err, w, r)
		return
	}

	serviceErr := h.service.Delete(id)

	if serviceErr != nil {
		handleError(convertError(serviceErr), w, r)
	} else {
		addJsonHeader(w)
		w.WriteHeader(http.StatusNoContent) // action enacted, response does not include an entity
	}
}
