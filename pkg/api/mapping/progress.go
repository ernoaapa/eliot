package mapping

import (
	pb "github.com/ernoaapa/eliot/pkg/api/services/pods/v1"
	"github.com/ernoaapa/eliot/pkg/progress"
)

// MapImageFetchProgressToAPIModel maps image fetch progress information to API model
func MapImageFetchProgressToAPIModel(progresses []*progress.ImageFetch) (result []*pb.ImageFetch) {
	for _, progress := range progresses {
		Layers := []*pb.ImageLayerStatus{}
		for _, layer := range progress.GetLayers() {
			Layers = append(Layers, &pb.ImageLayerStatus{
				Ref:    layer.Ref,
				Digest: layer.Digest,
				Offset: layer.Offset,
				Total:  layer.Total,
			})
		}
		result = append(result, &pb.ImageFetch{
			ContainerID: progress.ContainerID,
			Image:       progress.Image,
			Resolved:    progress.Resolved,
			Failed:      progress.Failed,
			Layers:      Layers,
		})
	}
	return result
}

// MapAPIModelToImageFetchProgress maps image fetch progress information to API model
func MapAPIModelToImageFetchProgress(progresses []*pb.ImageFetch) (result []*progress.ImageFetch) {
	for _, image := range progresses {
		statuses := map[string]*progress.Status{}
		for _, layer := range image.Layers {
			statuses[layer.Ref] = &progress.Status{
				Ref:    layer.Ref,
				Digest: layer.Digest,
				Offset: layer.Offset,
				Total:  layer.Total,
			}
		}
		result = append(result, progress.CreateImageFetch(
			image.ContainerID,
			image.Image,
			image.Resolved,
			statuses,
		))
	}
	return result
}
