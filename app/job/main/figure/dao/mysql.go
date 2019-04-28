package dao

import (
	"context"
	"fmt"

	"go-common/app/job/main/figure/model"
	"go-common/library/database/sql"
	"go-common/library/log"

	"github.com/pkg/errors"
)

const (
	_shard        = 100
	_insertFigure = "INSERT INTO figure_user_%02d (mid,score,lawful_score,wide_score,friendly_score,bounty_score,creativity_score,ver,ctime,mtime) VALUES (?,?,?,?,?,?,?,?,?,?)"
	_existFigure  = "SELECT id FROM figure_user_%02d WHERE mid=? LIMIT 1"
)

func hit(mid int64) int64 {
	return mid % _shard
}

// ExistFigure exist user figure info
func (d *Dao) ExistFigure(c context.Context, mid int64) (id int64, err error) {
	res := d.db.QueryRow(c, fmt.Sprintf(_existFigure, hit(mid)), mid)
	if err = res.Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			return
		}
		errors.Wrapf(err, "row.Scan(%d) error", mid)
	}

	return
}

// SaveFigure init user figure info
func (d *Dao) SaveFigure(c context.Context, f *model.Figure) (id int64, err error) {
	res, err := d.db.Exec(c, fmt.Sprintf(_insertFigure, hit(f.Mid)), f.Mid, f.Score, f.LawfulScore, f.WideScore, f.FriendlyScore, f.BountyScore, f.CreativityScore, f.Ver, f.Ctime, f.Mtime)
	if err != nil {
		errors.Wrapf(err, "init user(%d) Figure info error(%v)", f.Mid, err)
		return
	}
	if id, err = res.LastInsertId(); err != nil {
		log.Error("res.LastInsertId() error(%v)", err)
	}
	return
}
