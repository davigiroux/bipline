//go:build integration

package generator

import (
	"context"
	"strings"
	"testing"

	"github.com/davigiroux/bipline/internal/eventsource"
)

// Replace these fixtures with your own past releases before running.
// These are example events — the content below is illustrative, not real.
var fixtureRelease = eventsource.Event{
	Type:  "release",
	Repo:  "owner/your-repo",
	URL:   "https://github.com/owner/your-repo/releases/tag/v0.1.0",
	Title: "v0.1.0 - Initial release",
	Body:  "First public release. Adds core functionality for X. Fixes Y edge case.",
	Tag:   "v0.1.0",
}

var fixtureRelease2 = eventsource.Event{
	Type:  "release",
	Repo:  "owner/another-repo",
	URL:   "https://github.com/owner/another-repo/releases/tag/v0.2.0",
	Title: "v0.2.0 - New feature",
	Body:  "Adds support for Z. Improves performance of the main loop by 30%.",
	Tag:   "v0.2.0",
}

func assertVoiceInvariants(t *testing.T, draft string, event eventsource.Event) {
	t.Helper()
	if strings.Contains(draft, "—") {
		t.Errorf("draft contains em dash\ndraft: %q", draft)
	}
	if !strings.Contains(draft, event.URL) {
		t.Errorf("draft missing event URL\nwant URL: %q\ndraft: %q", event.URL, draft)
	}
	if len(draft) > 500 {
		t.Errorf("draft too long: %d chars (max 500)\ndraft: %q", len(draft), draft)
	}
}

func TestDraft_Release(t *testing.T) {
	draft, err := Draft(context.Background(), fixtureRelease, "../../voice.md")
	if err != nil {
		t.Fatalf("Draft: %v", err)
	}
	t.Logf("draft:\n%s", draft)
	assertVoiceInvariants(t, draft, fixtureRelease)
}

func TestDraft_Release2(t *testing.T) {
	draft, err := Draft(context.Background(), fixtureRelease2, "../../voice.md")
	if err != nil {
		t.Fatalf("Draft: %v", err)
	}
	t.Logf("draft:\n%s", draft)
	assertVoiceInvariants(t, draft, fixtureRelease2)
}
