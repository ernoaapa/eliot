package containers

// Defaults set default values to container definitions
func Defaults(containers []*Container) (result []*Container) {
	for _, container := range containers {
		result = append(result, Default(container))
	}
	return result
}

// Default set default values to Container model
func Default(container *Container) *Container {
	// Set here default values if needed...
	return container
}
