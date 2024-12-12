package model

type QueryInput map[QueryID]any

type QueryOutput struct {
	Alert []AlertBody `json:"alert"`
}
