package patcher

import (
	"IconUpdater/pkg/config"
	"IconUpdater/pkg/file"
	"fmt"
	"path"
)

type (
	patch struct {
		App      App
		Icon     Icon
		Priority int
	}
	App struct {
		Provider config.AppProvider
		config.App
	}
	Icon struct {
		Provider config.IconProvider
		config.Icon
	}
)

func Apply(components *config.Components) error {
	apps, icons := buildAppsAndIcons(components)

	patches := map[string]patch{}
	refresh := false
	for _, patcher := range components.Patchers {
		builtPatches, err := buildPatches(patcher, apps, icons)
		if err != nil {
			return err
		}
		for _, currentPatch := range builtPatches {
			patchToCompare, ok := patches[currentPatch.App.Path]
			if !ok || patchToCompare.Priority < currentPatch.Priority {
				patches[currentPatch.App.Path] = currentPatch
			}
		}
		if patcher.Specifications.RefreshDock {
			refresh = true
		}
	}

	for _, currentPatch := range patches {
		err := currentPatch.Apply()
		if err != nil {
			return err
		}
	}

	if refresh {
		err := RefreshDock()
		if err != nil {
			return err
		}
	}

	return nil
}

func (patch *patch) Apply() error {
	appPaths, err := file.Find(patch.App.Path, patch.App.AppPattern)
	if err != nil {
		return err
	}
	for _, appPath := range appPaths {
		icns, err := file.Find(path.Join(appPath, "Contents", "Resources"), patch.App.IcnsPattern)
		if err != nil {
			return err
		}
		for _, icn := range icns {

			wasBackupCreated, backupPath, err := createBackupFile(icn)
			if err != nil {
				return err
			}

			_, err = file.Copy(patch.Icon.Path, icn)
			if err != nil {
				return err
			}

			fmt.Println(patch.Icon.Path, "->", icn)
			if wasBackupCreated {
				fmt.Println("\tCreated backup:", backupPath)
			}
		}
		err = file.Touch(appPath)
		if err != nil {
			return err
		}
	}
	return nil
}

// createBackupFile creates a copy with "_backup" as extension, if the backup file does not exist. Returns true if a backup was created
func createBackupFile(icn string) (bool, string, error) {
	backupPath := icn + "_backup"
	exists, err := file.Exists(backupPath)
	if err != nil || exists {
		return false, backupPath, err
	}
	_, err = file.Copy(icn, backupPath)
	if err != nil {
		return false, backupPath, err
	}
	return true, backupPath, nil
}
