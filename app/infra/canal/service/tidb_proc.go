package service

import (
	"context"
	"fmt"
	"time"

	"go-common/app/infra/canal/conf"
	"go-common/library/log"
)

func (ins *tidbInstance) proc(ch chan *msg) {
	defer ins.waitTable.Done()
	for {
		msg, ok := <-ch
		if !ok {
			return
		}
		data, err := tidbMakeData(msg)
		if err != nil {
			log.Error("tidb MakeData(%+v) err: %+v", msg, err)
			continue
		}
		if data == nil {
			continue
		}
		for _, target := range ins.targets[msg.db] {
			for {
				if err = target.Send(context.TODO(), data.Key, data); err == nil {
					stats.Incr("send_counter", target.Name(), msg.db, msg.tableRegexp, data.Action)
					break
				}
				stats.Incr("retry_counter", target.Name(), msg.db, msg.tableRegexp, data.Action)
				log.Error("tidb %s scheme(%s) pub fail,add to retry", target.Name(), msg.db)
				time.Sleep(time.Second)
			}
			log.Info("tidb %s pub(key:%s, value:%+v) succeed", target.Name(), data.Key, data)
		}
	}
}

func (ins *tidbInstance) syncproc() {
	defer ins.waitConsume.Done()
	for {
		if ins.closed {
			return
		}
		time.Sleep(time.Second * 5)
		ins.sync()
	}
}

func (c *Canal) tidbEventproc() {
	ech := conf.TiDBEvent()
	for {
		insc := <-ech
		if insc == nil {
			continue
		}
		c.tidbInsl.Lock()
		if old, ok := c.tidbInstances[insc.Name]; ok {
			old.close()
		}
		c.tidbInsl.Unlock()
		ins, err := newTiDBInstance(c, insc)
		if err != nil {
			log.Error("new instance error(%v)", err)
			c.sendWx(fmt.Sprintf("reload tidb canal instance(%s) failed error(%v)", ins.String(), err))
			continue
		}
		c.tidbInsl.Lock()
		c.tidbInstances[insc.Name] = ins
		c.tidbInsl.Unlock()
		go ins.start()
		log.Info("reload tidb canal instance(%s) success", ins.String())
		c.sendWx(fmt.Sprintf("reload tidb canal instance(%s) success", ins.String()))
	}
}
