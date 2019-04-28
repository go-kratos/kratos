package tidb

import (
	"time"

	"go-common/library/log"
)

func (db *DB) nodeproc(e <-chan struct{}) {
	if db.dis == nil {
		return
	}
	for {
		<-e
		nodes := db.nodeList()
		if len(nodes) == 0 {
			continue
		}
		cm := make(map[string]*conn)
		var conns []*conn
		for _, conn := range db.conns {
			cm[conn.addr] = conn
		}
		for _, node := range nodes {
			if cm[node] != nil {
				conns = append(conns, cm[node])
				continue
			}
			c, err := db.connectDSN(genDSN(db.conf.DSN, node))
			if err == nil {
				conns = append(conns, c)
			} else {
				log.Error("tidb: connect addr: %s err: %+v", node, err)
			}
		}
		if len(conns) == 0 {
			log.Error("tidb: no nodes ignore event")
			continue
		}
		oldConns := db.conns
		db.mutex.Lock()
		db.conns = conns
		db.mutex.Unlock()
		log.Info("tidb: new nodes: %v", nodes)
		var removedConn []*conn
		for _, conn := range oldConns {
			var exist bool
			for _, c := range conns {
				if c.addr == conn.addr {
					exist = true
					break
				}
			}
			if !exist {
				removedConn = append(removedConn, conn)
			}
		}
		go db.closeConns(removedConn)
	}
}

func (db *DB) closeConns(conns []*conn) {
	if len(conns) == 0 {
		return
	}
	du := db.conf.QueryTimeout
	if db.conf.ExecTimeout > du {
		du = db.conf.ExecTimeout
	}
	if db.conf.TranTimeout > du {
		du = db.conf.TranTimeout
	}
	time.Sleep(time.Duration(du))
	for _, conn := range conns {
		err := conn.Close()
		if err != nil {
			log.Error("tidb: close removed conn: %s err: %v", conn.addr, err)
		} else {
			log.Info("tidb: close removed conn: %s", conn.addr)
		}
	}
}
