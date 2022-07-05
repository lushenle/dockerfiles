package main

import (
	"bufio"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
)

var dockerRegistryUserID = "manunkind"

type ErrorLine struct {
	Error       string      `json:"error"`
	ErrorDetail ErrorDetail `json:"errorDetail"`
}

type ErrorDetail struct {
	Message string `json:"message"`
}

func print(rd io.Reader) error {
	var lastLine string

	scanner := bufio.NewScanner(rd)
	for scanner.Scan() {
		lastLine = scanner.Text()
		fmt.Println(scanner.Text())
	}

	errLine := &ErrorLine{}
	json.Unmarshal([]byte(lastLine), errLine)
	if errLine.Error != "" {
		return errors.New(errLine.Error)
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func imageBuild(dockerClient *client.Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*120)
	defer cancel()

	tar, err := archive.TarWithOptions(".", &archive.TarOptions{})
	if err != nil {
		log.Fatal(err)
	}
	opts := types.ImageBuildOptions{
		Dockerfile: "Dockerfile",
		Tags:       []string{dockerRegistryUserID + "/wireguard_exporter:go"},
		Remove:     true,
	}
	res, err := dockerClient.ImageBuild(ctx, tar, opts)
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()

	err = print(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func imagePush(dockerClient *client.Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*120)
	defer cancel()

	var authConfig = types.AuthConfig{
		Username:      "xxxxxxx",
		Password:      "xxxxxxx", //Docker Hub Password or Access Token
		ServerAddress: "https://index.docker.io/v1/",
	}

	authConfigBytes, _ := json.Marshal(authConfig)
	authConfigEncoded := base64.URLEncoding.EncodeToString(authConfigBytes)

	tag := dockerRegistryUserID + "/wireguard_exporter"
	opts := types.ImagePushOptions{RegistryAuth: authConfigEncoded}
	rd, err := dockerClient.ImagePush(ctx, tag, opts)
	if err != nil {
		log.Fatal(err)
	}

	defer rd.Close()

	err = print(rd)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func main() {
	log.SetFlags(log.Flags() | log.Lshortfile)
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(cli.ClientVersion())

	err = imageBuild(cli)
	if err != nil {
		log.Fatal(err)
	}

	err = imagePush(cli)
	if err != nil {
		log.Fatal(err)
	}
}
