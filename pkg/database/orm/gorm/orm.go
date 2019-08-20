package gorm

import (
	"context"
	"fmt"
	"time"

	"github.com/bilibili/kratos/pkg/net/netutil/breaker"
	"github.com/bilibili/kratos/pkg/net/trace"
	"github.com/pkg/errors"

	"github.com/bilibili/kratos/pkg/log"
	"github.com/jinzhu/gorm"
	// mysql driver
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

const (
	_family          = "sql_client"
	_slowLogDuration = time.Millisecond * 250
)

// SQLCommon ...
type (
	// SQLCommon alias of gorm.SQLCommon
	SQLCommon = gorm.SQLCommon
	// Callback alias of gorm.Callback
	Callback = gorm.Callback
	// CallbackProcessor alias of gorm.CallbackProcessor
	CallbackProcessor = gorm.CallbackProcessor
	// Dialect alias of gorm.Dialect
	Dialect = gorm.Dialect
	// Scope ...
	Scope = gorm.Scope
	// DB ...
	DB = gorm.DB
	// Model ...
	Model = gorm.Model
	// ModelStruct ...
	ModelStruct = gorm.ModelStruct
	// Field ...
	Field = gorm.Field
	// FieldStruct ...
	StructField = gorm.StructField
	// RowQueryResult ...
	RowQueryResult = gorm.RowQueryResult
	// RowsQueryResult ...
	RowsQueryResult = gorm.RowsQueryResult
	// Association ...
	Association = gorm.Association
	// Errors ...
	Errors = gorm.Errors
	// Logger ...
	Logger = gorm.Logger
)

var (
	errSlowCommand = errors.New("mysql slow command")

	// IsRecordNotFoundError ...
	IsRecordNotFoundError = gorm.IsRecordNotFoundError

	// ErrRecordNotFound returns a "record not found error". Occurs only when attempting to query the database with a struct; querying with a slice won't return this error
	ErrRecordNotFound = gorm.ErrRecordNotFound
	// ErrInvalidSQL occurs when you attempt a query with invalid SQL
	ErrInvalidSQL = gorm.ErrInvalidSQL
	// ErrInvalidTransaction occurs when you are trying to `Commit` or `Rollback`
	ErrInvalidTransaction = gorm.ErrInvalidTransaction
	// ErrCantStartTransaction can't start transaction when you are trying to start one with `Begin`
	ErrCantStartTransaction = gorm.ErrCantStartTransaction
	// ErrUnaddressable unaddressable value
	ErrUnaddressable = gorm.ErrUnaddressable
)

// WithContext ...
func WithContext(ctx context.Context, db *DB) *DB {
	db.InstantSet("_context", ctx)
	return db
}

// Open ...
func Open(dialect string, config *Config) (*DB, error) {
	inner, err := gorm.Open(dialect, config.DSN)
	if err != nil {
		return nil, err
	}

	inner.LogMode(config.Debug)
	inner.DB().SetMaxIdleConns(config.Idle)
	inner.DB().SetMaxOpenConns(config.Active)

	if config.IdleTimeout != 0 {
		inner.DB().SetConnMaxLifetime(config.IdleTimeout)
	}

	brkGroup := breaker.NewGroup(config.Breaker)

	d, err := ParseDSN(config.DSN)
	if err != nil {
		return nil, err
	}

	label := fmt.Sprintf("%s_%s", d.DBName, d.Addr)

	brk := brkGroup.Get(d.Addr)

	replace := func(processor func() *gorm.CallbackProcessor, callbackName string, wrapper func(func(*Scope), string) func(*Scope)) {
		old := processor().Get(callbackName)
		processor().Replace(callbackName, wrapper(old, label))
	}

	invoke := func(op string) func(callback func(scope *Scope), label string) func(scope *Scope) {
		return func(callback func(scope *Scope), label string) func(scope *Scope) {
			return func(scope *Scope) {
				fn := func() {
					beg := time.Now()
					callback(scope)

					// todo(gorexlv): desensitization required
					slowLog(scope.CombinedConditionSql(), beg)

					// error metric
					if scope.HasError() {
						_incConnMetric(d.DBName+"."+scope.TableName(), d.Addr, "ERR")
					} else {
						_incConnMetric(d.DBName+"."+scope.TableName(), d.Addr, "OK")
					}
					_observeConn(time.Since(beg), d.DBName+"."+scope.TableName(), d.Addr)
				}

				if val, ok := scope.Get("_context"); ok {
					if ctx, ok := val.(context.Context); ok {
						if t, ok := trace.FromContext(ctx); ok {
							t = t.Fork(_family, "query")
							t.SetTag(trace.String(trace.TagAddress, d.Addr), trace.String(trace.TagComment, scope.CombinedConditionSql()))
							defer t.Finish(&scope.DB().Error)
						}
					}
				}

				if err := brk.Allow(); err != nil {
					_incConnMetric(d.DBName+"."+scope.TableName(), d.Addr, "BAN")
					scope.DB().Error = err
					return
				}

				// todo(gorexlv) add tracing
				fn()
			}
		}
	}

	replace(inner.Callback().Delete, "gorm:delete", invoke("delete"))
	replace(inner.Callback().Update, "gorm:update", invoke("update"))
	replace(inner.Callback().Create, "gorm:create", invoke("create"))
	replace(inner.Callback().Query, "gorm:query", invoke("query"))
	replace(inner.Callback().RowQuery, "gorm:row_query", invoke("row_query"))

	return inner, err
}

func slowLog(statement string, now time.Time) {
	du := time.Since(now)
	if du > _slowLogDuration {
		log.Warn("%s slow log statement: %s time: %v", _family, statement, du)
	}
}
