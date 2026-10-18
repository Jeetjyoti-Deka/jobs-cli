package providers

import (
	"github.com/Jeetjyoti-Deka/jobs-cli/internal/config"
	"github.com/Jeetjyoti-Deka/jobs-cli/internal/models"
)

type Provider interface {
	Scrape(company config.Company) ([]models.Job, error)
}
