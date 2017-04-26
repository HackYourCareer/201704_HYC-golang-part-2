package web

import (
	"github.com/SAPHybrisGliwice/golang-part-2/tts-service/service"
	"net/http"
	"net/http/httptest"
	"testing"

	"bytes"
	"encoding/json"
	. "github.com/smartystreets/goconvey/convey"
	"strings"
)

func TestRestController(t *testing.T) {
	Convey("Rest Controller", t, func(c C) {

		const selfUrl = "http://localhost:3000"
		const rootUrl = "/voiceMessages"

		Convey("when handling GET request on /tts", func() {

			Convey("should respond with 405 (Method Not Allowed) status code", func() {
				req, err := http.NewRequest("GET", rootUrl, nil)

				if err != nil {
					t.Fatal(err)
				}

				mux := http.NewServeMux()
				New(mux, defaultMockService(), nil, selfUrl)

				//Test the request
				rr := httptest.NewRecorder()

				mux.ServeHTTP(rr, req)

				So(rr.Code, ShouldEqual, http.StatusMethodNotAllowed)

			})
		})

		Convey("when handling POST request on /tts", func() {

			Convey("should validate request Content-Type", func() {
				req, err := http.NewRequest("POST", rootUrl, nil)
				if err != nil {
					t.Fatal(err)
				}

				mux := http.NewServeMux()
				New(mux, defaultMockService(), nil, selfUrl)

				//Test the request
				rr := httptest.NewRecorder()

				mux.ServeHTTP(rr, req)

				So(rr.Code, ShouldEqual, http.StatusUnsupportedMediaType)
			})

			Convey("should validate missing body", func() {
				req, err := http.NewRequest("POST", rootUrl, nil)
				if err != nil {
					t.Fatal(err)
				}
				req.Header.Set("Content-Type", "application/json")

				mux := http.NewServeMux()
				New(mux, defaultMockService(), nil, selfUrl)

				//Test the request
				rr := httptest.NewRecorder()

				mux.ServeHTTP(rr, req)

				So(rr.Code, ShouldEqual, http.StatusBadRequest)
			})

			Convey("should validate empty body", func() {
				req, err := http.NewRequest("POST", rootUrl, strings.NewReader(" "))
				if err != nil {
					t.Fatal(err)
				}
				req.Header.Set("Content-Type", "application/json")

				mux := http.NewServeMux()
				New(mux, defaultMockService(), nil, selfUrl)

				//Test the request
				rr := httptest.NewRecorder()

				mux.ServeHTTP(rr, req)

				So(rr.Code, ShouldEqual, http.StatusBadRequest)
			})

			Convey("should validate language", func() {

				//Encode JSON
				u := CreateDTO{Text: "abcdef", Language: "DE"}
				b := new(bytes.Buffer)
				json.NewEncoder(b).Encode(u)

				//Prepare request
				req, err := http.NewRequest("POST", rootUrl, b)
				if err != nil {
					t.Fatal(err)
				}
				req.Header.Set("Content-Type", "application/json")

				mux := http.NewServeMux()
				New(mux, defaultMockService(), nil, selfUrl)

				//Test the request
				rr := httptest.NewRecorder()

				mux.ServeHTTP(rr, req)

				So(rr.Code, ShouldEqual, http.StatusBadRequest)
				const expected = `{"status":400,"message":"Invalid payload","details":["Unsupported Language: DE"]}` + "\n"
				So(string(rr.Body.String()), ShouldEqual, expected)
			})

			Convey("should require Text", func() {

				//Encode JSON
				u := CreateDTO{Text: "", Language: "PL"}
				b := new(bytes.Buffer)
				json.NewEncoder(b).Encode(u)

				//Prepare request
				req, err := http.NewRequest("POST", rootUrl, b)
				if err != nil {
					t.Fatal(err)
				}
				req.Header.Set("Content-Type", "application/json")

				mux := http.NewServeMux()
				New(mux, defaultMockService(), nil, selfUrl)

				//Test the request
				rr := httptest.NewRecorder()

				mux.ServeHTTP(rr, req)

				So(rr.Code, ShouldEqual, http.StatusBadRequest)
				const expected = `{"status":400,"message":"Invalid payload","details":["Text is empty"]}` + "\n"
				So(string(rr.Body.String()), ShouldEqual, expected)
			})

			Convey("should validate language and text with Content-Type", func() {

				//Encode JSON
				u := CreateDTO{Text: "", Language: "DE"}
				b := new(bytes.Buffer)
				json.NewEncoder(b).Encode(u)

				//Prepare request
				req, err := http.NewRequest("POST", rootUrl, b)
				if err != nil {
					t.Fatal(err)
				}
				req.Header.Set("Content-Type", "application/json")

				mux := http.NewServeMux()
				New(mux, defaultMockService(), nil, selfUrl)

				//Test the request
				rr := httptest.NewRecorder()

				mux.ServeHTTP(rr, req)

				So(rr.Code, ShouldEqual, http.StatusBadRequest)
				So(rr.Header().Get("Content-Type"), ShouldEqual, "application/json")

				const expected = `{"status":400,"message":"Invalid payload","details":["Text is empty","Unsupported Language: DE"]}` + "\n"
				So(string(rr.Body.String()), ShouldEqual, expected)
			})

			Convey("should correctly handle proper json body", func() {

				//Encode JSON
				u := CreateDTO{Text: "abcdef", Language: "EN"}
				b := new(bytes.Buffer)
				json.NewEncoder(b).Encode(u)

				//Prepare request
				req, err := http.NewRequest("POST", rootUrl, b)
				if err != nil {
					t.Fatal(err)
				}
				req.Header.Set("Content-Type", "application/json")

				mux := http.NewServeMux()
				New(mux, defaultMockService(), nil, selfUrl)

				//Test the request
				rr := httptest.NewRecorder()

				mux.ServeHTTP(rr, req)

				So(rr.Code, ShouldEqual, http.StatusAccepted)
				const expected = `{"id":"abc123","text":"Received: abcdef","language":"EN","status":"PENDING"}` + "\n"
				So(string(rr.Body.String()), ShouldEqual, expected)
			})

			Convey("should correctly return json with MediaUrl and correct Content-Type", func() {

				//Encode JSON
				u := CreateDTO{Text: "abcdef", Language: "EN"}
				b := new(bytes.Buffer)
				json.NewEncoder(b).Encode(u)

				//Prepare request
				req, err := http.NewRequest("POST", rootUrl, b)
				if err != nil {
					t.Fatal(err)
				}
				req.Header.Set("Content-Type", "application/json")

				mux := http.NewServeMux()
				New(mux, getMockService("123", service.StatusReady), nil, selfUrl)

				//Test the request
				rr := httptest.NewRecorder()

				mux.ServeHTTP(rr, req)

				So(rr.Code, ShouldEqual, http.StatusAccepted)
				So(rr.Header().Get("Content-Type"), ShouldEqual, "application/json")

				const expected = `{"id":"abc123","text":"Received: abcdef","language":"EN","status":"READY","mediaUrl":"` + selfUrl + `/media/123"}` + "\n"
				So(string(rr.Body.String()), ShouldEqual, expected)
			})

		})

		Convey("when handling GET request on /voiceMessages/{ID}", func() {

			Convey("should require ID value", func() {
				req, err := http.NewRequest("GET", rootUrl+"/", nil)
				if err != nil {
					t.Fatal(err)
				}

				mux := http.NewServeMux()
				New(mux, defaultMockService(), nil, selfUrl)

				//Test the request
				rr := httptest.NewRecorder()

				mux.ServeHTTP(rr, req)

				So(rr.Code, ShouldEqual, http.StatusBadRequest)
				const expected = `{"status":400,"message":"Missing ID in request path"}` + "\n"
				So(string(rr.Body.String()), ShouldEqual, expected)
			})

			Convey("should return result by ID", func() {
				req, err := http.NewRequest("GET", rootUrl+"/cafe", nil)
				if err != nil {
					t.Fatal(err)
				}

				mux := http.NewServeMux()
				New(mux, defaultMockService(), nil, selfUrl)

				//Test the request
				rr := httptest.NewRecorder()

				mux.ServeHTTP(rr, req)

				So(rr.Code, ShouldEqual, http.StatusOK)
				const expected = `{"id":"cafe","text":"coffee'h good","language":"EN","status":"PENDING"}` + "\n"
				So(string(rr.Body.String()), ShouldEqual, expected)
			})

			Convey("should return 404 for non-existing TTS", func() {
				req, err := http.NewRequest("GET", rootUrl+"/tea", nil)
				if err != nil {
					t.Fatal(err)
				}

				mux := http.NewServeMux()
				New(mux, defaultMockService(), nil, selfUrl)

				//Test the request
				rr := httptest.NewRecorder()

				mux.ServeHTTP(rr, req)

				So(rr.Code, ShouldEqual, http.StatusNotFound)
				const expected = `{"status":404,"message":"TTS with ID: 'tea' doesn't exist"}` + "\n"
				So(string(rr.Body.String()), ShouldEqual, expected)
			})

		})

		Convey("when handling DELETE request on /voiceMessages/{ID}", func() {

			Convey("should require ID value", func() {
				req, err := http.NewRequest("DELETE", rootUrl+"/", nil)
				if err != nil {
					t.Fatal(err)
				}

				mux := http.NewServeMux()
				New(mux, defaultMockService(), nil, selfUrl)

				//Test the request
				rr := httptest.NewRecorder()

				mux.ServeHTTP(rr, req)

				So(rr.Code, ShouldEqual, http.StatusBadRequest)
				const expected = `{"status":400,"message":"Missing ID in request path"}` + "\n"
				So(string(rr.Body.String()), ShouldEqual, expected)
			})

			Convey("should return 402 after enacting deletion", func() {
				req, err := http.NewRequest("DELETE", rootUrl+"/cafe", nil)
				if err != nil {
					t.Fatal(err)
				}

				mux := http.NewServeMux()
				New(mux, defaultMockService(), nil, selfUrl)

				//Test the request
				rr := httptest.NewRecorder()

				mux.ServeHTTP(rr, req)

				So(rr.Code, ShouldEqual, http.StatusNoContent)
				const expected = ``
				So(string(rr.Body.String()), ShouldEqual, expected)
			})

			Convey("should return 404 for non-existing TTS", func() {
				req, err := http.NewRequest("DELETE", rootUrl+"/tea", nil)
				if err != nil {
					t.Fatal(err)
				}

				mux := http.NewServeMux()
				New(mux, defaultMockService(), nil, selfUrl)

				//Test the request
				rr := httptest.NewRecorder()

				mux.ServeHTTP(rr, req)

				So(rr.Code, ShouldEqual, http.StatusNotFound)
				const expected = `{"status":404,"message":"TTS with ID: 'tea' doesn't exist"}` + "\n"
				So(string(rr.Body.String()), ShouldEqual, expected)
			})

		})
	})
}

// Mocks for service.TtsService
func defaultMockService() service.TtsService {
	return getMockService("", service.StatusPending)
}

func getMockService(mediaId string, status service.StatusEnum) service.TtsService {
	return mockService{status, mediaId}
}

type mockService struct {
	status  service.StatusEnum
	mediaId string
}

func (s mockService) Create(create *service.TtsCreate) (*service.TtsResult, error) {
	res := service.TtsResult{
		Id:       "abc123",
		Text:     "Received: " + create.Text,
		Language: create.Language,
		Status:   s.status,
		MediaId:  s.mediaId,
	}
	return &res, nil
}

func (s mockService) Get(id string) (*service.TtsResult, error) {

	if id == "cafe" {
		res := service.TtsResult{
			Id:       id,
			Text:     "coffee'h good",
			Language: service.EN,
			Status:   service.StatusPending,
		}
		return &res, nil
	} else {
		return nil, service.NotFound(id)
	}
}

func (s mockService) Delete(id string) error {

	if id == "cafe" {
		return nil
	} else {
		return service.NotFound(id)
	}
}
