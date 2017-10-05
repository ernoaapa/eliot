package progress

import "sync"

// ImageFetch stores container pull status
type ImageFetch struct {
	ContainerID string
	Image       string
	Resolved    bool
	layers      map[string]*Status
	mu          sync.Mutex
}

// Status represents single layer ref current progress
type Status struct {
	Ref    string
	Digest string
	Status string
	Offset int64
	Total  int64
}

// NewImageFetch creates new ImageFetch for given name
func NewImageFetch(containerID, image string) *ImageFetch {
	return CreateImageFetch(containerID, image, false, map[string]*Status{})
}

// CreateImageFetch creates new ImageFetch for given name
func CreateImageFetch(containerID, image string, resolved bool, layers map[string]*Status) *ImageFetch {
	return &ImageFetch{
		ContainerID: containerID,
		Image:       image,
		Resolved:    resolved,
		layers:      layers,
	}
}

// Add new layer ref to the progress list
func (s *ImageFetch) Add(ref, digest string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	// Started to fetch layers so image is resolved
	s.Resolved = true

	if _, ok := s.layers[ref]; ok {
		return // Already added
	}

	s.layers[ref] = &Status{
		Ref:    ref,
		Digest: digest,
		Status: "waiting",
	}
}

// SetToWaiting updates layer ref to the waiting state
func (s *ImageFetch) SetToWaiting(ref string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.layers[ref]; !ok {
		return // not added yet
	}

	s.layers[ref].Status = "waiting"
}

// SetToDownloading updates layer ref to the waiting state
func (s *ImageFetch) SetToDownloading(ref string, offset, total int64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.layers[ref]; !ok {
		return // not added yet
	}

	s.layers[ref].Status = "downloading"
	s.layers[ref].Offset = offset
	s.layers[ref].Total = total
}

// SetToDone updates layer ref to the done state
func (s *ImageFetch) SetToDone(ref string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.layers[ref]; !ok {
		return // not added yet
	}

	s.layers[ref].Status = "done"
}

// GetLayers return list of layers
func (s *ImageFetch) GetLayers() (result []Status) {
	for _, status := range s.layers {
		result = append(result, *status)
	}
	return result
}
