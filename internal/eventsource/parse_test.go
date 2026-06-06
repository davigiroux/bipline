package eventsource

import (
	"os"
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	releaseJSON, err := os.ReadFile("testdata/release.json")
	if err != nil {
		t.Fatalf("read release fixture: %v", err)
	}
	prJSON, err := os.ReadFile("testdata/pr.json")
	if err != nil {
		t.Fatalf("read pr fixture: %v", err)
	}

	tests := []struct {
		name      string
		eventName string
		payload   []byte
		want      Event
		wantErr   bool
		errSubstr string
	}{
		{
			name:      "release published",
			eventName: "release",
			payload:   releaseJSON,
			want: Event{
				Type:  "release",
				Repo:  "devgiroux/safenudge",
				URL:   "https://github.com/devgiroux/safenudge/releases/tag/v1.2.0",
				Title: "v1.2.0 - Add batching",
				Body:  "Groups notifications by sender. Cuts interruptions by ~60%.",
				Tag:   "v1.2.0",
			},
		},
		{
			name:      "merged pull request",
			eventName: "pull_request",
			payload:   prJSON,
			want: Event{
				Type:  "pr",
				Repo:  "devgiroux/safenudge",
				URL:   "https://github.com/devgiroux/safenudge/pull/42",
				Title: "Add quiet hours setting",
				Body:  "Lets users configure a window where notifications are suppressed.",
				Tag:   "",
			},
		},
		{
			name:      "unmerged pull request",
			eventName: "pull_request",
			payload:   []byte(`{"action":"closed","pull_request":{"merged":false,"title":"x","body":"y","html_url":"https://github.com/o/r/pull/1"},"repository":{"full_name":"o/r"}}`),
			wantErr:   true,
			errSubstr: "not merged",
		},
		{
			name:      "unsupported event",
			eventName: "push",
			payload:   []byte(`{}`),
			wantErr:   true,
			errSubstr: "unsupported event",
		},
		{
			name:      "invalid JSON for release",
			eventName: "release",
			payload:   []byte(`not json`),
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.eventName, tt.payload)
			if tt.wantErr {
				if err == nil {
					t.Fatal("want error, got nil")
				}
				if tt.errSubstr != "" && !strings.Contains(err.Error(), tt.errSubstr) {
					t.Errorf("error %q does not contain %q", err.Error(), tt.errSubstr)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("got  %+v\nwant %+v", got, tt.want)
			}
		})
	}
}
