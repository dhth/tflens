package view

import "time"

type HTMLConfig struct {
	CustomTemplate *string
	Title          string
}

type HTMLData struct {
	Title     string
	Columns   []string
	Rows      []HTMLRow
	Timestamp string
}

type HTMLRow struct {
	Data   []string
	Status string
}

func NewHTMLData(title string, referenceTime time.Time) HTMLData {
	return HTMLData{
		Title:     title,
		Timestamp: referenceTime.UTC().Format("2006-01-02 15:04:05 UTC"),
	}
}
