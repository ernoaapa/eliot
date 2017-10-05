package main

import (
	"math/rand"
	"time"

	"github.com/ernoaapa/can/pkg/progress"
)

// Simple app to test out how the progress rendering looks like
func main() {
	steps := int64(100)
	total := int64(876433)
	status := &progress.Status{
		Ref:    "layer-sha256:f83e14495c19e0bb1c7187b76571e3c6a2125dae2926678da165cfc6c7da0670",
		Digest: "f83e14495c19e0bb1c7187b76571e3c6a2125dae2926678da165cfc6c7da0670",
		Status: "downloading",
		Offset: 0,
		Total:  total,
	}
	progresses := []*progress.ImageFetch{
		progress.CreateImageFetch(
			"my-pod",
			"the-image",
			true,
			map[string]*progress.Status{
				"layer-sha256:f83e14495c19e0bb1c7187b76571e3c6a2125dae2926678da165cfc6c7da0670": status,
			},
		),
	}
	renderer := progress.NewRenderer()

	for i := int64(0); i < steps; i++ {
		status.Offset = (total / steps) * i
		time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
		renderer.Update(progresses)
	}

	renderer.Stop()
}
