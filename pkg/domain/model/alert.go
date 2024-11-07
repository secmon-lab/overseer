package model

import (
	"time"

	"github.com/m-mizutani/goerr"
)

type AlertID string

type Alert struct {
	ID          AlertID   `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Timestamp   time.Time `json:"timestamp"`
	Attrs       Attrs     `json:"attrs"`
}

type Attrs map[string]any

func (x Alert) Validate() error {
	if x.ID == "" {
		return goerr.New("id is required")
	}
	if x.Title == "" {
		return goerr.New("title is required")
	}

	if x.Timestamp.IsZero() {
		return goerr.New("timestamp is required")
	}

	return nil
}
