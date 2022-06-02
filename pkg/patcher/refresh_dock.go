package patcher

import (
	"IconUpdater/pkg/file"
	"os/exec"
	"path/filepath"
)

// RefreshDock clears the icon cache and kills the Dock process
func RefreshDock() error {
	err := clearIconCache()
	if err != nil {
		return err
	}

	err = exec.Command("killall", "Dock", "Finder").Run()
	if err != nil {
		return err
	}

	return nil
}

// clearIconCache finds the dock iconcache and deletes the file. This enforces reloading the icons the next time the dock refreshes.
func clearIconCache() error {
	paths, err := filepath.Glob("/var/folders/*/*/*/com.apple.dock.iconcache*")
	if err != nil {
		return err
	}
	for _, path := range paths {
		err := file.Remove(path)
		if err != nil {
			return err
		}
	}
	return nil
}
