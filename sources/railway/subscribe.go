package railway

type operationMessage struct {
	Id      string           `json:"id"`
	Type    string           `json:"type"`
	Payload subscribePayload `json:"payload"`
}

type subscribePayload struct {
	Query     string     `json:"query"`
	Variables *variables `json:"variables"`
}

type variables struct {
	EnvironmentId string `json:"environmentId"`
	Filter        string `json:"filter"`
	BeforeLimit   int64  `json:"beforeLimit"`
	BeforeDate    string `json:"beforeDate"`
}

var (
	connectionInit = []byte(`{"type":"connection_init"}`)
	connectionAck  = []byte(`{"type":"connection_ack"}`)
)

/*
func (gql *GraphQLConfig) createSubscription(ctx context.Context) (*websocket.Conn, error) {
	subscribePayload := &subscribePayload{
		Query: streamEnvironmentLogsQuery,
		Variables: &variables{
			EnvironmentId: ,
		},
	}
}*/
