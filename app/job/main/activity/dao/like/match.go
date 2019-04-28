package like

import (
	"context"
	"database/sql"
	"fmt"

	"go-common/app/job/main/activity/model/match"
	"go-common/library/log"
	"go-common/library/xstr"
)

const (
	_unDoMatchSQL   = "SELECT mid,stake,result FROM act_match_user_log WHERE m_o_id = ? AND status = 0 LIMIT ?"
	_upMatchUserSQL = "UPDATE act_match_user_log SET status = 1 WHERE m_o_id = ? AND mid IN (%s)"
	_matchObjSQL    = "SELECT id,match_id,sid,result FROM act_matchs_object WHERE status = 0 AND id = ?"
)

// UnDoMatchUsers un finish users.
func (d *Dao) UnDoMatchUsers(c context.Context, matchObjID int64, limit int) (list []*match.ActMatchUser, err error) {
	rows, err := d.db.Query(c, _unDoMatchSQL, matchObjID, limit)
	if err != nil {
		log.Error("UnDoMatchUsers.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		user := new(match.ActMatchUser)
		if err = rows.Scan(&user.Mid, &user.Stake, &user.Result); err != nil {
			log.Error("UnDoMatchUsers row.Scan error(%v)", err)
			return
		}
		list = append(list, user)
	}
	return
}

// UpMatchUserResult update match user result.
func (d *Dao) UpMatchUserResult(c context.Context, matchObjID int64, mids []int64) (err error) {
	if _, err = d.db.Exec(c, fmt.Sprintf(_upMatchUserSQL, xstr.JoinInts(mids)), matchObjID); err != nil {
		log.Error("UpMatchUserResult d.db.Exec mids(%v) error(%v)", mids, err)
	}
	return
}

// MatchObjInfo get match object info from db.
func (d *Dao) MatchObjInfo(c context.Context, matchObjID int64) (data *match.ActMatchObj, err error) {
	row := d.db.QueryRow(c, _matchObjSQL, matchObjID)
	data = new(match.ActMatchObj)
	if err = row.Scan(&data.ID, &data.MatchID, &data.SID, &data.Result); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			log.Error("MatchObjInfo row.Scan() error(%v)", err)
		}
	}
	return
}
