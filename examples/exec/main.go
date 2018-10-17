package main

import (
	"fmt"
	"io/ioutil"

	"github.com/PMoneda/whaler"
)

func main() {
	out, err := whaler.RunCommand("<container_id>", "ls")
	buf, _ := ioutil.ReadAll(out)
	fmt.Println(string(buf))
	fmt.Println(err)
}
