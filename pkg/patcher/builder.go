package patcher

import (
	"fmt"
	"github.com/Gebes/IconUpdater/pkg/config"
	"github.com/Gebes/IconUpdater/pkg/file"
	"path"
	"strings"
)

func buildPatches(patcher config.Patcher, apps []App, icons []Icon) ([]patch, error) {
	patches := make([]patch, 0)

	for _, match := range patcher.Specifications.Matches {
		appNames := map[string]App{}
		for _, currentApp := range apps {
			if currentApp.Provider.Metadata.Name == match.AppProvider {
				appNames[currentApp.Name] = currentApp
			}
		}

		for _, currentIcon := range icons {
			name := currentIcon.Icon.Name
			currentApp, ok := appNames[name]
			if currentIcon.Provider.Metadata.Name != match.IconProvider || !ok {
				continue
			}
			patches = append(patches, patch{
				App:      currentApp,
				Icon:     currentIcon,
				Priority: match.Priority,
				DryRun:   patcher.Specifications.DryRun,
			})
		}
	}

	for _, currentPatch := range patcher.Specifications.Patches {
		appSplit := strings.SplitN(currentPatch.App, ".", 2)
		if len(appSplit) != 2 {
			return nil, fmt.Errorf("invalid App for patch: %s", currentPatch.App)
		}
		var app *App
		appProvider, appName := appSplit[0], appSplit[1]
		for i, a := range apps {
			if a.Provider.Metadata.Name == appProvider && a.Name == appName {
				app = &apps[i]
				break
			}
		}
		if app == nil {
			return nil, fmt.Errorf("app not found for patch: %s", currentPatch.App)
		}

		iconSplit := strings.SplitN(currentPatch.Icon, ".", 2)
		if len(iconSplit) != 2 {
			return nil, fmt.Errorf("invalid icon for path: %s", currentPatch.Icon)
		}

		var icon *Icon
		iconProvider, iconName := iconSplit[0], iconSplit[1]
		for _, i := range icons {
			if i.Provider.Metadata.Name == iconProvider && i.Name == iconName {
				icon = &i
				break
			}
		}
		if icon == nil {
			return nil, fmt.Errorf("icon not found for patch: %s", currentPatch.Icon)
		}

		patches = append(patches, patch{
			App:      *app,
			Icon:     *icon,
			Priority: currentPatch.Priority,
			DryRun:   patcher.Specifications.DryRun,
		})
	}

	return patches, nil
}

func buildAppsAndIcons(configurations *config.Components) ([]App, []Icon, error) {
	apps := make([]App, 0)
	for _, provider := range configurations.AppProviders {
		for _, currentApp := range provider.Specifications.Apps {
			apps = append(apps, App{
				Provider: provider,
				App:      currentApp,
			})
		}
	}
	icons := make([]Icon, 0)
	for _, provider := range configurations.IconProviders {
		for _, currentIcon := range provider.Specifications.Icons {
			icons = append(icons, Icon{
				Provider: provider,
				Icon:     currentIcon,
			})
		}
		for _, folder := range provider.Specifications.IconFolders {
			paths, err := file.Find(folder.Path, folder.IcnsPattern)
			if err != nil {
				return nil, nil, err
			}

			for _, icnsPath := range paths {
				fileWithExt := path.Base(icnsPath)
				file := strings.TrimSuffix(fileWithExt, path.Ext(fileWithExt))
				icons = append(icons, Icon{
					Provider: provider,
					Icon: config.Icon{
						Name: file,
						Path: icnsPath,
					},
				})
			}
		}
	}

	return apps, icons, nil
}
