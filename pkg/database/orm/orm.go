package orm

import (
	"strings"
	"time"

	"github.com/go-kratos/kratos/pkg/log"
	xtime "github.com/go-kratos/kratos/pkg/time"

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
	DbPrev      string         //表前缀
}

type ormLog struct{}

func (l ormLog) Print(v ...interface{}) {
	log.Info(strings.Repeat("%v ", len(v)), v...)
}

// NewMySQL new db and retry connection when has error.
func NewMySQL(c *Config) (db *gorm.DB) {
	db, err := gorm.Open("mysql", c.DSN)
	if err != nil {
		log.Error("db dsn(%s) error: %v", c.DSN, err)
		panic(err)
	}

	if c.DbPrev != "" {
		gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
			return c.DbPrev + defaultTableName
		}
		db.SingularTable(true)
	}
	//SetMaxIdleConns设置空闲连接池中的最大连接数。
	db.DB().SetMaxIdleConns(c.Idle)
	// SetMaxOpenConns设置到数据库的最大打开连接数。
	db.DB().SetMaxOpenConns(c.Active)
	// SetConnMaxLifetime设置连接可以重用的最长时间。
	db.DB().SetConnMaxLifetime(time.Duration(c.IdleTimeout) / time.Second)
	db.SetLogger(ormLog{})
	return
}
