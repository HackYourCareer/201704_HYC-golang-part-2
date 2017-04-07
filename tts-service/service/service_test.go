package service

import (
	"errors"
	"github.com/SAPHybrisGliwice/golang-part-2/tts-service/tts"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestService(t *testing.T) {
	Convey("TTS Service", t, func(c C) {

		Convey("'generateId' function should generate ID based on text and language", func() {
			text1 := "Hello , World"
			text2 := "  hELLO,wORLD  "
			text3 := "Hello World"

			res1en := generateId(text1, "EN")
			res1pl := generateId(text1, "PL")
			res2en := generateId(text2, "EN")
			res3en := generateId(text3, "EN")

			So(res1en, ShouldNotEqual, res1pl)
			So(res2en, ShouldEqual, res1en)
			So(res3en, ShouldNotEqual, res1en)
		})

		Convey("Get by Id should return an error if not exists", func() {
			//given
			mock := mock("abc", ttsData{"Hello,World", "EN", StatusPending.String(), ""})
			s := New(mock, mock)

			//when
			data, err := s.Get("def")

			//then
			So(data, ShouldBeNil)
			So(err, ShouldNotBeNil)
			notFound, ok := err.(ObjectNotFoundError)

			So(ok, ShouldBeTrue)
			So(notFound.Message, ShouldEqual, "TTS with ID: 'def' doesn't exist")

		})

		Convey("Get by Id should return an object if exists", func() {
			//given
			mock := mock("abc", ttsData{"Hello,World", "EN", StatusPending.String(), ""})
			s := New(mock, mock)

			//when
			data, err := s.Get("abc")

			//then
			So(err, ShouldBeNil)
			So(data, ShouldNotBeNil)

			So(data.Id, ShouldEqual, "abc")
			So(data.Text, ShouldEqual, "Hello,World")
			So(data.Language, ShouldEqual, EN)
			So(data.Status, ShouldEqual, StatusPending)
			So(data.MediaId, ShouldEqual, "")
		})

		Convey("Create should store the object, then invoke the engine, then update mediaId and status", func() {
			const mediaId = "newAudio345"
			const text = "Hello, TTS"
			actions := []string{}

			//given
			mock := mock("", ttsData{}) //Notice no initial data
			mock.mediaIdToGenerate = mediaId
			s := New(mock, mock)

			//when
			res, err := s.Create(&TtsCreate{text, EN})

			//then after Create
			So(err, ShouldBeNil)
			So(res, ShouldNotBeNil)
			id := res.Id
			So(id, ShouldNotBeEmpty)
			assertCommonValues(res, text, EN, StatusPending, "")

			//Ensure all operations in the backgrounds completed...
			actions = readBlocking(actions, mock.recordChan)
			actions = readBlocking(actions, mock.recordChan)
			actions = readBlocking(actions, mock.recordChan)

			//then after Get
			res, err = s.Get(id)
			So(err, ShouldBeNil)
			So(res, ShouldNotBeNil)
			So(res.Id, ShouldEqual, id)
			assertCommonValues(res, text, EN, StatusReady, mediaId)

			//Verify interaction
			So(actions[0], ShouldEqual, "persistence.create")
			So(actions[1], ShouldEqual, "tts.Engine.Process")
			So(actions[2], ShouldEqual, "persistence.update")
		})

		Convey("Create should update status on media generation failure", func() {
			const text = "Hello, TTS"
			actions := []string{}

			//given
			mock := mock("", ttsData{}) //Notice no initial data
			mock.mediaIdToGenerate = "" //Indicates that mock media engine should generate an error
			s := New(mock, mock)

			//when
			res, err := s.Create(&TtsCreate{text, EN})

			//then after Create
			So(err, ShouldBeNil)
			So(res, ShouldNotBeNil)
			id := res.Id
			So(id, ShouldNotBeEmpty)
			assertCommonValues(res, text, EN, StatusPending, "")

			//Ensure all operations in the backgrounds completed...
			actions = readBlocking(actions, mock.recordChan)
			actions = readBlocking(actions, mock.recordChan)
			actions = readBlocking(actions, mock.recordChan)

			//then after Get
			res, err = s.Get(id)
			So(err, ShouldBeNil)
			So(res, ShouldNotBeNil)
			So(res.Id, ShouldEqual, id)
			assertCommonValues(res, text, EN, StatusError, "")

			//Verify interaction
			So(actions[0], ShouldEqual, "persistence.create")
			So(actions[1], ShouldEqual, "tts.Engine.Process")
			So(actions[2], ShouldEqual, "persistence.update")
		})

		Convey("Create should not start media generation on create failure", func() {
			const text = "Boom!"
			actions := []string{}

			//given
			mock := mock("", ttsData{}) //Notice no initial data
			mock.ttsTextThatFails = text
			s := New(mock, mock)

			//when
			res, err := s.Create(&TtsCreate{text, EN})

			//then after Create
			So(err, ShouldNotBeNil)
			So(res, ShouldBeNil)

			//Ensure only "create" operation has been performed
			actions = readBlocking(actions, mock.recordChan)
			So(actions[0], ShouldEqual, "persistence.create")

			actions, ok := readNonBlocking(actions, mock.recordChan)
			So(ok, ShouldBeFalse)
			So(len(actions), ShouldEqual, 1)
		})

		Convey("Create should return existing object on conflict", func() {
			const text = "Hello, TTS"
			const id = "15f3f83eec955266793622006b0f66a47398f3b1"
			actions := []string{}

			//given
			mock := mock(id, ttsData{text, "EN", StatusReady.String(), "mediaId#123"})
			mock.ttsTextThatConflicts = text
			s := New(mock, mock)

			//when
			res, err := s.Create(&TtsCreate{text, EN})

			//then after Create
			So(err, ShouldBeNil)
			So(res, ShouldNotBeNil)
			So(res.Id, ShouldEqual, id)
			assertCommonValues(res, text, EN, StatusReady, "mediaId#123")

			//Ensure all operations in the backgrounds completed...
			actions = readBlocking(actions, mock.recordChan)
			actions = readBlocking(actions, mock.recordChan)

			//Verify interaction
			So(actions[0], ShouldEqual, "persistence.create")
			So(actions[1], ShouldEqual, "persistence.get")
		})
	})
}

//Creates a mock that fulfills the contract of both: service.ttsPersistence and tts.Engine to test interaction.
//id, data - params to pre-fill the persistence mock
func mock(id string, data ttsData) *interactionMock {
	m := interactionMock{id: id, data: data}
	m.recordChan = make(chan string, 5)

	return &m
}

//Mock object used to verify correct interaction between service.persistence and tts.Engine inside service
//This mock implement both service.TtsPersistence and MediaEngine interfaces.
type interactionMock struct {
	id   string  //tts id
	data ttsData //tts persistence data

	mediaIdToGenerate    string //if empty, return error from tts.Engine.Process
	ttsTextThatFails     string //if invoked with this text, simulate persistence failure
	ttsTextThatConflicts string //if invoked with this text, return ObjectAlreadyExistsError

	recordChan chan string
}

//service.TtsPersistence contract
func (mp *interactionMock) create(id string, data ttsData) error {
	mp.recordChan <- "persistence.create"

	if data.Text == mp.ttsTextThatFails {
		return errors.New("Persistence Failure")
	}

	if data.Text == mp.ttsTextThatConflicts {
		return AlreadyExists(id)
	}

	mp.id = id
	mp.data = data
	return nil
}
func (mp *interactionMock) get(id string) (*ttsData, error) {
	mp.recordChan <- "persistence.get"

	if mp.id != id {
		return nil, NotFound(id)
	} else {
		return &mp.data, nil
	}
}
func (mp *interactionMock) update(id string, status string, mediaId string) error {
	mp.recordChan <- "persistence.update"

	if mp.id != id {
		return NotFound(id)
	} else {
		mp.data.Status = status
		mp.data.MediaId = mediaId

		return nil
	}
}

func (mp *interactionMock) del(id string) error {
	if mp.id != id {
		return NotFound(id)
	} else {
		return nil
	}
}

//Implements MediaEngine interface
func (mp *interactionMock) Process(text string, meta tts.Metadata) (string, error) {
	mp.recordChan <- "tts.Engine.Process"

	if mp.mediaIdToGenerate == "" {
		//Simulate error
		return "", errors.New("Network Unreachable")
	}
	return mp.mediaIdToGenerate, nil
}

func readBlocking(source []string, recordChan chan string) []string {
	s := <-recordChan
	return append(source, s)
}

func readNonBlocking(source []string, recordChan chan string) ([]string, bool) {
	select {
	case s := <-recordChan:
		res := append(source, s)
		return res, true
	default:
		return source, false
	}
}

func assertCommonValues(res *TtsResult, text string, lang LangEnum, status StatusEnum, mediaId string) {
	So(res.Text, ShouldEqual, text)
	So(res.Language, ShouldEqual, lang)
	So(res.Status, ShouldEqual, status)
	So(res.MediaId, ShouldEqual, mediaId)
}
