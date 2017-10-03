package containerd

import (
	"context"
	"strings"

	"github.com/containerd/containerd"
	"github.com/containerd/containerd/containers"
	"github.com/ernoaapa/can/pkg/model"
	"github.com/ernoaapa/can/pkg/runtime/containerd/mapping"
	specs "github.com/opencontainers/runtime-spec/specs-go"
)

// WithCwd spets the container current working directory (cwd)
func WithCwd(cwd string) containerd.SpecOpts {
	return func(_ context.Context, _ *containerd.Client, _ *containers.Container, s *specs.Spec) error {
		s.Process.Cwd = cwd
		return nil
	}
}

// WithEnv you can add or override process environment variables
// overrides should be list of strings in format 'KEY=value'
func WithEnv(overrides []string) containerd.SpecOpts {
	return func(_ context.Context, _ *containerd.Client, _ *containers.Container, s *specs.Spec) error {
		if len(overrides) > 0 {
			s.Process.Env = replaceOrAppendEnvValues(s.Process.Env, overrides)
		}
		return nil
	}
}

// replaceOrAppendEnvValues returns the defaults with the overrides either
// replaced by env key or appended to the list
func replaceOrAppendEnvValues(defaults, overrides []string) []string {
	cache := make(map[string]int, len(defaults))
	for i, e := range defaults {
		parts := strings.SplitN(e, "=", 2)
		cache[parts[0]] = i
	}

	for _, value := range overrides {
		// Values w/o = means they want this env to be removed/unset.
		if !strings.Contains(value, "=") {
			if i, exists := cache[value]; exists {
				defaults[i] = "" // Used to indicate it should be removed
			}
			continue
		}

		// Just do a normal set/update
		parts := strings.SplitN(value, "=", 2)
		if i, exists := cache[parts[0]]; exists {
			defaults[i] = value
		} else {
			defaults = append(defaults, value)
		}
	}

	// Now remove all entries that we want to "unset"
	for i := 0; i < len(defaults); i++ {
		if defaults[i] == "" {
			defaults = append(defaults[:i], defaults[i+1:]...)
			i--
		}
	}

	return defaults
}

// WithMounts you can add mount points to the container
func WithMounts(mounts []model.Mount) containerd.SpecOpts {
	return func(_ context.Context, _ *containerd.Client, _ *containers.Container, s *specs.Spec) error {
		for _, mount := range mounts {
			s.Mounts = append(s.Mounts, mapping.MapMountToContainerdModel(mount))
		}
		return nil
	}
}
