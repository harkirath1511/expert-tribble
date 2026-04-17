package docker

import (
	"archive/tar"
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/moby/moby/client"
)

// get func's for images

func ListImages(apiclient *client.Client) client.ImageListResult {
	res, err := apiclient.ImageList(context.Background(), client.ImageListOptions{
		Manifests: true,
	})
	if err != nil {
		log.Fatal("Some err listing all imgs : ", err)
	}
	return res
}

func InspectImg(apiclient *client.Client, id string) client.ImageInspectResult {

	if id == "" {
		log.Fatal("Id or name is missing")
	}

	res, err := apiclient.ImageInspect(context.Background(), id)
	if err != nil {
		log.Fatal("Some err inspecting the img : ", err)
	}

	return res
}

func SearchForImg(apiclient *client.Client, name string) client.ImageSearchResult {

	if name == "" {
		log.Fatal("Name of img is missing")
	}
	res, err := apiclient.ImageSearch(context.Background(), name, client.ImageSearchOptions{})
	if err != nil {
		log.Fatal("Some err inspecting the img : ", err)
	}

	return res
}

//post/del func's for images

func DeleteImg(apiclient *client.Client, id string) {
	if id == "" {
		log.Fatal("Name or id of img is missing")
	}
	_, err := apiclient.ImageRemove(context.Background(), id, client.ImageRemoveOptions{})
	if err != nil {
		log.Fatal("Some err removing the img : ", err)
	}

	fmt.Println("Successfully removed img!")
}

func BuildImg(apiclient *client.Client, path string, tag string) {
	root := filepath.Join(path, "internal", "docker")
	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)
	filepath.Walk(root, func(p string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}
		rel, _ := filepath.Rel(root, p)
		h, _ := tar.FileInfoHeader(info, "")
		h.Name = filepath.ToSlash(rel)
		_ = tw.WriteHeader(h)
		f, _ := os.Open(p)
		defer f.Close()
		_, _ = io.Copy(tw, f)
		return nil
	})
	_ = tw.Close()
	tar := io.NopCloser(bytes.NewReader(buf.Bytes()))
	defer tar.Close()

	res, err := apiclient.ImageBuild(context.Background(), tar, client.ImageBuildOptions{
		Dockerfile: "Dockerfile",
		Remove:     true,
		Tags:       []string{tag},
	})
	if err != nil {
		log.Fatal("some err : ", err)
	}

	defer res.Body.Close()

	io.Copy(os.Stdout, res.Body)

	fmt.Println("The response : ", res)
}

func CreateImg(apiclient *client.Client, image string) {
	res, err := apiclient.ImagePull(context.Background(), image, client.ImagePullOptions{})
	if err != nil {
		log.Fatal("Some err : ", err)
	}
	defer res.Close()

	data, _ := io.Copy(os.Stdout, res)
	fmt.Println("res : ", data)
}
