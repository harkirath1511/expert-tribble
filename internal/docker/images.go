package docker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/moby/moby/client"
)

func ListImages(apiclient *client.Client){
	res, err := apiclient.ImageList(context.Background(), client.ImageListOptions{})
	if err!=nil{
		log.Fatal("Some err listing all imgs : ",err)
	}
	//var data map[string]any

	data,_ := json.MarshalIndent(res, "", " ")
	fmt.Println("Res : ",string(data))
}