package config

import (
	"fmt"
	"strings"
)

func (c *Components) Validate() error {

	// name:kind
	components := map[string]string{}
	for _, component := range c.VariableProviders {
		name := component.Metadata.Name
		if len(strings.TrimSpace(name)) == 0 {
			return fmt.Errorf("name of VariableProvider is empty")
		}
		kind, ok := components[name]
		if ok {
			return fmt.Errorf("name of VariableProvider is already used for %s component: %s", kind, name)
		}
		components[name] = string(component.Kind)
	}
	for _, component := range c.AppProviders {
		name := component.Metadata.Name
		if len(strings.TrimSpace(name)) == 0 {
			return fmt.Errorf("name of AppProvider is empty")
		}
		kind, ok := components[name]
		if ok {
			return fmt.Errorf("name of AppProvider is already used for %s component: %s", kind, name)
		}
		components[name] = string(component.Kind)
	}
	for _, component := range c.IconProviders {
		name := component.Metadata.Name
		if len(strings.TrimSpace(name)) == 0 {
			return fmt.Errorf("name of IconProvider is empty")
		}
		kind, ok := components[name]
		if ok {
			return fmt.Errorf("name of IconProvider is already used for %s component: %s", kind, name)
		}
		components[name] = string(component.Kind)
	}
	for _, component := range c.Patchers {
		name := component.Metadata.Name
		if len(strings.TrimSpace(name)) == 0 {
			return fmt.Errorf("name of Patcher is empty")
		}
		kind, ok := components[name]
		if ok {
			return fmt.Errorf("name of Patcher is already used for %s component: %s", kind, name)
		}
		components[name] = string(component.Kind)
	}

	return nil

}
