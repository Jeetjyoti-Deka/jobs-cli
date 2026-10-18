package cmd

import (
	"fmt"
	"os"

	"github.com/Jeetjyoti-Deka/jobs-cli/internal/config"
	"github.com/Jeetjyoti-Deka/jobs-cli/internal/providers"
	"github.com/Jeetjyoti-Deka/jobs-cli/internal/httpclient"
	"github.com/Jeetjyoti-Deka/jobs-cli/internal/scraper"
	"github.com/Jeetjyoti-Deka/jobs-cli/internal/output"
	"github.com/Jeetjyoti-Deka/jobs-cli/internal/filter"
	"github.com/spf13/cobra"
)

var (
	outputType string
	outputFile string
	titleFilter []string
	locationFilter []string
)

func createOutput() (output.Output, *os.File, error) {
	switch outputType {
	case "console":
		return output.NewConsoleOutput(os.Stdout), nil, nil

	case "markdown":
		if outputFile == "" {
			outputFile = "jobs.md"
		}

		file, err := os.Create(outputFile)
		if err != nil {
			return nil, nil, fmt.Errorf(
				"create output file: %w",
				err,
			)
		}

		return output.NewMarkdownOutput(file), file, nil

	default:
		return nil, nil, fmt.Errorf(
			"unsupported output type %q: expected console or markdown",
			outputType,
		)
	}
}

func createFilters() []filter.Filter {
	var filters []filter.Filter

	if len(titleFilter) > 0 {
		filters = append(
			filters,
			filter.NewTitleFilter(titleFilter...),
		)
	}

	if len(locationFilter) > 0 {
		filters = append(
			filters,
			filter.NewLocationFilter(locationFilter...),
		)
	}

	return filters
}

var scrapeCmd = &cobra.Command{
	Use:   "scrape",
	Short: "Scrape configured companies",
	RunE: func(cmd *cobra.Command, args []string) error {

		cfg, err := config.Load(configPath)
		if err != nil {
			return err
		}
		client := httpclient.New()

		factory := providers.NewFactory(client)

		sc := scraper.New(factory)

		results := sc.Run(cfg.Companies)

		// Create filters from CLI flags.
		filters := createFilters()

		// Apply filters to each successful scrape result.
		engine := filter.NewEngine(filters...)

		for i := range results {
			if results[i].Err != nil {
				continue
			}

			results[i].Jobs = engine.Apply(results[i].Jobs)
		}

		var out output.Output

		out, _, err = createOutput()
		if err != nil{
			return fmt.Errorf("error creating output formatter: %w", err)
		}

		if err := out.Write(results); err != nil {
			return fmt.Errorf("write output: %w", err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(scrapeCmd)

	scrapeCmd.Flags().StringVarP(
		&outputType,
		"output",
		"o",
		"console",
		"Output type: console or markdown",
	)

	scrapeCmd.Flags().StringVarP(
		&outputFile,
		"file",
		"f",
		"jobs.md",
		"Output file",
	)

	scrapeCmd.Flags().StringSliceVarP(
		&titleFilter,
		"title",
		"t",
		nil,
		"Filter jobs by title",
	)

	scrapeCmd.Flags().StringSliceVarP(
		&locationFilter,
		"location",
		"l",
		nil,
		"Filter jobs by location",
	)
}