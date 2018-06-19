package main

import (
	"fmt"

	"github.com/PMoneda/whaler"
)

func main() {
	conf := whaler.BuildImageConfig{
		Tag: "localhost:5000/company/app:0.0.1",
	}
	str, err := whaler.BuildImageWithDockerfile(conf)

	if err != nil {
		fmt.Println(err)
		return
	} else {
		fmt.Println(str)
	}
	if out, err := whaler.Publish(conf.Tag, "docker", "docker"); err != nil {
		fmt.Println(out)
		fmt.Println(err)
	} else {
		fmt.Println(out)
	}
	id, err := whaler.CreateContainer(whaler.CreateContainerConfig{
		Image: "localhost:5000/company/app:0.0.1",
		Name:  "my-container",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(id)
	err = whaler.StartContainer(id)
	if err != nil {
		fmt.Println(err)
	}
}
