package tts

import "io"

// Engine aggregates converter and storage types.
// It is supposed to be used in other packages.
type Engine struct {
	crt converter
	str storage
}

// Process converts a given data to an audio media.
// It returns a media ID or an error, if any.
func (e Engine) Process(text string, meta Metadata) (string, error) {

	r, err := e.crt.Convert(text, meta)
	if err != nil {
		return "", err
	}

	defer r.Close()

	id, err := e.str.Save(r)
	if err != nil {
		return "", err
	}

	return id, nil
}

// Delete removes audio media from storage.
func (e Engine) Delete(id string) error {
	return e.str.Delete(id)
}

// Result returns the processing result based on its ID.
// It returns an io.ReadCloser of an error, if any.
func (e Engine) Result(id string) (io.ReadCloser, error) {

	return e.str.Get(id)
}

// https://golang.org/doc/effective_go.html#composite_literals
func NewEngine() *Engine {

	return &Engine{crt: newVoiceRssConverter(), str: newFileSystemStorage()}
}

type Metadata struct {
	Lang string
}
