package service

import (
	"fmt"
	"time"

	"go-common/app/infra/canal/conf"
	"go-common/library/log"

	"github.com/pkg/errors"
	"github.com/siddontang/go-mysql/canal"
	"github.com/siddontang/go-mysql/mysql"
	"github.com/siddontang/go-mysql/replication"
)

// Instance canal instance
type Instance struct {
	*canal.Canal
	config *conf.InsConf
	// one instance can have lots of target different by schema and table
	targets []*Target
	master  *dbMasterInfo

	// binlog latest timestamp, be used check delay
	latestTimestamp int64
	err             error
	closed          bool
}

// NewInstance new canal instance
func NewInstance(c *conf.InsConf) (ins *Instance, err error) {
	// new instance
	ins = &Instance{config: c}
	// check and modify config
	if c.MasterInfo == nil {
		err = errors.New("no masterinfo config")
		ins.err = err
		return
	}
	if c.ReadTimeout == 0 {
		c.ReadTimeout = 90
	}
	c.ReadTimeout = c.ReadTimeout * time.Second
	c.HeartbeatPeriod = c.HeartbeatPeriod * time.Second

	if ins.master, err = newDBMasterInfo(ins.config.Addr, ins.config.MasterInfo); err != nil {
		log.Error("init master info error(%v)", err)
		ins.err = errors.Wrap(err, "init master info")
		return
	}
	if ins.targets, err = newTargets(c); err != nil {
		log.Error("db.init() error(%v)", err)
		ins.err = errors.Wrap(err, "db init")
		return
	}
	ins.latestTimestamp = time.Now().Unix()
	// new canal
	if ins.Canal, err = canal.NewCanal(c.Config); err != nil {
		log.Error("canal.NewCanal(%v) error(%v)", c.Config, err)
		ins.err = errors.Wrapf(err, "canal NewCanal(%v)", c.Config)
		return
	}
	// implement self as canal's event handler
	ins.Canal.SetEventHandler(ins)

	return
}

func newTargets(c *conf.InsConf) (targets []*Target, err error) {
	targets = make([]*Target, 0, len(c.Databases))
	for _, db := range c.Databases {
		if err = db.CheckTable(c.Addr, c.User, c.Password); err != nil {
			log.Error("db.CheckTable() error(%v)", err)
			return
		}
		targets = append(targets, NewTarget(db))
	}
	return
}

// Start start binlog receive
func (ins *Instance) Start() {
	pos := ins.master.Pos()
	if pos.Name == "" || pos.Pos == 0 {
		var err error
		if pos, err = ins.Canal.GetMasterPos(); err != nil {
			log.Error("c.MasterPos error(%v)", err)
			ins.err = errors.Wrap(err, "canal get master pos when start")
			return
		}
	}
	ins.err = ins.Canal.RunFrom(pos)
}

// Close close instance
func (ins *Instance) Close() {
	if ins.err != nil {
		return
	}
	ins.Canal.Close()
	for _, t := range ins.targets {
		t.close()
	}
	ins.closed = true
	ins.err = errors.New("canal closed")
}

// Check filter row event
func (ins *Instance) Check(ev *canal.RowsEvent) (ts []*Target) {
	for _, t := range ins.targets {
		if t.compare(ev.Table.Schema, ev.Table.Name, ev.Action) {
			ts = append(ts, t)
		}
	}
	return
}

func (ins *Instance) String() string {
	return ins.config.Addr
}

// OnRotate OnRotate
func (ins *Instance) OnRotate(re *replication.RotateEvent) error {
	log.Info("OnRotate binlog addr(%s) rotate binname(%s) pos(%d)", ins.config.Addr, re.NextLogName, re.Position)
	return nil
}

// OnDDL OnDDL
func (ins *Instance) OnDDL(pos mysql.Position, qe *replication.QueryEvent) error {
	log.Info("OnDDL binlog addr(%s) ddl binname(%s) pos(%d)", ins.config.Addr, pos.Name, pos.Pos)
	return nil
}

// OnXID OnXID
func (ins *Instance) OnXID(mysql.Position) error {
	return nil
}

//OnGTID OnGTID
func (ins *Instance) OnGTID(mysql.GTIDSet) error {
	return nil
}

// OnPosSynced OnPosSynced
func (ins *Instance) OnPosSynced(pos mysql.Position, force bool) error {
	return ins.master.Save(pos, force)
}

// OnRow send the envent to table
func (ins *Instance) OnRow(ev *canal.RowsEvent) error {
	for _, t := range ins.Check(ev) {
		t.send(ev)
	}
	if stats != nil {
		stats.Incr("syncer_counter", ins.String(), ev.Table.Schema, tblReplacer.ReplaceAllString(ev.Table.Name, ""), ev.Action)
		stats.State("delay_syncer", ins.delay(), ins.String(), ev.Table.Schema, "", "")
	}
	ins.latestTimestamp = time.Now().Unix()
	return nil
}

// Error returns instance error.
func (ins *Instance) Error() string {
	if ins.err == nil {
		return ""
	}
	return fmt.Sprintf("+%v", ins.err)
}

func (ins *Instance) delay() int64 {
	return time.Now().Unix() - ins.latestTimestamp
}
