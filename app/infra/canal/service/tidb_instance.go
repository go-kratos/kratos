package service

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"go-common/app/infra/canal/conf"
	"go-common/app/infra/canal/infoc"
	"go-common/app/infra/canal/model"
	"go-common/app/infra/canal/service/reader"
	"go-common/library/log"
	"go-common/library/queue/databus"

	pb "github.com/pingcap/tidb-tools/tidb_binlog/slave_binlog_proto/go-binlog"
)

type tidbInstance struct {
	canal  *Canal
	config *conf.TiDBInsConf
	// one instance can have lots of target different by schema and table
	targets      map[string][]producer
	latestOffset int64
	// scan latest timestamp, be used check delay
	latestTimestamp int64
	reader          *reader.Reader
	err             error
	closed          bool
	tables          map[string]map[string]*Table
	ignoreTables    map[string]map[string]bool
	waitConsume     sync.WaitGroup
	waitTable       sync.WaitGroup
}

// Table db table
type Table struct {
	PrimaryKey []string        // kafka msg key
	OmitField  map[string]bool // field will be ignored in table
	OmitAction map[string]bool // action will be ignored in table
	name       string
	ch         chan *msg
}

type msg struct {
	db          string
	table       string
	tableRegexp string
	mu          *pb.TableMutation
	ignore      map[string]bool
	keys        []string
	columns     []*pb.ColumnInfo
}

// newTiDBInstance new canal instance
func newTiDBInstance(cl *Canal, c *conf.TiDBInsConf) (ins *tidbInstance, err error) {
	// new instance
	ins = &tidbInstance{
		config:       c,
		canal:        cl,
		targets:      make(map[string][]producer, len(c.Databases)),
		tables:       make(map[string]map[string]*Table),
		ignoreTables: make(map[string]map[string]bool),
	}
	cfg := &reader.Config{
		Name:      c.Name,
		KafkaAddr: c.Addrs,
		Offset:    c.Offset,
		CommitTS:  c.CommitTS,
		ClusterID: c.ClusterID,
	}
	if err = ins.check(); err != nil {
		return
	}
	position, err := cl.dao.TiDBPosition(context.Background(), c.Name)
	if err == nil && position != nil {
		cfg.Offset = position.Offset
		cfg.CommitTS = position.CommitTS
	}
	ins.latestOffset = cfg.Offset
	ins.latestTimestamp = cfg.CommitTS
	for _, db := range c.Databases {
		if db.Databus != nil {
			if ins.targets == nil {
				ins.targets = make(map[string][]producer)
			}
			ins.targets[db.Schema] = append(ins.targets[db.Schema], &databusP{group: db.Databus.Group, topic: db.Databus.Topic, Databus: databus.New(db.Databus)})
		}
		if db.Infoc != nil {
			ins.targets[db.Schema] = append(ins.targets[db.Schema], &infocP{taskID: db.Infoc.TaskID, Infoc: infoc.New(db.Infoc)})
		}
	}
	ins.reader, ins.err = reader.NewReader(cfg)
	return ins, ins.err
}

// start start binlog receive
func (ins *tidbInstance) start() {
	defer ins.waitConsume.Done()
	ins.waitConsume.Add(2)
	go ins.reader.Run()
	go ins.syncproc()
	for msg := range ins.reader.Messages() {
		ins.process(msg)
	}
}

// close close instance
func (ins *tidbInstance) close() {
	if ins.closed {
		return
	}
	ins.closed = true
	ins.reader.Close()
	ins.waitConsume.Wait()
	for _, tables := range ins.tables {
		for _, table := range tables {
			close(table.ch)
		}
	}
	ins.waitTable.Wait()
	ins.sync()
}

// String .
func (ins *tidbInstance) String() string {
	return fmt.Sprintf("%s-%s", ins.config.Name, ins.config.ClusterID)
}

func (ins *tidbInstance) process(m *reader.Message) (err error) {
	if m.Binlog.Type == pb.BinlogType_DDL {
		log.Info("tidb %s got ddl: %s", ins.String(), m.Binlog.DdlData.String())
		return
	}
	for _, table := range m.Binlog.DmlData.Tables {
		tb := ins.getTable(table.GetSchemaName(), table.GetTableName())
		if tb == nil {
			continue
		}
		for _, mu := range table.Mutations {
			action := strings.ToLower(mu.GetType().String())
			if tb.OmitAction[action] {
				continue
			}
			tb.ch <- &msg{
				db:          table.GetSchemaName(),
				table:       table.GetTableName(),
				mu:          mu,
				ignore:      tb.OmitField,
				keys:        tb.PrimaryKey,
				columns:     table.ColumnInfo,
				tableRegexp: tb.name,
			}
			if stats != nil {
				stats.Incr("syncer_counter", ins.String(), table.GetSchemaName(), tb.name, action)
				stats.State("delay_syncer", ins.delay(), ins.String(), tb.name, "", "")
			}
		}
	}
	ins.latestOffset = m.Offset
	ins.latestTimestamp = m.Binlog.CommitTs
	return nil
}

// Error returns instance error.
func (ins *tidbInstance) Error() string {
	if ins.err == nil {
		return ""
	}
	return fmt.Sprintf("+%v", ins.err)
}

func (ins *tidbInstance) delay() int64 {
	return time.Now().Unix() - ins.latestTimestamp
}

func (ins *tidbInstance) sync() {
	info := &model.TiDBInfo{
		Name:      ins.config.Name,
		ClusterID: ins.config.ClusterID,
		Offset:    ins.latestOffset,
		CommitTS:  ins.latestTimestamp,
	}
	ins.canal.dao.UpdateTiDBPosition(context.Background(), info)
}
