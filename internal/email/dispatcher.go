package email

import (
	"bytes"
	"encoding/json"
	"log"
	"os"
	"os/exec"
)

func findBinary() string {
	candidates := []string{
		"./email/target/release/cast-onion-email",
		"./email/target/debug/cast-onion-email",
	}
	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	return ""
}

type Dispatcher struct {
	enabled bool
}

func NewDispatcher() *Dispatcher {
	path := findBinary()
	if path == "" {
		log.Println("email: cast-onion-email binary not found, emails disabled — build it in the email/ folder")
		return &Dispatcher{enabled: false}
	}
	log.Printf("email: using binary at %s", path)
	return &Dispatcher{enabled: true}
}

func (d *Dispatcher) send(payload any) {
	if !d.enabled {
		return
	}
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("email: marshal error: %v", err)
		return
	}
	go func() {
		cmd := exec.Command(findBinary())
		cmd.Stdin = bytes.NewReader(data)
		out, err := cmd.CombinedOutput()
		if err != nil {
			log.Printf("email: send failed: %v — %s", err, string(out))
		} else {
			log.Printf("email: sent — %s", string(out))
		}
	}()
}

func (d *Dispatcher) ApplicationReceived(to, name, stationName, genre, applicationID string) {
	d.send(map[string]any{
		"type":           "application_received",
		"to":             to,
		"name":           name,
		"station_name":   stationName,
		"genre":          genre,
		"application_id": applicationID,
	})
}

func (d *Dispatcher) ApplicationApproved(to, name, stationName, stationKey, accessToken string) {
	d.send(map[string]any{
		"type":         "application_approved",
		"to":           to,
		"name":         name,
		"station_name": stationName,
		"station_key":  stationKey,
		"access_token": accessToken,
	})
}

func (d *Dispatcher) ApplicationDenied(to, name, stationName, reason string) {
	d.send(map[string]any{
		"type":         "application_denied",
		"to":           to,
		"name":         name,
		"station_name": stationName,
		"reason":       reason,
	})
}

func (d *Dispatcher) StationSuspended(to, name, stationName, reason string) {
	d.send(map[string]any{
		"type":         "station_suspended",
		"to":           to,
		"name":         name,
		"station_name": stationName,
		"reason":       reason,
	})
}

func (d *Dispatcher) StationRevoked(to, name, stationName, reason string) {
	d.send(map[string]any{
		"type":         "station_revoked",
		"to":           to,
		"name":         name,
		"station_name": stationName,
		"reason":       reason,
	})
}
