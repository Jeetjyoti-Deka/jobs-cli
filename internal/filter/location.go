package filter

import (
	"strings"

	"github.com/Jeetjyoti-Deka/jobs-cli/internal/models"
)

type LocationFilter struct {
	queries []string
}

func NewLocationFilter(queries ...string) *LocationFilter {
	normalized := make([]string, 0, len(queries))

	for _, query := range queries {
		query = strings.ToLower(strings.TrimSpace(query))

		if query != "" {
			normalized = append(normalized, query)
		}
	}

	return &LocationFilter{
		queries: normalized,
	}
}

func (f *LocationFilter) Match(job models.Job) bool {
	location := strings.ToLower(job.Location)

	for _, query := range f.queries {
		if strings.Contains(location, query) {
			return true
		}
	}

	return false
}