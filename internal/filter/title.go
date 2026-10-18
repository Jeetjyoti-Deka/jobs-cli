package filter

import (
	"strings"

	"github.com/Jeetjyoti-Deka/jobs-cli/internal/models"
)

type TitleFilter struct {
	queries []string
}

func NewTitleFilter(queries ...string) *TitleFilter {
	normalized := make([]string, 0, len(queries))

	for _, query := range queries {
		query = strings.ToLower(strings.TrimSpace(query))

		if query != "" {
			normalized = append(normalized, query)
		}
	}

	return &TitleFilter{
		queries: normalized,
	}
}

func (f *TitleFilter) Match(job models.Job) bool {
	title := strings.ToLower(job.Title)

	for _, query := range f.queries {
		if strings.Contains(title, query) {
			return true
		}
	}

	return false
}