package main

import (
	"fmt"

	"github.com/PMoneda/whaler"
)

func main() {
	conf := whaler.BuildImageConfig{
		Tag: "company/app:0.0.1",
	}
	str, err := whaler.BuildImageWithDockerfile(conf)

	if err != nil {
		fmt.Println(err)
		return
	} else {
		fmt.Println(str)
	}
	id, err := whaler.CreateContainer(whaler.CreateContainerConfig{
		Image: "company/app:0.0.1",
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
