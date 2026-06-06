package eventsource

import (
	"encoding/json"
	"fmt"
)

type releasePayload struct {
	Release struct {
		TagName string `json:"tag_name"`
		Name    string `json:"name"`
		Body    string `json:"body"`
		HTMLURL string `json:"html_url"`
	} `json:"release"`
	Repository struct {
		FullName string `json:"full_name"`
	} `json:"repository"`
}

type pullRequestPayload struct {
	PullRequest struct {
		Merged  bool   `json:"merged"`
		Title   string `json:"title"`
		Body    string `json:"body"`
		HTMLURL string `json:"html_url"`
	} `json:"pull_request"`
	Repository struct {
		FullName string `json:"full_name"`
	} `json:"repository"`
}

// Parse parses a GitHub event payload into a normalized Event.
// eventName is the GitHub event name ("release" or "pull_request").
// payload is the raw JSON bytes from GITHUB_EVENT_PATH.
func Parse(eventName string, payload []byte) (Event, error) {
	switch eventName {
	case "release":
		return parseRelease(payload)
	case "pull_request":
		return parsePR(payload)
	default:
		return Event{}, fmt.Errorf("eventsource: unsupported event %q", eventName)
	}
}

func parseRelease(payload []byte) (Event, error) {
	var p releasePayload
	if err := json.Unmarshal(payload, &p); err != nil {
		return Event{}, fmt.Errorf("eventsource: parse release: %w", err)
	}
	return Event{
		Type:  "release",
		Repo:  p.Repository.FullName,
		URL:   p.Release.HTMLURL,
		Title: p.Release.Name,
		Body:  p.Release.Body,
		Tag:   p.Release.TagName,
	}, nil
}

func parsePR(payload []byte) (Event, error) {
	var p pullRequestPayload
	if err := json.Unmarshal(payload, &p); err != nil {
		return Event{}, fmt.Errorf("eventsource: parse pull_request: %w", err)
	}
	if !p.PullRequest.Merged {
		return Event{}, fmt.Errorf("eventsource: pull_request is not merged")
	}
	return Event{
		Type:  "pr",
		Repo:  p.Repository.FullName,
		URL:   p.PullRequest.HTMLURL,
		Title: p.PullRequest.Title,
		Body:  p.PullRequest.Body,
		Tag:   "",
	}, nil
}
