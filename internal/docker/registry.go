package docker

import (
	"github.com/harkirath1511/docker-cli/internal/utils"
	"github.com/moby/moby/client"
)

type Tool struct {
	Name        string `json:"name"`
	Description string `json:"desc"`
	Execute     func(apiClient *client.Client, args map[string]interface{}) (string, error)
}

var DockerTools = map[string]Tool{
	"list_containers": Tool{
		Name:        "list_containers",
		Description: "List all Docker containers with their basic information",
		Execute: func(apiClient *client.Client, args map[string]interface{}) (string, error) {
			containerRes := ListContainers(apiClient)
			res, err := FormatContList(containerRes)
			if err != nil {
				return "", err
			}
			return res, nil
		},
	},
	"inspect_container": Tool{
		Name:        "inspect_container",
		Description: "Get detailed information about a specific container including configuration and state",
		Execute: func(apiClient *client.Client, args map[string]interface{}) (string, error) {
			name, err := utils.GetArg(args, "name")
			if err != nil {
				return "", err
			}

			containerRes := InspectContainer(apiClient, name)
			res, err2 := FormatContInspect(containerRes)
			if err2 != nil {
				return "", err2
			}
			return res, nil
		},
	},
	"container_processes": Tool{
		Name:        "container_processes",
		Description: "List all running processes inside a specific container",
		Execute: func(apiClient *client.Client, args map[string]interface{}) (string, error) {
			name, err := utils.GetArg(args, "name")
			if err != nil {
				return "", err
			}

			containerRes := ProcInsideContainer(apiClient, name)
			res, err2 := FormatProcRes(containerRes)
			if err2 != nil {
				return "", err2
			}
			return res, nil
		},
	},
	"container_logs": Tool{
		Name:        "container_logs",
		Description: "Retrieve and display logs from a specific container",
		Execute: func(apiClient *client.Client, args map[string]interface{}) (string, error) {
			name, err := utils.GetArg(args, "name")
			if err != nil {
				return "", err
			}

			containerRes := GetContainerLogs(apiClient, name)
			res, err2 := FormatContLogs(containerRes)
			if err2 != nil {
				return "", err2
			}
			return res, nil
		},
	},
	"container_stats": Tool{
		Name:        "container_stats",
		Description: "Display live resource usage statistics for a specific container",
		Execute: func(apiClient *client.Client, args map[string]interface{}) (string, error) {
			name, err := utils.GetArg(args, "name")
			if err != nil {
				return "", err
			}

			res := GetContainerStats(apiClient, name)
			return res, nil
		},
	},
	"start_container": Tool{
		Name:        "start_container",
		Description: "Start a stopped container",
		Execute: func(apiClient *client.Client, args map[string]interface{}) (string, error) {
			name, err := utils.GetArg(args, "name")
			if err != nil {
				return "", err
			}

			res := StartContainer(apiClient, name)

			return res, nil
		},
	},
	"stop_container": Tool{
		Name:        "stop_container",
		Description: "Stop a running container gracefully",
		Execute: func(apiClient *client.Client, args map[string]interface{}) (string, error) {
			name, err := utils.GetArg(args, "name")
			if err != nil {
				return "", err
			}

			res := StopContainer(apiClient, name)

			return res, nil
		},
	},
	"restart_container": Tool{
		Name:        "restart_container",
		Description: "Restart a running or stopped container",
		Execute: func(apiClient *client.Client, args map[string]interface{}) (string, error) {
			name, err := utils.GetArg(args, "name")
			if err != nil {
				return "", err
			}

			res := RestartContainer(apiClient, name)

			return res, nil
		},
	},
	"rename_container": Tool{
		Name:        "rename_container",
		Description: "Rename an existing container",
		Execute: func(apiClient *client.Client, args map[string]interface{}) (string, error) {
			oldName, err := utils.GetArg(args, "oldName")
			if err != nil {
				return "", err
			}

			newName, err := utils.GetArg(args, "newName")
			if err != nil {
				return "", err
			}

			res := RenameContainer(apiClient, oldName, newName)

			return res, nil
		},
	},
	"pause_container": Tool{
		Name:        "pause_container",
		Description: "Pause all processes in a running container",
		Execute: func(apiClient *client.Client, args map[string]interface{}) (string, error) {
			name, err := utils.GetArg(args, "name")
			if err != nil {
				return "", err
			}

			res := PauseContainer(apiClient, name)

			return res, nil
		},
	},
	"unpause_container": Tool{
		Name:        "unpause_container",
		Description: "Resume a paused container",
		Execute: func(apiClient *client.Client, args map[string]interface{}) (string, error) {
			name, err := utils.GetArg(args, "name")
			if err != nil {
				return "", err
			}

			res := UnpauseContainer(apiClient, name)

			return res, nil
		},
	},
	"kill_container": Tool{
		Name:        "kill_container",
		Description: "Forcefully stop a running container",
		Execute: func(apiClient *client.Client, args map[string]interface{}) (string, error) {
			name, err := utils.GetArg(args, "name")
			if err != nil {
				return "", err
			}

			res := KillContainer(apiClient, name)

			return res, nil
		},
	},
	"delete_container": Tool{
		Name:        "delete_container",
		Description: "Remove a container",
		Execute: func(apiClient *client.Client, args map[string]interface{}) (string, error) {
			name, err := utils.GetArg(args, "name")
			if err != nil {
				return "", err
			}

			res := DeleteContainer(apiClient, name)

			return res, nil
		},
	},

	//img funcs
	"list_images": Tool{
		Name:        "list_images",
		Description: "List all Docker images available on the local machine",
		Execute: func(apiClient *client.Client, args map[string]interface{}) (string, error) {
			imgRes := ListImages(apiClient)

			res, err := FormatImgList(imgRes)
			if err != nil {
				return "", err
			}
			return res, nil
		},
	},
	"inspect_image": Tool{
		Name:        "inspect_image",
		Description: "Inspect a Docker image and return its detailed metadata",
		Execute: func(apiClient *client.Client, args map[string]interface{}) (string, error) {
			name, err := utils.GetArg(args, "name")
			if err != nil {
				return "", err
			}

			imgRes := InspectImg(apiClient, name)
			res, err2 := FormatImgInspect(imgRes)
			if err2 != nil {
				return "", err2
			}

			return res, nil
		},
	},
	"search_image": Tool{
		Name:        "search_image",
		Description: "Search Docker Hub or configured registries for images matching a name",
		Execute: func(apiClient *client.Client, args map[string]interface{}) (string, error) {
			name, err := utils.GetArg(args, "name")
			if err != nil {
				return "", err
			}

			imgRes := SearchForImg(apiClient, name)
			res, err2 := FormatImgSrchRes(imgRes)
			if err2 != nil {
				return "", err2
			}

			return res, nil
		},
	},
	"delete_image": Tool{
		Name:        "delete_image",
		Description: "Remove a Docker image from the local machine",
		Execute: func(apiClient *client.Client, args map[string]interface{}) (string, error) {
			name, err := utils.GetArg(args, "name")
			if err != nil {
				return "", err
			}

			res := DeleteImg(apiClient, name)
			return res, nil
		},
	},
	"build_image": Tool{
		Name:        "build_image",
		Description: "Build a Docker image from a Dockerfile in the provided path and tag it",
		Execute: func(apiClient *client.Client, args map[string]interface{}) (string, error) {
			path, err := utils.GetArg(args, "path")
			if err != nil {
				return "", err
			}

			tag, err := utils.GetArg(args, "tag")
			if err != nil {
				return "", err
			}

			res := BuildImg(apiClient, path, tag)
			return res, nil
		},
	},
	"create_image": Tool{
		Name:        "create_image",
		Description: "Pull a Docker image from a registry using its image reference",
		Execute: func(apiClient *client.Client, args map[string]interface{}) (string, error) {
			name, err := utils.GetArg(args, "name")
			if err != nil {
				return "", err
			}

			res := CreateImg(apiClient, name)
			return res, nil
		},
	},
}
