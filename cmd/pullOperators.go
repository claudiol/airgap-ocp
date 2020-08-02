/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"os"
        "strings"
	"bytes"
	
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
)


// pullOperatorsCmd represents the pullOperators command
var pullOperatorsCmd = &cobra.Command{
	Use:   "pullOperators",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
	     pullOperators(AirGapConfig.OcpDisconnectedOperators)
	     fmt.Println("pullOperators called")
	},
}

func init() {
	rootCmd.AddCommand(pullOperatorsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// pullOperatorsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// pullOperatorsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}

func pullOperators(operatorList [] string) {

     for _, operator := range operatorList {

	ctx      := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	
	if err != nil {
		panic(err)
	}

	// Needed to ensure API version compatibility
        cli.NegotiateAPIVersion(ctx)

	authConfig := types.AuthConfig{
		Username: AirGapConfig.UserName,
		Password: AirGapConfig.Password,
	}
	encodedJSON, err := json.Marshal(authConfig)
	if err != nil {
		panic(err)
	}
	authStr   := base64.URLEncoding.EncodeToString(encodedJSON)
	imageName := operator

	out, err  := cli.ImagePull(ctx, imageName, types.ImagePullOptions{RegistryAuth: authStr})
	
	if err != nil {
		fmt.Printf("ERROR: %s \n Could not pull image [%s] ... continuing. \n", err.Error(), imageName)
		continue
	}

	defer out.Close()
	io.Copy(os.Stdout, out)

	err = cli.ImageTag(ctx, imageName, imageName)
	if err != nil {
		panic(err)
	}
	
        imageIds := getImages(cli)
	saveImages(cli, imageIds, true)
     }

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

	for _, image := range images {
		fmt.Printf("Saving: %s | %s | %d\n", image.ID, image.RepoTags, image.Size)
                imageIds = append(imageIds, image.RepoTags[0])
	}
        //fmt.Println(imageIds)
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

func saveImages(cli *client.Client, imageIds [] string, remove bool) {

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

	imageName    := strings.Split(imageId,":")
	imageVersion := imageName[1]
	imageName     = strings.Split(imageName[0], "/")

	var name string
	if AirGapConfig.DestDir != "" {
	   name = AirGapConfig.DestDir + imageName[len(imageName)-1] + "-" + imageVersion + ".tar"
	} else {
	   name = imageName[len(imageName)-1] + "-" + imageVersion + ".tar"
	}
	
	f, err := os.Create(name)
	if err != nil {
	  panic(err)
	}

        fileBytes := buf.Bytes()

	// Write the file - This should be a file that can be loaded using
	// docker load -i file.tar or podman load -i file.tar
	f.Write(fileBytes)
	f.Close()
	if remove == true {
	   cli.ImageRemove(context.Background(), imageId, types.ImageRemoveOptions{ Force: true })
	}
    }
}

