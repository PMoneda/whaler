package whaler

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

//Container is a basic representation of a docker container
type Container struct {
	Name    string
	ID      string
	Image   string
	Volumes []string
	Env     []string
	Ports   []string
}

//Image is a basic representation of a docker image
type Image struct {
	ID   string
	Name string
}

//Network is a basic representation of a docker network
type Network struct {
	Name string
	ID   string
}

//CreateContainerConfig is a basic configuration structure to create a docker container
type CreateContainerConfig struct {
	Name        string
	Image       string
	Volumes     []string
	Env         []string
	Ports       []string
	NetworkName string
}

//BuildImageConfig is a basic configuration to build an image
type BuildImageConfig struct {
	PathContext string
	Dockerfile  string
	Tag         string
}

//GetContainers from docker
func GetContainers() ([]Container, error) {
	cli, err := client.NewEnvClient()
	if err != nil {
		return nil, err
	}
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		return nil, err
	}
	_containers := make([]Container, 0, 0)
	for _, container := range containers {
		_containers = append(_containers, Container{ID: container.ID, Name: container.Names[0], Image: container.Image})
	}
	return _containers, nil
}

//BuildImageWithDockerfile builds an image using a specific dockerfile
func BuildImageWithDockerfile(config BuildImageConfig) (string, error) {
	buf := bytes.NewBuffer(nil)
	if config.PathContext == "" {
		if wd, err := os.Getwd(); err != nil {
			return "", err
		} else {
			config.PathContext = wd
		}
	}
	err := compress(config.PathContext, buf)
	if err != nil {
		return "", err
	}
	if config.Dockerfile == "" {
		config.Dockerfile = "Dockerfile"
	}

	return buildDockerImage(config.Dockerfile, config.Tag, buf)
}

func buildDockerImage(dockerfile, tag string, ctx io.Reader) (string, error) {
	cli, err := client.NewEnvClient()
	if err != nil {
		return "", err
	}
	resp, err := cli.ImageBuild(context.Background(), ctx, types.ImageBuildOptions{
		Dockerfile:  dockerfile,
		Tags:        []string{tag},
		NetworkMode: "bridge",
		NoCache:     true,
	})
	if err != nil {
		return "", err
	}
	b, _ := ioutil.ReadAll(resp.Body)
	return string(b), nil
}

//CreateContainer creates a new container on Docker
func CreateContainer(config CreateContainerConfig) (string, error) {
	cli, err := client.NewEnvClient()
	if err != nil {
		return "", err
	}
	containerConfig := &container.Config{
		Env:             config.Env,
		Image:           config.Image,
		AttachStderr:    false,
		AttachStdin:     false,
		AttachStdout:    false,
		NetworkDisabled: false,
		ExposedPorts:    nat.PortSet{},
	}
	net := &Network{}
	if config.NetworkName != "" {
		var ex error
		net, ex = FindNetworkByName(config.NetworkName)
		if ex != nil {
			return "", err
		}
	}

	if config.Volumes == nil {
		config.Volumes = make([]string, 0)
	}
	hostConfig := &container.HostConfig{
		NetworkMode:  container.NetworkMode(config.NetworkName),
		Binds:        config.Volumes,
		PortBindings: make(nat.PortMap),
	}
	for _, port := range config.Ports {
		parts := strings.Split(port, ":")
		pb := nat.PortBinding{}
		pb.HostIP = ""
		pb.HostPort = parts[0]

		var s = nat.Port(parts[1] + "/tcp")
		var a struct{}
		containerConfig.ExposedPorts[s] = a
		list := make([]nat.PortBinding, 0, 1)
		list = append(list, pb)
		hostConfig.PortBindings[s] = list
	}
	networkConfig := &network.NetworkingConfig{
		EndpointsConfig: map[string]*network.EndpointSettings{config.NetworkName: &network.EndpointSettings{
			NetworkID: net.ID,
			Aliases:   []string{config.Name},
		}},
	}
	if config.NetworkName == "" {
		networkConfig = nil
	}
	resp, err := cli.ContainerCreate(context.Background(), containerConfig, hostConfig, networkConfig, config.Name)
	if err != nil {
		return "", err
	}
	return resp.ID, nil
}

//FindNetworkByName look for a docker network with specific name
func FindNetworkByName(name string) (*Network, error) {
	networks, err := GetNetworks()
	if err != nil {
		return nil, err
	}
	for _, net := range networks {
		if net.Name == name {
			return &net, nil
		}
	}
	return nil, fmt.Errorf("network %s not found", name)
}

//GetNetworks from docker
func GetNetworks() ([]Network, error) {
	nets := make([]Network, 0)
	cli, err := client.NewEnvClient()
	if err != nil {
		return nil, err
	}
	networks, err := cli.NetworkList(context.Background(), types.NetworkListOptions{})
	if err != nil {
		return nil, err
	}
	for _, net := range networks {
		nets = append(nets, Network{
			Name: net.Name,
			ID:   net.ID,
		})
	}
	return nets, nil
}

//StartContainer starts container by id
func StartContainer(id string) error {
	cli, err := client.NewEnvClient()
	if err != nil {
		return err
	}
	return cli.ContainerStart(context.Background(), id, types.ContainerStartOptions{})

}
