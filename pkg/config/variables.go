package config

import (
	"strings"
)

func (c *Components) ApplyVariables() {
	for _, variableProvider := range c.VariableProviders {
		for i := range c.IconProviders {
			c.IconProviders[i].ApplyVariables(variableProvider)
		}
		for i := range c.AppProviders {
			c.AppProviders[i].ApplyVariables(variableProvider)
		}

	}
}

func (p *IconProvider) ApplyVariables(provider VariableProvider) {
	for _, v := range provider.Specifications.Variables {
		key, value := "$"+provider.Metadata.Name+"."+v.Name, v.Value
		for i := range p.Specifications.Icons {
			applyVariable(key, value, &p.Specifications.Icons[i].Name)
			applyVariable(key, value, &p.Specifications.Icons[i].Path)
		}
		for i := range p.Specifications.IconFolders {
			applyVariable(key, value, &p.Specifications.IconFolders[i].Path)
			applyVariable(key, value, &p.Specifications.IconFolders[i].IcnsPattern)
		}
	}
}

func (p *AppProvider) ApplyVariables(provider VariableProvider) {
	for _, v := range provider.Specifications.Variables {
		key, value := "$"+provider.Metadata.Name+"."+v.Name, v.Value
		for i := range p.Specifications.Apps {
			applyVariable(key, value, &p.Specifications.Apps[i].Name)
			applyVariable(key, value, &p.Specifications.Apps[i].Path)
			applyVariable(key, value, &p.Specifications.Apps[i].AppPattern)
			applyVariable(key, value, &p.Specifications.Apps[i].IcnsPattern)
			applyVariable(key, value, &p.Specifications.Apps[i].IcnsPath)
		}
	}
}

func (p *Patcher) ApplyVariables(provider VariableProvider) {
	for _, v := range provider.Specifications.Variables {
		key, value := "$"+provider.Metadata.Name+"."+v.Name, v.Value
		for i := range p.Specifications.Matches {
			applyVariable(key, value, &p.Specifications.Matches[i].AppProvider)
			applyVariable(key, value, &p.Specifications.Matches[i].IconProvider)
		}
		for i := range p.Specifications.Patches {
			applyVariable(key, value, &p.Specifications.Patches[i].App)
			applyVariable(key, value, &p.Specifications.Patches[i].Icon)
		}
	}
}

func applyVariable(key, value string, destination *string) {
	*destination = strings.ReplaceAll(*destination, key, value)
}
