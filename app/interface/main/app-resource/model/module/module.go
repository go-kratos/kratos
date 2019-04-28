package module

import (
	xtime "go-common/library/time"
)

const (
	Total       = 0
	Incremental = 1

	EnvRelease = "1"
	EnvTest    = "2"
	EnvDefault = "3"

	NotValid = int8(0)
	Valid    = int8(1)
)

type ResourcePool struct {
	ID        int         `json:"-"`
	Name      string      `json:"name"`
	Resources []*Resource `json:"resources,omitempty"`
}

type Resource struct {
	ID           int        `json:"-"`
	ResID        int        `json:"-"`
	Name         string     `json:"name"`
	Compresstype int        `json:"compresstype"`
	Type         string     `json:"type"`
	URL          string     `json:"url"`
	MD5          string     `json:"md5"`
	TotalMD5     string     `json:"total_md5"`
	Size         int        `json:"size"`
	Version      int        `json:"ver"`
	Increment    int        `json:"increment"`
	FromVer      int        `json:"-"`
	Condition    *Condition `json:"-"`
	Level        int        `json:"level,omitempty"`
	IsWifi       int8       `json:"is_wifi"`
}

type Condition struct {
	ID        int                  `json:"-"`
	ResID     int                  `json:"-"`
	STime     xtime.Time           `json:"stime"`
	ETime     xtime.Time           `json:"etime"`
	Valid     int8                 `json:"valid"`
	ValidTest int8                 `json:"valid_test"`
	Default   int                  `json:"-"`
	Columns   map[string][]*Column `json:"columns"`
	IsWifi    int8                 `json:"-"`
}

type Column struct {
	Condition string `json:"condition"`
	Value     string `json:"value"`
}

type Versions struct {
	PoolName string `json:"name"`
	Resource []struct {
		ResourceName string      `json:"name"`
		Version      interface{} `json:"ver"`
	} `json:"resources"`
}
