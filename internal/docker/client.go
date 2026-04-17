package docker

import (
	"context"
	"fmt"
	"log"

	"github.com/moby/moby/client"
)

func Init() (*client.Client, error) {

	apiClient, err := client.New(
		client.FromEnv,
		client.WithUserAgent("docker-ai/1.0.0"),
		client.WithAPIVersionNegotiation(),
	)
	if err != nil {
		log.Fatal("Noo there is was an err : ", err)
	}

	res, err := apiClient.Ping(context.TODO(), client.PingOptions{
		NegotiateAPIVersion: true,
	})
	if err != nil {
		log.Fatal("Noo there is an err : ", err)
	}

	fmt.Println("RESULT : ", res)
	fmt.Println("Success!")
	return apiClient, nil
}
