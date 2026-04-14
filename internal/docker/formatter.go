package docker

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/moby/moby/client"
)

type Container struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Image   string `json:"image"`
	Created string `json:"created"`
	State   string `json:"state"`
	Status  string `json:"status"`
}

type Image struct {
	ID         string  `json:"id"`
	Created    string  `json:"created"`
	Size       float64 `json:"size"`
	Tag        string  `json:"tag"`
	Containers int     `json:"containers"`
}

func FormatContainerList(containerList client.ContainerListResult) (string, error) {

	if len(containerList.Items) == 0 {
		fmt.Println("Empty args provided!")
		return "[]", nil
	}

	var res []Container

	for _, ctr := range containerList.Items {
		name := "unknown"
		if len(ctr.Names) > 0 {
			name = strings.TrimPrefix(ctr.Names[0], "/")
		}

		res = append(res, Container{
			ID:      ctr.ID,
			Name:    name,
			Image:   ctr.Image,
			Created: strconv.FormatInt(ctr.Created, 10),
			State:   string(ctr.State),
			Status:  ctr.Status,
		})
	}

	jsonBytes, err := json.MarshalIndent(res, "", " ")
	if err != nil {
		log.Fatal("some err : ", err)
	}

	fmt.Println("RES : ", string(jsonBytes))
	return string(jsonBytes), nil
}

func FormatImgList(images client.ImageListResult) (string, error) {
	if len(images.Items) == 0 {
		fmt.Println("Empty args provided!")
		return "[]", nil
	}

	var res []Image

	for _, img := range images.Items {

		tags := "unknown"
		if len(img.RepoTags) > 0 {
			tags = img.RepoTags[0]
		}

		rawId := strings.TrimPrefix(img.ID, "sha256:")
		if len(rawId) > 12 {
			rawId = rawId[:12]
		}

		res = append(res, Image{
			ID:         rawId,
			Created:    strconv.FormatInt(img.Created, 10),
			Size:       float64(img.Size) / 1024 / 1024,
			Tag:        tags,
			Containers: int(img.Containers),
		})
	}

	jsonByte, err := json.MarshalIndent(res, "", " ")
	if err != nil {
		log.Fatal("Somer err : ", err)
	}

	fmt.Println("RES : ", string(jsonByte))

	return string(jsonByte), nil
}
