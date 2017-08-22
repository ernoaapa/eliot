package model

type Pod struct {
	Name string `yaml:"name"`
	Spec Spec   `yaml:"spec"`
}

type Spec struct {
	Containers []Container `yaml:"containers"`
}

type Container struct {
	Name  string `yaml:"name"`
	Image string `yaml:"image"`
}

type NodeInfo struct {
}
