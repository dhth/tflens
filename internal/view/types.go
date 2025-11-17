package view

import (
	"html/template"
	"time"
)

type HTMLConfig struct {
	CustomTemplate *string
	Title          string
}

type HTMLData struct {
	Title     string
	Columns   []string
	Rows      []HTMLRow
	Diffs     []HTMLDiff
	Timestamp string
}

type HTMLRow struct {
	Data   []string
	Status string
}

type HTMLDiff struct {
	ModuleName string
	Output     template.HTML
	BaseLabel  string
	HeadLabel  string
	BaseRef    string
	HeadRef    string
}

func NewHTMLData(title string, referenceTime time.Time) HTMLData {
	return HTMLData{
		Title:     title,
		Timestamp: referenceTime.UTC().Format("2006-01-02 15:04:05 UTC"),
	}
}
