package scraper

import (
	"sync"

	"github.com/Jeetjyoti-Deka/jobs-cli/internal/config"
	"github.com/Jeetjyoti-Deka/jobs-cli/internal/models"
	"github.com/Jeetjyoti-Deka/jobs-cli/internal/providers"
)

const DefaultConcurrency = 8

type Scraper struct {
	factory     *providers.Factory
	concurrency int
}

func New(factory *providers.Factory, concurrency int) *Scraper {
	return &Scraper{
		factory:     factory,
		concurrency: max(concurrency, 1),
	}
}

func (s *Scraper) Run(companies []config.Company) []models.ScrapeResult {
	queue := make(chan config.Company, len(companies))
	for _, company := range companies {
		queue <- company
	}
	close(queue)

	results := make(chan models.ScrapeResult)
	var wg sync.WaitGroup

	for range min(s.concurrency, len(companies)) {
		wg.Add(1)
		go s.worker(queue, results, &wg)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	scrapeResults := make([]models.ScrapeResult, 0, len(companies))

	for result := range results {
		scrapeResults = append(scrapeResults, result)
	}

	return scrapeResults
}

func (s *Scraper) worker(
	companies <-chan config.Company,
	results chan<- models.ScrapeResult,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	for company := range companies {
		provider, err := s.factory.Get(company.Provider)
		if err != nil {
			results <- models.ScrapeResult{
				Company: company.Name,
				Err:     err,
			}
			continue
		}

		jobs, err := provider.Scrape(company)

		results <- models.ScrapeResult{
			Company: company.Name,
			Jobs:    jobs,
			Err:     err,
		}
	}
}
