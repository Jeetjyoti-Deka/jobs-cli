package filter

import "github.com/Jeetjyoti-Deka/jobs-cli/internal/models"

type Filter interface {
	Match(job models.Job) bool
}

type Engine struct {
	filters []Filter
}

func NewEngine(filters ...Filter) *Engine {
	return &Engine{
		filters: filters,
	}
}

func (e *Engine) Apply(jobs []models.Job) []models.Job {
	if len(e.filters) == 0 {
		return jobs
	}

	filtered := make([]models.Job, 0, len(jobs))

	for _, job := range jobs {
		if e.matches(job) {
			filtered = append(filtered, job)
		}
	}

	return filtered
}

func (e *Engine) matches(job models.Job) bool {
	for _, filter := range e.filters {
		if !filter.Match(job) {
			return false
		}
	}

	return true
}
