package providers

import (
	"fmt"

	"github.com/Jeetjyoti-Deka/jobs-cli/internal/config"
	"github.com/Jeetjyoti-Deka/jobs-cli/internal/models"
	"github.com/Jeetjyoti-Deka/jobs-cli/internal/httpclient"
)

const greenhouseAPI = "https://boards-api.greenhouse.io/v1/boards/%s/jobs"


type greenhouseResponse struct {
	Jobs []greenhouseJob `json:"jobs"`
}

type greenhouseJob struct {
	Title       string `json:"title"`
	AbsoluteURL string `json:"absolute_url"`

	Location struct {
		Name string `json:"name"`
	} `json:"location"`
}

type GreenhouseProvider struct {
	client *httpclient.Client
}

func NewGreenhouseProvider(client *httpclient.Client) *GreenhouseProvider {
    return &GreenhouseProvider{
        client: client,
    }
}

func (p *GreenhouseProvider) Scrape(company config.Company) ([]models.Job, error) {

	url := fmt.Sprintf(
		greenhouseAPI,
		company.Board,
	)

	var gh greenhouseResponse

	if err := p.client.GetJSON(url, &gh); err != nil {
		return nil, err
	}

	jobs := make([]models.Job, 0, len(gh.Jobs))

	for _, job := range gh.Jobs {

		jobs = append(jobs, models.Job{
			Company: company.Name,
			Title: job.Title,
			Location: job.Location.Name,
			URL: job.AbsoluteURL,
		})
	}

	return jobs, nil
}