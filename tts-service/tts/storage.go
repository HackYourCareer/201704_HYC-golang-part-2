package tts

import (
	"io"
	"log"
	"os"
	"strconv"
	"time"
)

type storage interface {

	// Save saves the media.
	// It returns the new media ID and an error, if any.
	Save(data io.Reader) (string, error)

	// Get retrieves the media by its ID.
	// It returns an io.ReadCloser and an error, if any.
	Get(id string) (io.ReadCloser, error)

	// Delete removes the media by its ID.
	// It returns an error, if any.
	Delete(id string) error
}

// Local file system based implementation of the storage interface
type fileSystemStorage struct {
	baseDir string
}

// https://golang.org/doc/faq#methods_on_values_or_pointers
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

func (s fileSystemStorage) Get(id string) (io.ReadCloser, error) {

	path := s.createPathFor(id)

	file, err := os.Open(path)

	if err != nil {

		return nil, err
	}

	return file, nil
}

func (s fileSystemStorage) Delete(id string) error {

	path := s.createPathFor(id)
	return os.Remove(path)
}

// Constructor for the fileSystemStorage
func newFileSystemStorage() *fileSystemStorage {

	value := os.Getenv("TTS_BASE_DIR")

	if len(value) == 0 {

		value = os.TempDir()
		log.Printf("TTS_BASE_DIR not provided. Using %s", value)
	}

	return &fileSystemStorage{baseDir: value}
}

// We need to distinguish between different path separators (Windows, Linux)
const separator = string(os.PathSeparator)

func (s fileSystemStorage) createPathFor(id string) string {

	return s.baseDir + separator + id
}
