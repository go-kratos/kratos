package canal

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/juju/errors"
	"github.com/siddontang/go-mysql/client"
	"github.com/siddontang/go-mysql/dump"
	"github.com/siddontang/go-mysql/mysql"
	"github.com/siddontang/go-mysql/replication"
	"github.com/siddontang/go-mysql/schema"
	log "github.com/sirupsen/logrus"
)

// Canal can sync your MySQL data into everywhere, like Elasticsearch, Redis, etc...
// MySQL must open row format for binlog
type Canal struct {
	m sync.Mutex

	cfg *Config

	useGTID bool

	master     *masterInfo
	dumper     *dump.Dumper
	dumpDoneCh chan struct{}
	syncer     *replication.BinlogSyncer

	eventHandler EventHandler

	connLock sync.Mutex
	conn     *client.Conn

	tableLock          sync.RWMutex
	tables             map[string]*schema.Table
	errorTablesGetTime map[string]time.Time

	tableMatchCache   map[string]bool
	includeTableRegex []*regexp.Regexp
	excludeTableRegex []*regexp.Regexp

	ctx    context.Context
	cancel context.CancelFunc
}

// canal will retry fetching unknown table's meta after UnknownTableRetryPeriod
var UnknownTableRetryPeriod = time.Second * time.Duration(10)
var ErrExcludedTable = errors.New("excluded table meta")

func NewCanal(cfg *Config) (*Canal, error) {
	c := new(Canal)
	c.cfg = cfg

	c.ctx, c.cancel = context.WithCancel(context.Background())

	c.dumpDoneCh = make(chan struct{})
	c.eventHandler = &DummyEventHandler{}

	c.tables = make(map[string]*schema.Table)
	if c.cfg.DiscardNoMetaRowEvent {
		c.errorTablesGetTime = make(map[string]time.Time)
	}
	c.master = &masterInfo{}

	var err error

	if err = c.prepareDumper(); err != nil {
		return nil, errors.Trace(err)
	}

	if err = c.prepareSyncer(); err != nil {
		return nil, errors.Trace(err)
	}

	if err := c.checkBinlogRowFormat(); err != nil {
		return nil, errors.Trace(err)
	}

	// init table filter
	if n := len(c.cfg.IncludeTableRegex); n > 0 {
		c.includeTableRegex = make([]*regexp.Regexp, n)
		for i, val := range c.cfg.IncludeTableRegex {
			reg, err := regexp.Compile(val)
			if err != nil {
				return nil, errors.Trace(err)
			}
			c.includeTableRegex[i] = reg
		}
	}

	if n := len(c.cfg.ExcludeTableRegex); n > 0 {
		c.excludeTableRegex = make([]*regexp.Regexp, n)
		for i, val := range c.cfg.ExcludeTableRegex {
			reg, err := regexp.Compile(val)
			if err != nil {
				return nil, errors.Trace(err)
			}
			c.excludeTableRegex[i] = reg
		}
	}

	if c.includeTableRegex != nil || c.excludeTableRegex != nil {
		c.tableMatchCache = make(map[string]bool)
	}

	return c, nil
}

func (c *Canal) prepareDumper() error {
	var err error
	dumpPath := c.cfg.Dump.ExecutionPath
	if len(dumpPath) == 0 {
		// ignore mysqldump, use binlog only
		return nil
	}

	if c.dumper, err = dump.NewDumper(dumpPath,
		c.cfg.Addr, c.cfg.User, c.cfg.Password); err != nil {
		return errors.Trace(err)
	}

	if c.dumper == nil {
		//no mysqldump, use binlog only
		return nil
	}

	dbs := c.cfg.Dump.Databases
	tables := c.cfg.Dump.Tables
	tableDB := c.cfg.Dump.TableDB

	if len(tables) == 0 {
		c.dumper.AddDatabases(dbs...)
	} else {
		c.dumper.AddTables(tableDB, tables...)
	}

	charset := c.cfg.Charset
	c.dumper.SetCharset(charset)

	c.dumper.SkipMasterData(c.cfg.Dump.SkipMasterData)
	c.dumper.SetMaxAllowedPacket(c.cfg.Dump.MaxAllowedPacketMB)

	for _, ignoreTable := range c.cfg.Dump.IgnoreTables {
		if seps := strings.Split(ignoreTable, ","); len(seps) == 2 {
			c.dumper.AddIgnoreTables(seps[0], seps[1])
		}
	}

	if c.cfg.Dump.DiscardErr {
		c.dumper.SetErrOut(ioutil.Discard)
	} else {
		c.dumper.SetErrOut(os.Stderr)
	}

	return nil
}

// Run will first try to dump all data from MySQL master `mysqldump`,
// then sync from the binlog position in the dump data.
// It will run forever until meeting an error or Canal closed.
func (c *Canal) Run() error {
	return c.run()
}

// RunFrom will sync from the binlog position directly, ignore mysqldump.
func (c *Canal) RunFrom(pos mysql.Position) error {
	c.useGTID = false
	c.master.Update(pos)

	return c.Run()
}

func (c *Canal) StartFromGTID(set mysql.GTIDSet) error {
	c.useGTID = true
	c.master.UpdateGTID(set)

	return c.Run()
}

func (c *Canal) run() error {
	defer func() {
		c.cancel()
	}()

	err := c.tryDump()
	close(c.dumpDoneCh)

	if err != nil {
		log.Errorf("canal dump mysql err: %v", err)
		return errors.Trace(err)
	}

	if err = c.runSyncBinlog(); err != nil {
		log.Errorf("canal start sync binlog err: %v", err)
		return errors.Trace(err)
	}

	return nil
}

func (c *Canal) Close() {
	log.Infof("closing canal")

	c.m.Lock()
	defer c.m.Unlock()

	c.cancel()
	c.connLock.Lock()
	c.conn.Close()
	c.conn = nil
	c.connLock.Unlock()
	c.syncer.Close()

	c.eventHandler.OnPosSynced(c.master.Position(), true)
}

func (c *Canal) WaitDumpDone() <-chan struct{} {
	return c.dumpDoneCh
}

func (c *Canal) Ctx() context.Context {
	return c.ctx
}

func (c *Canal) checkTableMatch(key string) bool {
	// no filter, return true
	if c.tableMatchCache == nil {
		return true
	}

	c.tableLock.RLock()
	rst, ok := c.tableMatchCache[key]
	c.tableLock.RUnlock()
	if ok {
		// cache hit
		return rst
	}
	matchFlag := false
	// check include
	if c.includeTableRegex != nil {
		for _, reg := range c.includeTableRegex {
			if reg.MatchString(key) {
				matchFlag = true
				break
			}
		}
	}
	// check exclude
	if matchFlag && c.excludeTableRegex != nil {
		for _, reg := range c.excludeTableRegex {
			if reg.MatchString(key) {
				matchFlag = false
				break
			}
		}
	}
	c.tableLock.Lock()
	c.tableMatchCache[key] = matchFlag
	c.tableLock.Unlock()
	return matchFlag
}

func (c *Canal) GetTable(db string, table string) (*schema.Table, error) {
	key := fmt.Sprintf("%s.%s", db, table)
	// if table is excluded, return error and skip parsing event or dump
	if !c.checkTableMatch(key) {
		return nil, ErrExcludedTable
	}
	c.tableLock.RLock()
	t, ok := c.tables[key]
	c.tableLock.RUnlock()

	if ok {
		return t, nil
	}

	if c.cfg.DiscardNoMetaRowEvent {
		c.tableLock.RLock()
		lastTime, ok := c.errorTablesGetTime[key]
		c.tableLock.RUnlock()
		if ok && time.Now().Sub(lastTime) < UnknownTableRetryPeriod {
			return nil, schema.ErrMissingTableMeta
		}
	}

	t, err := schema.NewTable(c, db, table)
	if err != nil {
		// check table not exists
		if ok, err1 := schema.IsTableExist(c, db, table); err1 == nil && !ok {
			return nil, schema.ErrTableNotExist
		}
		// work around : RDS HAHeartBeat
		// ref : https://github.com/alibaba/canal/blob/master/parse/src/main/java/com/alibaba/otter/canal/parse/inbound/mysql/dbsync/LogEventConvert.java#L385
		// issue : https://github.com/alibaba/canal/issues/222
		// This is a common error in RDS that canal can't get HAHealthCheckSchema's meta, so we mock a table meta.
		// If canal just skip and log error, as RDS HA heartbeat interval is very short, so too many HAHeartBeat errors will be logged.
		if key == schema.HAHealthCheckSchema {
			// mock ha_health_check meta
			ta := &schema.Table{
				Schema:  db,
				Name:    table,
				Columns: make([]schema.TableColumn, 0, 2),
				Indexes: make([]*schema.Index, 0),
			}
			ta.AddColumn("id", "bigint(20)", "", "")
			ta.AddColumn("type", "char(1)", "", "")
			c.tableLock.Lock()
			c.tables[key] = ta
			c.tableLock.Unlock()
			return ta, nil
		}
		// if DiscardNoMetaRowEvent is true, we just log this error
		if c.cfg.DiscardNoMetaRowEvent {
			c.tableLock.Lock()
			c.errorTablesGetTime[key] = time.Now()
			c.tableLock.Unlock()
			// log error and return ErrMissingTableMeta
			log.Errorf("canal get table meta err: %v", errors.Trace(err))
			return nil, schema.ErrMissingTableMeta
		}
		return nil, err
	}

	c.tableLock.Lock()
	c.tables[key] = t
	if c.cfg.DiscardNoMetaRowEvent {
		// if get table info success, delete this key from errorTablesGetTime
		delete(c.errorTablesGetTime, key)
	}
	c.tableLock.Unlock()

	return t, nil
}

// ClearTableCache clear table cache
func (c *Canal) ClearTableCache(db []byte, table []byte) {
	key := fmt.Sprintf("%s.%s", db, table)
	c.tableLock.Lock()
	delete(c.tables, key)
	if c.cfg.DiscardNoMetaRowEvent {
		delete(c.errorTablesGetTime, key)
	}
	c.tableLock.Unlock()
}

// Check MySQL binlog row image, must be in FULL, MINIMAL, NOBLOB
func (c *Canal) CheckBinlogRowImage(image string) error {
	// need to check MySQL binlog row image? full, minimal or noblob?
	// now only log
	if c.cfg.Flavor == mysql.MySQLFlavor {
		if res, err := c.Execute(`SHOW GLOBAL VARIABLES LIKE "binlog_row_image"`); err != nil {
			return errors.Trace(err)
		} else {
			// MySQL has binlog row image from 5.6, so older will return empty
			rowImage, _ := res.GetString(0, 1)
			if rowImage != "" && !strings.EqualFold(rowImage, image) {
				return errors.Errorf("MySQL uses %s binlog row image, but we want %s", rowImage, image)
			}
		}
	}

	return nil
}

func (c *Canal) checkBinlogRowFormat() error {
	res, err := c.Execute(`SHOW GLOBAL VARIABLES LIKE "binlog_format";`)
	if err != nil {
		return errors.Trace(err)
	} else if f, _ := res.GetString(0, 1); f != "ROW" {
		return errors.Errorf("binlog must ROW format, but %s now", f)
	}

	return nil
}

func (c *Canal) prepareSyncer() error {
	seps := strings.Split(c.cfg.Addr, ":")
	if len(seps) != 2 {
		return errors.Errorf("invalid mysql addr format %s, must host:port", c.cfg.Addr)
	}

	port, err := strconv.ParseUint(seps[1], 10, 16)
	if err != nil {
		return errors.Trace(err)
	}

	cfg := replication.BinlogSyncerConfig{
		ServerID:        c.cfg.ServerID,
		Flavor:          c.cfg.Flavor,
		Host:            seps[0],
		Port:            uint16(port),
		User:            c.cfg.User,
		Password:        c.cfg.Password,
		Charset:         c.cfg.Charset,
		HeartbeatPeriod: c.cfg.HeartbeatPeriod,
		ReadTimeout:     c.cfg.ReadTimeout,
		UseDecimal:      c.cfg.UseDecimal,
	}

	c.syncer = replication.NewBinlogSyncer(cfg)

	return nil
}

// Execute a SQL
func (c *Canal) Execute(cmd string, args ...interface{}) (rr *mysql.Result, err error) {
	c.connLock.Lock()
	defer c.connLock.Unlock()

	retryNum := 3
	for i := 0; i < retryNum; i++ {
		if c.conn == nil {
			c.conn, err = client.Connect(c.cfg.Addr, c.cfg.User, c.cfg.Password, "")
			if err != nil {
				return nil, errors.Trace(err)
			}
		}

		rr, err = c.conn.Execute(cmd, args...)
		if err != nil && !mysql.ErrorEqual(err, mysql.ErrBadConn) {
			return
		} else if mysql.ErrorEqual(err, mysql.ErrBadConn) {
			c.conn.Close()
			c.conn = nil
			continue
		} else {
			return
		}
	}
	return
}

func (c *Canal) SyncedPosition() mysql.Position {
	return c.master.Position()
}
