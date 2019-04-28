package backend

import (
	"context"
	"fmt"
	"net"
	"strings"
)

var factoryMap map[string]Factory

func init() {
	factoryMap = make(map[string]Factory)
}

// Factory backend factory
type Factory func(map[string]interface{}) (Backend, error)

// Registry registry backdend
func Registry(name string, factory Factory) {
	if _, ok := factoryMap[name]; ok {
		panic(fmt.Sprintf("backend %s already exists", name))
	}
	factoryMap[name] = factory
}

// New backend
func New(name string, conf map[string]interface{}) (Backend, error) {
	if factory, ok := factoryMap[name]; ok {
		return factory(conf)
	}
	return nil, fmt.Errorf("backend %s not exists", name)
}

// Selector selector
type Selector struct {
	Env      string
	Region   string
	Zone     string
	Hostname string
}

func (s Selector) String() string {
	strs := make([]string, 0, 4)
	if s.Env != "" {
		strs = append(strs, s.Env)
	}
	if s.Region != "" {
		strs = append(strs, s.Region)
	}
	if s.Zone != "" {
		strs = append(strs, s.Zone)
	}
	if s.Hostname != "" {
		strs = append(strs, s.Hostname)
	}
	return strings.Join(strs, "-")
}

// Target global unique application identifier
type Target struct {
	Name string
}

func (t Target) String() string {
	return t.Name
}

// ParseName parse qname get name and selector
func ParseName(name string, defaultSel Selector) (target Target, sel Selector, err error) {
	// TODO: support selector
	return Target{Name: name}, defaultSel, nil
}

// Metadata metadata contain env, zone, region e.g.
type Metadata struct {
	ClientHost       string
	LatestTimestamps string
}

// Instance service instance struct
type Instance struct {
	Region      string `json:"region"`
	Zone        string `json:"zone"`
	Env         string `json:"env"`
	Hostname    string `json:"hostname"`
	DiscoveryID string `json:"discovery_id"`
	TreeID      int64  `json:"tree_id"`
	IPAddr      net.IP `json:"ip_addr,omitempty"` // hacked field
}

// Backend provide service query
type Backend interface {
	Ping(ctx context.Context) error
	Query(ctx context.Context, target Target, sel Selector, md Metadata) ([]*Instance, error)
	Close(ctx context.Context) error
}
