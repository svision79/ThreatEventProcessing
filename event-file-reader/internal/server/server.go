package server

import (
	"event-file-reader/internal/handlers"
	"github.com/gin-gonic/gin"
)

func Run() error {
	r := gin.Default()
	r.GET("/events", handlers.ReadEventsFromFile) // single API to fetch events from file
	return r.Run(":9090")                         // listens on port 9090
}
