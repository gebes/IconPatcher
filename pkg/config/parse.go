package config

import (
	"bytes"
	"fmt"
	"github.com/Gebes/IconUpdater/pkg/file"
	"gopkg.in/yaml.v3"
	"io"
)

type Components struct {
	VariableProviders []VariableProvider
	AppProviders      []AppProvider
	IconProviders     []IconProvider
	Patchers          []Patcher
}

func Load(path string, components *Components) error {
	data, err := file.Read(path)
	if err != nil {
		return err
	}

	if components == nil {
		*components = Components{}
	}

	dec := yaml.NewDecoder(bytes.NewReader(data))
	for {
		doc, err := nextDocument(dec)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		var config BaseComponent
		err = yaml.Unmarshal(doc, &config)
		if err != nil {
			return err
		}

		if config.ApiVersion != ApiVersion1 {
			return fmt.Errorf("unknown version: %s", config.ApiVersion)
		}
		switch config.Kind {
		case VariableProviderKind:
			err = unmarshalAndAppend(doc, &components.VariableProviders)
			break
		case AppProviderKind:
			err = unmarshalAndAppend(doc, &components.AppProviders)
			break
		case IconProviderKind:
			err = unmarshalAndAppend(doc, &components.IconProviders)
			break
		case PatcherKind:
			err = unmarshalAndAppend(doc, &components.Patchers)
			break
		default:
			return fmt.Errorf("unknown kind with version %s: %s", config.Kind, config.ApiVersion)

		}
	}
	return err
}

func unmarshalAndAppend[T components](data []byte, toAppend *[]T) error {
	var t T
	err := yaml.Unmarshal(data, &t)
	if err != nil {
		return err
	}
	*toAppend = append(*toAppend, t)
	return nil
}

func nextDocument(dec *yaml.Decoder) ([]byte, error) {
	var doc any
	err := dec.Decode(&doc)
	if err != nil {
		return nil, err
	}
	return yaml.Marshal(doc)
}
