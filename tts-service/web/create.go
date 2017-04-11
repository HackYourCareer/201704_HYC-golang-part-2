package web

import (
	"encoding/json"
	"net/http"

	"github.com/SAPHybrisGliwice/golang-part-2/tts-service/service"
	"strings"
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
	result, serviceErr := h.service.Create(ttsCreate)
	if serviceErr != nil {
		message := ErrorDTO{http.StatusInternalServerError, serviceErr.Error(), nil}
		handleError(message, w, r)
	} else {
		addJsonHeader(w)
		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(toResultDTO(result, h.mediaUrl))
	}

}

func readCreateDTO(r *http.Request) (*CreateDTO, error) {
	var createDTO CreateDTO

	contentType := r.Header.Get("Content-Type")

	if !strings.HasPrefix(contentType, "application/json") {
		return nil, ErrorDTO{http.StatusUnsupportedMediaType, errInvalidContentType, nil}
	}

	if r.Body == nil {
		err := ErrorDTO{400, errEmptyBody, nil}
		return nil, err
	} else {
		err := json.NewDecoder(r.Body).Decode(&createDTO)
		if err == nil {
			return &createDTO, nil
		} else {
			return nil, ErrorDTO{http.StatusBadRequest, errJsonParse + err.Error(), nil}
		}
	}
}

func validateCreateDTO(dto *CreateDTO) (*service.TtsCreate, error) {
	var details []string

	if dto.Text == "" {
		details = append(details, errEmptyText)
	}

	var langEnum service.LangEnum = nil

	switch dto.Language {
	case "EN":
		langEnum = service.EN
	case "PL":
		langEnum = service.PL
	default:
		details = append(details, errUnsupportedLang+dto.Language)
	}

	if len(details) == 0 {
		return &service.TtsCreate{dto.Text, langEnum}, nil
	} else {
		return nil, ErrorDTO{http.StatusBadRequest, errInvalidPayload, details}
	}
}

const errInvalidContentType = "Invalid Content-Type. Only application/json is supported"
const errEmptyBody = "Request body must not be empty"
const errJsonParse = "Can't read json data: "
const errEmptyText = "Text is empty"
const errUnsupportedLang = "Unsupported Language: "
const errInvalidPayload = "Invalid payload"
