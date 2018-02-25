package affinity

import (
	"fmt"
	device "github.com/ernoaapa/eliot/pkg/api/services/device/v1"
	pods "github.com/ernoaapa/eliot/pkg/api/services/pods/v1"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNodeMatchAllRequirementsMatch(t *testing.T) {
	node := &device.Info{
		Labels: []*device.Label{
			{Key: "eliot.io/hostname", Value: "foobar"},
			{Key: "eliot.io/area", Value: "kitchen"},
			{Key: "eliot.io/version", Value: "1.2.3"},
		},
	}
	affinity := pods.NodeAffinity{
		Selectors: []*pods.NodeSelector{
			{Match: []*pods.Matcher{
				{NodeLabel: []*pods.Requirement{
					{Key: "eliot.io/hostname", Operator: pods.Requirement_IN, Values: []string{"foobar"}},
					{Key: "eliot.io/area", Operator: pods.Requirement_IN, Values: []string{"kitchen"}},
				}},
			}},
		},
	}
	assert.True(t, NodeMatch(node, affinity))
}

func TestNodeNotMatchIfNoSelectors(t *testing.T) {
	node := &device.Info{
		Labels: []*device.Label{
			{Key: "eliot.io/hostname", Value: "foobar"},
		},
	}
	affinity := pods.NodeAffinity{
		Selectors: []*pods.NodeSelector{},
	}
	assert.False(t, NodeMatch(node, affinity))
}

func TestNodeNotMatchWithLabels(t *testing.T) {
	node := &device.Info{
		Labels: []*device.Label{
			{Key: "eliot.io/hostname", Value: "foobar"},
		},
	}
	affinity := pods.NodeAffinity{
		Selectors: []*pods.NodeSelector{
			{Match: []*pods.Matcher{
				{NodeLabel: []*pods.Requirement{
					{Key: "eliot.io/hostname", Operator: pods.Requirement_IN, Values: []string{"something-else"}},
				}},
			}},
		},
	}
	assert.False(t, NodeMatch(node, affinity))
}

func TestNodeNotMatchIfNoLabelLabels(t *testing.T) {
	node := &device.Info{}
	affinity := pods.NodeAffinity{
		Selectors: []*pods.NodeSelector{
			{Match: []*pods.Matcher{
				{NodeLabel: []*pods.Requirement{
					{Key: "eliot.io/hostname", Operator: pods.Requirement_IN, Values: []string{"something-else"}},
				}},
			}},
		},
	}
	assert.False(t, NodeMatch(node, affinity))
}

func TestEvaluateToTrueWhenIsInRequirement(t *testing.T) {
	requirement := &pods.Requirement{
		Key:      "eliot.io/hostname",
		Operator: pods.Requirement_IN,
		Values:   []string{"foo", "bar", "baz"},
	}
	match, _ := evaluateRequirement(requirement, map[string]string{
		"eliot.io/hostname": "foo",
	})
	assert.True(t, match)
}

func TestEvaluateToFalseWhenIsNotRequirement(t *testing.T) {
	requirement := &pods.Requirement{
		Key:      "eliot.io/hostname",
		Operator: pods.Requirement_IN,
		Values:   []string{"foo", "bar", "baz"},
	}
	match, _ := evaluateRequirement(requirement, map[string]string{
		"eliot.io/hostname": "something-else",
	})
	assert.False(t, match)
}

func TestEvaluateToTrueWhenIsNotInRequirement(t *testing.T) {
	requirement := &pods.Requirement{
		Key:      "eliot.io/hostname",
		Operator: pods.Requirement_NOT_IN,
		Values:   []string{"foo", "bar"},
	}
	match, _ := evaluateRequirement(requirement, map[string]string{
		"eliot.io/hostname": "baz",
	})
	assert.True(t, match)
}

func TestEvaluateToFalseWhenIsInRequirement(t *testing.T) {
	requirement := &pods.Requirement{
		Key:      "eliot.io/hostname",
		Operator: pods.Requirement_NOT_IN,
		Values:   []string{"foo", "bar"},
	}
	match, _ := evaluateRequirement(requirement, map[string]string{
		"eliot.io/hostname": "foo",
	})
	assert.False(t, match)
}

func TestEvaluateToFalseWhenNotInOperatorDontHaveGivenLabel(t *testing.T) {
	requirement := &pods.Requirement{
		Key:      "eliot.io/hostname",
		Operator: pods.Requirement_NOT_IN,
		Values:   []string{"foo", "bar"},
	}
	match, _ := evaluateRequirement(requirement, map[string]string{})

	assert.False(t, match)
}

func TestEvaluateToTrueWhenRequirementExists(t *testing.T) {
	requirement := &pods.Requirement{
		Key:      "eliot.io/hostname",
		Operator: pods.Requirement_EXISTS,
	}
	match, _ := evaluateRequirement(requirement, map[string]string{
		"eliot.io/hostname": "baz",
	})
	assert.True(t, match)
}

func TestEvaluateReturnErrorWhenExistsHaveValues(t *testing.T) {
	requirement := &pods.Requirement{
		Key:      "eliot.io/hostname",
		Operator: pods.Requirement_EXISTS,
		Values:   []string{"EXISTS operator should not include values"},
	}
	_, err := evaluateRequirement(requirement, map[string]string{})
	assert.Error(t, err)
}

func TestEvaluateToFalseWhenRequirementDoesNotExist(t *testing.T) {
	requirement := &pods.Requirement{
		Key:      "eliot.io/hostname",
		Operator: pods.Requirement_DOES_NOT_EXIST,
	}
	match, _ := evaluateRequirement(requirement, map[string]string{
		"eliot.io/version": "1.2.3",
	})
	assert.True(t, match)
}

func TestEvaluateToFalseWhenRequirementExist(t *testing.T) {
	requirement := &pods.Requirement{
		Key:      "eliot.io/hostname",
		Operator: pods.Requirement_DOES_NOT_EXIST,
	}
	match, _ := evaluateRequirement(requirement, map[string]string{
		"eliot.io/hostname": "1.2.3",
	})
	assert.False(t, match)
}

func TestEvaluateReturnErrorWhenDoesNotExistsHaveValues(t *testing.T) {
	requirement := &pods.Requirement{
		Key:      "eliot.io/hostname",
		Operator: pods.Requirement_DOES_NOT_EXIST,
		Values:   []string{"DOES_NOT_EXIST operator should not include values"},
	}
	_, err := evaluateRequirement(requirement, map[string]string{})
	assert.Error(t, err)
}

func TestEvaluateToTrueWhenRequirementIsGreaterThan(t *testing.T) {
	requirement := &pods.Requirement{
		Key:      "eliot.io/count",
		Operator: pods.Requirement_GT,
		Values:   []string{"10"},
	}
	match, _ := evaluateRequirement(requirement, map[string]string{
		"eliot.io/count": "11",
	})
	assert.True(t, match)
}

func TestEvaluateToFalseWhenRequirementIsNotGreaterThan(t *testing.T) {
	requirement := &pods.Requirement{
		Key:      "eliot.io/count",
		Operator: pods.Requirement_GT,
		Values:   []string{"10"},
	}
	match, _ := evaluateRequirement(requirement, map[string]string{
		"eliot.io/count": "9",
	})
	assert.False(t, match)
}

func TestEvaluateToFalseWhenGtOperatorDontHaveGivenLabel(t *testing.T) {
	requirement := &pods.Requirement{
		Key:      "eliot.io/count",
		Operator: pods.Requirement_GT,
		Values:   []string{"10"},
	}
	match, _ := evaluateRequirement(requirement, map[string]string{})

	assert.False(t, match)
}

func TestEvaluateToTrueWhenRequirementIsLessThan(t *testing.T) {
	requirement := &pods.Requirement{
		Key:      "eliot.io/count",
		Operator: pods.Requirement_LT,
		Values:   []string{"10"},
	}
	match, _ := evaluateRequirement(requirement, map[string]string{
		"eliot.io/count": "9",
	})
	assert.True(t, match)
}

func TestEvaluateToFalseWhenRequirementIsNotLessThan(t *testing.T) {
	requirement := &pods.Requirement{
		Key:      "eliot.io/count",
		Operator: pods.Requirement_LT,
		Values:   []string{"10"},
	}
	match, _ := evaluateRequirement(requirement, map[string]string{
		"eliot.io/count": "11",
	})
	assert.False(t, match)
}

func TestEvaluateToFalseWhenLtOperatorDontHaveGivenLabel(t *testing.T) {
	requirement := &pods.Requirement{
		Key:      "eliot.io/count",
		Operator: pods.Requirement_LT,
		Values:   []string{"10"},
	}
	match, _ := evaluateRequirement(requirement, map[string]string{})

	assert.False(t, match)
}

func TestEvaluateReturnErrorIfMultipleValuesGivenToNumberOperators(t *testing.T) {
	for _, operator := range []pods.Requirement_Operator{pods.Requirement_GT, pods.Requirement_LT} {
		var err error
		_, err = evaluateRequirement(&pods.Requirement{
			Key:      "eliot.io/count",
			Operator: operator,
			Values:   []string{"10", "99"},
		}, map[string]string{
			"eliot.io/count": "0",
		})
		assert.Error(t, err, fmt.Sprintf("evaluateRequirement didn't return error with operator %s even though there were multiple values in requirement", operator))

		_, err = evaluateRequirement(&pods.Requirement{
			Key:      "eliot.io/count",
			Operator: operator,
			Values:   []string{"this is a text"},
		}, map[string]string{
			"eliot.io/count": "0",
		})
		assert.Error(t, err, fmt.Sprintf("evaluateRequirement didn't return error with operator %s even though value is not number", operator))
	}
}

func TestEvaluateImplementsAllOperators(t *testing.T) {
	for operator, name := range pods.Requirement_Operator_name {
		requirement := &pods.Requirement{
			Key:      "eliot.io/hostname",
			Operator: pods.Requirement_Operator(operator),
			Values:   []string{"foo", "bar", "baz"},
		}

		_, err := evaluateRequirement(requirement, map[string]string{
			"eliot.io/hostname": "something-else",
		})
		assert.False(t, IsNotImplemented(err), fmt.Sprintf("evaluateRequirement does not implement '%s' Requirement Operator", name))
	}
}
