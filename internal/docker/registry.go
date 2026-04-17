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
		Name: "list_containers",
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
		Name: "inspect_container",
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
		Name: "container_processes",
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
		Name: "container_logs",
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
		Name: "container_stats",
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
		Name: "start_container",
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
		Name: "stop_container",
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
		Name: "restart_container",
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
		Name: "rename_container",
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
		Name: "pause_container",
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
		Name: "unpause_container",
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
		Name: "kill_container",
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
		Name: "delete_container",
		Execute: func(apiClient *client.Client, args map[string]interface{}) (string, error) {
			name, err := utils.GetArg(args, "name")
			if err != nil {
				return "", err
			}

			res := DeleteContainer(apiClient, name)

			return res, nil
		},
	},
}
