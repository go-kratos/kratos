package dao

import (
	"context"
	"database/sql"
	"time"

	"go-common/app/service/main/antispam/conf"
	"go-common/library/cache/redis"
	xsql "go-common/library/database/sql"
)

// Executer .
type Executer interface {
	Exec(ctx context.Context, SQL string, args ...interface{}) (sql.Result, error)
}

// Querier .
type Querier interface {
	QueryRow(ctx context.Context, SQL string, args ...interface{}) *xsql.Row
	Query(ctx context.Context, SQL string, args ...interface{}) (*xsql.Rows, error)
}

// KeywordDao .
type KeywordDao interface {
	GetByID(context.Context, int64) (*Keyword, error)
	GetByIDs(context.Context, []int64) ([]*Keyword, error)
	GetByCond(context.Context, *Condition) ([]*Keyword, int64, error)
	GetByOffsetLimit(context.Context, *Condition) ([]*Keyword, error)
	GetByAreaAndContents(context.Context, *Condition) ([]*Keyword, error)
	GetByAreaAndContent(context.Context, *Condition) (*Keyword, error)
	GetRubbish(context.Context, *Condition) ([]*Keyword, error)

	Insert(context.Context, *Keyword) (*Keyword, error)
	Update(context.Context, *Keyword) (*Keyword, error)
	DeleteByIDs(context.Context, []int64) ([]*Keyword, error)
}

// RuleDao .
type RuleDao interface {
	GetByID(context.Context, int64) (*Rule, error)
	GetByIDs(context.Context, []int64) ([]*Rule, error)
	GetByCond(context.Context, *Condition) ([]*Rule, int64, error)
	GetByArea(context.Context, *Condition) ([]*Rule, error)
	GetByAreaAndTypeAndScope(context.Context, *Condition) (*Rule, error)
	GetByAreaAndLimitType(context.Context, *Condition) ([]*Rule, error)

	Insert(context.Context, *Rule) (*Rule, error)
	Update(context.Context, *Rule) (*Rule, error)
}

// RegexpDao .
type RegexpDao interface {
	GetByID(context.Context, int64) (*Regexp, error)
	GetByIDs(context.Context, []int64) ([]*Regexp, error)
	GetByCond(context.Context, *Condition) ([]*Regexp, int64, error)
	GetByAreaAndContent(context.Context, *Condition) (*Regexp, error)
	GetByContents(context.Context, []string) ([]*Regexp, error)

	Insert(context.Context, *Regexp) (*Regexp, error)
	Update(context.Context, *Regexp) (*Regexp, error)
}

// Tx .
type Tx interface {
	Executer

	RegexpTx
	KeywordTx
	RuleTx

	Commit() error
	Rollback() error
}

// KeywordTx .
type KeywordTx interface {
	InsertKeyword(*Keyword) error
	UpdateKeyword(*Keyword) error
}

// RegexpTx .
type RegexpTx interface {
	InsertRegexp(*Regexp) error
	UpdateRegexp(*Regexp) error
}

// RuleTx .
type RuleTx interface {
	InsertRule(*Rule) error
	UpdateRule(*Rule) error
}

// Dao .
type Dao struct {
	// db          *xsql.DB
	redis       *redis.Pool
	redisExpire int
}

// New a dao and return.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		// db
		// db: xsql.NewMySQL(c.MySQL.AntiSpam),
		// redis
		redis:       redis.NewPool(c.Redis.Config),
		redisExpire: int(time.Duration(c.Redis.IndexExpire) / time.Second),
	}
	return
}

// Ping check connection used in dao
func (d *Dao) Ping(c context.Context) (err error) {
	if err = d.pingRedis(c); err != nil {
		return
	}
	err = PingMySQL(c)
	return
}

// Close close all connection.
func (d *Dao) Close() {
	if d.redis != nil {
		d.redis.Close()
	}
	Close()
}
