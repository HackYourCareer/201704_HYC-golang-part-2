# Go in practice

## Text to speech microservice

### Prerequisites

1. Install Go (https://golang.org/doc/install)
2. Install Git (https://git-scm.com/book/en/v2/Getting-Started-Installing-Git)
3. Clone repository into your Go workspace directory

`mkdir SAPHybrisGliwice && cd $_ && git clone https://github.com/SAPHybrisGliwice/golang-part-2.git`

### How to run

1. Setup following environment variables

Environment variable | Description | Mandatory 
--- | --- | --- 
PUBLIC_DIR | UI main file location | true 
VOICE_RSS_API_KEY | API key for VoiceRSS API | true 
SERVICE_SELF_URL | Service URL used to produce media URLs. If not provided, localhost will be used | false 
TTS_BASE_DIR | Location for storing media. If not provided, temporary directory will be used | false 
PERSISTENCE_BASE_DIR | Location for storing text metadata. If not provided, temporary directory will be used | false

2. Run `go run app.go`

3. If you want to use UI, enter the following URL: `http://localhost:8080/public/index.html`
