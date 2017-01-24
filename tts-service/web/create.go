package web

import (
	"net/http"

	"github.com/SAPHybrisGliwice/golang-part-2/tts-service/service"
	"errors"
	"encoding/json"
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
	} else {
		resultDTO := toResultDTO(result, h.mediaUrl)
		json.NewEncoder(w).Encode(resultDTO)
	}
}

//Parses input data. Returns an error if unable to parse for any reason.
func readCreateDTO(r *http.Request) (*CreateDTO, error) {
	return nil, errors.New("Not implemented")
}

//Performs validation of the data.
//If data is valid, returns service.TtsCreate object.
//An error is returned otherwise.
func validateCreateDTO(dto *CreateDTO) (*service.TtsCreate, error) {
	var details []string

	langEnum := service.EN

	if len(details) == 0 {
		return &service.TtsCreate{dto.Text, langEnum}, nil
	} else {
               return nil, ErrorDTO{400, "Invalid payload", details}
	}
}
