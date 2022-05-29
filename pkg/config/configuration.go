package config

type preConfiguration struct {
	RefreshDock bool              `yaml:"refreshDock"`
	Variables   map[string]string `yaml:"variables"`
	Icons       any               `yaml:"icons"`
	Patches     any               `yaml:"patches"`
}

type Configuration struct {
	RefreshDock bool              `yaml:"refreshDock"`
	Icons       map[string]string `yaml:"icons"`
	Patches     []Patch           `yaml:"patches"`
}

type Patch struct {
	Icon string
	App  App `yaml:"app"`
}

type App struct {
	Folder      string `yaml:"folder"`
	AppPattern  string `yaml:"appPattern"`
	IcnsPattern string `yaml:"icnsPattern"`
}
