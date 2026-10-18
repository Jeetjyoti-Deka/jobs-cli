package providers

import (
	"fmt"
	"strings"

	"github.com/Jeetjyoti-Deka/jobs-cli/internal/config"
	"github.com/Jeetjyoti-Deka/jobs-cli/internal/models"
	"github.com/Jeetjyoti-Deka/jobs-cli/internal/httpclient"
)

const ashbyAPI = "https://api.ashbyhq.com/posting-api/job-board/%s"

type ashbyResponse struct {
	Jobs []ashbyJob `json:"jobs"`
}

type ashbyJob struct {
	Title    string `json:"title"`
	Location string `json:"location"`
	JobURL   string `json:"jobUrl"`
	IsListed bool   `json:"isListed"`
	IsRemote bool   `json:"isRemote"`

	SecondaryLocations []struct {
		Location string `json:"location"`
	} `json:"secondaryLocations"`
}

type AshbyProvider struct {
	client *httpclient.Client
}

func NewAshbyProvider(client *httpclient.Client) *AshbyProvider {
	return &AshbyProvider{
		client: client,
	}
}

func (p *AshbyProvider) Scrape(company config.Company) ([]models.Job, error) {

	url := fmt.Sprintf(
		ashbyAPI,
		company.Board,
	)

	var response ashbyResponse

	if err := p.client.GetJSON(url, &response); err != nil {
		return nil, err
	}

	jobs := make([]models.Job, 0, len(response.Jobs))

	for _, job := range response.Jobs {
		if !job.IsListed {
			continue
		}

		jobs = append(jobs, models.Job{
			Company: company.Name,
			Title: job.Title,
			Location: ashbyLocation(job),
			URL: job.JobURL,
		})
	}

	return jobs, nil
}

// ashbyLocation joins the primary and secondary locations so the location
// filter matches any of them, and marks remote jobs that don't say so.
func ashbyLocation(job ashbyJob) string {
	locations := make([]string, 0, 1+len(job.SecondaryLocations))

	if job.Location != "" {
		locations = append(locations, job.Location)
	}

	for _, secondary := range job.SecondaryLocations {
		if secondary.Location != "" {
			locations = append(locations, secondary.Location)
		}
	}

	location := strings.Join(locations, "; ")

	if job.IsRemote && !strings.Contains(strings.ToLower(location), "remote") {
		if location == "" {
			return "Remote"
		}

		location += " (Remote)"
	}

	return location
}
