package main

import (
	"context"
	"fmt"
	"os"

	"github.com/davigiroux/bipline/internal/buffer"
)

const (
	orgID   = "6a073bf15d53897094b8a76b"
	service = "twitter"
)

func main() {
	apiKey := os.Getenv("BUFFER_API_KEY")
	if apiKey == "" {
		fmt.Fprintln(os.Stderr, "BUFFER_API_KEY must be set")
		os.Exit(1)
	}

	ctx := context.Background()
	client := buffer.New(apiKey, orgID)

	ch, err := client.FindChannel(ctx, service)
	if err != nil {
		fmt.Fprintf(os.Stderr, "FindChannel: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("channel: %s (%s) id=%s\n", ch.Name, ch.Service, ch.ID)

	if err := client.CreateIdea(ctx, "bipline probe", "testing the bipline buffer client — ignore this idea"); err != nil {
		fmt.Fprintf(os.Stderr, "CreateIdea: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("idea created — check Buffer Ideas board")
}
