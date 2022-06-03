package main

import (
	"fmt"
	"github.com/Gebes/IconUpdater/pkg/config"
	"github.com/Gebes/IconUpdater/pkg/patcher"
	"os"
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Println("Provide at least one configuration file")
		return
	}

	components := &config.Components{}
	for _, path := range args {
		err := config.Load(path, components)
		if err != nil {
			fmt.Printf("Could not load configuration file \"%s\": %s\n", path, err)
			return
		}
	}

	err := components.Validate()
	if err != nil {
		fmt.Printf("Validation failed: %v\n", err)
		return
	}

	components.ApplyVariables()

	err = patcher.Apply(components)
	if err != nil {
		fmt.Printf("Patching failed: %v\n", err)
		return
	}

}
