package main

import (
	"github.com/jdewinne/kotsd/cli"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	cli.InitAndExecute(version)
}
