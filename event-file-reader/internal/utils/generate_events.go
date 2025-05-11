package utils

import (
	"encoding/json"
	"event-file-reader/internal/models"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

var sources = []string{"Firewall", "IDS", "Endpoint", "SIEM", "WAF"}
var threats = []string{"Malware", "Phishing", "DDoS", "Ransomware", "Brute Force"}

func GenerateDummyEvents() error {
	rand.Seed(time.Now().UnixNano())

	var events []models.Event
	for i := 1; i <= 200; i++ {
		detected := time.Now().Add(-time.Duration(rand.Intn(1000)) * time.Minute)
		processed := detected.Add(time.Duration(rand.Intn(120)) * time.Minute)

		event := models.Event{
			ID:          int(i),
			Source:      sources[rand.Intn(len(sources))],
			ThreatType:  threats[rand.Intn(len(threats))],
			DetectedAt:  detected,
			ProcessedAt: processed,
			Details:     fmt.Sprintf("Suspicious activity detected on host %d", rand.Intn(1000)),
		}
		events = append(events, event)
	}

	if err := os.MkdirAll("data", os.ModePerm); err != nil {
		return fmt.Errorf("failed to create data directory: %w", err)
	}

	path := filepath.Join("data", "events.json")
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(events); err != nil {
		return fmt.Errorf("failed to write events to file: %w", err)
	}
	return nil
}
