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
	ContCnt      int      `json:"contcnt"`
	Container    []string `json:"containers,omitempty"`
	OS           string   `json:"os,omitempty"`
	Cmd          []string `json:"cmd"`
	Entrypoint   []string `json:"entrypoint"`
	ExposedPorts []string `json:"exposed_ports"`
	EnvVars      []string `json:"env_vars"`
	Architecture string   `json:"architecture,omitempty"`
}

type ImgSrch struct {
	Description string `json:"desc,omitempty"`
	IsOfficial  bool   `json:"isofficial"`
	Name        string `json:"name"`
	StarCnt     int    `json:"starcnt"`
}

//Container funcs

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

func FormatContLogs(data string) string {
	maxChars := 3000

	if len(data) <= maxChars {
		fmt.Println("Logs : ", data)
		return data
	}

	fmt.Println("reduced Logs : ", data[len(data)-maxChars:])
	return "...RLOeduced Length logs...\n" + data[len(data)-maxChars:]
}

func FormatProcRes(rawData client.ContainerTopResult) string {

	res := make(map[string]string)

	titles := rawData.Titles
	proc := rawData.Processes

	if len(titles) > 0 && len(proc) > 0 {
		for i := 0; i < len(titles); i++ {
			res[titles[i]] = fmt.Sprintf("%s , %s", proc[0][i], proc[1][i])
		}
	}

	jsonBytes, err := json.MarshalIndent(res, "", " ")
	if err != nil {
		log.Fatal("Some err : ", err)
	}

	fmt.Println("res : ", string(jsonBytes))
	return string(jsonBytes)
}

//Imgs funcs

func FormatImgList(images client.ImageListResult) (string, error) {
	if len(images.Items) == 0 {
		fmt.Println("Empty args provided!")
		return "[]", nil
	}

	var res []Image

	for _, img := range images.Items {
		var imgTags []string

		if len(img.RepoTags) > 3 {
			imgTags = append(imgTags, string(img.RepoTags[0]), string(img.RepoTags[1]), string(img.RepoTags[3]))
		}

		rawId := strings.TrimPrefix(img.ID, "sha256:")

		res = append(res, Image{
			ID:      rawId,
			Created: strconv.FormatInt(img.Created, 10),
			Size:    float64(img.Size) / 1024 / 1024,
			Tags:    imgTags,
			ContCnt: int(img.Containers),
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
		imgTags = append(imgTags, string(img.RepoTags[0]), string(img.RepoTags[1]), string(img.RepoTags[3]))
	}

	var exposedPorts []string

	for _, ports := range img.Config.ExposedPorts {
		exposedPorts = append(exposedPorts, fmt.Sprintf("%s , ", ports))
	}

	var containers = []string{""}

	if len(img.Manifests) > 0 {
		containers = img.Manifests[0].ImageData.Containers
	}

	res = Image{
		ID:           img.ID,
		Created:      img.Created,
		Size:         float64(img.Size),
		Tags:         imgTags,
		Container:    containers,
		OS:           img.Os,
		Cmd:          img.Config.Cmd,
		Entrypoint:   img.Config.Entrypoint,
		ExposedPorts: exposedPorts,
		EnvVars:      img.Config.Env,
	}

	jsonBytes, err := json.MarshalIndent(res, "", " ")
	if err != nil {
		log.Fatal("There's some err : ", err)
	}

	fmt.Println("Res -> ", string(jsonBytes))

	return string(jsonBytes), nil
}

func FormatImgSrchRes(imgRes client.ImageSearchResult) (string, error) {

	var topRes []ImgSrch

	for cnt, res := range imgRes.Items {
		
		if cnt>5 {
			break
		}

		el := ImgSrch{
			Description: res.Description,
			IsOfficial:  res.IsOfficial,
			Name:        res.Name,
			StarCnt:     res.StarCount,
		}

		topRes = append(topRes, el)
	}

	jsonByte, err := json.MarshalIndent(topRes, "", " ")
	if err != nil {
		log.Fatal("there's some err : ", err)
	}

	fmt.Println("res : ", string(jsonByte))

	return string(jsonByte), nil
}
