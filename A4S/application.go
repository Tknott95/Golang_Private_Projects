package main

import (
	"log"

	router "github.com/tknott95/Private_Go_Projects/Concurrency4Go/Controllers"
)

func main() {
	router.InitServer()
	log.Println(`Chat server launched on port 8080 🚀`)

}
