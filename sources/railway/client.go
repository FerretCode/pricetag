package railway

import (
	"errors"
	"net/http"

	"github.com/hasura/go-graphql-client"
)

type authedTransport struct {
	token   string
	wrapped http.RoundTripper
}

func (t *authedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", "Bearer "+t.token)
	req.Header.Set("Content-Type", "application/json")
	return t.wrapped.RoundTrip(req)
}

func NewClient(gqlConfig *GraphQLConfig) (*GraphQLConfig, error) {
	if gqlConfig == nil {
		return nil, errors.New("gql config must not be nil")
	}

	if gqlConfig.AuthToken == "" {
		return nil, errors.New("auth token cannot be empty")
	}

	httpClient := &http.Client{
		Transport: &authedTransport{
			token:   gqlConfig.AuthToken,
			wrapped: http.DefaultTransport,
		},
	}

	if gqlConfig.BaseURL != "" {
		gqlConfig.client = graphql.NewClient(gqlConfig.BaseURL, httpClient)
	}

	return gqlConfig, nil
}
