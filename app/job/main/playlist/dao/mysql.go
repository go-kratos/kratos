package dao

import (
	"context"

	plmdl "go-common/app/interface/main/playlist/model"
	"go-common/app/job/main/playlist/model"
)

const (
	_statSQL    = "SELECT id,mid,fid,view,reply,fav,share,mtime FROM `playlist_stat` WHERE id = ?"
	_upViewSQL  = "UPDATE playlist_stat SET  view= ?, mtime=mtime  WHERE id = ?"
	_upFavSQL   = "UPDATE playlist_stat SET  fav= ?, mtime=mtime  WHERE id = ?"
	_upReplySQL = "UPDATE playlist_stat SET  reply= ?, mtime=mtime  WHERE id = ?"
	_upShareSQL = "UPDATE playlist_stat SET  share= ?, mtime=mtime  WHERE id = ?"
)

// Update updates stat  in db.
func (d *Dao) Update(c context.Context, stat *model.StatM, tp string) (rows int64, err error) {
	var tmpSQL string
	switch tp {
	case model.ViewCountType:
		tmpSQL = _upViewSQL
	case model.FavCountType:
		tmpSQL = _upFavSQL
	case model.ReplyCountType:
		tmpSQL = _upReplySQL
	case model.ShareCountType:
		tmpSQL = _upShareSQL
	}
	res, err := d.db.Exec(c, tmpSQL, *stat.Count, stat.ID)
	if err != nil {
		PromError("db:更新计数", tp+" Update(%d,%+v) error(%v)", stat.ID, stat, err)
		return
	}
	rows, err = res.RowsAffected()
	return
}

// Stat returns stat info.
func (d *Dao) Stat(c context.Context, pid int64) (stat *plmdl.PlStat, err error) {
	row := d.db.QueryRow(c, _statSQL, pid)
	stat = &plmdl.PlStat{}
	if err = row.Scan(&stat.ID, &stat.Mid, &stat.Fid, &stat.View, &stat.Reply, &stat.Fav, &stat.Share, &stat.MTime); err != nil {
		PromError("db:读取计数", "Stat(%v) error(%v)", pid, err)
	}
	return
}
