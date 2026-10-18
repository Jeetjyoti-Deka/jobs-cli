package providers

import (
	"github.com/Jeetjyoti-Deka/jobs-cli/internal/config"
	"github.com/Jeetjyoti-Deka/jobs-cli/internal/models"
)

type DummyProvider struct{}

func NewDummyProvider() *DummyProvider {
	return &DummyProvider{}
}

func (p *DummyProvider) Scrape(company config.Company) ([]models.Job, error) {
	return []models.Job{
		{
			Company:  company.Name,
			Title:    "Software Engineer",
			Location: "Remote",
			URL:      company.Board,
		},
	}, nil
}
