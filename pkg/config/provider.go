package config

// Provider works as Facade for underlying config with capability to override some values
type Provider struct {
	config            *Config
	namespaceOverride string
	endpointsOverride []Endpoint
}

// NewProvider creates new Provider around config instance
func NewProvider(config *Config) *Provider {
	return &Provider{
		config: config,
	}
}

// OverrideNamespace set namespace to be overrided with given value
func (c *Provider) OverrideNamespace(namespace string) {
	c.namespaceOverride = namespace
}

// GetNamespace return current namespace
func (c *Provider) GetNamespace() string {
	if c.namespaceOverride != "" {
		return c.namespaceOverride
	}
	return c.config.Namespace
}

// GetEndpointByName finds Endpoint by name or return ok=false
func (c *Provider) GetEndpointByName(name string) (endpoint Endpoint, ok bool) {
	for _, endpoint := range c.GetEndpoints() {
		if endpoint.Name == name {
			return endpoint, true
		}
	}
	return endpoint, false
}

// OverrideEndpoints set endpoint to be overrided with given value
func (c *Provider) OverrideEndpoints(endpoints []Endpoint) {
	c.endpointsOverride = endpoints
}

// GetEndpoints return current namespace
func (c *Provider) GetEndpoints() []Endpoint {
	if len(c.endpointsOverride) != 0 {
		return c.endpointsOverride
	}
	return c.config.Endpoints
}
