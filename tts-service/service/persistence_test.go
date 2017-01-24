package service

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestPersistence(t *testing.T) {
	Convey("TTS Persistence", t, func(c C) {

		Convey("should create, read and delete the file", func() {
			persistence := NewPersistence()
			id := "test1"

			err := persistence.create(id, ttsData{
				Text:     "test text 1",
				Language: EN.String(),
				Status:   StatusPending.String(),
				MediaId:  "",
			})
			defer persistence.del(id)

			So(err, ShouldBeNil)

			data, err := persistence.get(id)
			So(err, ShouldBeNil)
			So(data, ShouldNotBeNil)
			So(data.Text, ShouldEqual, "test text 1")
			So(data.Language, ShouldEqual, EN.String())
			So(data.Status, ShouldEqual, StatusPending.String())
			So(data.MediaId, ShouldEqual, "")
		})

		Convey("should return ObjectNotFoundError for non-existing files", func() {
			persistence := NewPersistence()
			id := "test2"

			data, err := persistence.get(id)
			So(err, ShouldNotBeNil)
			So(data, ShouldBeNil)

			_, ok := err.(ObjectNotFoundError)
			So(ok, ShouldBeTrue)
			So(err.Error(), ShouldEqual, "TTS with ID: '"+id+"' doesn't exist")
		})

		Convey("should return ObjectAlreadyExistsError for existing files", func() {
			persistence := NewPersistence()
			id := "test3"

			data := ttsData{
				Text:     "test text 3",
				Language: EN.String(),
				Status:   StatusPending.String(),
				MediaId:  "audio123",
			}

			err := persistence.create(id, data)

			defer persistence.del(id)
			So(err, ShouldBeNil)

			err = persistence.create(id, ttsData{
				Text:     "test text 3",
				Language: EN.String(),
				Status:   StatusPending.String(),
				MediaId:  "audio123",
			})

			So(err, ShouldNotBeNil)

			_, ok := err.(ObjectAlreadyExistsError)
			So(ok, ShouldBeTrue)
			So(err.Error(), ShouldEqual, "TTS with ID: '"+id+"' already exists")
		})

		Convey("should update existing file", func() {
			persistence := NewPersistence()
			id := "test4"

			//Create
			err := persistence.create(id, ttsData{
				Text:     "test text 4",
				Language: EN.String(),
				Status:   StatusPending.String(),
				MediaId:  "",
			})
			defer persistence.del(id)

			So(err, ShouldBeNil)

			//Get to verify
			data, err := persistence.get(id)
			So(err, ShouldBeNil)
			So(data, ShouldNotBeNil)
			So(data.Text, ShouldEqual, "test text 4")
			So(data.Language, ShouldEqual, EN.String())
			So(data.Status, ShouldEqual, StatusPending.String())
			So(data.MediaId, ShouldEqual, "")

			//Update
			err = persistence.update(id, StatusReady.String(), "media123")
			So(err, ShouldBeNil)

			//Get to verify once again
			data, err = persistence.get(id)
			So(err, ShouldBeNil)
			So(data, ShouldNotBeNil)
			So(data.Text, ShouldEqual, "test text 4")
			So(data.Language, ShouldEqual, EN.String())
			So(data.Status, ShouldEqual, StatusReady.String())
			So(data.MediaId, ShouldEqual, "media123")
		})
	})
}
