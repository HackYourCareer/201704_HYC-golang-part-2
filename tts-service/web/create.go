package web

import (
	"net/http"
	"github.com/SAPHybrisGliwice/golang-part-2/tts-service/service"
	"errors"
)

func onCreateRequest(h createHandling, w http.ResponseWriter, r *http.Request) {

	//1 read create DTO
	createDTO, parseErr := readCreateDTO(r)
	if parseErr != nil {
		handleError(parseErr, w, r)
		return
	}

	//2 validate create DTO
	ttsCreate, validateErr := validateCreateDTO(createDTO)
	if validateErr != nil {
		handleError(validateErr, w,r)
		return
	}

	//3 Invoke service
	_, serviceErr := h.service.Create(ttsCreate)
	if serviceErr != nil {
		handleError(serviceErr, w, r)
		return
	}

	//4 send the result
}

func readCreateDTO(r *http.Request) (*CreateDTO, error) {
	//1 validate content-type
	//2 validate body is not null
	//3 parse input as json

	return nil, errors.New("not implemented")
}

func validateCreateDTO(dto *CreateDTO) (*service.TtsCreate, error) {
	//1 validate text is not empty
	//2 validate language

	return nil, errors.New("not implemented")
}


const errInvalidContentType = "Invalid Content-Type. Only application/json is supported"
const errEmptyBody = "Request body must not be empty"
const errJsonParse = "Can't read json data: "
const errEmptyText = "Text is empty"
const errUnsupportedLang = "Unsupported Language: "
const errInvalidPayload = "Invalid payload"