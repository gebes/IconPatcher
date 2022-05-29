package patcher

import (
	"IconUpdater/pkg/config"
	"IconUpdater/pkg/file"
	"fmt"
	"path"
)

func Patch(configuration config.Configuration) error {
	for _, patch := range configuration.Patches {
		files, err := file.Find(patch.App.Folder, patch.App.AppPattern)
		if err != nil {
			return err
		}
		err = patchIcon(files, configuration.Icons, patch)
		if err != nil {
			return err
		}
	}

	return nil
}

func patchIcon(appPaths []string, icons map[string]string, patch config.Patch) error {
	for _, appPath := range appPaths {
		icns, err := file.Find(path.Join(appPath, "Contents", "Resources"), patch.App.IcnsPattern)
		if err != nil {
			return err
		}
		for _, icn := range icns {
			_, err = file.Copy(icons[patch.Icon], icn)
			if err != nil {
				return err
			}
			fmt.Println("Patched patch", icn)
		}
		err = file.Touch(appPath)
		if err != nil {
			return err
		}
	}
	return nil
}

func RefreshDock() error {
	fmt.Println("\nRun this command to refresh your Dock. You need to close the udpated Apps too.\n\trm /var/folders/*/*/*/com.apple.dock.iconcache* && killall Finder && killall Dock")
	return nil
}
