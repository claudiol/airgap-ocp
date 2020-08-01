package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"os"
        "fmt"
        "strings"
	"bytes"
	
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func main() {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	// Needed to ensure API version compatibility
        cli.NegotiateAPIVersion(ctx)

	authConfig := types.AuthConfig{
		Username: "rhn-gps-claudiol",
		Password: "Azsx$123",
	}
	encodedJSON, err := json.Marshal(authConfig)
	if err != nil {
		panic(err)
	}
	authStr := base64.URLEncoding.EncodeToString(encodedJSON)

        //imageName := "registry.redhat.io/distributed-tracing/jaeger-ingester-rhel7"
	imageName := "registry.redhat.io/rhpam-7/rhpam-rhel8-operator"

	out, err := cli.ImagePull(ctx, imageName, types.ImagePullOptions{RegistryAuth: authStr})
	if err != nil {
		panic(err)
	}

	defer out.Close()
	io.Copy(os.Stdout, out)

	err = cli.ImageTag(ctx, imageName, imageName)
	if err != nil {
		panic(err)
	}
	
        imageIds := getImages(cli)
	saveImages(cli, imageIds)
}

func listImages(cli *client.Client) {
	//List all images available locally
	images, err := cli.ImageList(context.Background(), types.ImageListOptions{})
	if err != nil {
		panic(err)
	}

	fmt.Println("LIST IMAGES\n-----------------------")
	fmt.Println("Image ID | Repo Tags | Size")
	for _, image := range images {
		fmt.Printf("Saving: %s | %s | %d\n", image.ID, image.RepoTags, image.Size)
	}
}

func getImages(cli *client.Client) []string {
        var imageIds []string

	//List all images available locally
	images, err := cli.ImageList(context.Background(), types.ImageListOptions{})
	if err != nil {
		panic(err)
	}

	fmt.Println("LIST IMAGES\n-----------------------")
	fmt.Println("Image ID | Repo Tags | Size")
	for _, image := range images {
		fmt.Printf("Saving: %s | %s | %d\n", image.ID, image.RepoTags, image.Size)
                imageIds = append(imageIds, image.RepoTags[0])
	}
        fmt.Println(imageIds)
	return imageIds
}

func listCointainers(cli *client.Client) {
	//Retrieve a list of containers
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	fmt.Print("\n\n\n")
	fmt.Println("LIST CONTAINERS\n-----------------------")
	fmt.Println("Container Names | Image | Mounts")
	//Iterate through all containers and display each container's properties
	for _, container := range containers {
		fmt.Printf("%s | %s | %s\n", container.Names, container.Image, container.Mounts)
	}

}

func saveImages(cli *client.Client, imageIds [] string) {

     // We are going to iterate through all the imaes that we retrieved from the Docker Daemon
     for _, imageId := range imageIds {

        var tmpIds [] string

	tmpIds       = append(tmpIds, imageId)
        reader, err := cli.ImageSave(context.Background(), tmpIds)
        if err != nil {
           fmt.Println(err.Error())
        }

	buf := new(bytes.Buffer)
	buf.ReadFrom(reader)

	imageName := strings.Split(imageId,":")
	imageName  = strings.Split(imageName[0], "/")
	
	name := imageName[len(imageName)-1] + ".tar"
	f, err := os.Create(name)
	if err != nil {
	  panic(err)
	}

        fileBytes := buf.Bytes()

	// Write the file - This should be a file that can be loaded using
	// docker load -i file.tar or podman load -i file.tar
	f.Write(fileBytes)
	f.Close()
    }
}

