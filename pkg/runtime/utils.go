package runtime

import (
	"os"

	"github.com/ernoaapa/elliot/pkg/fs"
	"github.com/ernoaapa/elliot/pkg/model"
	"github.com/pkg/errors"
)

func ensureMountSourceDirExists(mounts []model.Mount) error {
	for _, mount := range mounts {
		if fs.FileExist(mount.Source) {
			// it's a file...
			continue
		}
		if !fs.DirExist(mount.Source) {
			err := os.MkdirAll(mount.Source, os.ModePerm)
			if err != nil {
				return errors.Wrapf(err, "Error while mkdir recursively mount source [%s]", mount.Source)
			}
		}
	}
	return nil
}
