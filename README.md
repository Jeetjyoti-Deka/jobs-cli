# jobs-cli

A small Go command-line tool that pulls current job openings from company careers boards and prints them to the terminal or writes them to a Markdown file.

You list the companies you care about in a YAML config, along with the applicant-tracking system (ATS) each one uses. `jobs-cli` then:

1. Queries every company's public job-board API **concurrently**.
2. Optionally filters the postings by **title** and/or **location**.
3. Writes the results to the **console** or to a **Markdown** file.

A failure for one company (bad board name, network error, unsupported provider) is reported for that company and does not stop the others.

### Supported providers

| Provider     | `provider` value | API used                                                  | Status              |
|--------------|------------------|-----------------------------------------------------------|---------------------|
| Greenhouse   | `greenhouse`     | `https://boards-api.greenhouse.io/v1/boards/<board>/jobs` | ✅ Working          |
| Lever        | `lever`          | `https://api.lever.co/v0/postings/<board>?mode=json`      | ✅ Working          |
| Ashby        | `ashby`          | `https://api.ashbyhq.com/posting-api/job-board/<board>`   | ✅ Working          |

Any other `provider` value returns `unknown provider "<value>"`.

---

## Requirements

- **Go 1.25+** (see `go.mod`)
- Internet access to reach the Greenhouse and Lever APIs

## Installation

```bash
# Clone and enter the repo
git clone https://github.com/Jeetjyoti-Deka/jobs-cli.git
cd jobs-cli

# Download dependencies
go mod tidy

# Create your config from the example (config.yaml is git-ignored),
# then replace the placeholder boards with real companies
cp config.example.yaml config.yaml

# Option A: run directly
go run . scrape

# Option B: build a binary
go build -o jobs-cli .
./jobs-cli scrape

# Option C: install into $GOPATH/bin
go install .
```

---

## Configuration and files

### `config.yaml` (required)

The tool reads its list of companies from a YAML file. By default it looks for `config.yaml` in the **current working directory**. Use `--config` / `-c` to point somewhere else.

```yaml
companies:
  - name: Crunchyroll        # Display name used in the output
    provider: greenhouse     # greenhouse | lever | ashby
    board: crunchyroll       # The company's board identifier on that provider

  - name: Spotify
    provider: lever
    board: spotify

  # Lines starting with # are ignored, so you can comment companies in and out
  # - name: OpenAI
  #   provider: greenhouse
  #   board: openai
```

| Field      | Description |
|------------|-------------|
| `name`     | Any label you like. Shown as the company heading in the output. |
| `provider` | The ATS the company uses: `greenhouse`, `lever`, or `ashby`. |
| `board`    | The company's slug on that ATS. |

**Finding the `board` value:** look at the company's careers page URL.

- Greenhouse: `https://boards.greenhouse.io/<board>` or `https://job-boards.greenhouse.io/<board>`
- Lever: `https://jobs.lever.co/<board>`
- Ashby: `https://jobs.ashbyhq.com/<board>`

If the file is missing or isn't valid YAML, the command exits with an error.

### `jobs.md` (generated)

When you use `--output markdown`, results are written to `jobs.md` in the current directory by default, or to the path you give with `--file`. **The file is overwritten on every run.**

---

## Commands and flags

```
jobs-cli [command] [flags]
```

### Global flags

These work with every command.

| Flag       | Short | Default       | Description |
|------------|-------|---------------|-------------|
| `--config` | `-c`  | `config.yaml` | Path to the YAML configuration file. |
| `--help`   | `-h`  |               | Show help for the command. |

### `scrape`

Fetches job postings for every company in the config, applies any filters, and writes the output.

```
jobs-cli scrape [flags]
```

| Flag         | Short | Default   | Description |
|--------------|-------|-----------|-------------|
| `--output`   | `-o`  | `console` | Output format: `console` or `markdown`. |
| `--file`     | `-f`  | `jobs.md` | Output file path. Only used with `--output markdown`. |
| `--title`    | `-t`  | *(none)*  | Keep only jobs whose title contains this text. Repeatable / comma-separated. |
| `--location` | `-l`  | *(none)*  | Keep only jobs whose location contains this text. Repeatable / comma-separated. |

#### How filters work

- Matching is a **case-insensitive substring** match (`-t engineer` matches "Senior Software Engineer").
- Multiple values for the **same** flag are **OR**ed: `-t backend -t platform` keeps jobs whose title contains *backend* **or** *platform*.
- `--title` and `--location` together are **AND**ed: a job must match a title value **and** a location value.
- You can pass several values at once with commas (`-t backend,platform`) or by repeating the flag.
- Leading/trailing spaces around values are trimmed.
- Jobs with no location listed never match a `--location` filter.

### Other commands

| Command      | Description |
|--------------|-------------|
| `help`       | Help about any command (`jobs-cli help scrape`). |
| `completion` | Generate a shell autocompletion script: `jobs-cli completion [bash\|zsh\|fish\|powershell]`. |

---

## Examples

```bash
# Scrape everything in ./config.yaml and print to the terminal
jobs-cli scrape

# Use a different config file
jobs-cli scrape --config ~/configs/companies.yaml
jobs-cli -c ~/configs/companies.yaml scrape

# Only engineering roles
jobs-cli scrape -t engineer

# Backend OR platform roles that are remote OR in India
jobs-cli scrape -t backend,platform -l remote -l india

# Write results to the default jobs.md
jobs-cli scrape -o markdown

# Write filtered results to a custom Markdown file
jobs-cli scrape -o markdown -f remote-jobs.md -l remote

# Set up zsh completion for the current session
source <(jobs-cli completion zsh)
```

### Sample console output

```
Crunchyroll (3 jobs)
----------------------------------------
Senior Software Engineer
Location: Remote - US
https://job-boards.greenhouse.io/crunchyroll/jobs/1234567

...
----------------------------------------
✔ 1 companies succeeded
✖ 0 companies failed
```

### Sample Markdown output

```markdown
## Crunchyroll (3 jobs)

### Senior Software Engineer

**Location:** Remote - US

[Apply](https://job-boards.greenhouse.io/crunchyroll/jobs/1234567)
```

---

## Project structure

```
.
├── main.go                    # Entry point; calls cmd.Execute()
├── config.example.yaml        # Example company list; copy to config.yaml
├── cmd/
│   ├── root.go                # Root command and --config flag
│   └── scrape.go              # scrape command, its flags, output/filter wiring
└── internal/
    ├── config/                # YAML config loading
    ├── httpclient/            # Thin HTTP client with a GetJSON helper
    ├── providers/             # Greenhouse, Lever, Ashby and the provider factory
    ├── scraper/               # Runs every company's scrape concurrently
    ├── filter/                # Title and location filters + filter engine
    ├── models/                # Job and ScrapeResult types
    └── output/                # Console and Markdown writers
```

### Adding a new provider

1. Create `internal/providers/<name>.go` with a type that implements:
   ```go
   Scrape(company config.Company) ([]models.Job, error)
   ```
2. Add a constant for its name in `internal/providers/types.go`.
3. Add a `case` for it in `Factory.Get` in `internal/providers/factory.go`.

---

## Known limitations

- **Output order isn't fixed.** Companies are scraped in parallel, so they appear in whatever order their requests finish.
- **No request timeout.** The HTTP client uses Go's defaults, so a slow board can hold up the whole run.
