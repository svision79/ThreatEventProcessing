package server

import (
	"ThreatEventProcessingService/teps/internal/handlers"
	"ThreatEventProcessingService/teps/internal/repository"
	"ThreatEventProcessingService/teps/internal/schedulers"
	"ThreatEventProcessingService/teps/internal/service"
	"github.com/gin-gonic/gin"
	"log"
)

/*
Run function for :
 1. injecting dependencies
 2. Define routes
 3. Start schedulers
 4. start server

Called in main method
returns error
*/
func Run() error {
	r := gin.Default()
	db := repository.NewDB()
	cache := repository.NewCache()
	eventService := service.NewEventService(db, cache)
	log.Println("Loading data from event reader and saving")
	eventService.FetchEventsAndSave()()
	eventHandler := handlers.NewEventHandler(eventService)
	eventScheduler := schedulers.InitSchedulers(eventService)
	err1, err2 := eventScheduler.StartSchedulers()
	if err1 != nil || err2 != nil {
		log.Panic("Failed to start schedulers", err1, err2)
	}
	eventHandler.RegisterRoutes(r)

	return r.Run(":8080")
}
