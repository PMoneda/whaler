package whaler

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

//Network is a basic representation of a docker network
type Network struct {
	Name string
	ID   string
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
