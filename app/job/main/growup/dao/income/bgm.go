package income

import (
	"context"
	"fmt"

	model "go-common/app/job/main/growup/model/income"
	"go-common/library/log"
)

const (
	_insertBGMSQL = "INSERT INTO background_music(mid,sid,aid,cid,join_at,title) VALUES %s"
	_getBGMSQL    = "SELECT id,mid,sid,aid,cid,join_at FROM background_music WHERE id > ? ORDER BY id LIMIT ?"
	_delBGMSQL    = "DELETE FROM background_music LIMIT ?"
)

// InsertBGM insert bgm from data platform
func (d *Dao) InsertBGM(c context.Context, values string) (rows int64, err error) {
	res, err := d.db.Exec(c, fmt.Sprintf(_insertBGMSQL, values))
	if err != nil {
		log.Error("insert bgm error(%v)", err)
		return
	}
	return res.RowsAffected()
}

// GetBGM get bgms
func (d *Dao) GetBGM(c context.Context, id int64, limit int64) (bs []*model.BGM, last int64, err error) {
	rows, err := d.db.Query(c, _getBGMSQL, id, limit)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		b := &model.BGM{}
		err = rows.Scan(&last, &b.MID, &b.SID, &b.AID, &b.CID, &b.JoinAt)
		if err != nil {
			return
		}
		bs = append(bs, b)
	}
	return
}

// DelBGM del bgm
func (d *Dao) DelBGM(c context.Context, limit int64) (rows int64, err error) {
	res, err := d.db.Exec(c, _delBGMSQL, limit)
	if err != nil {
		log.Error("del bgm error(%v)", err)
		return
	}
	return res.RowsAffected()
}
