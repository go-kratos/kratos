package dao

import (
	"context"
	"fmt"

	"go-common/app/service/main/figure/model"
	xsql "go-common/library/database/sql"

	"github.com/pkg/errors"
)

const (
	_shard      = 100
	_figureInfo = "SELECT id, mid, score, lawful_score, wide_score, friendly_score, bounty_score, creativity_score, ver, ctime, mtime FROM figure_user_%02d WHERE mid=? ORDER BY id DESC LIMIT 1"
	_rank       = `SELECT score_from,score_to,percentage FROM figure_rank ORDER BY percentage ASC`
)

func hit(mid int64) int64 {
	return mid % _shard
}

// FigureInfo get user figure info
func (d *Dao) FigureInfo(c context.Context, mid int64) (res *model.Figure, err error) {
	row := d.db.QueryRow(c, fmt.Sprintf(_figureInfo, hit(mid)), mid)
	res = &model.Figure{}
	if err = row.Scan(&res.ID, &res.Mid, &res.Score, &res.LawfulScore, &res.WideScore, &res.FriendlyScore, &res.BountyScore, &res.CreativityScore, &res.Ver, &res.Ctime, &res.Mtime); err != nil {
		if err == xsql.ErrNoRows {
			err = nil
			res = nil
			return
		}
		err = errors.WithStack(err)
	}
	return
}

// Ranks get figure score rank by ver
func (d *Dao) Ranks(c context.Context) (ranks []*model.Rank, err error) {
	var (
		rows *xsql.Rows
	)
	if rows, err = d.db.Query(c, _rank); err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var rank = &model.Rank{}
		if err = rows.Scan(&rank.ScoreFrom, &rank.ScoreTo, &rank.Percentage); err != nil {
			ranks = nil
			return
		}
		ranks = append(ranks, rank)
	}
	err = rows.Err()
	return
}
