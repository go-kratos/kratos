package canal

import (
	"github.com/siddontang/go-mysql/mysql"
	"github.com/siddontang/go-mysql/replication"
)

type EventHandler interface {
	OnRotate(roateEvent *replication.RotateEvent) error
	OnDDL(nextPos mysql.Position, queryEvent *replication.QueryEvent) error
	OnRow(e *RowsEvent) error
	OnXID(nextPos mysql.Position) error
	OnGTID(gtid mysql.GTIDSet) error
	// OnPosSynced Use your own way to sync position. When force is true, sync position immediately.
	OnPosSynced(pos mysql.Position, force bool) error
	String() string
}

type DummyEventHandler struct {
}

func (h *DummyEventHandler) OnRotate(*replication.RotateEvent) error { return nil }
func (h *DummyEventHandler) OnDDL(mysql.Position, *replication.QueryEvent) error {
	return nil
}
func (h *DummyEventHandler) OnRow(*RowsEvent) error                 { return nil }
func (h *DummyEventHandler) OnXID(mysql.Position) error             { return nil }
func (h *DummyEventHandler) OnGTID(mysql.GTIDSet) error             { return nil }
func (h *DummyEventHandler) OnPosSynced(mysql.Position, bool) error { return nil }
func (h *DummyEventHandler) String() string                         { return "DummyEventHandler" }

// `SetEventHandler` registers the sync handler, you must register your
// own handler before starting Canal.
func (c *Canal) SetEventHandler(h EventHandler) {
	c.eventHandler = h
}
