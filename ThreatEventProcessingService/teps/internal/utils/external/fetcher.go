package external

import (
	"ThreatEventProcessingService/teps/internal/models"
	"errors"
	"fmt"
	"github.com/goccy/go-json"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func FetchEvents() ([]models.Event, error) {
	eventFeedApi := os.Getenv("FETCH_FEED_API")
	if eventFeedApi == "" {
		log.Println("FETCH_FEED_API not set")
		return nil, errors.New("FETCH_FEED_API not set")
	}
	resp, err := http.Get(eventFeedApi)
	if err != nil {
		log.Println("Failed to fetch events:", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 response: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var events []models.Event
	err = json.Unmarshal(body, &events)
	if err != nil {
		return nil, err
	}
	for i := range events {
		if events[i].ProcessedAt.IsZero() {
			events[i].ProcessedAt = time.Now()
		}
	}
	return events, nil

}
