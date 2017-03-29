package tts

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestConverter(t *testing.T) {

	Convey("VoiceRSS based converter", t, func(c C) {

		Convey("should convert a text to a speech", func() {

			f, _ := os.Open("testdata" + string(os.PathSeparator) + "test")
			defer f.Close()

			fc, _ := ioutil.ReadAll(f)

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

				io.Copy(w, bytes.NewBuffer(fc))
			}))
			defer server.Close()

			converter := &voiceRssConverter{apiUrl: server.URL}

			r, err := converter.Convert("This is just a simple test", Metadata{})

			So(err, ShouldBeNil)

			rc, _ := ioutil.ReadAll(r)

			So(rc, ShouldResemble, fc)
		})

		Convey("should return an error in case of an internal error", func() {

			converter := &voiceRssConverter{apiUrl: ""}

			_, err := converter.Convert("whatever", Metadata{})

			So(err, ShouldNotBeNil)
		})

		Convey("should return an error in case of an unexpected response code", func() {

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusBadRequest)
			}))
			defer server.Close()

			converter := &voiceRssConverter{apiUrl: server.URL}

			_, err := converter.Convert("whatever", Metadata{})

			So(err, ShouldNotBeNil)
		})

		Convey("should return an error in case of an unexpected response content type", func() {

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

				io.Copy(w, strings.NewReader("whatever"))
			}))
			defer server.Close()

			converter := &voiceRssConverter{apiUrl: server.URL}

			_, err := converter.Convert("whatever", Metadata{})

			So(err, ShouldNotBeNil)
		})
	})
}
