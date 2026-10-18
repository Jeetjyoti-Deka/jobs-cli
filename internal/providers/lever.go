package providers

import (
    "fmt"

    "github.com/Jeetjyoti-Deka/jobs-cli/internal/config"
	"github.com/Jeetjyoti-Deka/jobs-cli/internal/models"
	"github.com/Jeetjyoti-Deka/jobs-cli/internal/httpclient"
)

const leverAPI = "https://api.lever.co/v0/postings/%s?mode=json"

type leverJob struct {
	Text string `json:"text"`

	HostedURL string `json:"hostedUrl"`

	Categories struct {
		Location string `json:"location"`
	} `json:"categories"`
}

type LeverProvider struct{
    client *httpclient.Client
}

func NewLeverProvider(client *httpclient.Client) *LeverProvider {
    return &LeverProvider{
        client: client,
    }
}

func (p *LeverProvider) Scrape(company config.Company) ([]models.Job, error) {
    url := fmt.Sprintf(
		leverAPI,
		company.Board,
	)

	var response []leverJob

	if err := p.client.GetJSON(url, &response); err != nil {
		return nil, err
	}

	jobs := make([]models.Job, 0, len(response))

	for _, posting := range response {

		jobs = append(jobs, models.Job{
			Company: company.Name,
			Title: posting.Text,
			Location: posting.Categories.Location,
			URL: posting.HostedURL,
		})
	}

	return jobs, nil
}