package canal

import (
	"sync"

	"github.com/siddontang/go-mysql/mysql"
	log "github.com/sirupsen/logrus"
)

type masterInfo struct {
	sync.RWMutex

	pos mysql.Position

	gtid mysql.GTIDSet
}

func (m *masterInfo) Update(pos mysql.Position) {
	log.Debugf("update master position %s", pos)

	m.Lock()
	m.pos = pos
	m.Unlock()
}

func (m *masterInfo) UpdateGTID(gtid mysql.GTIDSet) {
	log.Debugf("update master gtid %s", gtid.String())

	m.Lock()
	m.gtid = gtid
	m.Unlock()
}

func (m *masterInfo) Position() mysql.Position {
	m.RLock()
	defer m.RUnlock()

	return m.pos
}

func (m *masterInfo) GTID() mysql.GTIDSet {
	m.RLock()
	defer m.RUnlock()

	return m.gtid
}
