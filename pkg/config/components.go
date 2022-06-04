package config

type Kind string

const (
	VariableProviderKind Kind = "VariableProvider"
	AppProviderKind      Kind = "AppProvider"
	IconProviderKind     Kind = "IconProvider"
	PatcherKind          Kind = "Patcher"
)

type ApiVersion string

const ApiVersion1 = "v1"

type BaseComponent struct {
	ApiVersion string `yaml:"apiVersion"`
	Kind       Kind   `yaml:"kind"`
}

type components interface {
	VariableProvider | AppProvider | IconProvider | Patcher
}

type Variable struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

type VariableProvider struct {
	BaseComponent `yaml:",inline"`
	Metadata      struct {
		Name string `yaml:"name"`
	} `yaml:"metadata"`
	Specifications struct {
		Variables []Variable `yaml:"variables"`
	} `yaml:"spec"`
}

type (
	AppProvider struct {
		BaseComponent `yaml:",inline"`
		Metadata      struct {
			Name string `yaml:"name"`
		} `yaml:"metadata"`
		Specifications struct {
			Apps []App `yaml:"apps"`
		} `yaml:"spec"`
	}
	App struct {
		Name        string `yaml:"name"`
		Path        string `yaml:"path"`
		AppPattern  string `yaml:"appPattern"`
		IcnsPattern string `yaml:"icnsPattern"`
		IcnsPath    string `yaml:"icnsPath"`
	}
)

type (
	IconProvider struct {
		BaseComponent `yaml:",inline"`
		Metadata      struct {
			Name string `yaml:"name"`
		} `yaml:"metadata"`
		Specifications struct {
			Icons       []Icon       `yaml:"icons"`
			IconFolders []IconFolder `yaml:"folders"`
		} `yaml:"spec"`
	}

	Icon struct {
		Name string `yaml:"name"`
		Path string `yaml:"path"`
	}
	IconFolder struct {
		Path        string `yaml:"path"`
		IcnsPattern string `yaml:"icnsPattern"`
	}
)

type (
	Patcher struct {
		BaseComponent `yaml:",inline"`
		Metadata      struct {
			Name       string   `yaml:"name"`
			Components []string `yaml:"components"`
		} `yaml:"metadata"`
		Specifications struct {
			RefreshDock bool    `yaml:"refreshDock"`
			DryRun      bool    `yaml:"dryRun"`
			Matches     []Match `yaml:"matches"`
			Patches     []Patch `yaml:"patches"`
		} `yaml:"spec"`
	}
	Match struct {
		AppProvider  string `yaml:"appProvider"`
		IconProvider string `yaml:"iconProvider"`
		Priority     int    `yaml:"priority"`
	}
	Patch struct {
		Icon     string `yaml:"icon"`
		App      string `yaml:"app"`
		Priority int    `yaml:"priority"`
	}
)
