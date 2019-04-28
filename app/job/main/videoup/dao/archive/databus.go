package archive

import (
	"context"

	"go-common/app/job/main/videoup/model/archive"
	"go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_dbusSQL   = "SELECT id,gp,topic,part,last_offset FROM archive_databus WHERE gp=? AND topic=? AND part=?"
	_inDBusSQL = "INSERT INTO archive_databus(gp,topic,part,last_offset) VALUES(?,?,?,?)"
	_upDBusSQL = "UPDATE archive_databus SET last_offset=? WHERE gp=? AND topic=? AND part=?"
)

// DBus get DBus by group+topic+partition
func (d *Dao) DBus(c context.Context, group, topic string, partition int32) (dbus *archive.Databus, err error) {
	row := d.db.QueryRow(c, _dbusSQL, group, topic, partition)
	dbus = &archive.Databus{}
	if err = row.Scan(&dbus.ID, &dbus.Group, &dbus.Topic, &dbus.Partition, &dbus.Offset); err != nil {
		if err == sql.ErrNoRows {
			dbus = nil
			err = nil
		} else {
			log.Error("row.Scan error(%v)", err)
		}
	}
	return
}

// AddDBus add databus
func (d *Dao) AddDBus(c context.Context, group, topic string, partition int32, offset int64) (rows int64, err error) {
	res, err := d.db.Exec(c, _inDBusSQL, group, topic, partition, offset)
	if err != nil {
		log.Error("d.db.Exec(%s, %s, %d, %d) error(%v)", group, topic, partition, offset, err)
		return
	}
	return res.RowsAffected()
}

// UpDBus update databus offset
func (d *Dao) UpDBus(c context.Context, group, topic string, partition int32, offset int64) (rows int64, err error) {
	res, err := d.db.Exec(c, _upDBusSQL, offset, group, topic, partition)
	if err != nil {
		log.Error("d.db.Exec(%d, %s, %s, %d) error(%v)", offset, group, topic, partition, err)
		return
	}
	return res.RowsAffected()
}
