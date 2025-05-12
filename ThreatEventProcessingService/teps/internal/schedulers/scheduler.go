package schedulers

import (
	"ThreatEventProcessingService/teps/internal/service"
	"github.com/robfig/cron/v3"
	"log"
	"os"
	"time"
)

var c *cron.Cron

type Scheduler struct {
	service service.EventService
}

func InitSchedulers(eventService service.EventService) *Scheduler {
	return &Scheduler{
		service: eventService,
	}
}

func (scheduler *Scheduler) StartSchedulers() (error, error) {
	c = cron.New()

	err1 := scheduler.dailyCleanUpSchedulers()
	err2 := scheduler.fetchEventsScheduler()

	if err1 != nil || err2 != nil {
		return err1, err2
	}
	c.Start()
	log.Println("Started schedulers..")
	return nil, nil
}

func (scheduler *Scheduler) fetchEventsScheduler() error {
	spec := os.Getenv("FETCH_API_TIME")
	_, err := c.AddFunc(spec, scheduler.service.FetchEventsAndSave())
	if err != nil {
		log.Println("Failed to schedule FetchEvents job:", err)
	} else {
		log.Println("Events Fetch job Scheduled: ")
	}
	return err
}

func (scheduler *Scheduler) dailyCleanUpSchedulers() error {

	spec := "59 23 * * *"

	batchId, err := c.AddFunc(spec, func() {
		events, err := scheduler.service.GetOlderEvents(24 * time.Hour)
		if err != nil {
			log.Println("Error getting older events")
			return
		}
		err = scheduler.service.DeleteEvents(*events)
		if err != nil {
			log.Println("Error deleting old events")
			return
		}
		//Commenting for now
		//go async.UploadEventsToS3(events)
	})
	if err != nil {
		log.Println("Failed to schedule cleanup job:", err)
	} else {
		log.Println("Cleanup job Scheduled: ", batchId)
	}
	return err
}
