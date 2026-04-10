package docker

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"

	"github.com/harkirath1511/docker-cli/internal/utils"
	"github.com/moby/moby/client"
)

//all needed get func's for containers

func ListContainers(apiClient *client.Client) {

	res, err := apiClient.ContainerList(context.Background(), client.ContainerListOptions{
		All: true,
	})

	if err != nil {
		panic(err)
	}
	fmt.Printf("%s  %-22s  %s\n", "ID", "STATUS", "IMAGE")

	for _, ctr := range res.Items {
		fmt.Printf("%s  %-22s  %s\n", ctr.ID, ctr.Status, ctr.Image)
	}
}

func InspectContainer(apiclient *client.Client, id string) {
	if id == "" {
		log.Fatal("You need to provide an id")
	}

	res, err := apiclient.ContainerInspect(context.Background(), id, client.ContainerInspectOptions{})
	if err != nil {
		log.Fatal("some err in inspecting container: ", err)
	}

	fmt.Println("CONTAINER DETAILS : ")
	fmt.Printf("Id : %s \n", res.Container.ID)
	fmt.Printf("Created : %s \n", res.Container.Created)
	fmt.Printf("Path : %s \n", res.Container.Path)
	fmt.Printf("Img : %s \n", res.Container.Image)
	fmt.Printf("Name : %s \n", res.Container.Name)
	fmt.Printf("Platform : %s \n", res.Container.Platform)
	fmt.Printf("Args : %s \n", res.Container.Args)
	fmt.Println("State : ", res.Container.State)
}

func ProcInsideContainer(apiclient *client.Client, id string) {
	if id == "" {
		log.Fatal("You need to provide an id")
	}

	res, err := apiclient.ContainerTop(context.Background(), id, client.ContainerTopOptions{})
	if err != nil {
		log.Fatal("some err in inspecting container: ", err)
	}

	fmt.Println(res.Processes)
}

func GetContainerLogs(apiclient *client.Client, id string) {
	if id == "" {
		log.Fatal("You need to provide an id")
	}

	res, err := apiclient.ContainerLogs(context.Background(), id, client.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
	})
	if err != nil {
		log.Fatal("some err in inspecting container: ", err)
	}

	defer res.Close()

	data, err := io.ReadAll(res)
	if err != nil {
		log.Fatal("Some err : ", err)
	}

	fmt.Println("CONTAINER LOGS : ")
	fmt.Println(string(data))

}

func GetContainerStats(apiclient *client.Client, id string) {
	if id == "" {
		log.Fatal("You need to provide an id")
	}

	res, err := apiclient.ContainerStats(context.Background(), id, client.ContainerStatsOptions{})
	if err != nil {
		log.Fatal("Some err : ", err)
	}
	defer res.Body.Close()

	jsonData, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal("err while reading stats body: ", err)
	}

	var mapData map[string]any

	if !json.Valid(jsonData) {
		log.Fatal("invalid JSON from stats response")
	}

	if err := json.Unmarshal(jsonData, &mapData); err != nil {
		log.Fatal("err : ", err)
	}

	data, _ := json.MarshalIndent(mapData, "", " ")
	fmt.Println("DATA : ", string(data))

	memoryStats, ok := utils.GetMap(mapData, "memory_stats")
	if !ok {
		log.Fatal("memory_stats missing")
	}

	memoryStatsInner, ok := utils.GetMap(memoryStats, "stats")
	if !ok {
		log.Fatal("memory_stats.stats missing")
	}

	memoryUsage, ok := utils.GetFloat(memoryStats, "usage")
	if !ok {
		log.Fatal("memory_stats.usage missing")
	}

	inactiveFile, ok := utils.GetFloat(memoryStatsInner, "inactive_file")
	if !ok {
		log.Fatal("memory_stats.stats.inactive_file missing")
	}

	memoryLimit, ok := utils.GetFloat(memoryStats, "limit")
	if !ok {
		log.Fatal("memory_stats.limit missing")
	}

	usedMemory := memoryUsage - inactiveFile
	availableMemory := memoryLimit

	memoryUsagePct := 0.0
	if availableMemory > 0 {
		memoryUsagePct = (usedMemory / availableMemory) * 100.0
	}

	cpuStats, ok := utils.GetMap(mapData, "cpu_stats")
	if !ok {
		log.Fatal("cpu_stats missing")
	}

	preCPUStats, ok := utils.GetMap(mapData, "precpu_stats")
	if !ok {
		log.Fatal("precpu_stats missing")
	}

	cpuUsage, ok := utils.GetMap(cpuStats, "cpu_usage")
	if !ok {
		log.Fatal("cpu_stats.cpu_usage missing")
	}

	preCPUUsage, ok := utils.GetMap(preCPUStats, "cpu_usage")
	if !ok {
		log.Fatal("precpu_stats.cpu_usage missing")
	}

	totalUsage, ok := utils.GetFloat(cpuUsage, "total_usage")
	if !ok {
		log.Fatal("cpu_stats.cpu_usage.total_usage missing")
	}

	preTotalUsage, ok := utils.GetFloat(preCPUUsage, "total_usage")
	if !ok {
		log.Fatal("precpu_stats.cpu_usage.total_usage missing")
	}

	systemCPUUsage, ok := utils.GetFloat(cpuStats, "system_cpu_usage")
	if !ok {
		log.Fatal("cpu_stats.system_cpu_usage missing")
	}

	preSystemCPUUsage, ok := utils.GetFloat(preCPUStats, "system_cpu_usage")
	if !ok {
		log.Fatal("precpu_stats.system_cpu_usage missing")
	}

	cpuDelta := totalUsage - preTotalUsage
	systemCPUDelta := systemCPUUsage - preSystemCPUUsage

	numberCPUs := 0.0
	if percpuRaw, ok := cpuUsage["percpu_usage"]; ok {
		if percpu, ok := percpuRaw.([]any); ok && len(percpu) > 0 {
			numberCPUs = float64(len(percpu))
		}
	}
	if numberCPUs == 0 {
		if onlineCPUs, ok := utils.GetFloat(cpuStats, "online_cpus"); ok {
			numberCPUs = onlineCPUs
		}
	}

	cpuUsagePct := 0.0
	if systemCPUDelta > 0 && cpuDelta > 0 && numberCPUs > 0 {
		cpuUsagePct = (cpuDelta / systemCPUDelta) * numberCPUs * 100.0
	}

	fmt.Printf("used_memory: %.0f bytes\n", usedMemory)
	fmt.Printf("available_memory: %.0f bytes\n", availableMemory)
	fmt.Printf("memory_usage_pct: %.2f%%\n", memoryUsagePct)
	fmt.Printf("cpu_delta: %.0f\n", cpuDelta)
	fmt.Printf("system_cpu_delta: %.0f\n", systemCPUDelta)
	fmt.Printf("number_cpus: %.0f\n", numberCPUs)
	fmt.Printf("cpu_usage_pct: %.2f%%\n", cpuUsagePct)
}
