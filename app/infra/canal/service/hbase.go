package service

import (
	"context"
	"fmt"
	"io"
	"time"

	"go-common/app/infra/canal/conf"
	"go-common/app/infra/canal/model"
	"go-common/library/log"
	"go-common/library/queue/databus"

	hbase "go-common/library/database/hbase.v2"

	"github.com/pkg/errors"
	"github.com/tsuna/gohbase/hrpc"
)

// HBaseInstance hbase canal instance
type HBaseInstance struct {
	config *conf.HBaseInsConf
	client *hbase.Client
	// one instance can have lots of target different by schema and table
	targets []*databus.Databus
	hinfo   *hbaseInfo

	// scan latest timestamp, be used check delay
	latestTimestamp int64

	err    error
	closed bool
}

// NewHBaseInstance new canal instance
func NewHBaseInstance(c *conf.HBaseInsConf) (ins *HBaseInstance, err error) {
	// new instance
	ins = &HBaseInstance{config: c}
	// check and modify config
	if c.MasterInfo == nil {
		err = errors.New("no masterinfo config")
		ins.err = err
		return
	}
	// new hbase info
	if ins.hinfo, err = newHBaseInfo(ins.config.Cluster, ins.config.MasterInfo); err != nil {
		log.Error("init hbase info error(%v)", err)
		ins.err = errors.Wrap(err, "init hbase info")
		return
	}
	// new hbase client
	ins.client = hbase.NewClient(&hbase.Config{Zookeeper: &hbase.ZKConfig{
		Root:  c.Root,
		Addrs: c.Addrs,
	}})
	return
}

// Start start binlog receive
func (ins *HBaseInstance) Start() {
	for _, db := range ins.config.Databases {
		databus := databus.New(db.Databus)
		ins.targets = append(ins.targets, databus)
		for _, tb := range db.Tables {
			go ins.scanproc(db, databus, tb)
		}
	}
}

func (ins *HBaseInstance) scanproc(c *conf.HBaseDatabase, databus *databus.Databus, table *conf.HBaseTable) {
	latestTs := ins.hinfo.LatestTs(table.Name)
	if latestTs == 0 {
		latestTs = uint64(time.Now().UnixNano() / 1e6)
	}
AGAIN:
	for {
		time.Sleep(300 * time.Millisecond)
		nowTs := uint64(time.Now().UnixNano() / 1e6)
		// scan
		scanner, err := ins.client.Scan(context.TODO(), []byte(table.Name), hrpc.TimeRangeUint64(latestTs, nowTs))
		if err != nil {
			time.Sleep(time.Second)
			continue
		}
		ins.latestTimestamp = time.Now().Unix()
		for {
			res, err := scanner.Next()
			if err == io.EOF {
				latestTs = nowTs
				ins.hinfo.Save(table.Name, latestTs)
				continue AGAIN
			} else if err != nil {
				time.Sleep(time.Second)
				continue AGAIN
			}
			var key string
			data := map[string]interface{}{}
			for _, cell := range res.Cells {
				row := string(cell.Row)
				if key == "" {
					key = row
				}
				if _, ok := data[row]; !ok {
					data[row] = map[string]string{}
				}
				column := string(cell.Qualifier)
				omit := false
				for _, field := range table.OmitField {
					if field == column {
						omit = true
						break
					}
				}
				if !omit {
					cm := data[row].(map[string]string)
					cm[column] = string(cell.Value)
				}
			}
			real := &model.Data{Table: table.Name, New: data}
			databus.Send(context.TODO(), key, real)
			log.Info("hbase databus:group(%s)topic(%s) pub(key:%s, value:%+v) succeed", c.Databus.Group, c.Databus.Topic, key, real)
		}
	}
}

// Close close instance
func (ins *HBaseInstance) Close() {
	if ins.err != nil {
		return
	}
	ins.client.Close()
	for _, t := range ins.targets {
		t.Close()
	}
	ins.closed = true
	ins.err = errors.New("canal hbase instance closed")
}

func (ins *HBaseInstance) String() string {
	return ins.config.Cluster
}

// Error returns instance error.
func (ins *HBaseInstance) Error() string {
	if ins.err == nil {
		return ""
	}
	return fmt.Sprintf("+%v", ins.err)
}

func (ins *HBaseInstance) delay() int64 {
	return time.Now().Unix() - ins.latestTimestamp
}
