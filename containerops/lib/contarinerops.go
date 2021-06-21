package lib

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/docker/docker/client"

	"github.com/docker/docker/api/types"
	"github.com/olekukonko/tablewriter"
)

var (
	Host string
	Port string
	Name string
	CTX  = context.Background()
)

func dropError(err error) error {
	if err != nil {
		log.Fatal(err)
	}
	return err
}

func cli() *client.Client {
	var (
		cli *client.Client
		err error
	)

	if Host == "" {
		cli, err = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
		dropError(err)
	} else {
		cli, err = client.NewClientWithOpts(client.WithHost("http://"+Host+":"+Port), client.WithTimeout(30*time.Second), client.WithScheme("http"), client.WithAPIVersionNegotiation())
		dropError(err)
	}
	return cli
}

func ListContainer() {
	//ctx := context.Background()
	containers, err := cli().ContainerList(CTX, types.ContainerListOptions{All: true})
	dropError(err)
	table := tablewriter.NewWriter(os.Stdout)
	table.SetRowLine(true)
	table.SetHeader([]string{"Name", "Image", "State", "Status"})
	table.SetCenterSeparator("*")
	table.SetColumnSeparator("╪")
	table.SetRowSeparator("-")

	for _, container := range containers {
		tableRow := []string{strings.TrimRight(strings.Trim(fmt.Sprint(container.Names), "[/"), "]"), container.Image, container.State, container.Status}
		table.Append(tableRow)
	}
	table.Render()
}

func StopContainer() {
	fmt.Println("Stop container:", Name)
	err := cli().ContainerStop(CTX, Name, nil)
	dropError(err)
}

func StopALLContainers() {
	containers, err := cli().ContainerList(CTX, types.ContainerListOptions{All: true, Quiet: true})
	dropError(err)
	for _, container := range containers {
		fmt.Println("Stop container:", strings.TrimRight(strings.Trim(fmt.Sprint(container.Names), "[/"), "]"))
		err := cli().ContainerStop(CTX, container.ID, nil)
		dropError(err)
	}
}

func RemoveContainer() {
	fmt.Println("Remove container:", Name)
	err := cli().ContainerRemove(CTX, Name, types.ContainerRemoveOptions{Force: true, RemoveVolumes: false, RemoveLinks: false})
	dropError(err)
}

func RemoveAllContainers() {
	containers, err := cli().ContainerList(CTX, types.ContainerListOptions{All: true, Quiet: true})
	dropError(err)
	for _, container := range containers {
		fmt.Println("Remove container:", strings.TrimRight(strings.Trim(fmt.Sprint(container.Names), "[/"), "]"))
		err := cli().ContainerRemove(CTX, container.ID, types.ContainerRemoveOptions{Force: true, RemoveVolumes: false, RemoveLinks: false})
		dropError(err)
	}
}
