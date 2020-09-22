package protocol

// Driver is the otlp protocol
type Driver interface {
	Open(name string) (Conn, error)
}
