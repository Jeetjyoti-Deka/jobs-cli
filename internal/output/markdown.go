package output

import (
	"fmt"
	"io"

	"github.com/Jeetjyoti-Deka/jobs-cli/internal/models"
)

type MarkdownOutput struct {
	writer io.Writer
}

func NewMarkdownOutput(writer io.Writer) *MarkdownOutput {
	return &MarkdownOutput{
		writer: writer,
	}
}

func (o *MarkdownOutput) Write(results []models.ScrapeResult) error {
	for _, result := range results {
		if result.Err != nil {
			fmt.Fprintf(
				o.writer,
				"## %s\n\n",
				result.Company,
			)

			fmt.Fprintf(
				o.writer,
				"**Error:** %v\n\n",
				result.Err,
			)

			continue
		}

		fmt.Fprintf(
			o.writer,
			"## %s (%d jobs)\n\n",
			result.Company,
			len(result.Jobs),
		)

		for _, job := range result.Jobs {
			fmt.Fprintf(
				o.writer,
				"### %s\n\n",
				job.Title,
			)

			if job.Location != "" {
				fmt.Fprintf(
					o.writer,
					"**Location:** %s\n\n",
					job.Location,
				)
			}

			fmt.Fprintf(
				o.writer,
				"[Apply](%s)\n\n",
				job.URL,
			)
		}
	}

	return nil
}