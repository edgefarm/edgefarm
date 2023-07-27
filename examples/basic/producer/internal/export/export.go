package export

// Exporter is the interface to publish exported messages
type Exporter interface {
	PubExport(subject string, data []byte) error
}
