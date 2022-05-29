package main

import (
	"IconUpdater/pkg/config"
	"IconUpdater/pkg/patcher"
	"fmt"
	"os"
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Println("Provide at least one configuration file")
		return
	}

	var configs []config.Configuration
	for _, path := range args {
		currentConfig, err := config.Load(path)
		if err != nil {
			fmt.Printf("Could not load configuration file \"%s\": %s\n", path, err)
			return
		}
		configs = append(configs, *currentConfig)
	}

	refreshDock := false
	for i, configuration := range configs {
		if configuration.RefreshDock {
			refreshDock = true
		}
		err := patcher.Patch(configuration)
		if err != nil {
			fmt.Printf("Could not apply configuration %d: %s\n", i+1, err)
			return
		}
	}

	if refreshDock {
		err := patcher.RefreshDock()
		if err != nil {
			fmt.Printf("Could not refresh dock: %s\n", err)
			return
		}
	}

}
