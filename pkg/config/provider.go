package config

import "strings"

// Provider works as Facade for underlying config with capability to override some values
type Provider struct {
	config            *Config
	namespaceOverride string
	endpointOverride  string
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
	return c.config.GetCurrentContext().Namespace
}

// OverrideEndpoint set endpoint to be overrided with given value
func (c *Provider) OverrideEndpoint(endpoint string) {
	c.endpointOverride = endpoint
}

// GetEndpoint return current namespace
func (c *Provider) GetEndpoint() string {
	if c.endpointOverride != "" {
		return c.endpointOverride
	}
	return c.config.GetCurrentEndpoint().URL
}

// GetEndpointHost return only the host (or ip) of endpoint
func (c *Provider) GetEndpointHost() string {
	endpoint := c.GetEndpoint()
	parts := strings.SplitN(endpoint, ":", 2)
	return parts[0]
}
