package like

import (
	"context"
	"database/sql"

	l "go-common/app/interface/main/activity/model/like"

	"github.com/pkg/errors"
)

const (
	_userAchieveSQL        = "select id,aid,sid,mid,award from act_like_user_achievements where id = ?"
	_userAchieveUpSQL      = "update act_like_user_achievements set award = ? where id = ? and  award = 0"
	_userAchievementAddSQL = "insert into act_like_user_achievements (`aid`,`mid`,`sid`,`award`) values(?,?,?,?)"
	_userAchievementSQL    = "select id,aid,sid,mid,award from act_like_user_achievements where sid = ? and mid = ? and del = 0"
	// AwardNotChange .
	AwardNotChange = 0
	// AwardHasChange .
	AwardHasChange = 1
	// AwardNoGet .
	AwardNoGet = 2
)

// AddUserAchievment .
func (d *Dao) AddUserAchievment(c context.Context, userAchi *l.ActLikeUserAchievement) (ID int64, err error) {
	var res sql.Result
	if res, err = d.db.Exec(c, _userAchievementAddSQL, userAchi.Aid, userAchi.Mid, userAchi.Sid, userAchi.Award); err != nil {
		err = errors.Wrapf(err, "d.db.Exec(%s)", _userAchievementAddSQL)
		return
	}
	return res.LastInsertId()
}

// UserAchievement .
func (d *Dao) UserAchievement(c context.Context, sid, mid int64) (res []*l.ActLikeUserAchievement, err error) {
	rows, err := d.db.Query(c, _userAchievementSQL, sid, mid)
	if err != nil {
		err = errors.Wrapf(err, "d.db.Query(%s)", _userAchievementSQL)
		return
	}
	for rows.Next() {
		n := &l.ActLikeUserAchievement{}
		if err = rows.Scan(&n.ID, &n.Aid, &n.Sid, &n.Mid, &n.Award); err != nil {
			err = errors.Wrapf(err, "d.db.Scan(%s)", _userAchievementSQL)
			return
		}
		res = append(res, n)
	}
	if err = rows.Err(); err != nil {
		err = errors.Wrapf(err, "rows.Err(%s)", _userAchievementSQL)
	}
	return
}

// RawActUserAchieve .
func (d *Dao) RawActUserAchieve(c context.Context, id int64) (res *l.ActLikeUserAchievement, err error) {
	res = &l.ActLikeUserAchievement{}
	if err = d.db.QueryRow(c, _userAchieveSQL, id).Scan(&res.ID, &res.Aid, &res.Sid, &res.Mid, &res.Award); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			err = errors.Wrapf(err, "d.db.QueryRow(%s)", _userAchieveSQL)
		}
	}
	return
}

// ActUserAchieveChange .
func (d *Dao) ActUserAchieveChange(c context.Context, id, award int64) (upID int64, err error) {
	var res sql.Result
	if res, err = d.db.Exec(c, _userAchieveUpSQL, award, id); err != nil {
		err = errors.Wrapf(err, "d.db.Exec(%s)", _userAchieveUpSQL)
		return
	}
	return res.RowsAffected()
}
