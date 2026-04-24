package main

import (
	"log"

	"github.com/harkirath1511/docker-cli/internal/docker"
	"github.com/harkirath1511/docker-cli/internal/llm"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	apiClient, err := docker.Init()
	if err != nil {
		log.Fatal("There's an err : ", err)
	}

	defer apiClient.Close()

	llm.GenerateRes("Hello , hi there!")

	// containerRes := docker.ListContainers(apiClient)
	// docker.FormatContList(containerRes)

	// inspecRes := docker.InspectContainer(apiClient, "chatbot-backend-1")
	// docker.FormatContInspect(inspecRes)

	// procRes := docker.ProcInsideContainer(apiClient, "3ea4b760a05a1e064f1f3fbdf4399d8292b7d890ecf8521402f2a6d7f5ee1ffd")
	// docker.FormatProcRes(procRes)

	// logs := docker.GetContainerLogs(apiClient, "chatbot-backend-1")
	// docker.FormatContLogs(logs)

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

	// inspectImg := docker.InspectImg(apiClient, "sha256:e35e37cf82b6894049e69e168e1135cbc2a084f3b63bdccc7908afff1bdc57d6")
	// docker.FormatImgInspect(inspectImg)

	// srchRes := docker.SearchForImg(apiClient, "node")
	// docker.FormatImgSrchRes(srchRes)

	//docker.DeleteImg(apiClient, "sha256:b8cf5e598b72087903acce6c5ca4292cf991a5dbc729f6e7f5783163117f5513")
	//docker.BuildImg(apiClient, "../..", "ai-docker-harkirat")
	//docker.CreateImg(apiClient, "alpine")
}
