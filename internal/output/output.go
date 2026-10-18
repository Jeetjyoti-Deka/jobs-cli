package output

import (
    "github.com/Jeetjyoti-Deka/jobs-cli/internal/models"
)

type Output interface {
    Write([]models.ScrapeResult) error
}