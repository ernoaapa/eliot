package manifest

// Source is interface for all state sources
type Source interface {
	Start()
	Stop()
}
