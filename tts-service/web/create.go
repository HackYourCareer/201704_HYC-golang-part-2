package web

import (
	"net/http"

	"github.com/SAPHybrisGliwice/golang-part-2/tts-service/service"
	"errors"
	"log"
)

func onCreateRequest(h createHandling, w http.ResponseWriter, r *http.Request) {

	createDTO, inputErr := readCreateDTO(r)
	if inputErr != nil {
		handleError(inputErr, w, r)
		return
	}

	ttsCreate, validationErr := validateCreateDTO(createDTO)
	if validationErr != nil {
		handleError(validationErr, w, r)
		return
	}

	//Invoke service
	result, serviceErr := h.service.Create(*ttsCreate)
	if serviceErr != nil {
		message := ErrorDTO{500, serviceErr.Error(), nil}
		handleError(message, w, r)
		return
	}

	log.Printf("result: %v", result)
}

func readCreateDTO(r *http.Request) (*CreateDTO, error) {
	return nil, errors.New("Not implemented")
}

func validateCreateDTO(dto *CreateDTO) (*service.TtsCreate, error) {
	return nil, errors.New("Not implemented")
}
