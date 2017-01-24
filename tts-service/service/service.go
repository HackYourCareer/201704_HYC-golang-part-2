package service

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"github.com/SAPHybrisGliwice/golang-part-2/tts-service/tts"
	"strings"
	"errors"
)

//Public API

type TtsService interface {
	Create(create TtsCreate) (*TtsResult, error)
	Get(ID string) (*TtsResult, error)
}

//Interface abstracting over tts.Engine
type MediaEngine interface {
	Process(text string, meta tts.Metadata) (string, error)
}

func New(persistence TtsPersistence, engine MediaEngine) TtsService {
	return impl{
		persistence: persistence,
		ttsEngine:   engine,
	}
}

//Implementation

type impl struct {
	persistence TtsPersistence
	ttsEngine   MediaEngine
}

func (srv impl) Create(create TtsCreate) (*TtsResult, error) {

	if create.Text == "" {
		return nil, errors.New("Cannot create: Text is empty")
	}

	id := generateId(create.Text, create.Language.String())

	initialStatus := StatusPending
	mediaId := ""

	//Save TTS definition data in the persistent store
	err := srv.persistence.create(id, ttsData{
		Text:     create.Text,
		Language: create.Language.String(),
		Status:   initialStatus.String(),
		MediaId:  mediaId,
	})

	if err != nil {
		//In case of conflict, just return already existing object
		_, ok := err.(ObjectAlreadyExistsError)
		if ok {
			return srv.Get(id)
		}

		//Propagate other errors
		return nil, err
	}

	res := TtsResult{
		Id:       id,
		Text:     create.Text,
		Language: create.Language,
		Status:   initialStatus,
		MediaId:  mediaId,
	}

	//Generate Media in the background
	go srv.generateMedia(res.Id, res.Text, res.Language)

	return &res, nil
}

func (srv impl) Get(id string) (*TtsResult, error) {

	data, err := srv.persistence.get(id)
	if err != nil {
		return nil, err
	}

	return &TtsResult{
		Id:       id,
		Text:     data.Text,
		Language: lang(data.Language),
		Status:   status(data.Status),
		MediaId:  data.MediaId,
	}, nil
}

func (srv impl) generateMedia(id, text string, language LangEnum) {

	metadata := tts.Metadata{
		Lang: language.String(),
	}

	mediaId, mediaErr := srv.ttsEngine.Process(text, metadata)

	if mediaErr == nil {
		srv.persistence.update(id, StatusReady.String(), mediaId)
	} else {
		fmt.Printf("Problem with TTS(id: %v) - an Error occured during media generation: %v\n", id, mediaErr)
		srv.persistence.update(id, StatusError.String(), "")
	}
}

func generateId(text string, language string) string {
	baseStr := strings.ToLower(strings.Replace(text, " ", "", -1) + language)
	sha1Sum := sha1.Sum([]byte(baseStr))
	encoded := hex.EncodeToString(sha1Sum[:])
	return encoded
}
