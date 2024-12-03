package model

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/m-mizutani/goerr"
)

type AlertID string

func NewAlertID() AlertID {
	id, err := uuid.NewV7()
	if err != nil {
		panic(err)
	}

	return AlertID(id.String())
}

const (
	AlertSchemaVersion = "v0"
)

type Alert struct {
	// ID is unique identifier of alert. It would be used for aggregation. If ID is empty, it will be generated automatically.
	ID AlertID `json:"id"`

	// Version is schema version of alert. It should be overwritten with AlertSchemaVersion.
	Version string `json:"version"`

	// Title is short description of alert. It's required.
	Title string `json:"title"`

	// Description is detailed description of alert.
	Description string `json:"description"`

	// Timestamp is time when alert is generated. If Timestamp is zero, it will be set to current time.
	Timestamp time.Time `json:"timestamp"`

	// Attrs is key-value pairs of additional information of alert.
	Attrs Attrs `json:"attrs"`
}

type Attrs map[string]any

func (x *Alert) Complete(ctx context.Context) {
	x.Version = AlertSchemaVersion

	if x.ID == "" {
		x.ID = NewAlertID()
	}
	if x.Timestamp.IsZero() {
		x.Timestamp = time.Now()
	}
}

func (x Alert) Validate() error {
	if x.Title == "" {
		return goerr.New("title is required")
	}

	return nil
}
