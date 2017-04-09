package tts

import (
	"errors"
	. "github.com/smartystreets/goconvey/convey"
	"io"
	"io/ioutil"
	"strings"
	"testing"
)

func TestEngine(t *testing.T) {

	Convey("TTS engine", t, func(c C) {

		Convey("Process method", func() {

			Convey("should not blow if there are no errors", func() {

				engine := &Engine{mockConverter{false}, mockStorage{false}}

				id, err := engine.Process("", Metadata{})

				So(err, ShouldBeNil)
				So(id, ShouldNotBeNil)
			})

			Convey("should pass error from converter", func() {

				engine := &Engine{mockConverter{true}, mockStorage{false}}

				id, err := engine.Process("", Metadata{})

				So(id, ShouldBeNil)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, converterErrorMessage)
			})

			Convey("should pass error from storage", func() {

				engine := &Engine{mockConverter{false}, mockStorage{true}}

				id, err := engine.Process("", Metadata{})

				So(id, ShouldBeNil)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, storageErrorMessage)
			})
		})

		Convey("GetResult method", func() {

			Convey("should pass error from storage", func() {

				engine := &Engine{mockConverter{false}, mockStorage{true}}

				_, err := engine.Process("", Metadata{})

				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, storageErrorMessage)
			})

			Convey("should not blow if there are no errors", func() {

				engine := &Engine{mockConverter{false}, mockStorage{false}}

				r, err := engine.Process("", Metadata{})

				So(r, ShouldNotBeNil)
				So(err, ShouldBeNil)
			})
		})
	})
}

const converterErrorMessage = "Converter. Unexpected error."
const storageErrorMessage = "Storage. Unexpected error."

type mockConverter struct {
	failing bool
}

func (mc mockConverter) Convert(text string, metadata Metadata) (io.ReadCloser, error) {

	if mc.failing {

		return nil, errors.New(converterErrorMessage)
	}

	return ioutil.NopCloser(strings.NewReader("test")), nil
}

type mockStorage struct {
	failing bool
}

func (ms mockStorage) Save(data io.Reader) (string, error) {

	if ms.failing {

		return "", errors.New(storageErrorMessage)
	}

	return "dummyID", nil
}

func (ms mockStorage) Get(id string) (io.ReadCloser, error) {

	if ms.failing {

		return nil, errors.New(storageErrorMessage)
	}

	return ioutil.NopCloser(strings.NewReader("test")), nil
}

func (ms mockStorage) Delete(id string) error {

	return nil
}
