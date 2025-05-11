package main

import (
	"ThreatEventProcessingService/teps/internal/server"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	err = server.Run()
	if err != nil {
		log.Fatalf("Error starting server %v", err)
	}
}
