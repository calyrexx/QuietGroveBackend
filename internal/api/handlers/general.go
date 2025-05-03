package handlers

import (
	"log"
	"net/http"
)

const (
	statusMessage = "status OK"
)

type General struct {
	version string
}

func NewGeneral(version string) (*General, error) {
	return &General{version: version}, nil
}

func (h *General) Version(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	_, err := w.Write([]byte(h.version))
	if err != nil {
		log.Printf("Error: %s", err)
	}
}

func (h *General) Health(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte(statusMessage))
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
}
