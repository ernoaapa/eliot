package containerd

import (
	"context"
	"time"

	"github.com/containerd/containerd"
	"github.com/containerd/containerd/errdefs"
	"github.com/ernoaapa/can/pkg/progress"
	digest "github.com/opencontainers/go-digest"
	log "github.com/sirupsen/logrus"
)

// UpdateFetchProgress start goroutine to update the fetch status until done channel closes
func UpdateFetchProgress(done <-chan struct{}, client *containerd.Client, progress *progress.ImageFetch) {
	var (
		ctx   = context.Background()
		start = time.Now()
	)

	for {
		select {
		case <-done:
			return
		case <-time.After(100 * time.Millisecond):
			active, err := client.ContentStore().ListStatuses(ctx, "")
			if err != nil {
				log.Errorf("Error while listing active content digestions: %s", err)
				continue
			}

			activeDownloads := []string{}

			for _, active := range active {
				activeDownloads = append(activeDownloads, active.Ref)
				log.Printf("Update %s to %d %d", active.Ref, active.Offset, active.Total)
				progress.SetToDownloading(active.Ref, active.Offset, active.Total)
			}

			for _, layer := range filter(progress.GetLayers(), activeDownloads) {
				info, err := client.ContentStore().Info(ctx, digest.FromString(layer.Digest))

				if err != nil {
					if errdefs.IsNotFound(err) {
						progress.SetToWaiting(layer.Ref)
					} else {
						log.Errorf("Error while fetching [%s] image layer: %s", layer.Ref, err)
					}
				}

				if info.CreatedAt.After(start) {
					progress.SetToDone(layer.Ref)
				}
			}
		}
	}
}

func filter(statuses []progress.Status, refs []string) (result []progress.Status) {
	for _, status := range statuses {
		if !contains(refs, status.Ref) {
			result = append(result, status)
		}
	}
	return result
}

func contains(source []string, value string) bool {
	for _, item := range source {
		if item == value {
			return true
		}
	}
	return false
}
