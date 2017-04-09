package tts

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type converter interface {

	// Convert converts the text to a media.
	// It returns an io.ReadCloser of an error, if any.
	Convert(text string, metadata Metadata) (io.ReadCloser, error)
}

// VoiceRss based implementation of the converter interface //
type voiceRssConverter struct {
	apiKey string
	apiUrl string
}

//  VoiceRSS mandatory parameters http://www.voicerss.org/api/documentation.aspx
//  key - The API key (mandatory)
//  src - The textual content for converting to speech (length limited by 100KB) (mandatory)
//  hl  - The textual content language. Allows values: see Languages (mandatory)
//  f   - The speech audio formats. Allows values: see Audio Formats. Default value: 8khz_8bit_mono. (optional)
//  r   - The speech rate (speed). Allows values: from -10 (slowest speed) up to 10 (fastest speed). Default value: 0 (normal speed). (optional)
func (c voiceRssConverter) Convert(text string, meta Metadata) (io.ReadCloser, error) {

	response, err := http.PostForm(c.apiUrl, url.Values{
		"key": {c.apiKey},
		"src": {text},
		"hl":  {c.resolveLang(meta.Lang)},
		"f":   {audioFormat},
		"r":   {speechRate},
	})

	if err != nil {
		return nil, err
	}

	switch response.StatusCode {

	case http.StatusOK:

		if strings.HasPrefix(response.Header.Get("Content-Type"), "audio") {
			return response.Body, nil
		}

		defer response.Body.Close()
		body, _ := ioutil.ReadAll(response.Body)

		return nil, fmt.Errorf("Unexpected response: %s", string(body))

	default:

		return nil, fmt.Errorf("Unexpected response: %d", response.StatusCode)
	}
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

const audioFormat = "16khz_16bit_stereo"
const speechRate = "-2"
