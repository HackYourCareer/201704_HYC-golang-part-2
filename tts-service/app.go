package main

import (
	"log"
	"net/http"

	"github.com/SAPHybrisGliwice/golang-part-2/tts-service/service"
	"github.com/SAPHybrisGliwice/golang-part-2/tts-service/tts"
	"github.com/SAPHybrisGliwice/golang-part-2/tts-service/web"
	"os"
	"strconv"
)

func main() {

	const port = 8080

	portStr := strconv.Itoa(port)

	engine := tts.NewEngine()
	persistence := service.NewPersistence()

	controller := service.New(persistence, engine)

	web.New(http.DefaultServeMux, controller, engine, selfUrl(portStr))

	log.Printf("Listening on port: %v", portStr)
	log.Fatal(http.ListenAndServe(":"+portStr, nil))
}

func selfUrl(port string) string {
	selfUrl := os.Getenv("SERVICE_SELF_URL")

	if len(selfUrl) == 0 {
		selfUrl = "http://localhost:" + port
	}
	return selfUrl
}
