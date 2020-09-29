package protocol

import (
	"fmt"
	"sync"
)

var (
	driversMu sync.RWMutex
	drivers   = make(map[string]Driver)
)

type ClientConn interface {
}

type DB struct {
	connector ClientConn
}

// Driver is the otlp protocol
type Driver interface {
	Open(name string) (ClientConn, error)
}

func Register(name string, driver Driver) {

	driversMu.Lock()
	defer driversMu.Unlock()
	if driver == nil {
		panic("sql: Register driver is nil")
	}
	if _, dup := drivers[name]; dup {
		panic("sql: Register called twice for driver " + name)
	}
	drivers[name] = driver
}

func Open(driverName, dataSourceName string) (*DB, error) {
	driversMu.RLock()
	driveri, ok := drivers[driverName]
	driversMu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("sql: unknown driver %q (forgotten import?)", driverName)
	}

	connector, err := driveri.Open(dataSourceName)
	if err != nil {
		return nil, err
	}
	return OpenDB(connector), nil
}

func OpenDB(c ClientConn) *DB {
	db := &DB{
		connector: c,
	}
	return db
}
