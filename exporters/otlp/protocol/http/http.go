package http

import (
	"fmt"

	"go.opentelemetry.io/otel/exporters/otlp/protocol"
)

// MySQLDriver is exported to make the driver directly accessible.
// In general the driver is used via the database/sql package.
type MySQLDriver struct{}

func init() {
	protocol.Register("http", &MySQLDriver{})
	fmt.Println("init export otlp http")
}

func (d MySQLDriver) Open(dsn string) (protocol.ClientConn, error) {
	return nil, nil
}
