package main

import (
	"github.com/maxiiot/devicebridge/cmd/devicebridge/cmd"
)

var version = "0.1.1"

func main() {
	cmd.Execute(version)
}
