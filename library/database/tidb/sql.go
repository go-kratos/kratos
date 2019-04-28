package tidb

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"go-common/library/naming"
	"go-common/library/net/netutil/breaker"
	"go-common/library/net/trace"

	"github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
)

const (
	_family = "tidb_client"
)

var (
	// ErrStmtNil prepared stmt error
	ErrStmtNil = errors.New("sql: prepare failed and stmt nil")
	// ErrNoRows is returned by Scan when QueryRow doesn't return a row.
	// In such a case, QueryRow returns a placeholder *Row value that defers
	// this error until a Scan.
	ErrNoRows = sql.ErrNoRows
	// ErrTxDone transaction done.
	ErrTxDone = sql.ErrTxDone
)

// DB database.
type DB struct {
	conf         *Config
	conns        []*conn
	idx          int64
	dis          naming.Resolver
	appid        string
	mutex        sync.RWMutex
	breakerGroup *breaker.Group
}

// conn database connection
type conn struct {
	*sql.DB
	breaker breaker.Breaker
	conf    *Config
	addr    string
}

// Tx transaction.
type Tx struct {
	db     *conn
	tx     *sql.Tx
	t      trace.Trace
	c      context.Context
	cancel func()
}

// Row row.
type Row struct {
	err error
	*sql.Row
	db     *conn
	query  string
	args   []interface{}
	t      trace.Trace
	cancel func()
}

// Scan copies the columns from the matched row into the values pointed at by dest.
func (r *Row) Scan(dest ...interface{}) (err error) {
	if r.t != nil {
		defer r.t.Finish(&err)
	}
	if r.err != nil {
		err = r.err
	} else if r.Row == nil {
		err = ErrStmtNil
	}
	if err != nil {
		return
	}
	err = r.Row.Scan(dest...)
	if r.cancel != nil {
		r.cancel()
	}
	r.db.onBreaker(&err)
	if err != ErrNoRows {
		err = errors.Wrapf(err, "addr: %s, query %s args %+v", r.db.addr, r.query, r.args)
	}
	return
}

// Rows rows.
type Rows struct {
	*sql.Rows
	cancel func()
}

// Close closes the Rows, preventing further enumeration. If Next is called
// and returns false and there are no further result sets,
// the Rows are closed automatically and it will suffice to check the
// result of Err. Close is idempotent and does not affect the result of Err.
func (rs *Rows) Close() (err error) {
	err = errors.WithStack(rs.Rows.Close())
	if rs.cancel != nil {
		rs.cancel()
	}
	return
}

// Stmt prepared stmt.
type Stmt struct {
	db    *conn
	tx    bool
	query string
	stmt  atomic.Value
	t     trace.Trace
}

// Stmts random prepared stmt.
type Stmts struct {
	query string
	sts   map[string]*Stmt
	mu    sync.RWMutex
	db    *DB
}

// Open opens a database specified by its database driver name and a
// driver-specific data source name, usually consisting of at least a database
// name and connection information.
func Open(c *Config) (db *DB, err error) {
	db = &DB{conf: c, breakerGroup: breaker.NewGroup(c.Breaker)}
	cfg, err := mysql.ParseDSN(c.DSN)
	if err != nil {
		return
	}
	var dsns []string
	if cfg.Net == "discovery" {
		db.appid = cfg.Addr
		for _, addr := range db.disc() {
			dsns = append(dsns, genDSN(c.DSN, addr))
		}
	} else {
		dsns = append(dsns, c.DSN)
	}

	cs := make([]*conn, 0, len(dsns))
	for _, dsn := range dsns {
		r, err := db.connectDSN(dsn)
		if err != nil {
			return db, err
		}
		cs = append(cs, r)
	}
	db.conns = cs
	return
}

func (db *DB) connectDSN(dsn string) (c *conn, err error) {
	d, err := connect(db.conf, dsn)
	if err != nil {
		return
	}
	addr := parseDSNAddr(dsn)
	brk := db.breakerGroup.Get(addr)
	c = &conn{DB: d, breaker: brk, conf: db.conf, addr: addr}
	return
}

func connect(c *Config, dataSourceName string) (*sql.DB, error) {
	d, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		err = errors.WithStack(err)
		return nil, err
	}
	d.SetMaxOpenConns(c.Active)
	d.SetMaxIdleConns(c.Idle)
	d.SetConnMaxLifetime(time.Duration(c.IdleTimeout))
	return d, nil
}

func (db *DB) conn() (c *conn) {
	db.mutex.RLock()
	c = db.conns[db.index()]
	db.mutex.RUnlock()
	return
}

// Begin starts a transaction. The isolation level is dependent on the driver.
func (db *DB) Begin(c context.Context) (tx *Tx, err error) {
	return db.conn().begin(c)
}

// Exec executes a query without returning any rows.
// The args are for any placeholder parameters in the query.
func (db *DB) Exec(c context.Context, query string, args ...interface{}) (res sql.Result, err error) {
	return db.conn().exec(c, query, args...)
}

// Prepare creates a prepared statement for later queries or executions.
// Multiple queries or executions may be run concurrently from the returned
// statement. The caller must call the statement's Close method when the
// statement is no longer needed.
func (db *DB) Prepare(query string) (*Stmt, error) {
	return db.conn().prepare(query)
}

// Prepared creates a prepared statement for later queries or executions.
// Multiple queries or executions may be run concurrently from the returned
// statement. The caller must call the statement's Close method when the
// statement is no longer needed.
func (db *DB) Prepared(query string) (s *Stmts) {
	s = &Stmts{query: query, sts: make(map[string]*Stmt), db: db}
	for _, c := range db.conns {
		st := c.prepared(query)
		s.mu.Lock()
		s.sts[c.addr] = st
		s.mu.Unlock()
	}
	return
}

// Query executes a query that returns rows, typically a SELECT. The args are
// for any placeholder parameters in the query.
func (db *DB) Query(c context.Context, query string, args ...interface{}) (rows *Rows, err error) {
	return db.conn().query(c, query, args...)
}

// QueryRow executes a query that is expected to return at most one row.
// QueryRow always returns a non-nil value. Errors are deferred until Row's
// Scan method is called.
func (db *DB) QueryRow(c context.Context, query string, args ...interface{}) *Row {
	return db.conn().queryRow(c, query, args...)
}

func (db *DB) index() int {
	if len(db.conns) == 1 {
		return 0
	}
	v := atomic.AddInt64(&db.idx, 1)
	return int(v) % len(db.conns)
}

// Close closes the databases, releasing any open resources.
func (db *DB) Close() (err error) {
	db.mutex.RLock()
	defer db.mutex.RUnlock()
	for _, d := range db.conns {
		if e := d.Close(); e != nil {
			err = errors.WithStack(e)
		}
	}
	return
}

// Ping verifies a connection to the database is still alive, establishing a
// connection if necessary.
func (db *DB) Ping(c context.Context) (err error) {
	if err = db.conn().ping(c); err != nil {
		return
	}
	return
}

func (db *conn) onBreaker(err *error) {
	if err != nil && *err != nil && *err != sql.ErrNoRows && *err != sql.ErrTxDone {
		db.breaker.MarkFailed()
	} else {
		db.breaker.MarkSuccess()
	}
}

func (db *conn) begin(c context.Context) (tx *Tx, err error) {
	now := time.Now()
	t, ok := trace.FromContext(c)
	if ok {
		t = t.Fork(_family, "begin")
		t.SetTag(trace.String(trace.TagAddress, db.addr), trace.String(trace.TagComment, ""))
		defer func() {
			if err != nil {
				t.Finish(&err)
			}
		}()
	}
	if err = db.breaker.Allow(); err != nil {
		stats.Incr("tidb:begin", "breaker")
		return
	}
	_, c, cancel := db.conf.TranTimeout.Shrink(c)
	rtx, err := db.BeginTx(c, nil)
	stats.Timing("tidb:begin", int64(time.Since(now)/time.Millisecond))
	if err != nil {
		err = errors.WithStack(err)
		cancel()
		return
	}
	tx = &Tx{tx: rtx, t: t, db: db, c: c, cancel: cancel}
	return
}

func (db *conn) exec(c context.Context, query string, args ...interface{}) (res sql.Result, err error) {
	now := time.Now()
	if t, ok := trace.FromContext(c); ok {
		t = t.Fork(_family, "exec")
		t.SetTag(trace.String(trace.TagAddress, db.addr), trace.String(trace.TagComment, query))
		defer t.Finish(&err)
	}
	if err = db.breaker.Allow(); err != nil {
		stats.Incr("tidb:exec", "breaker")
		return
	}
	_, c, cancel := db.conf.ExecTimeout.Shrink(c)
	res, err = db.ExecContext(c, query, args...)
	cancel()
	db.onBreaker(&err)
	stats.Timing("tidb:exec", int64(time.Since(now)/time.Millisecond))
	if err != nil {
		err = errors.Wrapf(err, "addr: %s exec:%s, args:%+v", db.addr, query, args)
	}
	return
}

func (db *conn) ping(c context.Context) (err error) {
	now := time.Now()
	if t, ok := trace.FromContext(c); ok {
		t = t.Fork(_family, "ping")
		t.SetTag(trace.String(trace.TagAddress, db.addr), trace.String(trace.TagComment, ""))
		defer t.Finish(&err)
	}
	if err = db.breaker.Allow(); err != nil {
		stats.Incr("tidb:ping", "breaker")
		return
	}
	_, c, cancel := db.conf.ExecTimeout.Shrink(c)
	err = db.PingContext(c)
	cancel()
	db.onBreaker(&err)
	stats.Timing("tidb:ping", int64(time.Since(now)/time.Millisecond))
	if err != nil {
		err = errors.WithStack(err)
	}
	return
}

func (db *conn) prepare(query string) (*Stmt, error) {
	stmt, err := db.Prepare(query)
	if err != nil {
		err = errors.Wrapf(err, "addr: %s prepare %s", db.addr, query)
		return nil, err
	}
	st := &Stmt{query: query, db: db}
	st.stmt.Store(stmt)
	return st, nil
}

func (db *conn) prepared(query string) (stmt *Stmt) {
	stmt = &Stmt{query: query, db: db}
	s, err := db.Prepare(query)
	if err == nil {
		stmt.stmt.Store(s)
		return
	}
	return
}

func (db *conn) query(c context.Context, query string, args ...interface{}) (rows *Rows, err error) {
	now := time.Now()
	if t, ok := trace.FromContext(c); ok {
		t = t.Fork(_family, "query")
		t.SetTag(trace.String(trace.TagAddress, db.addr), trace.String(trace.TagComment, query))
		defer t.Finish(&err)
	}
	if err = db.breaker.Allow(); err != nil {
		stats.Incr("tidb:query", "breaker")
		return
	}
	_, c, cancel := db.conf.QueryTimeout.Shrink(c)
	rs, err := db.DB.QueryContext(c, query, args...)
	db.onBreaker(&err)
	stats.Timing("tidb:query", int64(time.Since(now)/time.Millisecond))
	if err != nil {
		err = errors.Wrapf(err, "addr: %s, query:%s, args:%+v", db.addr, query, args)
		cancel()
		return
	}
	rows = &Rows{Rows: rs, cancel: cancel}
	return
}

func (db *conn) queryRow(c context.Context, query string, args ...interface{}) *Row {
	now := time.Now()
	t, ok := trace.FromContext(c)
	if ok {
		t = t.Fork(_family, "queryrow")
		t.SetTag(trace.String(trace.TagAddress, db.addr), trace.String(trace.TagComment, query))
	}
	if err := db.breaker.Allow(); err != nil {
		stats.Incr("tidb:queryrow", "breaker")
		return &Row{db: db, t: t, err: err}
	}
	_, c, cancel := db.conf.QueryTimeout.Shrink(c)
	r := db.DB.QueryRowContext(c, query, args...)
	stats.Timing("tidb:queryrow", int64(time.Since(now)/time.Millisecond))
	return &Row{db: db, Row: r, query: query, args: args, t: t, cancel: cancel}
}

// Close closes the statement.
func (s *Stmt) Close() (err error) {
	stmt, ok := s.stmt.Load().(*sql.Stmt)
	if ok {
		err = errors.WithStack(stmt.Close())
	}
	return
}

func (s *Stmt) prepare() (st *sql.Stmt) {
	var ok bool
	if st, ok = s.stmt.Load().(*sql.Stmt); ok {
		return
	}
	var err error
	if st, err = s.db.Prepare(s.query); err == nil {
		s.stmt.Store(st)
	}
	return
}

// Exec executes a prepared statement with the given arguments and returns a
// Result summarizing the effect of the statement.
func (s *Stmt) Exec(c context.Context, args ...interface{}) (res sql.Result, err error) {
	now := time.Now()
	if s.tx {
		if s.t != nil {
			s.t.SetTag(trace.String(trace.TagAnnotation, s.query))
		}
	} else if t, ok := trace.FromContext(c); ok {
		t = t.Fork(_family, "exec")
		t.SetTag(trace.String(trace.TagAddress, s.db.addr), trace.String(trace.TagComment, s.query))
		defer t.Finish(&err)
	}
	if err = s.db.breaker.Allow(); err != nil {
		stats.Incr("tidb:stmt:exec", "breaker")
		return
	}
	stmt := s.prepare()
	if stmt == nil {
		err = ErrStmtNil
		return
	}
	_, c, cancel := s.db.conf.ExecTimeout.Shrink(c)
	res, err = stmt.ExecContext(c, args...)
	cancel()
	s.db.onBreaker(&err)
	stats.Timing("tidb:stmt:exec", int64(time.Since(now)/time.Millisecond))
	if err != nil {
		err = errors.Wrapf(err, "addr: %s exec:%s, args:%+v", s.db.addr, s.query, args)
	}
	return
}

// Query executes a prepared query statement with the given arguments and
// returns the query results as a *Rows.
func (s *Stmt) Query(c context.Context, args ...interface{}) (rows *Rows, err error) {
	if s.tx {
		if s.t != nil {
			s.t.SetTag(trace.String(trace.TagAnnotation, s.query))
		}
	} else if t, ok := trace.FromContext(c); ok {
		t = t.Fork(_family, "query")
		t.SetTag(trace.String(trace.TagAddress, s.db.addr), trace.String(trace.TagComment, s.query))
		defer t.Finish(&err)
	}
	if err = s.db.breaker.Allow(); err != nil {
		stats.Incr("tidb:stmt:query", "breaker")
		return
	}
	stmt := s.prepare()
	if stmt == nil {
		err = ErrStmtNil
		return
	}
	now := time.Now()
	_, c, cancel := s.db.conf.QueryTimeout.Shrink(c)
	rs, err := stmt.QueryContext(c, args...)
	s.db.onBreaker(&err)
	stats.Timing("tidb:stmt:query", int64(time.Since(now)/time.Millisecond))
	if err != nil {
		err = errors.Wrapf(err, "addr: %s, query:%s, args:%+v", s.db.addr, s.query, args)
		cancel()
		return
	}
	rows = &Rows{Rows: rs, cancel: cancel}
	return
}

// QueryRow executes a prepared query statement with the given arguments.
// If an error occurs during the execution of the statement, that error will
// be returned by a call to Scan on the returned *Row, which is always non-nil.
// If the query selects no rows, the *Row's Scan will return ErrNoRows.
// Otherwise, the *Row's Scan scans the first selected row and discards the rest.
func (s *Stmt) QueryRow(c context.Context, args ...interface{}) (row *Row) {
	now := time.Now()
	row = &Row{db: s.db, query: s.query, args: args}
	if s.tx {
		if s.t != nil {
			s.t.SetTag(trace.String(trace.TagAnnotation, s.query))
		}
	} else if t, ok := trace.FromContext(c); ok {
		t = t.Fork(_family, "queryrow")
		t.SetTag(trace.String(trace.TagAddress, s.db.addr), trace.String(trace.TagComment, s.query))
		row.t = t
	}
	if row.err = s.db.breaker.Allow(); row.err != nil {
		stats.Incr("tidb:stmt:queryrow", "breaker")
		return
	}
	stmt := s.prepare()
	if stmt == nil {
		return
	}
	_, c, cancel := s.db.conf.QueryTimeout.Shrink(c)
	row.Row = stmt.QueryRowContext(c, args...)
	row.cancel = cancel
	stats.Timing("tidb:stmt:queryrow", int64(time.Since(now)/time.Millisecond))
	return
}

func (s *Stmts) prepare(conn *conn) (st *Stmt) {
	if conn == nil {
		conn = s.db.conn()
	}
	s.mu.RLock()
	st = s.sts[conn.addr]
	s.mu.RUnlock()
	if st == nil {
		st = conn.prepared(s.query)
		s.mu.Lock()
		s.sts[conn.addr] = st
		s.mu.Unlock()
	}
	return
}

// Exec executes a prepared statement with the given arguments and returns a
// Result summarizing the effect of the statement.
func (s *Stmts) Exec(c context.Context, args ...interface{}) (res sql.Result, err error) {
	return s.prepare(nil).Exec(c, args...)
}

// Query executes a prepared query statement with the given arguments and
// returns the query results as a *Rows.
func (s *Stmts) Query(c context.Context, args ...interface{}) (rows *Rows, err error) {
	return s.prepare(nil).Query(c, args...)
}

// QueryRow executes a prepared query statement with the given arguments.
// If an error occurs during the execution of the statement, that error will
// be returned by a call to Scan on the returned *Row, which is always non-nil.
// If the query selects no rows, the *Row's Scan will return ErrNoRows.
// Otherwise, the *Row's Scan scans the first selected row and discards the rest.
func (s *Stmts) QueryRow(c context.Context, args ...interface{}) (row *Row) {
	return s.prepare(nil).QueryRow(c, args...)
}

// Close closes the statement.
func (s *Stmts) Close() (err error) {
	for _, st := range s.sts {
		if err = errors.WithStack(st.Close()); err != nil {
			return
		}
	}
	return
}

// Commit commits the transaction.
func (tx *Tx) Commit() (err error) {
	err = tx.tx.Commit()
	tx.cancel()
	tx.db.onBreaker(&err)
	if tx.t != nil {
		tx.t.Finish(&err)
	}
	if err != nil {
		err = errors.WithStack(err)
	}
	return
}

// Rollback aborts the transaction.
func (tx *Tx) Rollback() (err error) {
	err = tx.tx.Rollback()
	tx.cancel()
	tx.db.onBreaker(&err)
	if tx.t != nil {
		tx.t.Finish(&err)
	}
	if err != nil {
		err = errors.WithStack(err)
	}
	return
}

// Exec executes a query that doesn't return rows. For example: an INSERT and
// UPDATE.
func (tx *Tx) Exec(query string, args ...interface{}) (res sql.Result, err error) {
	now := time.Now()
	if tx.t != nil {
		tx.t.SetTag(trace.String(trace.TagAnnotation, fmt.Sprintf("exec %s", query)))
	}
	res, err = tx.tx.ExecContext(tx.c, query, args...)
	stats.Timing("tidb:tx:exec", int64(time.Since(now)/time.Millisecond))
	if err != nil {
		err = errors.Wrapf(err, "addr: %s exec:%s, args:%+v", tx.db.addr, query, args)
	}
	return
}

// Query executes a query that returns rows, typically a SELECT.
func (tx *Tx) Query(query string, args ...interface{}) (rows *Rows, err error) {
	if tx.t != nil {
		tx.t.SetTag(trace.String(trace.TagAnnotation, fmt.Sprintf("query %s", query)))
	}
	now := time.Now()
	defer func() {
		stats.Timing("tidb:tx:query", int64(time.Since(now)/time.Millisecond))
	}()
	rs, err := tx.tx.QueryContext(tx.c, query, args...)
	if err == nil {
		rows = &Rows{Rows: rs}
	} else {
		err = errors.Wrapf(err, "addr: %s, query:%s, args:%+v", tx.db.addr, query, args)
	}
	return
}

// QueryRow executes a query that is expected to return at most one row.
// QueryRow always returns a non-nil value. Errors are deferred until Row's
// Scan method is called.
func (tx *Tx) QueryRow(query string, args ...interface{}) *Row {
	if tx.t != nil {
		tx.t.SetTag(trace.String(trace.TagAnnotation, fmt.Sprintf("queryrow %s", query)))
	}
	now := time.Now()
	defer func() {
		stats.Timing("tidb:tx:queryrow", int64(time.Since(now)/time.Millisecond))
	}()
	r := tx.tx.QueryRowContext(tx.c, query, args...)
	return &Row{Row: r, db: tx.db, query: query, args: args}
}

// Stmt returns a transaction-specific prepared statement from an existing statement.
func (tx *Tx) Stmt(stmt *Stmt) *Stmt {
	if stmt == nil {
		return nil
	}
	as, ok := stmt.stmt.Load().(*sql.Stmt)
	if !ok {
		return nil
	}
	ts := tx.tx.StmtContext(tx.c, as)
	st := &Stmt{query: stmt.query, tx: true, t: tx.t, db: tx.db}
	st.stmt.Store(ts)
	return st
}

// Stmts returns a transaction-specific prepared statement from an existing statement.
func (tx *Tx) Stmts(stmt *Stmts) *Stmt {
	return tx.Stmt(stmt.prepare(tx.db))
}

// Prepare creates a prepared statement for use within a transaction.
// The returned statement operates within the transaction and can no longer be
// used once the transaction has been committed or rolled back.
// To use an existing prepared statement on this transaction, see Tx.Stmt.
func (tx *Tx) Prepare(query string) (*Stmt, error) {
	if tx.t != nil {
		tx.t.SetTag(trace.String(trace.TagAnnotation, fmt.Sprintf("prepare %s", query)))
	}
	stmt, err := tx.tx.Prepare(query)
	if err != nil {
		err = errors.Wrapf(err, "addr: %s prepare %s", tx.db.addr, query)
		return nil, err
	}
	st := &Stmt{query: query, tx: true, t: tx.t, db: tx.db}
	st.stmt.Store(stmt)
	return st, nil
}

// parseDSNAddr parse dsn name and return addr.
func parseDSNAddr(dsn string) (addr string) {
	if dsn == "" {
		return
	}
	cfg, err := mysql.ParseDSN(dsn)
	if err != nil {
		return
	}
	addr = cfg.Addr
	return
}

func genDSN(dsn, addr string) (res string) {
	cfg, err := mysql.ParseDSN(dsn)
	if err != nil {
		return
	}
	cfg.Addr = addr
	cfg.Net = "tcp"
	res = cfg.FormatDSN()
	return
}
