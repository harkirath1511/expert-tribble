package main

import (
	"log"

	"github.com/harkirath1511/docker-cli/internal/docker"
)

func main() {

	apiClient, err := docker.Init()
	if err != nil {
		log.Fatal("There's an err : ", err)
	}

	defer apiClient.Close()

	//docker.ListContainers(apiClient)
	//docker.InspectContainer(apiClient, "1642268ca3eed9f736e6ce343b4dc8ff6f8d41f69797c5d5e4f9b3169de7cd60")
	//docker.ProcInsideContainer(apiClient, "1642268ca3eed9f736e6ce343b4dc8ff6f8d41f69797c5d5e4f9b3169de7cd60")
	//docker.GetContainerLogs(apiClient, "1642268ca3eed9f736e6ce343b4dc8ff6f8d41f69797c5d5e4f9b3169de7cd60")
	//docker.GetContainerStats(apiClient, "1642268ca3eed9f736e6ce343b4dc8ff6f8d41f69797c5d5e4f9b3169de7cd60")
	//docker.StartContainer(apiClient, "1642268ca3eed9f736e6ce343b4dc8ff6f8d41f69797c5d5e4f9b3169de7cd60")
	//docker.StopContainer(apiClient, "1642268ca3eed9f736e6ce343b4dc8ff6f8d41f69797c5d5e4f9b3169de7cd60")
	//docker.RestartContainer(apiClient, "1642268ca3eed9f736e6ce343b4dc8ff6f8d41f69797c5d5e4f9b3169de7cd60")
	//docker.RenameContainer(apiClient, "1642268ca3eed9f736e6ce343b4dc8ff6f8d41f69797c5d5e4f9b3169de7cd60", "syncCord-cont")
	//docker.PauseContainer(apiClient, "1642268ca3eed9f736e6ce343b4dc8ff6f8d41f69797c5d5e4f9b3169de7cd60")
	//docker.UnpauseContainer(apiClient, "1642268ca3eed9f736e6ce343b4dc8ff6f8d41f69797c5d5e4f9b3169de7cd60")
	//docker.KillContainer(apiClient, "1642268ca3eed9f736e6ce343b4dc8ff6f8d41f69797c5d5e4f9b3169de7cd60")
	//docker.DeleteContainer(apiClient, "1642268ca3eed9f736e6ce343b4dc8ff6f8d41f69797c5d5e4f9b3169de7cd60")

	docker.ListImages(apiClient)
}
