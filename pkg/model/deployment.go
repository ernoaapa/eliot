package model

// Deployment model
type Deployment struct {
	Metadata
	ID   string         `validate:"required,gt=0"`
	Spec DeploymentSpec `validate:"required"`
}

// DeploymentSpec model
type DeploymentSpec struct {
	Selector map[string]string
	Template Pod `validate:"required"`
}
