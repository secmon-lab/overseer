package model

type QueryInput map[QueryID]any

type QueryOutput struct {
	Alert []Alert `json:"alert"`
}
