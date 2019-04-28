package like

import (
	"context"

	likemdl "go-common/app/interface/main/activity/model/like"
	xsql "go-common/library/database/sql"

	"github.com/pkg/errors"
)

const (
	_achievesSQL = "select `id`,`name`,`icon`,`dic`,`unlock`,`ctime`,`mtime`,`del`,`sid`,`image`,`award` from act_like_achievements where sid = ? and del = 0 order by `unlock` limit 100"
	// HaveAward award state
	HaveAward = 1
)

// RawActLikeAchieves .
func (d *Dao) RawActLikeAchieves(c context.Context, sid int64) (res *likemdl.Achievements, err error) {
	var rows *xsql.Rows
	if rows, err = d.db.Query(c, _achievesSQL, sid); err != nil {
		if err == xsql.ErrNoRows {
			err = nil
		} else {
			err = errors.Wrapf(err, "RawActLikeAchieves:Query(%s)", _achievesSQL)
			return
		}
	}
	defer rows.Close()
	list := make([]*likemdl.ActLikeAchievement, 0, 100)
	for rows.Next() {
		a := &likemdl.ActLikeAchievement{}
		if err = rows.Scan(&a.ID, &a.Name, &a.Icon, &a.Dic, &a.Unlock, &a.Ctime, &a.Mtime, &a.Del, &a.Sid, &a.Image, &a.Award); err != nil {
			err = errors.Wrap(err, "RawActLikeAchieves:scan()")
			return
		}
		list = append(list, a)
	}
	res = &likemdl.Achievements{Achievements: list}
	if err = rows.Err(); err != nil {
		err = errors.Wrap(err, "RawActLikeAchieves:rows.Err()")
	}
	return
}
