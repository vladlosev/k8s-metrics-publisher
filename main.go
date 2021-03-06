package main

import (
	"os"

	"github.com/vladlosev/k8s-metrics-publisher/pkg/cmd"
)

func main() {
	cmd := cmd.NewServerCommand()
	err := cmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
