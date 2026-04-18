package docker

import (
	"fmt"

	"github.com/moby/moby/client"
)

func Execute(apiClient *client.Client, name string, args map[string]any) (string, error) {
	tools, exists := DockerTools[name]
	if !exists {
		return "", fmt.Errorf("ai requested unknown tool: %s", name)
	}
	return tools.Execute(apiClient, args)
}
