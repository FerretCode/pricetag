package railway

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/charmbracelet/log"
	"github.com/coder/websocket"
	"github.com/google/uuid"
)

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

func (gql *GraphQLConfig) buildMetadataMap(ctx context.Context, config *Config) (map[string]string, error) {
	services, err := gql.GetServices(ctx, config)
	if err != nil {
		return nil, err
	}

	environments, err := gql.GetEnvironments(ctx, config)
	if err != nil {
		return nil, err
	}

	for k, v := range environments {
		services[k] = v
	}

	return services, nil
}

func (gql *GraphQLConfig) createSubscription(ctx context.Context, config *Config) (*websocket.Conn, error) {
	if config == nil {
		return nil, errors.New("config must be present")
	}

	subscribePayload := &subscribePayload{
		Query: streamEnvironmentLogsQuery,
		Variables: &variables{
			EnvironmentId: config.EnvironmentId,
			BeforeDate:    time.Now().UTC().Add(-5 * time.Minute).Format(time.RFC3339Nano),
			BeforeLimit:   500,
		}}

	operationMessage := operationMessage{
		Id:      uuid.Must(uuid.NewUUID()).String(),
		Type:    "subscribe",
		Payload: *subscribePayload,
	}

	operationMessageBytes, err := json.Marshal(&operationMessage)
	if err != nil {
		return nil, err
	}

	opts := &websocket.DialOptions{
		HTTPHeader: http.Header{
			"Authorization": []string{"Bearer" + gql.AuthToken},
			"Content-Type":  []string{"application/json"},
		},
		Subprotocols: []string{"graphql-transport-ws"},
	}

	timeout, cancel := context.WithTimeout(ctx, (10 * time.Second))
	defer cancel()

	conn, _, err := websocket.Dial(timeout, gql.BaseSubscriptionURL, opts)
	if err != nil {
		return nil, err
	}

	conn.SetReadLimit(-1)

	if err := conn.Write(ctx, websocket.MessageText, connectionInit); err != nil {
		return nil, err
	}

	_, ackMessage, err := conn.Read(ctx)
	if err != nil {
		return nil, err
	}

	if !bytes.Equal(ackMessage, connectionAck) {
		return nil, errors.New("did not receive connection acknowledgement from server")
	}

	if err := conn.Write(ctx, websocket.MessageText, operationMessageBytes); err != nil {
		return nil, err
	}

	return conn, nil
}

func (gql *GraphQLConfig) SubscribeToLogs(ctx context.Context, config *Config) error {
	idToNameMap, err := gql.buildMetadataMap(ctx, config)
	if err != nil {
		return err
	}

	conn, err := gql.createSubscription(ctx, config)
	if err != nil {
		return err
	}
	defer conn.CloseNow()

	LogTime := time.Now().UTC()

	for {
		_, logPayload, err := safeConnRead(conn, ctx)
		if err != nil {
			log.Error("resubscribing to logs endpoint", "reason", err)

			safeConnCloseNow(conn)

			conn, err = gql.createSubscription(ctx, config)
			if err != nil {
				return err
			}

			continue
		}

		logs := &logPayloadResponse{}

		if err := json.Unmarshal(logPayload, &logs); err != nil {
			return err
		}

		if logs.Type != TypeNext {
			log.Error("resubscribing to logs endpoint", "reason", fmt.Sprintf("log type not next: %s", logs.Type))

			safeConnCloseNow(conn)

			conn, err = gql.createSubscription(ctx, config)
			if err != nil {
				return err
			}

			continue
		}

		filteredLogs := []railwayLog{}

		for i := range logs.Payload.Data.EnvironmentLogs {
			LogTime = logs.Payload.Data.EnvironmentLogs[i].Timestamp

			if logs.Payload.Data.EnvironmentLogs[i].Timestamp.Before(LogTime) || LogTime == logs.Payload.Data.EnvironmentLogs[i].Timestamp {
				log.Debug("skipping stale log message")
				continue
			}

			serviceName, ok := idToNameMap[logs.Payload.Data.EnvironmentLogs[i].Tags.ServiceID]
			if !ok {
				log.Warn("service name not found")
				serviceName = "undefined"
			}

			logs.Payload.Data.EnvironmentLogs[i].Tags.ServiceName = serviceName

			environmentName, ok := idToNameMap[logs.Payload.Data.EnvironmentLogs[i].Tags.EnvironmentID]
			if !ok {
				log.Warn("environment name not found")
				environmentName = "undefined"
			}

			logs.Payload.Data.EnvironmentLogs[i].Tags.EnvironmentName = environmentName

			projectName, ok := idToNameMap[logs.Payload.Data.EnvironmentLogs[i].Tags.ProjectID]
			if !ok {
				log.Warn("project name not found")
				projectName = "undefined"
			}

			logs.Payload.Data.EnvironmentLogs[i].Tags.ProjectName = projectName

			filteredLogs = append(filteredLogs, logs.Payload.Data.EnvironmentLogs[i])
		}

		if len(filteredLogs) == 0 {
			continue
		}

		// TODO: reconstruct log message
	}
}

func safeConnRead(conn *websocket.Conn, ctx context.Context) (mT websocket.MessageType, b []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("recovered from panic: %v", r)
		}
	}()

	return conn.Read(ctx)
}

func safeConnCloseNow(conn *websocket.Conn) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("recovered from panic: %v", r)
		}
	}()

	return conn.CloseNow()
}
