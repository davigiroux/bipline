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
	b, _ := json.Marshal(f.data)
	return json.Unmarshal(b, resp.Data)
}

func TestFindChannel(t *testing.T) {
	fixture := map[string]interface{}{
		"channels": []map[string]interface{}{
			{"id": "6a073cec090476fb99230076", "name": "devgiroux", "service": "twitter"},
			{"id": "6a077599090476fb99243bdc", "name": "davi-alvarenga-028614119", "service": "linkedin"},
		},
	}

	tests := []struct {
		service string
		wantID  string
		wantErr bool
	}{
		{service: "twitter", wantID: "6a073cec090476fb99230076"},
		{service: "linkedin", wantID: "6a077599090476fb99243bdc"},
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
