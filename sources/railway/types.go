package railway

import (
	"github.com/ferretcode/pricetag/sink"
	"github.com/hasura/go-graphql-client"
)

type GraphQLConfig struct {
	AuthToken           string
	BaseSubscriptionURL string
	BaseURL             string
	client              *graphql.Client
	sink                *sink.Sink
}
