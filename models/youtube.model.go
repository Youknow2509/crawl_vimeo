package models

// model chapter
type Chapter struct {
	Time    string `json:"time"`    // Time in the format "00:00:00"
	Content string `json:"content"` // Content of the chapter
}
