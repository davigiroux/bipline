package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/davigiroux/bipline/internal/buffer"
	"github.com/davigiroux/bipline/internal/eventsource"
	"github.com/davigiroux/bipline/internal/generator"
)

func main() {
	vals, err := requireEnv(
		"BUFFER_API_KEY",
		"ANTHROPIC_API_KEY",
		"BUFFER_ORG_ID",
		"BIPLINE_VOICE_PATH",
		"GITHUB_EVENT_NAME",
		"GITHUB_EVENT_PATH",
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	payload, err := os.ReadFile(vals["GITHUB_EVENT_PATH"])
	if err != nil {
		fmt.Fprintf(os.Stderr, "bipline: read event: %v\n", err)
		os.Exit(1)
	}

	event, err := eventsource.Parse(vals["GITHUB_EVENT_NAME"], payload)
	if err != nil {
		// Unsupported event type or unmerged PR — not a failure, skip silently.
		fmt.Fprintf(os.Stderr, "bipline: skip: %v\n", err)
		os.Exit(0)
	}

	ctx := context.Background()

	draft, err := generator.Draft(ctx, event, vals["BIPLINE_VOICE_PATH"])
	if err != nil {
		fmt.Fprintf(os.Stderr, "bipline: generate: %v\n", err)
		os.Exit(1)
	}

	client := buffer.New(vals["BUFFER_API_KEY"], vals["BUFFER_ORG_ID"])
	if err := client.CreateIdea(ctx, event.Title, draft); err != nil {
		fmt.Fprintf(os.Stderr, "bipline: create idea: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("bipline: created draft idea for %q\n", event.Title)
}

func requireEnv(names ...string) (map[string]string, error) {
	vals := make(map[string]string, len(names))
	var missing []string
	for _, name := range names {
		v := os.Getenv(name)
		if v == "" {
			missing = append(missing, name)
		} else {
			vals[name] = v
		}
	}
	if len(missing) > 0 {
		return nil, fmt.Errorf("bipline: missing required env vars: %s", strings.Join(missing, ", "))
	}
	return vals, nil
}
