package tts

import (
	"bytes"
	"io"
	"io/ioutil"
)

type converter interface {

	// Convert converts the text to a media.
	// It returns an io.ReadCloser of an error, if any.
	Convert(text string, metadata Metadata) (io.ReadCloser, error)
}

// VoiceRss based implementation of the converter interface
type voiceRssConverter struct {
	apiUrl string
}

//  VoiceRSS http://www.voicerss.org/api/documentation.aspx
//
//  The API has the following URL: https://api.voicerss.org/
//
//  Parameters:
//  key - The API key (mandatory)
//  src - The textual content for converting to speech (length limited by 100KB) (mandatory)
//  hl  - The textual content language. Allows values: see Languages (mandatory)
//  f   - The speech audio formats. Allows values: see Audio Formats. Default value: 8khz_8bit_mono. (optional)
//  r   - The speech rate (speed). Allows values: from -10 (slowest speed) up to 10 (fastest speed). Default value: 0 (normal speed). (optional)
func (c voiceRssConverter) Convert(text string, meta Metadata) (io.ReadCloser, error) {

	return ioutil.NopCloser(new(bytes.Buffer)), nil
}

func newVoiceRssConverter() *voiceRssConverter {

	return &voiceRssConverter{
		apiUrl: "https://api.voicerss.org/",
	}
}

func (c voiceRssConverter) resolveLang(lang string) string {

	if lang == "PL" {
		return "pl-pl"
	}

	return "en-us"
}
