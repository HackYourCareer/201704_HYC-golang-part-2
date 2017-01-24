package web

import (
	"go-university/tts-service/service"
)

type CreateDTO struct {
	Text     string
	Language string
}

type ResultDTO struct {
	ID       string `json:"id"`
	Text     string `json:"text"`
	Language string `json:"language"`
	Status   string `json:"status"`
	MediaUrl string `json:"mediaUrl,omitempty"`
}

//Converts service result to REST response object
func (r *ResultDTO) createWith(s *service.TtsResult, mediaUrl mediaUrlFunc) {
	r.ID = s.Id
	r.Text = s.Text
	r.Language = s.Language.String()
	r.Status = s.Status.String()

	if s.MediaId != "" {
		r.MediaUrl = mediaUrl(s.MediaId)
	} //QUESTION: Why no else here?

}

// REST error object
type ErrorDTO struct {
	Status  int      `json:"status"`
	Message string   `json:"message"`
	Details []string `json:"details,omitempty"`
}

//ErrorDTO implements built-in  "error" interface
func (err ErrorDTO) Error() string {
	return err.Message
}
