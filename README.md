# cronlog

Structured logging wrapper for cron jobs with failure notifications via webhook.

## Installation

```bash
go get github.com/yourusername/cronlog
```

## Usage

```go
package main

import (
    "github.com/yourusername/cronlog"
)

func main() {
    logger := cronlog.New(cronlog.Config{
        JobName:    "daily-report",
        WebhookURL: "https://hooks.slack.com/services/your/webhook/url",
        LogLevel:   cronlog.INFO,
    })

    if err := runJob(logger); err != nil {
        logger.Failure("job failed", "error", err)
        return
    }

    logger.Success("job completed successfully")
}

func runJob(logger *cronlog.Logger) error {
    logger.Info("starting data export", "records", 1500)
    // your cron job logic here
    return nil
}
```

### Output

```json
{"time":"2024-01-15T08:00:00Z","level":"INFO","job":"daily-report","msg":"starting data export","records":1500}
{"time":"2024-01-15T08:00:01Z","level":"INFO","job":"daily-report","msg":"job completed successfully","status":"success","duration_ms":1243}
```

On failure, cronlog automatically sends a structured payload to the configured webhook endpoint with the job name, error details, and execution duration.

## Configuration

| Field        | Type   | Description                        |
|--------------|--------|------------------------------------|
| `JobName`    | string | Identifier for the cron job        |
| `WebhookURL` | string | Endpoint to POST failure alerts to |
| `LogLevel`   | Level  | Minimum log level (default: INFO)  |

## License

MIT © [yourusername](https://github.com/yourusername)