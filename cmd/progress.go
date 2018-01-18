package cmd

import (
	ui "github.com/ernoaapa/eliot/pkg/cmd/ui"
	"github.com/ernoaapa/eliot/pkg/progress"
)

// ShowDownloadProgress prints UI "downloading" lines and updates until
// the progress channel closes
func ShowDownloadProgress(progressc <-chan []*progress.ImageFetch) {
	lines := map[string]ui.Line{}
	for fetches := range progressc {
		for _, fetch := range fetches {
			if _, ok := lines[fetch.Image]; !ok {
				lines[fetch.Image] = ui.NewLine().Loadingf("Download %s", fetch.Image)
			}

			if fetch.IsDone() {
				if fetch.Failed {
					lines[fetch.Image].Errorf("Failed %s", fetch.Image)
				} else {
					lines[fetch.Image].Donef("Downloaded %s", fetch.Image)
				}
			} else {
				current, total := fetch.GetProgress()
				lines[fetch.Image].WithProgress(current, total)
			}
		}
	}

	for image, line := range lines {
		line.Donef("Completed %s", image)
	}
}
