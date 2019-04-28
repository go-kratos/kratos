package dao

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"go-common/app/job/main/figure-timer/model"
	xsql "go-common/library/database/sql"

	"github.com/pkg/errors"
)

const (
	_shard             = 100
	_selFigure         = `SELECT id,mid,score,lawful_score,wide_score,friendly_score,bounty_score,creativity_score,ver,ctime,mtime FROM figure_user_%02d WHERE mid=? LIMIT 1`
	_selFigures        = `SELECT id,mid,score,lawful_score,wide_score,friendly_score,bounty_score,creativity_score,ver,ctime,mtime FROM figure_user_%02d WHERE mid>? LIMIT ?`
	_upsertFigure      = `INSERT INTO figure_user_%02d (mid,score,lawful_score,wide_score,friendly_score,bounty_score,creativity_score,ver,ctime,mtime)  VALUES (?,?,?,?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE score=?,lawful_score=?,wide_score=?,friendly_score=?,bounty_score=?,creativity_score=?,ver=?,mtime=?`
	_insertRank        = `INSERT INTO figure_rank (score_from,score_to,percentage) VALUES(?,?,?) ON DUPLICATE KEY UPDATE score_from=?,score_to=?,percentage=?`
	_insertRankHistory = `INSERT INTO figure_rank_history (score_from,score_to,percentage,ver) VALUES(?,?,?,?)`
)

func hit(mid int64) int64 {
	return mid % _shard
}

// Figure get Figure from db
func (d *Dao) Figure(c context.Context, mid int64) (figure *model.Figure, err error) {
	row := d.mysql.QueryRow(c, fmt.Sprintf(_selFigure, hit(mid)), mid)
	figure = &model.Figure{}
	if err = row.Scan(&figure.ID, &figure.Mid, &figure.Score, &figure.LawfulScore, &figure.WideScore, &figure.FriendlyScore, &figure.BountyScore, &figure.CreativityScore, &figure.Ver, &figure.Ctime, &figure.Mtime); err != nil {
		if err == xsql.ErrNoRows {
			err = nil
			figure = nil
			return
		}
		err = errors.WithStack(err)
		return
	}
	return
}

// Figures get all figure info from formMid
func (d *Dao) Figures(c context.Context, fromMid int64, limit int) (figures []*model.Figure, end bool, err error) {
	if limit <= 0 {
		return
	}
	var (
		rows *xsql.Rows
	)
	if rows, err = d.mysql.Query(c, fmt.Sprintf(_selFigures, hit(fromMid)), fromMid, limit); err != nil {
		err = errors.WithStack(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		figure := &model.Figure{}
		if err = rows.Scan(&figure.ID, &figure.Mid, &figure.Score, &figure.LawfulScore, &figure.WideScore, &figure.FriendlyScore, &figure.BountyScore, &figure.CreativityScore, &figure.Ver, &figure.Ctime, &figure.Mtime); err != nil {
			if err == xsql.ErrNoRows {
				err = nil
				end = true
				return
			}
			return
		}
		figures = append(figures, figure)
	}
	if len(figures) < limit {
		end = true
	}
	err = errors.WithStack(rows.Err())
	return
}

// UpsertFigure insert or update(if mid duplicated) Figure
func (d *Dao) UpsertFigure(c context.Context, figure *model.Figure) (id int64, err error) {
	var (
		result sql.Result
		now    = time.Now()
	)
	if result, err = d.mysql.Exec(c, fmt.Sprintf(_upsertFigure, hit(figure.Mid)), figure.Mid, figure.Score, figure.LawfulScore, figure.WideScore, figure.FriendlyScore, figure.BountyScore, figure.CreativityScore, figure.Ver, now, now, figure.Score, figure.LawfulScore, figure.WideScore, figure.FriendlyScore, figure.BountyScore, figure.CreativityScore, figure.Ver, now); err != nil {
		return
	}
	if id, err = result.LastInsertId(); err != nil {
		err = errors.WithStack(err)
	}
	return
}

// InsertRankHistory insert figure rank history to db
func (d *Dao) InsertRankHistory(c context.Context, rank *model.Rank) (id int64, err error) {
	var (
		res sql.Result
	)
	if res, err = d.mysql.Exec(c, _insertRankHistory, rank.ScoreFrom, rank.ScoreTo, rank.Percentage, rank.Ver); err != nil {
		return
	}
	if id, err = res.LastInsertId(); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

// UpsertRank insert or update figure rank to db
func (d *Dao) UpsertRank(c context.Context, rank *model.Rank) (id int64, err error) {
	var (
		res sql.Result
	)
	if res, err = d.mysql.Exec(c, _insertRank, rank.ScoreFrom, rank.ScoreTo, rank.Percentage, rank.ScoreFrom, rank.ScoreTo, rank.Percentage); err != nil {
		return
	}
	if id, err = res.LastInsertId(); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}
