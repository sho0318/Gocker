package main

import (
	"fmt"
	"os"

	"mydocker/internal"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	if len(os.Args) < 2 {
		return fmt.Errorf("usage: %s <run|child> <command> [args...]", os.Args[0])
	}

	command := os.Args[1]
	
	switch command {
	case "run":
		if len(os.Args) < 3 {
			return fmt.Errorf("usage: %s run <command> [args...]", os.Args[0])
		}
		return runContainer(os.Args[2:])
		
	case "child":
		if len(os.Args) < 3 {
			return fmt.Errorf("usage: %s child <command> [args...]", os.Args[0])
		}
		return runChild(os.Args[2:])
		
	default:
		return fmt.Errorf("unknown command: %s\nAvailable commands: run, child", command)
	}
}

func runContainer(command []string) error {
	cfg := internal.NewDefaultConfig(command)
	container := internal.NewContainer(cfg)
	return container.Run()
}

func runChild(command []string) error {
	cfg := internal.NewDefaultConfig(command)
	container := internal.NewContainer(cfg)
	return container.RunChild()
}
