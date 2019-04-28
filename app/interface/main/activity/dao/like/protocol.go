package like

import (
	"context"

	lmdl "go-common/app/interface/main/activity/model/like"
	"go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_actProtocolSQL = "select id,sid,protocol,mtime,ctime,types,tags,pubtime,deltime,editime,hot,bgm_id,paster_id,oids,screen_set from act_subject_protocol where sid = ? limit 1"
)

// RawActSubjectProtocol .
func (dao *Dao) RawActSubjectProtocol(c context.Context, sid int64) (res *lmdl.ActSubjectProtocol, err error) {
	row := dao.db.QueryRow(c, _actProtocolSQL, sid)
	res = new(lmdl.ActSubjectProtocol)
	if err = row.Scan(&res.ID, &res.Sid, &res.Protocol, &res.Mtime, &res.Ctime, &res.Types, &res.Tags, &res.Pubtime, &res.Deltime, &res.Editime, &res.Hot, &res.BgmID, &res.PasterID, &res.Oids, &res.ScreenSet); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			log.Error("RawActSubjectProtocol:row.Scan error(%v)", err)
		}
	}
	return
}
