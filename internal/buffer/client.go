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

// CreateIdea creates a draft idea in the org's Ideas board.
// The mutation returns a CreateIdeaPayload union; all branches are handled.
func (c *BufferClient) CreateIdea(ctx context.Context, title, text string) error {
	resp, err := CreateIdea(ctx, c.gql, c.orgID, title, text)
	if err != nil {
		return fmt.Errorf("CreateIdea: %w", err)
	}
	switch r := resp.CreateIdea.(type) {
	case *CreateIdeaCreateIdea:
		if r.Id == "" {
			return fmt.Errorf("buffer: createIdea returned Idea with empty id")
		}
		return nil
	case *CreateIdeaCreateIdeaIdeaResponse:
		if r.Idea.Id == "" {
			return fmt.Errorf("buffer: createIdea returned IdeaResponse with empty idea id")
		}
		return nil
	case *CreateIdeaCreateIdeaInvalidInputError,
		*CreateIdeaCreateIdeaUnauthorizedError,
		*CreateIdeaCreateIdeaUnexpectedError,
		*CreateIdeaCreateIdeaLimitReachedError:
		return fmt.Errorf("buffer: %s", r.(interface{ GetMessage() string }).GetMessage())
	default:
		return fmt.Errorf("buffer: unexpected createIdea response: %T", resp.CreateIdea)
	}
}

// FindChannel returns the channel matching service (e.g. "twitter", "linkedin").
// Returns error if no channel is found for that service.
func (c *BufferClient) FindChannel(ctx context.Context, service string) (*Channel, error) {
	resp, err := GetChannels(ctx, c.gql, c.orgID)
	if err != nil {
		return nil, fmt.Errorf("GetChannels: %w", err)
	}
	svc := Service(service)
	for _, ch := range resp.Channels {
		if ch.Service == svc {
			return &Channel{ID: ch.Id, Name: ch.Name, Service: service}, nil
		}
	}
	return nil, fmt.Errorf("no channel for service %q in org %s", service, c.orgID)
}
