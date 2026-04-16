package docker

import (
	"encoding/json"
	"fmt"
	"log"
	"net/netip"
	"strconv"
	"strings"

	"github.com/moby/moby/client"
)

type Container struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Image    string   `json:"image"`
	Created  string   `json:"created"`
	State    string   `json:"state"`
	Status   string   `json:"status"`
	Path     string   `json:"path,omitempty"`
	ExitCode int      `json:"exitcode,omitempty"`
	EnvVars  []string `json:"envs,omitempty"`
	IpAdd    string   `json:"ip,omitempty"`
	Command  []string `json:"cmd,omitempty"`
	Ports    []string `json:"ports,omitempty"`
}

type Image struct {
	ID           string   `json:"id"`
	Created      string   `json:"created"`
	Size         float64  `json:"size"`
	Tags         []string `json:"tags"`
	Containers   int      `json:"containers"`
	OS           string   `json:"os,omitempty"`
	Cmd          []string `json:"cmd"`
	Entrypoint   []string `json:"entrypoint"`
	ExposedPorts []string `json:"exposed_ports"`
	EnvVars      []string `json:"env_vars"`
	Architecture string   `json:"architecture"`
}

func FormatContList(containerList client.ContainerListResult) (string, error) {

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

func FormatContInspect(ctr client.ContainerInspectResult) (string, error) {

	var container Container

	var ip netip.Addr

	if len(ctr.Container.NetworkSettings.Networks) > 0 {
		for _, net := range ctr.Container.NetworkSettings.Networks {
			ip = net.IPAddress
			break
		}
	}

	var ports []string

	for intPort, bindings := range ctr.Container.NetworkSettings.Ports {
		if len(bindings) > 0 {
			ports = append(ports, fmt.Sprintf("%s -> %s : %s ", intPort, bindings[0].HostIP, bindings[0].HostPort))
		} else {
			ports = append(ports, fmt.Sprintf("%s -> %s", intPort, ""))
		}
	}

	container = Container{
		ID:       ctr.Container.ID,
		Name:     ctr.Container.Name,
		Image:    ctr.Container.Image,
		Created:  ctr.Container.Created,
		State:    string(ctr.Container.State.Status),
		Path:     ctr.Container.Path,
		ExitCode: ctr.Container.State.ExitCode,
		EnvVars:  ctr.Container.Config.Env,
		IpAdd:    ip.String(),
		Command:  ctr.Container.Config.Cmd,
		Ports:    ports,
	}

	jsonByte, err := json.MarshalIndent(container, "", " ")
	if err != nil {
		log.Fatal("Some err :( -> ", err)
	}

	fmt.Println(string(jsonByte))

	return string(jsonByte), nil
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

func FormatImgInspect(img client.ImageInspectResult) (string, error) {
	var res Image

	var imgTags []string

	if len(img.RepoTags) > 3 {
		for _, tags := range img.RepoTags {
			imgTags = append(imgTags, string(tags[0]), string(tags[1]), string(tags[3]))
			break
		}
	}

	res = Image{
		ID:         img.ID,
		Created:    img.Created,
		Size:       float64(img.Size),
		Tags:       imgTags,
		Containers: img.Manifests[0].ImageData.Containers,
		OS:         img.Os,
		Cmd:        img.Config.Cmd,
		EnvVars:    img.Config.Env,
	}

}
