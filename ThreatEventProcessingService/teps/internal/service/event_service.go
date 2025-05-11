package service

import (
	"ThreatEventProcessingService/teps/internal/models"
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"log"
	"time"
)

type EventService interface {
	Create(event *models.Event) error
	GetById(id int) (*models.Event, error)
	ListEvents() (*[]models.Event, error)
	DeleteById(id int) error
	UpdateById(id int, event *models.Event) error
	GetOlderEvents(duration time.Duration) (*[]models.Event, error)
	DeleteEvents(oldEvent []models.Event) error
}

type eventService struct {
	db    *gorm.DB
	cache *redis.Client
}

func NewEventService(db *gorm.DB, cache *redis.Client) EventService {
	return &eventService{
		db:    db,
		cache: cache,
	}
}

func (service *eventService) Create(event *models.Event) error {
	return service.db.Create(&event).Error
}

func (service *eventService) GetById(id int) (*models.Event, error) {
	ctx := context.Background()
	redisKey := buildEventsCacheKey(id)
	eventBytes, cacheErr := service.cache.Get(ctx, redisKey).Result()
	if cacheErr == nil {
		var event models.Event
		err := json.Unmarshal([]byte(eventBytes), &event)
		if err == nil {
			return &event, nil
		}
	} else {
		log.Println("Error getting event from cache ", redisKey, cacheErr)
	}
	var event models.Event
	err := service.db.First(&event, id).Error
	if err != nil {
		return nil, err
	}
	cacheData, _ := json.Marshal(event)
	err = service.cache.Set(ctx, redisKey, cacheData, 0).Err()
	if err != nil {
		log.Println("Error adding event to cache")
	}
	return &event, nil
}

func (service *eventService) ListEvents() (*[]models.Event, error) {
	var events []models.Event
	err := service.db.Find(&events).Error
	if err != nil {
		return nil, err
	}
	return &events, nil
}

func (service *eventService) DeleteById(id int) error {
	err := service.db.Delete(&models.Event{}, id).Error
	if err != nil {
		return err
	}

	ctx := context.Background()
	_ = service.cache.Del(ctx, buildEventsCacheKey(id)).Err()
	return nil
}

func (service *eventService) UpdateById(id int, event *models.Event) error {
	_, err := service.GetById(id)
	if err != nil {
		return err
	}
	err = service.db.Save(&event).Error
	if err != nil {
		return err
	}
	ctx := context.Background()
	redisKey := buildEventsCacheKey(id)
	_ = service.cache.Del(ctx, redisKey).Err()
	return nil
}

func (service *eventService) DeleteEvents(oldEvent []models.Event) error {
	ids := make([]int, len(oldEvent))
	for i := 0; i < len(oldEvent); i++ {
		ids[i] = oldEvent[i].ID
	}
	tx := service.db.Begin()

	if tx.Error != nil {
		return tx.Error
	}
	err := tx.Delete(&models.Event{}, ids).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	ctx := context.Background()
	for id := range ids {
		service.cache.Del(ctx, buildEventsCacheKey(id))
	}
	return tx.Commit().Error

}

func buildEventsCacheKey(id int) string {
	return fmt.Sprintf("event:%d", id)
}

func (service *eventService) GetOlderEvents(duration time.Duration) (*[]models.Event, error) {
	var events []models.Event
	threshold := time.Now().Add(-duration)
	err := service.db.Where("detected_at < ?", threshold).Find(&events).Error
	return &events, err
}
