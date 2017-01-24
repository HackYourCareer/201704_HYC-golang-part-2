package tts

import (
	"io"
        "bytes"
)

type Converter interface {

	//Convert converts the text to a media.
	//It returns an io.Reader of an error, if any.
	Convert(text string, metadata Metadata) (io.Reader, error)
}

type voiceRssConverter struct {
	apiUrl string
}

func (c voiceRssConverter) Convert(text string, meta Metadata) (io.Reader, error) {

        return new(bytes.Buffer), nil
}

func newVoiceRssConverter() *voiceRssConverter {

	return &voiceRssConverter{
		apiUrl: "https://api.voicerss.org/",
	}
}

func (c voiceRssConverter) resolveLang(lang string) string {

	switch lang {

	case "PL":

		return "pl-pl"
	default:

		return "en-us"
	}
}