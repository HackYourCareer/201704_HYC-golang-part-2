package tts

import (
	"io"
	"log"
	"os"
	"strconv"
	"time"
)

type Storage interface {

	//Save saves the media.
	//It returns the new media ID and an error, if any.
	Save(data io.Reader) (string, error)

	//Get retrieves the media by its ID.
	//It returns an io.Reader and an error, if any.
	Get(id string) (io.Reader, error)

	//Delete removes the media by its ID.
	//It returns an error, if any.
	Delete(id string) error
}

type fileSystemStorage struct {
	baseDir string
}

const separator = string(os.PathSeparator)

func newFileSystemStorage() *fileSystemStorage {

	value := os.Getenv("TTS_BASE_DIR")

	if len(value) == 0 {

		value = os.TempDir()
		log.Printf("TTS_BASE_DIR not provided. Using %s", value)
	}

	return &fileSystemStorage{baseDir: value}
}

//https://golang.org/doc/faq#methods_on_values_or_pointers
func (s fileSystemStorage) Save(data io.Reader) (string, error) {

	id := strconv.FormatInt(time.Now().UnixNano(), 10)
	path := s.createPathFor(id)

	file, err := os.Create(path)

	if err != nil {

		return "", err
	}

	defer file.Close()

	_, err = io.Copy(file, data)

	if err != nil {

		return "", err
	}

	return id, nil
}

func (s fileSystemStorage) Get(id string) (io.Reader, error) {

	path := s.createPathFor(id)

	file, err := os.Open(path)

	if err != nil {

		return nil, err
	}

	r, w := io.Pipe()

	go func() {
		defer file.Close()
		defer w.Close()

		io.Copy(w, file)
	}()

	return r, nil
}

func (s fileSystemStorage) Delete(id string) error {

	path := s.createPathFor(id)
	return os.Remove(path)
}

func (s fileSystemStorage) createPathFor(id string) string {

	return s.baseDir + separator + id
}
