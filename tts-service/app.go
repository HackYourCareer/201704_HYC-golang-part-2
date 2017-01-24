package main

import (
	"fmt"
	"log"
	"net/http"

	"go-university/tts-service/service"
	"go-university/tts-service/web"
	"go-university/tts-service/tts"
	"strconv"
	"os"
)

func main() {

	const port = 8080

	portStr := strconv.Itoa(port)

	engine := tts.NewEngine()
	persistence := service.NewPersistence()

	controller := service.New(persistence, engine)

	web.New(http.DefaultServeMux, controller, engine, selfUrl(portStr))

	fmt.Println("Listening on port: " + portStr)
	log.Fatal(http.ListenAndServe(":" + portStr, nil))
}

func selfUrl(port string) string {
	selfUrl := os.Getenv("SERVICE_SELF_URL")

	if len(selfUrl) == 0 {
		selfUrl = "http://localhost:" + port
	}
	return selfUrl
}
