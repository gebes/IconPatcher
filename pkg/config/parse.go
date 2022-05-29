package config

import (
	"IconUpdater/pkg/file"
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"strings"
)

func Load(path string) (*Configuration, error) {
	data, err := file.Read(path)
	if err != nil {
		return nil, err
	}
	var preConfiguration preConfiguration
	err = yaml.Unmarshal(data, &preConfiguration)
	if err != nil {
		return nil, err
	}

	icons := map[string]string{}
	switch t := preConfiguration.Icons.(type) {
	case string:
		iconsData, err := file.Read(t)
		if err != nil {
			return nil, err
		}
		err = yaml.Unmarshal(iconsData, &icons)
		if err != nil {
			return nil, err
		}
	case map[string]any:
		for s, a := range t {
			icons[s] = fmt.Sprint(a)
		}
	default:
		return nil, errors.New("icons needs to be a string or map")
	}

	var patches []Patch
	switch t := preConfiguration.Patches.(type) {
	case string:
		patchesData, err := file.Read(t)
		if err != nil {
			return nil, err
		}
		err = yaml.Unmarshal(patchesData, &patches)
		if err != nil {
			return nil, err
		}
	case []any:
		patches = make([]Patch, len(t))
		for i, possiblePatch := range t {
			data, err = yaml.Marshal(possiblePatch)
			if err != nil {
				return nil, err
			}
			err = yaml.Unmarshal(data, &patches[i])
			if err != nil {
				return nil, fmt.Errorf("patch %d is not a valid patch: %w", i+1, err)
			}
		}
	default:
		return nil, errors.New("patches needs to be a string or list of patches")
	}

	configuration := &Configuration{
		RefreshDock: preConfiguration.RefreshDock,
		Icons:       icons,
		Patches:     patches,
	}

	for k, v := range configuration.Icons {
		configuration.Icons[k] = replace(v, preConfiguration.Variables)
	}
	for i, patch := range configuration.Patches {
		configuration.Patches[i].App.Folder = replace(patch.App.Folder, preConfiguration.Variables)
		configuration.Patches[i].App.AppPattern = replace(patch.App.AppPattern, preConfiguration.Variables)
		configuration.Patches[i].App.IcnsPattern = replace(patch.App.IcnsPattern, preConfiguration.Variables)
	}

	return configuration, nil
}

func replace(s string, variables map[string]string) string {
	for k, v := range variables {
		s = strings.ReplaceAll(s, "$"+k, v)
	}
	return s
}
