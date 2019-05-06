package tidb

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/bilibili/kratos/pkg/conf/env"
	"github.com/bilibili/kratos/pkg/log"
	"github.com/bilibili/kratos/pkg/naming"
	"github.com/bilibili/kratos/pkg/naming/discovery"
)

var _schema = "tidb://"

func (db *DB) nodeList() (nodes []string) {
	var (
		insZone *naming.InstancesInfo
		ins     []*naming.Instance
		ok      bool
	)
	if insZone, ok = db.dis.Fetch(context.Background()); !ok {
		return
	}
	if ins, ok = insZone.Instances[env.Zone]; !ok || len(ins) == 0 {
		return
	}
	for _, in := range ins {
		for _, addr := range in.Addrs {
			if strings.HasPrefix(addr, _schema) {
				addr = strings.Replace(addr, _schema, "", -1)
				nodes = append(nodes, addr)
			}
		}
	}
	log.Info("tidb get %s instances(%v)", db.appid, nodes)
	return
}

func (db *DB) disc() (nodes []string) {
	db.dis = discovery.Build(db.appid)
	e := db.dis.Watch()
	select {
	case <-e:
		nodes = db.nodeList()
	case <-time.After(10 * time.Second):
		panic("tidb init discovery err")
	}
	if len(nodes) == 0 {
		panic(fmt.Sprintf("tidb %s no instance", db.appid))
	}
	go db.nodeproc(e)
	log.Info("init tidb discvoery info successfully")
	return
}
