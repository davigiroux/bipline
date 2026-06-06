package generator

import (
	"context"
	"fmt"
	"os"

	anthropic "github.com/anthropics/anthropic-sdk-go"
	"github.com/davigiroux/bipline/internal/eventsource"
)

// Draft generates a social media post draft for the given event.
// voicePath is the path to the user's voice.md file, which becomes the Claude system prompt.
// ANTHROPIC_API_KEY must be set in the environment.
func Draft(ctx context.Context, event eventsource.Event, voicePath string) (string, error) {
	voice, err := os.ReadFile(voicePath)
	if err != nil {
		return "", fmt.Errorf("generator: read voice: %w", err)
	}

	client := anthropic.NewClient()

	userMessage := fmt.Sprintf(
		"Repo: %s\nTag: %s\nTitle: %s\nURL: %s\n\nRelease notes:\n%s",
		event.Repo, event.Tag, event.Title, event.URL, event.Body,
	)

	msg, err := client.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     anthropic.ModelClaudeSonnet4_6,
		MaxTokens: 500,
		System: []anthropic.TextBlockParam{
			{Text: string(voice)},
		},
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(userMessage)),
		},
	})
	if err != nil {
		return "", fmt.Errorf("generator: claude: %w", err)
	}
	if len(msg.Content) == 0 {
		return "", fmt.Errorf("generator: empty response from Claude")
	}
	textBlock, ok := msg.Content[0].AsAny().(anthropic.TextBlock)
	if !ok {
		return "", fmt.Errorf("generator: unexpected response block type: %T", msg.Content[0].AsAny())
	}
	return textBlock.Text, nil
}
