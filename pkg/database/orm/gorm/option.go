package gorm

import (
	"time"

	"github.com/bilibili/kratos/pkg/net/netutil/breaker"
)

// Config options
type Config struct {
	// mysql://root:secret@tcp(127.0.0.1:3307)/mysql?timeout=20s&readTimeout=20s
	DSN         string
	Debug       bool
	Idle        int
	Active      int
	IdleTimeout time.Duration
	DialTimeout time.Duration
	Breaker     *breaker.Config
}
