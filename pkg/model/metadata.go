package model

var (
	nameKey      = "name"
	namespaceKey = "namespace"
)

// Metadata is any random metadata
type Metadata map[string]interface{}

// NewMetadata creates new metadata with name and metadata fields
func NewMetadata(name, namespace string) Metadata {
	return Metadata{
		nameKey:      name,
		namespaceKey: namespace,
	}
}

// GetName returns name field from metadata
func (m Metadata) GetName() string {
	return m.GetMetadataValue(nameKey)
}

// SetName returns name field from metadata
func (m Metadata) SetName(name string) {
	m.SetMetadataValue(nameKey, name)
}

// GetNamespace returns name field from metadata
func (m Metadata) GetNamespace() string {
	return m.GetMetadataValue(namespaceKey)
}

// SetNamespace returns name field from metadata
func (m Metadata) SetNamespace(namespace string) {
	m.SetMetadataValue(namespaceKey, namespace)
}

// HasMetadataValue return true if metadata contains value
func (m Metadata) HasMetadataValue(key string) bool {
	_, ok := m[key]
	return ok
}

// GetMetadataValue return metadata value or empty string with key
func (m Metadata) GetMetadataValue(key string) string {
	value := m[key]
	if value == nil {
		return ""
	}
	return m[key].(string)
}

// SetMetadataValue updates metadata value by key
func (m Metadata) SetMetadataValue(key, value string) {
	m[key] = value
}
