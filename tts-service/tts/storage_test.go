package tts

import (
	. "github.com/smartystreets/goconvey/convey"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestStorage(t *testing.T) {

	Convey("File system based storage", t, func(c C) {

		baseDir, _ := ioutil.TempDir("", "test")
		defer os.Remove(baseDir)

		os.Setenv("TTS_BASE_DIR", baseDir)
		defer os.Unsetenv("TTS_BASE_DIR")

		storage := newFileSystemStorage()

		Convey("should save media", func() {

			id, err := storage.Save(strings.NewReader("test"))

			So(err, ShouldBeNil)
			So(id, ShouldNotBeEmpty)
		})

		Convey("should get media", func() {

			data := "This is just a simple test"

			id, _ := storage.Save(strings.NewReader(data))
			reader, err := storage.Get(id)

			content, _ := ioutil.ReadAll(reader)

			So(err, ShouldBeNil)
			So(string(content), ShouldEqual, data)
		})

		Convey("should return an error if media does not exist", func() {

			_, err := storage.Get("notExistingID")

			So(err, ShouldNotBeNil)
		})

		Convey("should remove media", func() {

			id, _ := storage.Save(strings.NewReader("test"))
			err := storage.Delete(id)

			So(err, ShouldBeNil)
		})
	})
}
