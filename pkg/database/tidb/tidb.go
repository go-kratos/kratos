package tidb

import (
	"github.com/bilibili/kratos/pkg/log"
	"github.com/bilibili/kratos/pkg/net/netutil/breaker"
	"github.com/bilibili/kratos/pkg/time"

	// database driver
	_ "github.com/go-sql-driver/mysql"
)

// Config mysql config.
type Config struct {
	DSN          string          // dsn
	Active       int             // pool
	Idle         int             // pool
	IdleTimeout  time.Duration   // connect max life time.
	QueryTimeout time.Duration   // query sql timeout
	ExecTimeout  time.Duration   // execute sql timeout
	TranTimeout  time.Duration   // transaction sql timeout
	Breaker      *breaker.Config // breaker
}

// NewTiDB new db and retry connection when has error.
func NewTiDB(c *Config) (db *DB) {
	if c.QueryTimeout == 0 || c.ExecTimeout == 0 || c.TranTimeout == 0 {
		panic("tidb must be set query/execute/transction timeout")
	}
	db, err := Open(c)
	if err != nil {
		log.Error("open tidb error(%v)", err)
		panic(err)
	}
	return
}
