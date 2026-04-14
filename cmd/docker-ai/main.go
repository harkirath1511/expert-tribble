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

	//containerRes := docker.ListContainers(apiClient)
	//docker.FormatContainerList(containerRes)

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

	// imgList := docker.ListImages(apiClient)
	// docker.FormatImgList(imgList)

	//docker.InspectImg(apiClient, "sha256:3cde66018e19cd9af6ae6dc4efd4d5174ffa73ade50ccc63fb6a710fc810d8b6")
	//docker.SearchForImg(apiClient, "node")
	//docker.DeleteImg(apiClient, "sha256:b8cf5e598b72087903acce6c5ca4292cf991a5dbc729f6e7f5783163117f5513")
	//docker.BuildImg(apiClient, "../..", "ai-docker-harkirat")
	//docker.CreateImg(apiClient, "alpine")
}
