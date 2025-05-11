package main

import (
	"event-file-reader/internal/server"
	"event-file-reader/internal/utils"
	"log"
)

func main() {

	go func() {
		if err := utils.GenerateDummyEvents(); err != nil {
			log.Fatalf("failed to generate dummy events: %v", err)
		}
	}()

	if err := server.Run(); err != nil {
		log.Fatalf("could not start server: %v", err)
	}
}
