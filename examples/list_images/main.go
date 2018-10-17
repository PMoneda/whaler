package main

import (
	"fmt"

	"github.com/PMoneda/whaler"
)

func main() {

	list, err := whaler.ListImages()
	if err != nil {
		panic(err)
	}
	for _, img := range list {
		fmt.Println(img)
	}
}
