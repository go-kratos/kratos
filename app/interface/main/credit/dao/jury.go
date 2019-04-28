package dao

import (
	"context"
	"fmt"
	"time"

	model "go-common/app/interface/main/credit/model"
	"go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/xstr"

	"github.com/pkg/errors"
)

const (
	_updateJurySQL      = "INSERT INTO blocked_jury(mid,status,expired) VALUES(?,?,?) ON DUPLICATE KEY UPDATE status=?,expired=?"
	_updateVoteTotalSQL = "UPDATE blocked_jury SET vote_total=vote_total+1 WHERE mid=?"
	_getJurySQL         = "SELECT id,mid,status,expired,invalid_reason,vote_total,vote_radio,vote_right,vote_total,black FROM blocked_jury WHERE mid=?"
	_juryInfosSQL       = "SELECT mid,status,expired,invalid_reason,vote_total,vote_radio,vote_right,total,black FROM blocked_jury WHERE mid IN (%s)"
)

//JuryApply user jury apply.
func (d *Dao) JuryApply(c context.Context, mid int64, expired time.Time) (err error) {
	if _, err = d.db.Exec(c, _updateJurySQL, mid, 1, expired, 1, expired); err != nil {
		log.Error("JuryApply: db.Exec(%d,%d,%v) error(%v)", mid, 1, expired, err)
	}
	return
}

// AddUserVoteTotal add user vote total.
func (d *Dao) AddUserVoteTotal(c context.Context, mid int64) (err error) {
	if _, err = d.db.Exec(c, _updateVoteTotalSQL, mid); err != nil {
		log.Error("AddUserVoteTotal: db.Exec(%d) error(%v)", mid, err)
	}
	return
}

// JuryInfo get user's jury info
func (d *Dao) JuryInfo(c context.Context, mid int64) (r *model.BlockedJury, err error) {
	row := d.db.QueryRow(c, _getJurySQL, mid)
	r = &model.BlockedJury{}
	if err = row.Scan(&r.ID, &r.MID, &r.Status, &r.Expired, &r.InvalidReason, &r.VoteTotal, &r.VoteRadio, &r.VoteRight, &r.CaseTotal, &r.Black); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			r = nil
		}
	}
	return
}

// JuryInfos get user's applying jury info
func (d *Dao) JuryInfos(c context.Context, mids []int64) (mbj map[int64]*model.BlockedJury, err error) {
	rows, err := d.db.Query(c, fmt.Sprintf(_juryInfosSQL, xstr.JoinInts(mids)))
	if err != nil {
		err = errors.Wrap(err, "JuryInfos")
		return
	}
	mbj = make(map[int64]*model.BlockedJury, len(mids))
	defer rows.Close()
	for rows.Next() {
		bj := new(model.BlockedJury)
		if err = rows.Scan(&bj.MID, &bj.Status, &bj.Expired, &bj.InvalidReason, &bj.VoteTotal, &bj.VoteRadio, &bj.VoteRight, &bj.CaseTotal, &bj.Black); err != nil {
			if err == sql.ErrNoRows {
				err = nil
				return
			}
			err = errors.Wrap(err, "JuryInfos")
			return
		}
		mbj[bj.MID] = bj
	}
	err = rows.Err()
	return
}
