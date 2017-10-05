package progress

import (
	"github.com/vbauerster/mpb"
	"github.com/vbauerster/mpb/decor"
)

// Renderer take care of rendering progress bars in the terminal
type Renderer struct {
	progress *mpb.Progress
	bars     map[string]*mpb.Bar
	totals   map[string]int64
}

// NewRenderer cretes new progress.Renderer
func NewRenderer() *Renderer {
	return &Renderer{
		progress: mpb.New(
			mpb.WithFormat("[▇▇-]"),
		),
		bars:   map[string]*mpb.Bar{},
		totals: map[string]int64{},
	}
}

// Update re-render current progress
func (r *Renderer) Update(images []*ImageFetch) {
	for _, image := range images {
		total := getTotal(image)
		if bar, ok := r.bars[image.ContainerID]; !ok || total != r.totals[image.ContainerID] {
			if bar != nil {
				r.progress.RemoveBar(bar)
			}
			r.bars[image.ContainerID] = r.newBar(image, total)
			r.totals[image.ContainerID] = total
		}

		completed := getCompleted(image)

		r.bars[image.ContainerID].Incr(int(completed - r.bars[image.ContainerID].Current()))
	}
}

// Done marks all bars to be completed
func (r *Renderer) Done() {
	for _, bar := range r.bars {
		bar.Complete()
	}
}

// Stop halts the progress bar update process
func (r *Renderer) Stop() {
	r.progress.Stop()
}

func (r *Renderer) newBar(image *ImageFetch, total int64) *mpb.Bar {
	return r.progress.AddBar(total,
		mpb.PrependDecorators(
			// StaticName decorator with minWidth and no extra config
			// If you need to change name while rendering, use DynamicName
			decor.StaticName(image.Image, len(image.Image), 0),
			decor.ETA(8, 0),
		),
		// Appending decorators
		mpb.AppendDecorators(
			// Percentage decorator with minWidth and no extra config
			decor.Percentage(5, 0),
		),
	)
}

func getTotal(image *ImageFetch) (result int64) {
	for _, layer := range image.layers {
		result = result + layer.Total
	}
	return result
}

func getCompleted(image *ImageFetch) (result int64) {
	for _, layer := range image.layers {
		result = result + layer.Offset
	}
	return result
}
