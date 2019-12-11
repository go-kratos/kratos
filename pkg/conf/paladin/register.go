package paladin

import (
	"fmt"
	"sort"
	"sync"
)

var (
	driversMu sync.RWMutex
	drivers   = make(map[string]Driver)
)

// Register makes a paladin driver available by the provided name.
// If Register is called twice with the same name or if driver is nil,
// it panics.
func Register(name string, driver Driver) {
	driversMu.Lock()
	defer driversMu.Unlock()

	if driver == nil {
		panic("paladin: driver is nil")
	}

	if _, dup := drivers[name]; dup {
		panic("paladin: Register called twice for driver " + name)
	}

	drivers[name] = driver
}

// Drivers returns a sorted list of the names of the registered paladin driver.
func Drivers() []string {
	driversMu.RLock()
	defer driversMu.RUnlock()

	var list []string
	for name := range drivers {
		list = append(list, name)
	}

	sort.Strings(list)
	return list
}

// GetDriver returns a driver implement by name.
func GetDriver(name string) (Driver, error) {
	driversMu.RLock()
	driveri, ok := drivers[name]
	driversMu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("paladin: unknown driver %q (forgotten import?)", name)
	}
	return driveri, nil
}
