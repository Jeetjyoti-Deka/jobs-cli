package scraper

import (
    "sync"

    "github.com/Jeetjyoti-Deka/jobs-cli/internal/providers"
    "github.com/Jeetjyoti-Deka/jobs-cli/internal/models"
    "github.com/Jeetjyoti-Deka/jobs-cli/internal/config"
)

type Scraper struct {
    factory *providers.Factory
}

func New(factory *providers.Factory) *Scraper {
    return &Scraper{
        factory: factory,
    }
}

func (s *Scraper) Run(companies []config.Company) []models.ScrapeResult {
   
		results := make(chan models.ScrapeResult)
		var wg sync.WaitGroup

		for _, company := range companies {
			wg.Add(1)

			go func(company config.Company) {
				defer wg.Done()

				provider, err := s.factory.Get(company.Provider)
				if err != nil {
					results <- models.ScrapeResult{
						Company: company.Name,
						Err: err,
					}
					return
				}

				jobs, err := provider.Scrape(company)

				results <- models.ScrapeResult{
					Company: company.Name,
					Jobs: jobs,
					Err: err,
				}

			}(company)
		}

		go func() {
			wg.Wait()
			close(results)
		}()

        scrapeResults := []models.ScrapeResult{}

        for result := range results {
            scrapeResults = append(scrapeResults, result)
        }

        return scrapeResults
}
