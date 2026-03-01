package main

import (
	"fmt"
	"os"

	"gocker/internal"
	"gocker/internal/container"
	"gocker/internal/host"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", r)
			os.Exit(1)
		}
	}()

	command := os.Args[1]
	
	switch command {
	case "run":
		runContainer(os.Args[2:])
	case "child":
		id := os.Args[2]
		runChild(id, os.Args[3:])
	}
}

func runContainer(command []string) {
	cfg := internal.NewDefaultConfig(command)
	launcher := host.NewLauncher(cfg)
	launcher.Start()
}

func runChild(id string, command []string) {
	cfg := internal.NewDefaultConfig(command)
	cfg.ContainerID = id
	runtime := container.NewRuntime(cfg)
	runtime.Start()
}
