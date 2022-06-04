package patcher

import (
	"fmt"
	"github.com/Gebes/IconUpdater/pkg/config"
	"github.com/Gebes/IconUpdater/pkg/file"
	"path"
)

type (
	patch struct {
		App      App
		Icon     Icon
		Priority int
		DryRun   bool
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
	apps, icons, err := buildAppsAndIcons(components)
	if err != nil {
		return err
	}

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
				patches[currentPatch.App.Name] = currentPatch
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

	if len(patch.App.IcnsPath) != 0 {
		icn := patch.App.IcnsPath
		wasBackupCreated, backupPath, err := createBackupFile(icn, patch.DryRun)
		if err != nil {
			return err
		}

		if !patch.DryRun {
			_, err = file.Copy(patch.Icon.Path, icn)
			if err != nil {
				return err
			}
			fmt.Println(patch.Icon.Path, "->", icn)
		} else {
			fmt.Println(patch.Icon.Path, "->", icn, "(DryRun)")
		}

		if wasBackupCreated {
			fmt.Println("\tCreated backup:", backupPath)
		}

		appPath := path.Dir(path.Dir(path.Dir(icn)))
		err = file.Touch(appPath)
		if err != nil {
			return err
		}
		return nil
	}

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

			wasBackupCreated, backupPath, err := createBackupFile(icn, patch.DryRun)
			if err != nil {
				return err
			}

			if !patch.DryRun {
				_, err = file.Copy(patch.Icon.Path, icn)
				if err != nil {
					return err
				}
				fmt.Println(patch.Icon.Path, "->", icn)

			} else {
				fmt.Println(patch.Icon.Path, "->", icn, "(DryRun)")
			}

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
func createBackupFile(icn string, dryRun bool) (bool, string, error) {
	backupPath := icn + "_backup"
	exists, err := file.Exists(backupPath)
	if err != nil || exists {
		return false, backupPath, err
	}

	if !dryRun {
		_, err = file.Copy(icn, backupPath)
		if err != nil {
			return false, backupPath, err
		}
	}
	return true, backupPath, nil
}
