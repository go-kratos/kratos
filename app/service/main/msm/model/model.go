package model

import (
	"container/list"
)

const (
	// CodePlatDefaut db common plat code.
	CodePlatDefaut = 1
	// CodePlatDefautMsg db common plat msg.
	CodePlatDefautMsg = "common"
	// CodeDelStatus db delete status.
	CodeDelStatus = 2
	// HostOffline host offline state.
	HostOffline = 0
)

// RPC rpc node value.
type RPC struct {
	Proto  string `json:"Proto"`
	Addr   string `json:"Addr"`
	Group  string `json:"Group"`
	Weight int    `json:"Weight"`
}

// Code ver and message.
type Code struct {
	Ver  int64
	Code int
	Msg  string
}

// Codes all codes local map cache.
type Codes struct {
	Ver  int64
	MD5  string
	Code map[int]string
}

// Version list and map.
type Version struct {
	List *list.List
	Map  map[int64]*list.Element
}

// Databus databus rule.
type Databus struct {
	Topic       string `json:"topic"`
	Group       string `json:"group"`
	Cluster     string `json:"cluster"`
	Business    string `json:"business"`
	Operation   int8   `json:"operation"`
	Leader      string `json:"leader"`
	Phone       string `json:"phone"`
	Email       string `json:"email"`
	AlarmSwitch int8   `json:"alarmSwitch"`
	Users       string `json:"users"`
	AlarmRule   string `json:"alarmRule"`
}

// Databuss databuss rules.
type Databuss struct {
	Rules []*Databus `json:"rules"`
	MD5   string     `json:"md5"`
}

// Limit limit.
type Limit struct {
	Burst int     `json:"burst"`
	Rate  float64 `json:"rate"`
}

// Limits limits.
type Limits struct {
	Apps map[string]*Limit `json:"apps"`
	MD5  string            `json:"md5"`
}

// Host host.
type Host struct {
	Name  string `json:"hostname"`
	State int    `json:"state"`
}

//CodesLangs ...
type CodesLangs struct {
	Ver  int64
	MD5  string
	Code map[int]map[string]string
}

// //Langs ...
// type Langs struct {
// 	Default  string `json:"default"`
// 	Localeds []*Locale
// }

//Locale ...
// type Locale struct {
// 	Locale  string
// 	Message string
// }

//CodeLangs ...
type CodeLangs struct {
	Ver  int64
	Code int
	Msg  map[string]string
}
