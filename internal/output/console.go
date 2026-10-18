package output

import (
    "io"
    "fmt"

    "github.com/Jeetjyoti-Deka/jobs-cli/internal/models"
)

type ConsoleOutput struct {
    writer io.Writer
}



func NewConsoleOutput(writer io.Writer) *ConsoleOutput {
    return &ConsoleOutput{
        writer: writer,
    }
}


func (o *ConsoleOutput) Write(results []models.ScrapeResult) error {
    failed := 0
    total := len(results)
    for _, result := range results {
		if result.Err != nil {
			failed++

			fmt.Fprintf(
				o.writer,
				"❌ %s: %v\n\n",
				result.Company,
				result.Err,
			)

			continue
		}

		fmt.Fprintf(
			o.writer,
			"%s (%d jobs)\n",
			result.Company,
			len(result.Jobs),
		)

		fmt.Fprintln(o.writer, "----------------------------------------")

		for _, job := range result.Jobs {
			fmt.Fprintf(o.writer, "%s\n", job.Title)

			if job.Location != "" {
				fmt.Fprintf(o.writer, "Location: %s\n", job.Location)
			}

			fmt.Fprintf(o.writer, "%s\n\n", job.URL)
		}
	}

	fmt.Fprintln(o.writer, "----------------------------------------")
	fmt.Fprintf(o.writer, "✔ %d companies succeeded \n", total-failed)
	fmt.Fprintf(o.writer, "✖ %d companies failed \n", failed)

	return nil
}