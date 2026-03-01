package internal

import (
	"fmt"
	"os/exec"
	"strings"
)

func Must(err error) {
	if err != nil {
		panic(fmt.Sprintf("Error: %v", err))
	}
}

func Run(name string, args ...string) {
	if err := exec.Command(name, args...).Run(); err != nil {
		panic(fmt.Sprintf("Command failed: %s %s\n%v", name, strings.Join(args, " "), err))
	}
}

func RunOutput(name string, args ...string) []byte {
	out, err := exec.Command(name, args...).CombinedOutput()
	if err != nil {
		panic(fmt.Sprintf("Command failed: %s %s\n%v\nOutput: %s", name, strings.Join(args, " "), err, string(out)))
	}
	return out
}
