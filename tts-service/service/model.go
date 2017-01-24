package service

//////////////////////////////////////// STRUCTS ////////////////////////////////////////

//Used to create new TTS data
type TtsCreate struct {
	Text     string
	Language LangEnum
}

//Defines Service result
//ID is the object unique identifier, derived from Text
//Text is the TTS source text
//MediaId is returned only if Status == Ready, and it's used to retrieve the data from Media Storage (outside of this Service)
type TtsResult struct {
	Id       string
	Text     string
	Language LangEnum
	Status   StatusEnum
	MediaId  string
}

//////////////////////////////////////// ENUMS ////////////////////////////////////////

////// LANGUAGE ENUM //////

//Exported interface provides type-safety
type LangEnum interface {
	//You can't implement this interface outside the package - return type is not exported
	Language() lang
	String() string
}

//Notice non-exported type
type lang string

// every lang must fullfill the LangEnum interface
func (l lang) Language() lang {
	return l
}
func (l lang) String() string {
	return string(l)
}

//Exported
const (
	EN lang = "EN"
	PL lang = "PL"
)

////// STATUS ENUM //////

//Exported interface provides type-safety
type StatusEnum interface {
	Status() status
	String() string
}

//Notice non-exported type
type status string

// every status must fullfill the StatusEnum interface
func (s status) Status() status {
	return s
}
func (s status) String() string {
	return string(s)
}

//Exported
const (
	StatusPending status = "PENDING"
	StatusReady   status = "READY"
	//"error" would collide with built-in error type
	StatusError status = "ERROR"
)
