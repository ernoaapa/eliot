package controller

import (
	"github.com/ernoaapa/can/pkg/model"
	"github.com/ernoaapa/can/pkg/runtime"
)

func getCurrentState(client runtime.Client) (result []model.Pod, err error) {
	result = []model.Pod{}
	namespaces, err := client.GetNamespaces()
	if err != nil {
		return result, err
	}

	for _, namespace := range namespaces {
		pods, err := client.GetPods(namespace)
		if err != nil {
			return result, err
		}

		result = append(result, pods...)
	}
	return result, nil
}
