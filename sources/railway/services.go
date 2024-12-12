package railway

import (
	"context"
	"errors"
)

func (gql *GraphQLConfig) getProjectInfo(ctx context.Context, config *Config) (*project, error) {
	if gql.client == nil {
		return nil, errors.New("client must not be nil")
	}

	environment := &environment{}

	variables := map[string]interface{}{
		"id": config.EnvironmentId,
	}

	if err := gql.client.Exec(ctx, environmentQuery, &environment, variables); err != nil {
		return nil, err
	}

	project := &project{}

	variables = map[string]interface{}{
		"id": environment.Environment.ProjectID,
	}

	if err := gql.client.Exec(ctx, projectQuery, &project, variables); err != nil {
		return nil, err
	}

	return project, nil
}

func (gql *GraphQLConfig) GetEnvironments(ctx context.Context, config *Config) (environments map[string]string, err error) {
	project, err := gql.getProjectInfo(ctx, config)
	if err != nil {
		return nil, err
	}

	for _, environment := range project.Project.Environments.Edges {
		environments[environment.Node.ID] = environment.Node.Name
	}

	return
}

func (gql *GraphQLConfig) GetServices(ctx context.Context, config *Config) (services map[string]string, err error) {
	project, err := gql.getProjectInfo(ctx, config)
	if err != nil {
		return nil, err
	}

	for _, service := range project.Project.Services.Edges {
		services[service.Node.ID] = service.Node.Name
	}

	return
}
