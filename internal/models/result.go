package models

type ScrapeResult struct {
	Company string
	Jobs    []Job
	Err     error
}
