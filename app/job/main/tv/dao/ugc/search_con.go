package ugc

import (
	"context"

	model "go-common/app/job/main/tv/model/pgc"
	"go-common/library/database/sql"
	"go-common/library/ecode"
	"go-common/library/log"
)

const (
	_UgcFromwhere = "FROM ugc_archive WHERE result=1 AND valid=1 AND deleted=0 "
	_UgcCont      = "SELECT aid,title,cover,`content`,pubtime,typeid " + _UgcFromwhere + "AND aid > ? ORDER BY aid ASC LIMIT ?"
	_UgcContCount = " SELECT count(*) " + _UgcFromwhere
)

// UgcCont is used for getting valid ugc archive data
func (d *Dao) UgcCont(ctx context.Context, aid int, limit int) (res []*model.SearUgcCon, maxID int, err error) {
	var (
		rows *sql.Rows
	)
	if rows, err = d.DB.Query(ctx, _UgcCont, aid, limit); err != nil {
		log.Error("d.UgcCont.Query: %s error(%v)", _UgcCont, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var r = &model.SearUgcCon{}
		if err = rows.Scan(&r.AID, &r.Title, &r.Cover, &r.Content, &r.Pubtime, &r.Typeid); err != nil {
			log.Error("UgcCont row.Scan() error(%v)", err)
			return
		}
		res = append(res, r)
	}
	if err = rows.Err(); err != nil {
		log.Error("d.UgcCont.Query error(%v)", err)
		return
	}
	if len(res) == 0 {
		err = ecode.NothingFound
		return
	}
	maxID = res[len(res)-1].AID
	return
}

// UgcCnt is used for getting valid data count
func (d *Dao) UgcCnt(ctx context.Context) (upCnt int, err error) {
	row := d.DB.QueryRow(ctx, _UgcContCount)
	if err = row.Scan(&upCnt); err != nil {
		log.Error("d.SeaContCount.Query: %s error(%v)", _UgcContCount, err)
	}
	return
}
