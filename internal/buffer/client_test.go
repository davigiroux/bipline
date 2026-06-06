package buffer

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/Khan/genqlient/graphql"
)

// fakeGQL implements graphql.Client by JSON-marshaling a fixed response into the
// resp.Data pointer that genqlient sets before calling MakeRequest.
type fakeGQL struct {
	data interface{}
	err  error
}

func (f *fakeGQL) MakeRequest(_ context.Context, _ *graphql.Request, resp *graphql.Response) error {
	if f.err != nil {
		return f.err
	}
	b, err := json.Marshal(f.data)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, resp.Data)
}

func TestCreateIdea_Success_Idea(t *testing.T) {
	fixture := map[string]interface{}{
		"createIdea": map[string]interface{}{
			"__typename": "Idea",
			"id":         "idea-123",
		},
	}
	c := &BufferClient{orgID: "test-org", gql: &fakeGQL{data: fixture}}
	if err := c.CreateIdea(context.Background(), "test title", "test text"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCreateIdea_Success_IdeaResponse(t *testing.T) {
	fixture := map[string]interface{}{
		"createIdea": map[string]interface{}{
			"__typename": "IdeaResponse",
			"idea":       map[string]interface{}{"id": "idea-456"},
		},
	}
	c := &BufferClient{orgID: "test-org", gql: &fakeGQL{data: fixture}}
	if err := c.CreateIdea(context.Background(), "test title", "test text"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCreateIdea_InvalidInputError(t *testing.T) {
	fixture := map[string]interface{}{
		"createIdea": map[string]interface{}{
			"__typename": "InvalidInputError",
			"message":    "invalid organizationId",
		},
	}
	c := &BufferClient{orgID: "test-org", gql: &fakeGQL{data: fixture}}
	err := c.CreateIdea(context.Background(), "test title", "test text")
	if err == nil {
		t.Fatal("want error, got nil")
	}
}

func TestFindChannel(t *testing.T) {
	fixture := map[string]interface{}{
		"channels": []map[string]interface{}{
			{"id": "6a073cec090476fb99230076", "name": "devgiroux", "service": string(ServiceTwitter)},
			{"id": "6a077599090476fb99243bdc", "name": "davi-alvarenga-028614119", "service": string(ServiceLinkedin)},
		},
	}

	tests := []struct {
		service string
		wantID  string
		wantErr bool
	}{
		{service: string(ServiceTwitter), wantID: "6a073cec090476fb99230076"},
		{service: string(ServiceLinkedin), wantID: "6a077599090476fb99243bdc"},
		{service: "tiktok", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.service, func(t *testing.T) {
			c := &BufferClient{orgID: "test-org", gql: &fakeGQL{data: fixture}}
			ch, err := c.FindChannel(context.Background(), tt.service)
			if tt.wantErr {
				if err == nil {
					t.Fatal("want error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if ch.ID != tt.wantID {
				t.Errorf("got ID %q, want %q", ch.ID, tt.wantID)
			}
			if ch.Service != tt.service {
				t.Errorf("got Service %q, want %q", ch.Service, tt.service)
			}
		})
	}
}
