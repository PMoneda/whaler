package whaler

import (
	"context"
	"fmt"
	"strings"
	"time"

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

//CreateContainerConfig is a basic configuration structure to create a docker container
type CreateContainerConfig struct {
	Name        string
	Image       string
	Volumes     []string
	Env         []string
	Ports       []string
	NetworkName string
}

//GetContainers from docker
func GetContainers(all bool) ([]Container, error) {
	cli, err := client.NewEnvClient()
	if err != nil {
		return nil, err
	}
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{
		All: all,
	})
	if err != nil {
		return nil, err
	}
	_containers := make([]Container, len(containers))
	i := 0
	for _, container := range containers {
		_containers[i] = Container{ID: container.ID, Name: container.Names[0], Image: container.Image}
		i++
	}
	return _containers, nil
}

//FindContainerByIdentifier finds container by id or name
func FindContainerByIdentifier(identifier string, all bool) (Container, error) {
	containers, err := GetContainers(all)
	if err != nil {
		return Container{}, err
	}
	for _, c := range containers {
		if c.Name == identifier || c.ID == identifier {
			return c, nil
		}
	}
	return Container{}, fmt.Errorf("container with identifer(name or id) %s not found", identifier)
}

//RemoveContainer by Id the container links and volumes will be removed too
func RemoveContainer(id string, force bool) error {
	cli, err := client.NewEnvClient()
	if err != nil {
		return err
	}
	return cli.ContainerRemove(context.Background(), id, types.ContainerRemoveOptions{
		Force: force,
	})
}

//RestartContainer by Id or Name
func RestartContainer(identifier string, timeout *time.Duration) error {
	cli, err := client.NewEnvClient()
	if err != nil {
		return err
	}
	if contaienr, err := FindContainerByIdentifier(identifier, false); err != nil {
		return err
	} else {
		return cli.ContainerRestart(context.Background(), contaienr.ID, timeout)
	}
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

//StartContainer starts container by id
func StartContainer(id string) error {
	cli, err := client.NewEnvClient()
	if err != nil {
		return err
	}
	return cli.ContainerStart(context.Background(), id, types.ContainerStartOptions{})

}
