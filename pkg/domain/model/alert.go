package model

type Alert struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Results     []Result `json:"results"`
}

type Result map[string]any
