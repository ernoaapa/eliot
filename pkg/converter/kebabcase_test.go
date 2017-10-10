package converter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKebabCaseToCamelCase(t *testing.T) {
	assert.Equal(t, "FooBarBaz", KebabCaseToCamelCase("foo-bar-baz"), "Should convert kebab-case to CamelCase")
	assert.Equal(t, "Foo", KebabCaseToCamelCase("foo"), "Should support also cases when no hypen")
	assert.Equal(t, "FooBarBaz", KebabCaseToCamelCase("FooBarBaz"), "If already CamelCase, should not do anything")
}
