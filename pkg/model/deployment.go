package model

// Deployment model
type Deployment struct {
	Metadata
	ID   string         `json:"id"`
	Spec DeploymentSpec `json:"spec"`
}

// DeploymentSpec model
type DeploymentSpec struct {
	Selector map[string]string `json:"selector"`
	Template Pod               `json:"template"`
}
