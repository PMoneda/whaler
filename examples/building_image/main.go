package main

import (
	"fmt"

	"github.com/PMoneda/whaler"
)

func main() {
	conf := whaler.BuildImageConfig{
		Tag: "my-image:0.0.1",
	}
	str, err := whaler.BuildImageWithDockerfile(conf)
	if err != nil {
		panic(err)
	} else {
		fmt.Println(str)
	}
}
