package dao

import (
	"context"

	"go-common/app/infra/canal/model"
	"go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_tidbPositonSQL       = "SELECT name, cluster_id, offset, tso FROM  tidb_info WHERE name = ?"
	_updateTidbPositonSQL = "INSERT INTO tidb_info(name, cluster_id, offset, tso) VALUES(?,?,?,?) ON DUPLICATE KEY UPDATE offset = ?, tso = ?"
)

// TiDBPosition get tidb positon
func (d *Dao) TiDBPosition(c context.Context, name string) (res *model.TiDBInfo, err error) {
	res = &model.TiDBInfo{}
	if err = d.db.QueryRow(c, _tidbPositonSQL, name).Scan(&res.Name, &res.ClusterID, &res.Offset, &res.CommitTS); err != nil {
		if err == sql.ErrNoRows {
			res = nil
			err = nil
			return
		}
		log.Error("db.TidbPosition.Query error(%v,%v,%v)", _tidbPositonSQL, name, err)
		return
	}
	return
}

// UpdateTiDBPosition update tidb position
func (d *Dao) UpdateTiDBPosition(c context.Context, info *model.TiDBInfo) (err error) {
	if info == nil {
		return
	}
	if _, err = d.db.Exec(c, _updateTidbPositonSQL, info.Name, info.ClusterID, info.Offset, info.CommitTS, info.Offset, info.CommitTS); err != nil {
		log.Error("db.UpdateTiDBPosition.Exec error(%v,%+v,%v)", _updateTidbPositonSQL, info, err)
	}
	return
}
