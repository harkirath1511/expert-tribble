package docker

import (
	"context"
	"fmt"

	"github.com/moby/moby/client"
)

func Init() (*client.Client, error) {

	apiClient, err := client.New(
		client.FromEnv,
		client.WithUserAgent("docker-ai/1.0.0"),
		client.WithAPIVersionNegotiation(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create docker client: %w", err)
	}

	_, err = apiClient.Ping(context.TODO(), client.PingOptions{
		NegotiateAPIVersion: true,
	})
	if err != nil {
		return nil, fmt.Errorf("docker daemon unreachable: %w\n\nIs Docker Desktop running?", err)
	}

	return apiClient, nil
}

