package dao

import (
	"context"
	"fmt"

	"go-common/app/job/main/dm/model"
	"go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_selSubSQL = "SELECT id,oid,type,pid,mid,state,attr,acount,count,mcount,move_count,maxlimit,childpool,ctime,mtime FROM dm_subject_%02d WHERE oid=? AND type=?"
)

// Subject get subject info from db.
func (d *Dao) Subject(c context.Context, tp int32, oid int64) (s *model.Subject, err error) {
	s = &model.Subject{}
	row := d.dmReader.QueryRow(c, fmt.Sprintf(_selSubSQL, d.hitSubject(oid)), oid, tp)
	if err = row.Scan(&s.ID, &s.Oid, &s.Type, &s.Pid, &s.Mid, &s.State, &s.Attr, &s.ACount, &s.Count, &s.MCount, &s.MoveCnt, &s.Maxlimit, &s.Childpool, &s.CTime, &s.MTime); err != nil {
		if err == sql.ErrNoRows {
			s = nil
			err = nil
		} else {
			log.Error("row.Scan() error(%v)", err)
		}
	}
	return
}
