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

	Convey("Process method", t, func(c C) {

		Convey("should return an ID in case of a success", func() {

			//TODO
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
