package buffer

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Khan/genqlient/graphql"
)

const endpoint = "https://api.buffer.com"

// Channel is a connected social media account.
type Channel struct {
	ID      string
	Name    string
	Service string
}

type authTransport struct {
	apiKey string
	base   http.RoundTripper
}

func (t *authTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	r = r.Clone(r.Context())
	r.Header.Set("Authorization", "Bearer "+t.apiKey)
	return t.base.RoundTrip(r)
}

// BufferClient is a thin wrapper over the genqlient-generated Buffer GraphQL client.
type BufferClient struct {
	orgID string
	gql   graphql.Client
}

// New creates a BufferClient authenticated with apiKey, scoped to orgID.
func New(apiKey, orgID string) *BufferClient {
	httpClient := &http.Client{
		Transport: &authTransport{apiKey: apiKey, base: http.DefaultTransport},
	}
	return &BufferClient{
		orgID: orgID,
		gql:   graphql.NewClient(endpoint, httpClient),
	}
}

// FindChannel returns the channel matching service (e.g. "twitter", "linkedin").
// Returns error if no channel is found for that service.
func (c *BufferClient) FindChannel(ctx context.Context, service string) (*Channel, error) {
	resp, err := GetChannels(ctx, c.gql, c.orgID)
	if err != nil {
		return nil, fmt.Errorf("GetChannels: %w", err)
	}
	for _, ch := range resp.Channels {
		if string(ch.Service) == service {
			return &Channel{ID: ch.Id, Name: ch.Name, Service: string(ch.Service)}, nil
		}
	}
	return nil, fmt.Errorf("no channel for service %q in org %s", service, c.orgID)
}
