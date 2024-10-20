package model

import "time"

type AlertID string

type Alert struct {
	ID        AlertID   `json:"id"`
	Title     string    `json:"title"`
	Timestamp time.Time `json:"timestamp"`
	Attrs     Attrs     `json:"attrs"`
}

type Attrs map[string]any
