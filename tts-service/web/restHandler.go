package web

import (
	"encoding/json"
	"net/http"

	"github.com/SAPHybrisGliwice/golang-part-2/tts-service/service"
	"github.com/SAPHybrisGliwice/golang-part-2/tts-service/tts"
	"strings"
)

//We pass ServeMux explicitly to be able to unit-test in isolation.
func New(mux *http.ServeMux, ttsService service.TtsService, engine *tts.Engine, selfUrl string) {

	const createPathPrefix = "/voiceMessages"
	const getPathPrefix = "/voiceMessages/"
	const mediaPathPrefix = "/media/"

	//Allows to construct URL to media given it's ID
	mediaUrl := func(mediaId string) string {
		return selfUrl + mediaPathPrefix + mediaId
	}

	create := createHandling{createPathPrefix, ttsService, mediaUrl}
	get := getHandling{getPathPrefix, ttsService, mediaUrl}
	media := mediaHandling{mediaPathPrefix, engine}

	//Second argument must be a http.HandlerFunc Function!
	mux.HandleFunc(create.pathPrefix, create.handle)
	mux.HandleFunc(get.pathPrefix, get.handle)
	mux.HandleFunc(media.pathPrefix, media.handle)

	//Handle simple UI
	mux.HandleFunc("/public/", uiHandler)
}

type mediaUrlFunc func(string) string

// CREATE HANDLING
type createHandling struct {
	pathPrefix string
	service    service.TtsService
	mediaUrl   mediaUrlFunc
}

func (h createHandling) handle(w http.ResponseWriter, r *http.Request) {
	onCreateRequest(h, w, r)
	//TODO WEB 1: Handle correct method
}

// GET HANDLING
type getHandling struct {
	pathPrefix string
	service    service.TtsService
	mediaUrl   mediaUrlFunc
}

func (h getHandling) handle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		onGetByIdRequest(h, w, r)
	default:
		onMethodNotSupported([]string{"GET"}, w, r)
	}
}

// MEDIA HANDLING
type mediaHandling struct {
	pathPrefix string
	engine     *tts.Engine
}

func (h mediaHandling) handle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		onGetMediaRequest(h, w, r)
	default:
		onMethodNotSupported([]string{"GET"}, w, r)
	}
}

// HELPER FUNCTIONS
func onMethodNotSupported(allowed []string, w http.ResponseWriter, r *http.Request) {

	err := ErrorDTO{
		Status:  http.StatusMethodNotAllowed,
		Message: "Method Not Allowed. Supported methods: [" + strings.Join(allowed, ",") + "]",
	}
	handleError(err, w, r)
}

//Name of the function emphasizes that we are explicitly handling the error in our code
func handleError(err error, w http.ResponseWriter, r *http.Request) {
	//Type assertion: Detect custom application error(h)
	message, ok := err.(ErrorDTO)
	if ok {
		addJsonHeader(w)
		w.WriteHeader(message.Status)
		json.NewEncoder(w).Encode(message)
	} else {
		//Generic error, should be handled somehow by our code. If "escaped" to this point, it'h an internal application bug.
		sendErrorResponse(500, "Internal Server Error: "+err.Error(), w)
	}
}

func sendErrorResponse(statusCode int, message string, w http.ResponseWriter) {
	w.WriteHeader(statusCode)
	w.Write([]byte(message))
}

func addJsonHeader(w http.ResponseWriter) {
	headers := w.Header()
	headers.Set("Content-Type", "application/json")
}

//Converts result object from the service into REST representation
func toResultDTO(s *service.TtsResult, mediaUrl mediaUrlFunc) *ResultDTO {
	//Convert result to REST format
	r := ResultDTO{}
	r.createWith(s, mediaUrl)
	return &r
}

//Tries to convert Service Error to ErrorDTO.
//If succesfull, returns ErrorDTO instance.
//Otherwise returns err argument unchanged
func convertError(err error) error {

	onf, ok := err.(service.ObjectNotFoundError)
	if ok {
		return ErrorDTO{
			Status:  404,
			Message: onf.Message,
		}
	}

	//Unknown error
	return err
}

//Extracts id from path
func getId(pathPrefix string, r *http.Request) (string, error) {
	path := r.URL.Path

	id := strings.TrimPrefix(path, pathPrefix)

	if id == "" {
		return "", ErrorDTO{
			Status:  400,
			Message: "Missing ID in request path",
		}
	}

	return id, nil
}
