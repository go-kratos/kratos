package config

import "time"

type BuildStatus struct {
	DB DB `json:"db, omitempty"`
}
type DB struct {
	IP          string        `json:"ip, omitempty"`
	Port        string        `json:"port, omitempty"`
	Name        string        `json:"name, omitempty"`
	Username    string        `json:"username, omitempty"`
	Password    string        `json:"password, omitempty"`
	Active      int           `json:"active, omitempty"`
	Idle        int           `json:"idle, omitempty"`
	IdleTimeout time.Duration `json:"idleTimeout, omitempty"`
}
