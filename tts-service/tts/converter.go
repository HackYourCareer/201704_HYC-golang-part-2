package tts

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
)

type converter interface {

	// Convert converts the text to a media.
	// It returns an io.ReadCloser of an error, if any.
	Convert(text string, metadata Metadata) (io.ReadCloser, error)
}

// VoiceRss based implementation of the converter interface
type voiceRssConverter struct {
	apiUrl string
	apiKey string
}

//  VoiceRSS http://www.voicerss.org/api/documentation.aspx
//
//  The API has the following URL: https://api.voicerss.org/
//
//  Parameters:
//  key - The API key (mandatory)
//  src - The textual content for converting to speech (length limited by 100KB) (mandatory)
//  hl  - The textual content language. Allows values: see Languages (mandatory)
func (c voiceRssConverter) Convert(text string, meta Metadata) (io.ReadCloser, error) {

	return ioutil.NopCloser(new(bytes.Buffer)), nil
}

func newVoiceRssConverter() *voiceRssConverter {

	return &voiceRssConverter{
		apiUrl: "https://api.voicerss.org/",
		apiKey: os.Getenv("VOICE_RSS_API_KEY"),
	}
}

func (c voiceRssConverter) resolveLang(lang string) string {

	if lang == "PL" {
		return "pl-pl"
	}

	return "en-us"
}
