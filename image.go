package whaler

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"io/ioutil"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

//Image is a basic representation of a docker image
type Image struct {
	ID   string
	Name string
}

//BuildImageConfig is a basic configuration to build an image
type BuildImageConfig struct {
	PathContext string
	Dockerfile  string
	Tag         string
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

//Publish image to registry
func Publish(image, username, password string) (string, error) {
	auth := types.AuthConfig{
		Username: username,
		Password: password,
	}
	encodedJSON, err := json.Marshal(auth)
	if err != nil {
		return "", err
	}
	encoded := base64.StdEncoding.EncodeToString(encodedJSON)
	cli, err := client.NewEnvClient()
	if err != nil {
		return "", err
	}

	out, err := cli.ImagePush(context.Background(), image, types.ImagePushOptions{
		All:          true,
		RegistryAuth: encoded,
	})
	if err != nil {
		return "", err
	}
	if b, err := ioutil.ReadAll(out); err != nil {
		return "", err
	} else {
		return string(b), nil
	}
}
