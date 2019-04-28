package orm

import (
	"strings"
	"time"

	"go-common/library/ecode"
	"go-common/library/log"
	xtime "go-common/library/time"

	// database driver
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

// Config mysql config.
type Config struct {
	DSN         string         // data source name.
	Active      int            // pool
	Idle        int            // pool
	IdleTimeout xtime.Duration // connect max life time.
}

type ormLog struct{}

func (l ormLog) Print(v ...interface{}) {
	log.Info(strings.Repeat("%v ", len(v)), v...)
}

func init() {
	gorm.ErrRecordNotFound = ecode.NothingFound
}

// NewMySQL new db and retry connection when has error.
func NewMySQL(c *Config) (db *gorm.DB) {
	db, err := gorm.Open("mysql", c.DSN)
	if err != nil {
		log.Error("db dsn(%s) error: %v", c.DSN, err)
		panic(err)
	}
	db.DB().SetMaxIdleConns(c.Idle)
	db.DB().SetMaxOpenConns(c.Active)
	db.DB().SetConnMaxLifetime(time.Duration(c.IdleTimeout) / time.Second)
	db.SetLogger(ormLog{})
	return
}
