package docker

import (
	"context"
	"fmt"
	"log"

	"github.com/moby/moby/client"
);


func ListContainers(apiClient *client.Client) {

	res, err := apiClient.ContainerList(context.Background(), client.ContainerListOptions{
		All: true,
	})

	if err != nil {
		panic(err)
	}
	fmt.Printf("%s  %-22s  %s\n", "ID", "STATUS", "IMAGE")

	for _, ctr := range res.Items{
		fmt.Printf("%s  %-22s  %s\n", ctr.ID, ctr.Status, ctr.Image)
	}
}

func InspectContainer(apiclient *client.Client, id string){
	if id==""{
		log.Fatal("You need to provide an id")
	}

	res, err := apiclient.ContainerInspect(context.Background(), id, client.ContainerInspectOptions{})
	if err!=nil{
		log.Fatal("some err in inspecting container: ",err)
	}

	fmt.Println("CONTAINER DETAILS : ")
	fmt.Printf("Id : %s \n",res.Container.ID)
	fmt.Printf("Created : %s \n",res.Container.Created)
	fmt.Printf("Path : %s \n",res.Container.Path)
	fmt.Printf("Img : %s \n",res.Container.Image)
	fmt.Printf("Name : %s \n",res.Container.Name)
	fmt.Printf("Platform : %s \n",res.Container.Platform)
	fmt.Printf("Args : %s \n",res.Container.Args)
	fmt.Println("State : ",res.Container.State)
}


