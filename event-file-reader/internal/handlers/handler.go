package handlers

import (
	"encoding/json"
	"event-file-reader/internal/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

func ReadEventsFromFile(c *gin.Context) {
	data, err := os.ReadFile("data/events.json")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not read file"})
		return
	}

	var events []models.Event
	if err := json.Unmarshal(data, &events); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid JSON format"})
		return
	}

	c.JSON(http.StatusOK, events)
}
