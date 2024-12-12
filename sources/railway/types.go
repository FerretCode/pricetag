package railway

import (
	"time"

	"github.com/ferretcode/pricetag/sink"
	"github.com/hasura/go-graphql-client"
)

type LogType string

const (
	TypeNext     LogType = "next"
	TypeComplete LogType = "complete"
)

type GraphQLConfig struct {
	AuthToken           string
	BaseSubscriptionURL string
	BaseURL             string
	client              *graphql.Client
	sink                *sink.Sink
}

type environment struct {
	Environment struct {
		ProjectID string `json:"projectId"`
	} `json:"environment"`
}

type railwayLog struct {
	Timestamp time.Time `json:"timestamp"`
	Message   string    `json:"message"`
	Severity  string    `json:"severity"`
	Tags      struct {
		ProjectID   string `json:"projectId"`
		ProjectName string `json:"projectName"`

		EnvironmentID   string `json:"environmentId"`
		EnvironmentName string `json:"environmentName"`

		ServiceID   string `json:"serviceId"`
		ServiceName string `json:"serviceName"`

		DeploymentID         string `json:"deploymentId"`
		DeploymentInstanceID string `json:"deploymentInstanceId"`
	} `json:"tags"`
	Attributes []attributes `json:"attributes"`
}

type attributes struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type project struct {
	Project struct {
		ID           string `json:"id"`
		Name         string `json:"name"`
		Description  string `json:"description"`
		Environments struct {
			Edges []struct {
				Node struct {
					ID   string `json:"id"`
					Name string `json:"name"`
				} `json:"node"`
			} `json:"edges"`
		} `json:"environments"`
		Services struct {
			Edges []struct {
				Node struct {
					ID               string `json:"id"`
					Name             string `json:"name"`
					ServiceInstances struct {
						Edges []struct {
							Node struct {
								EnvironmentID string `json:"environmentId"`
							} `json:"node"`
						} `json:"edges"`
					} `json:"serviceInstances"`
				} `json:"node"`
			} `json:"edges"`
		} `json:"services"`
	} `json:"project"`
}

type logPayloadResponse struct {
	Payload struct {
		Data struct {
			EnvironmentLogs []railwayLog `json:"environmentLogs"`
		} `json:"data"`
	} `json:"payload"`
	Type LogType `json:"type"`
}
