package tts

import (
	"io"
)

type Engine struct {
	crt Converter
	str Storage
}

//Process converts a given data to an audio media.
// It returns a media ID or an error, if any.
func (e Engine) Process(text string, meta Metadata) (string, error) {

	result, err := e.crt.Convert(text, meta)

	if err != nil {

		return "", err
	}

	id, err := e.str.Save(result)

	if err != nil {

		return "", err
	}

	return id, nil
}

//GetResult returns the processing result based on its ID.
//It returns an io.Reader of an error, if any.
func (e Engine) GetResult(id string) (io.Reader, error) {

	return e.str.Get(id)
}

//https://golang.org/doc/effective_go.html#composite_literals
func NewEngine() *Engine {

	c := newVoiceRssConverter()
	s := newFileSystemStorage()

	return newEngine(c, s)
}

func newEngine(c Converter, s Storage) *Engine {

	return &Engine{crt: c, str: s}
}

type Metadata struct {
	Lang string
}
