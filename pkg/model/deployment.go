package model

// Deployment model
type Deployment struct {
	ID       string         `json:"id"`
	Metadata Metadata       `json:"metadata"`
	Spec     DeploymentSpec `json:"spec"`
}

// DeploymentSpec model
type DeploymentSpec struct {
	Selector map[string]string `json:"selector"`
	Template Pod               `json:"template"`
}
