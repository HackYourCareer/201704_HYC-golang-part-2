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

type Converter interface {

	//Convert converts the text to a media.
	//It returns an io.Reader of an error, if any.
	Convert(text string, metadata Metadata) (io.Reader, error)
}

type voiceRssConverter struct {
	apiKey string
	apiUrl string
}

func (c voiceRssConverter) Convert(text string, meta Metadata) (io.Reader, error) {

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

			r, w := io.Pipe()

			go func() {
				defer response.Body.Close()
				defer w.Close()

				io.Copy(w, response.Body)
			}()

			return r, nil
		}

		defer response.Body.Close()
		body, _ := ioutil.ReadAll(response.Body)

		return nil, fmt.Errorf("Unexpected response: %s", string(body))

	default:

		return nil, fmt.Errorf("Unexpected response: %d", response.StatusCode)
	}
}

func newVoiceRssConverter() *voiceRssConverter {

	apiKey := os.Getenv("VOICE_RRS_API_KEY")

	return &voiceRssConverter{
		apiUrl: "https://api.voicerss.org/",
		apiKey: apiKey,
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

const audioFormat = "16khz_16bit_stereo"
const speechRate = "-2"
