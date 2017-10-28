package progress

import (
	"sync"
)

// ImageFetch stores container pull status
type ImageFetch struct {
	ContainerID string
	Image       string
	Resolved    bool
	Failed      bool
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

// NewStatus creates new status in waiting state
func NewStatus(ref, digest string) *Status {
	return &Status{
		Ref:    ref,
		Digest: digest,
		Status: "waiting",
	}
}

// Waiting marks status to be in waiting state
func (s *Status) Waiting() {
	s.Status = "waiting"
}

// Downloading updates Status to downloading
func (s *Status) Downloading(offset, total int64) {
	s.Status = "downloading"
	s.Offset = offset
	s.Total = total
}

// Done marks Status to done state
func (s *Status) Done() {
	s.Offset = s.Total
	s.Status = "done"
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

// IsDone return true if all bytes of all layers are downloaded
func (s *ImageFetch) IsDone() bool {
	current, total := s.GetProgress()
	return current == total && current != 0
}

// GetProgress calculates current and total bytes of all layers
func (s *ImageFetch) GetProgress() (current, total int64) {
	for _, layer := range s.layers {
		current += layer.Offset
		total += layer.Total
	}
	return current, total
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

	s.layers[ref] = NewStatus(ref, digest)
}

// SetToWaiting updates layer ref to the waiting state
func (s *ImageFetch) SetToWaiting(ref string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.layers[ref]; !ok {
		return // not added yet
	}

	s.layers[ref].Waiting()
}

// SetToDownloading updates layer ref to the waiting state
func (s *ImageFetch) SetToDownloading(ref string, offset, total int64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.layers[ref]; !ok {
		return // not added yet
	}

	s.layers[ref].Downloading(offset, total)
}

// SetToDone updates layer ref to the done state
func (s *ImageFetch) SetToDone(ref string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.layers[ref]; !ok {
		return // not added yet
	}

	s.layers[ref].Done()
}

// SetToFailed marks fetch to be failed
func (s *ImageFetch) SetToFailed() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Failed = true
}

// AllDone marks all layers downloaded
func (s *ImageFetch) AllDone() {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, layer := range s.layers {
		layer.Done()
	}
}

// GetLayers return list of layers
func (s *ImageFetch) GetLayers() (result []Status) {
	for _, status := range s.layers {
		result = append(result, *status)
	}
	return result
}
