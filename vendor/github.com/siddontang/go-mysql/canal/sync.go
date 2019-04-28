package canal

import (
	"fmt"
	"regexp"
	"time"

	"github.com/juju/errors"
	"github.com/satori/go.uuid"
	"github.com/siddontang/go-mysql/mysql"
	"github.com/siddontang/go-mysql/replication"
	"github.com/siddontang/go-mysql/schema"
	log "github.com/sirupsen/logrus"
)

var (
	expCreateTable = regexp.MustCompile("(?i)^CREATE\\sTABLE(\\sIF\\sNOT\\sEXISTS)?\\s`{0,1}(.*?)`{0,1}\\.{0,1}`{0,1}([^`\\.]+?)`{0,1}\\s.*")
	expAlterTable  = regexp.MustCompile("(?i)^ALTER\\sTABLE\\s.*?`{0,1}(.*?)`{0,1}\\.{0,1}`{0,1}([^`\\.]+?)`{0,1}\\s.*")
	expRenameTable = regexp.MustCompile("(?i)^RENAME\\sTABLE\\s.*?`{0,1}(.*?)`{0,1}\\.{0,1}`{0,1}([^`\\.]+?)`{0,1}\\s{1,}TO\\s.*?")
	expDropTable   = regexp.MustCompile("(?i)^DROP\\sTABLE(\\sIF\\sEXISTS){0,1}\\s`{0,1}(.*?)`{0,1}\\.{0,1}`{0,1}([^`\\.]+?)`{0,1}($|\\s)")
)

func (c *Canal) startSyncer() (*replication.BinlogStreamer, error) {
	if !c.useGTID {
		pos := c.master.Position()
		s, err := c.syncer.StartSync(pos)
		if err != nil {
			return nil, errors.Errorf("start sync replication at binlog %v error %v", pos, err)
		}
		log.Infof("start sync binlog at binlog file %v", pos)
		return s, nil
	} else {
		gset := c.master.GTID()
		s, err := c.syncer.StartSyncGTID(gset)
		if err != nil {
			return nil, errors.Errorf("start sync replication at GTID %v error %v", gset, err)
		}
		log.Infof("start sync binlog at GTID %v", gset)
		return s, nil
	}
}

func (c *Canal) runSyncBinlog() error {
	s, err := c.startSyncer()
	if err != nil {
		return err
	}

	savePos := false
	force := false
	for {
		ev, err := s.GetEvent(c.ctx)

		if err != nil {
			return errors.Trace(err)
		}
		savePos = false
		force = false
		pos := c.master.Position()

		curPos := pos.Pos
		//next binlog pos
		pos.Pos = ev.Header.LogPos

		// We only save position with RotateEvent and XIDEvent.
		// For RowsEvent, we can't save the position until meeting XIDEvent
		// which tells the whole transaction is over.
		// TODO: If we meet any DDL query, we must save too.
		switch e := ev.Event.(type) {
		case *replication.RotateEvent:
			pos.Name = string(e.NextLogName)
			pos.Pos = uint32(e.Position)
			log.Infof("rotate binlog to %s", pos)
			savePos = true
			force = true
			if err = c.eventHandler.OnRotate(e); err != nil {
				return errors.Trace(err)
			}
		case *replication.RowsEvent:
			// we only focus row based event
			err = c.handleRowsEvent(ev)
			if err != nil {
				e := errors.Cause(err)
				// if error is not ErrExcludedTable or ErrTableNotExist or ErrMissingTableMeta, stop canal
				if e != ErrExcludedTable &&
					e != schema.ErrTableNotExist &&
					e != schema.ErrMissingTableMeta {
					log.Errorf("handle rows event at (%s, %d) error %v", pos.Name, curPos, err)
					return errors.Trace(err)
				}
			}
			continue
		case *replication.XIDEvent:
			savePos = true
			// try to save the position later
			if err := c.eventHandler.OnXID(pos); err != nil {
				return errors.Trace(err)
			}
		case *replication.MariadbGTIDEvent:
			// try to save the GTID later
			gtid := &e.GTID
			c.master.UpdateGTID(gtid)
			if err := c.eventHandler.OnGTID(gtid); err != nil {
				return errors.Trace(err)
			}
		case *replication.GTIDEvent:
			u, _ := uuid.FromBytes(e.SID)
			gset, err := mysql.ParseMysqlGTIDSet(fmt.Sprintf("%s:%d", u.String(), e.GNO))
			if err != nil {
				return errors.Trace(err)
			}
			c.master.UpdateGTID(gset)
			if err := c.eventHandler.OnGTID(gset); err != nil {
				return errors.Trace(err)
			}
		case *replication.QueryEvent:
			var (
				mb     [][]byte
				schema []byte
				table  []byte
			)
			regexps := []regexp.Regexp{*expCreateTable, *expAlterTable, *expRenameTable, *expDropTable}
			for _, reg := range regexps {
				mb = reg.FindSubmatch(e.Query)
				if len(mb) != 0 {
					break
				}
			}
			mbLen := len(mb)
			if mbLen == 0 {
				continue
			}

			// the first last is table name, the second last is database name(if exists)
			if len(mb[mbLen-2]) == 0 {
				schema = e.Schema
			} else {
				schema = mb[mbLen-2]
			}
			table = mb[mbLen-1]

			savePos = true
			force = true
			c.ClearTableCache(schema, table)
			log.Infof("table structure changed, clear table cache: %s.%s\n", schema, table)
			if err = c.eventHandler.OnDDL(pos, e); err != nil {
				return errors.Trace(err)
			}

		default:
			continue
		}

		if savePos {
			c.master.Update(pos)
			c.eventHandler.OnPosSynced(pos, force)
		}
	}

	return nil
}

func (c *Canal) handleRowsEvent(e *replication.BinlogEvent) error {
	ev := e.Event.(*replication.RowsEvent)

	// Caveat: table may be altered at runtime.
	schema := string(ev.Table.Schema)
	table := string(ev.Table.Table)

	t, err := c.GetTable(schema, table)
	if err != nil {
		return err
	}
	var action string
	switch e.Header.EventType {
	case replication.WRITE_ROWS_EVENTv1, replication.WRITE_ROWS_EVENTv2:
		action = InsertAction
	case replication.DELETE_ROWS_EVENTv1, replication.DELETE_ROWS_EVENTv2:
		action = DeleteAction
	case replication.UPDATE_ROWS_EVENTv1, replication.UPDATE_ROWS_EVENTv2:
		action = UpdateAction
	default:
		return errors.Errorf("%s not supported now", e.Header.EventType)
	}
	events := newRowsEvent(t, action, ev.Rows)
	return c.eventHandler.OnRow(events)
}

func (c *Canal) WaitUntilPos(pos mysql.Position, timeout time.Duration) error {
	timer := time.NewTimer(timeout)
	for {
		select {
		case <-timer.C:
			return errors.Errorf("wait position %v too long > %s", pos, timeout)
		default:
			curPos := c.master.Position()
			if curPos.Compare(pos) >= 0 {
				return nil
			} else {
				log.Debugf("master pos is %v, wait catching %v", curPos, pos)
				time.Sleep(100 * time.Millisecond)
			}
		}
	}

	return nil
}

func (c *Canal) GetMasterPos() (mysql.Position, error) {
	rr, err := c.Execute("SHOW MASTER STATUS")
	if err != nil {
		return mysql.Position{"", 0}, errors.Trace(err)
	}

	name, _ := rr.GetString(0, 0)
	pos, _ := rr.GetInt(0, 1)

	return mysql.Position{name, uint32(pos)}, nil
}

func (c *Canal) CatchMasterPos(timeout time.Duration) error {
	pos, err := c.GetMasterPos()
	if err != nil {
		return errors.Trace(err)
	}

	return c.WaitUntilPos(pos, timeout)
}
