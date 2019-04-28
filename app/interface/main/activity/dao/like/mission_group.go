package like

import (
	"context"
	"database/sql"
	"fmt"

	l "go-common/app/interface/main/activity/model/like"
	"go-common/library/xstr"

	"github.com/pkg/errors"
)

const (
	_likeMissionBuffSQL  = "select id  from like_mission_group where sid = ? and mid = ?"
	_likeMissionAddSQL   = "insert into like_mission_group (`sid`,`mid`,`state`) values (?,?,?)"
	_likeMissionGroupSQL = "select id,sid,mid,state,ctime,mtime from like_mission_group where id in (%s)"
	// MissionStateInit the init state
	MissionStateInit = 0
)

// RawLikeMissionBuff get mid has .
func (d *Dao) RawLikeMissionBuff(c context.Context, sid, mid int64) (ID int64, err error) {
	res := &l.MissionGroup{}
	row := d.db.QueryRow(c, _likeMissionBuffSQL, sid, mid)
	if err = row.Scan(&res.ID); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			err = errors.Wrap(err, "RawLikeMissionBuff:QueryRow")
			return
		}
	}
	ID = res.ID
	return
}

// MissionGroupAdd add like_mission_group data .
func (d *Dao) MissionGroupAdd(c context.Context, group *l.MissionGroup) (misID int64, err error) {
	var res sql.Result
	if res, err = d.db.Exec(c, _likeMissionAddSQL, group.Sid, group.Mid, group.State); err != nil {
		err = errors.Wrapf(err, "d.db.Exec(%s)", _likeMissionAddSQL)
		return
	}
	return res.LastInsertId()
}

// RawMissionGroupItems get mission_group item by ids.
func (d *Dao) RawMissionGroupItems(c context.Context, lids []int64) (res map[int64]*l.MissionGroup, err error) {
	res = make(map[int64]*l.MissionGroup, len(lids))
	rows, err := d.db.Query(c, fmt.Sprintf(_likeMissionGroupSQL, xstr.JoinInts(lids)))
	if err != nil {
		err = errors.Wrapf(err, "d.db.Query(%s)", _likeMissionGroupSQL)
		return
	}
	defer rows.Close()
	for rows.Next() {
		n := &l.MissionGroup{}
		if err = rows.Scan(&n.ID, &n.Sid, &n.Mid, &n.State, &n.Ctime, &n.Mtime); err != nil {
			err = errors.Wrapf(err, "d.db.Scan(%s)", _likeMissionGroupSQL)
			return
		}
		res[n.ID] = n
	}
	if err = rows.Err(); err != nil {
		err = errors.Wrap(err, "RawMissionGroupItem:rows.Err()")
	}
	return
}
