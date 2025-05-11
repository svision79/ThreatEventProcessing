package handlers

import (
	"ThreatEventProcessingService/teps/internal/models"
	"ThreatEventProcessingService/teps/internal/service"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type EventHandler struct {
	service service.EventService
}

func NewEventHandler(s service.EventService) *EventHandler {
	return &EventHandler{
		service: s,
	}
}

func (h *EventHandler) RegisterRoutes(r *gin.Engine) {
	e := r.Group("/api/v1/events")
	e.POST("", h.CreateEvent)
	e.PUT("", h.UpdateEvent)
	e.GET("/:id", h.GetEvent)
	e.DELETE("/:id", h.DeleteEvent)
}

func (h *EventHandler) CreateEvent(context *gin.Context) {
	var event models.Event

	err := context.ShouldBindBodyWithJSON(&event)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.service.Create(&event)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusCreated, gin.H{"event": event})
}

func (h *EventHandler) UpdateEvent(context *gin.Context) {
	var event models.Event
	err := context.ShouldBindBodyWithJSON(&event)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = h.service.UpdateById(event.ID, &event)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		context.JSON(http.StatusNotFound, gin.H{"error": "Record not found for given id"})
		return
	} else if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"event": event})
	return
}

func (h *EventHandler) GetEvent(context *gin.Context) {
	var event *models.Event
	eventId := context.Param("id")
	eventIdInt, er := strconv.Atoi(eventId)
	if er != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": er.Error()})
		return
	}
	event, err := h.service.GetById(eventIdInt)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": "Record not found for given id"})
		return
	}
	context.JSON(http.StatusOK, gin.H{"event": event})
	return
}

func (h *EventHandler) DeleteEvent(context *gin.Context) {
	eventId := context.Param("id")
	eventIdInt, er := strconv.Atoi(eventId)
	if er != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": er.Error()})
		return
	}
	err := h.service.DeleteById(eventIdInt)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"event": eventIdInt})
	return
}
