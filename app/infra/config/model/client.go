package model

import (
	"go-common/library/time"
)

const (
	// HostOffline host offline state.
	HostOffline = 0
	// HostOnline host online state.
	HostOnline = 1
	// HostStateOK host state ok.
	HostStateOK = 2
	// UnknownVersion unknown version.
	UnknownVersion = -1
)

// Diff return to client.
type Diff struct {
	Version int64   `json:"version"`
	Diffs   []int64 `json:"diffs"`
}

// Version return to client.
type Version struct {
	Version int64 `json:"version"`
}

// ReVer reVer
type ReVer struct {
	Version int64  `json:"version"`
	Remark  string `json:"remark"`
}

// Versions versions
type Versions struct {
	Version []*ReVer `json:"version"`
	DefVer  int64    `json:"defver"`
}

// Content return to client.
type Content struct {
	Version int64  `json:"version"`
	Md5     string `json:"md5"`
	Content string `json:"content"`
}

// Namespace the key-value config object.
type Namespace struct {
	Name string            `json:"name"`
	Data map[string]string `json:"data"`
}

// Service service
type Service struct {
	Name         string
	BuildVersion string
	Env          string
	Token        string
	File         string
	Version      int64
	Host         string
	IP           string
	Appoint      int64
}

// NSValue config value.
type NSValue struct {
	ConfigID    int64  `json:"cid"`
	NamespaceID int64  `json:"nsid"`
	Name        string `json:"name"`
	Config      string `json:"config"`
}

// Value config value.
type Value struct {
	ConfigID int64  `json:"cid"`
	Name     string `json:"name"`
	Config   string `json:"config"`
}

// Host host.
type Host struct {
	Name          string    `json:"hostname"`
	Service       string    `json:"service"`
	BuildVersion  string    `json:"build"`
	IP            string    `json:"ip"`
	ConfigVersion int64     `json:"version"`
	HeartbeatTime time.Time `json:"heartbeat_time"`
	State         int       `json:"state"`
	Appoint       int64     `json:"appoint"`
	Customize     string    `json:"customize"`
	Force         int8      `json:"force"`
	ForceVersion  int64     `json:"force_version"`
}
