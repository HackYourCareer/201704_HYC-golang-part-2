package service

import (
	"encoding/json"
	"log"
	"os"
	"strings"
)

//The interface of TTS data persistence
type TtsPersistence interface {

	//Request to store tts data with given id
	//May return ObjectAlreadyExistsError
	create(id string, data ttsData) error

	//Returns tts data given it's id
	//May return ObjectNotFoundError
	get(id string) (*ttsData, error)

	//Updates status of tts data given it's id
	//May return ObjectNotFoundError
	update(id string, status string, mediaId string) error

	//Removes tts data given it's id
	//May return ObjectNotFoundError
	del(id string) error
}

type ttsData struct {
	Text     string
	Language string
	Status   string
	MediaId  string
}

//Initializes the persistence module
func NewPersistence() TtsPersistence {
	directory := os.Getenv("PERSISTENCE_BASE_DIR")

	if len(directory) == 0 {

		directory = os.TempDir()
		log.Printf("PERSISTENCE_BASE_DIR not provided. Using %s", directory)
	}

	return &fileBased{directory}
}

// Errors

//Returned on get/del
type ObjectNotFoundError struct {
	Message string
}

//ObjectNotFoundError implements built-in  "error" interface
func (err ObjectNotFoundError) Error() string {
	return err.Message
}

//Returned on create
type ObjectAlreadyExistsError struct {
	Message string
}

//ObjectAlreadyExistsError implements built-in  "error" interface
func (err ObjectAlreadyExistsError) Error() string {
	return err.Message
}

//Helper functions
func NotFound(id string) ObjectNotFoundError {
	return ObjectNotFoundError{"TTS with ID: '" + id + "' doesn't exist"}
}

func AlreadyExists(id string) ObjectAlreadyExistsError {
	return ObjectAlreadyExistsError{"TTS with ID: '" + id + "' already exists"}
}

// Implementation

const separator = string(os.PathSeparator)

type fileBased struct {
	directory string
}

func (fb fileBased) create(id string, data ttsData) error {
	path := fb.pathWithId(id)

	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)

	if err == nil {
		encoder := json.NewEncoder(file)
		encoder.Encode(data)
		file.Close()
		return nil
	}

	if os.IsExist(err) {
		err = AlreadyExists(id)
	}

	return err
}

func (fb fileBased) get(id string) (*ttsData, error) {
	path := fb.pathWithId(id)
	file, err := os.Open(path)

	if err == nil {
		data := &ttsData{}
		decoder := json.NewDecoder(file)
		err = decoder.Decode(&data)
		file.Close()
		return data, nil
	}

	if os.IsNotExist(err) {
		err = NotFound(id)
	}

	return nil, err

}

func (fb fileBased) update(id string, status string, mediaId string) error {
	//Read file
	data, err := fb.get(id)

	if err != nil {
		return err
	}

	//Update data
	data.Status = status
	data.MediaId = mediaId

	path := fb.pathWithId(id)

	//Write updated data
	file, err := os.OpenFile(path, os.O_RDWR, 0666)
	if err == nil {
		encoder := json.NewEncoder(file)
		encoder.Encode(data)
	}
	file.Close()
	return err
}

func (fb fileBased) del(id string) error {
	return os.Remove(fb.pathWithId(id))
}

func (fb fileBased) pathWithId(name string) string {
	var res string
	if strings.HasSuffix(fb.directory, separator) {
		res = fb.directory + name
	} else {
		res = fb.directory + separator + name
	}

	return res + ".json"
}
