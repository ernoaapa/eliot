package affinity

import (
	device "github.com/ernoaapa/eliot/pkg/api/services/device/v1"
	pods "github.com/ernoaapa/eliot/pkg/api/services/pods/v1"
	log "github.com/sirupsen/logrus"
	"strconv"
)

func NodeMatch(device *device.Info, affinity pods.NodeAffinity) bool {
	for _, selector := range affinity.Selectors {
		for _, matcher := range selector.Match {
			for _, requirement := range matcher.NodeLabel {
				match, err := evaluateRequirement(requirement, labelsToMap(device.Labels))
				if err != nil {
					if IsNotImplemented(err) {
						log.Fatal(err)
					}
				}
				if !match {
					return false
				}
			}
			return true
		}
	}
	return false
}

func labelsToMap(labels []*device.Label) map[string]string {
	result := make(map[string]string, len(labels))
	for _, label := range labels {
		result[label.Key] = label.Value
	}
	return result
}

func evaluateRequirement(requirement *pods.Requirement, labels map[string]string) (bool, error) {
	value, found := labels[requirement.Key]
	switch requirement.Operator {
	case pods.Requirement_IN:
		if !found {
			return false, nil
		}
		return contains(requirement.Values, value), nil
	case pods.Requirement_NOT_IN:
		if !found {
			return false, nil
		}
		return !contains(requirement.Values, value), nil
	case pods.Requirement_EXISTS:
		if len(requirement.Values) > 0 {
			return false, ErrWithMessagef(ErrInvalidValue, "requirement with operator 'EXISTS' should not include any values")
		}
		return found, nil
	case pods.Requirement_DOES_NOT_EXIST:
		if len(requirement.Values) > 0 {
			return false, ErrWithMessagef(ErrInvalidValue, "requirement with operator 'DOES_NOT_EXIST' should not include any values")
		}
		return !found, nil
	case pods.Requirement_GT:
		if !found {
			return false, nil
		}

		a, b, err := getNumbers(requirement, value)
		if err != nil {
			return false, err
		}

		return a > b, nil
	case pods.Requirement_LT:
		if !found {
			return false, nil
		}

		a, b, err := getNumbers(requirement, value)
		if err != nil {
			return false, err
		}

		return a < b, nil
	default:
		return false, ErrWithMessagef(ErrNotImplemented, "evaluateRequirement does not implement case for pods.Requirement operator %s", requirement.Operator)
	}
}

func getNumbers(requirement *pods.Requirement, value string) (int, int, error) {
	a, err := strconv.Atoi(value)
	if err != nil {
		return 0, 0, ErrWithMessagef(ErrInvalidValue, "[%s] value must be integer but it's [%s]", requirement.Key, value)
	}

	if len(requirement.Values) != 1 {
		return 0, 0, ErrWithMessagef(ErrInvalidValue, "requirement with operator 'GT' must have only one integer value but there's %d values", len(requirement.Values))
	}

	b, err := strconv.Atoi(requirement.Values[0])
	if err != nil {
		return 0, 0, ErrWithMessagef(ErrInvalidValue, "Requirement.Key must be integer with GT operator but it's [%s]", requirement.Values[0])
	}
	return a, b, nil
}

func contains(s []string, v string) bool {
	for _, e := range s {
		if e == v {
			return true
		}
	}
	return false
}
