package storage

// Interface describes storage interface.
type Interface interface {
	Process() error
}
