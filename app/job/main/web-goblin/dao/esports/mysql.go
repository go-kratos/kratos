package esports

import (
	"context"
	"fmt"
	"strconv"
	"time"

	mdlesp "go-common/app/job/main/web-goblin/model/esports"
	arcmdl "go-common/app/service/main/archive/api"
	xsql "go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/xstr"

	"github.com/pkg/errors"
)

const (
	_contestsSQL = "SELECT c.id,c.stime,c.live_room,c.home_id,c.away_id,c.success_team,c.special,c.special_name,c.special_tips,s.title,s.sub_title FROM `es_contests` as c INNER JOIN `es_seasons` as s ON c.sid=s.id  WHERE c.status = 0  AND  c.stime >= ? and c.stime < ? "
	_teamSQL     = "SELECT id,title,sub_title FROM `es_teams`  WHERE  is_deleted = 0  AND (id = ? or id = ?)"
	_arcSQL      = "SELECT id,aid,score,is_deleted FROM `es_archives`  WHERE  is_deleted = 0  AND id > ? ORDER BY id ASC LIMIT ? "
	_arcEditSQL  = "UPDATE es_archives SET score = CASE %s END WHERE aid IN (%s)"
)

// Contests  contests by time.
func (d *Dao) Contests(c context.Context, stime, etime int64) (res []*mdlesp.Contest, err error) {
	var (
		rows *xsql.Rows
	)
	if rows, err = d.db.Query(c, _contestsSQL, stime, etime); err != nil {
		log.Error("Contests:d.db.Query(%d) error(%v)", stime, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := new(mdlesp.Contest)
		if err = rows.Scan(&r.ID, &r.Stime, &r.LiveRoom, &r.HomeID, &r.AwayID, &r.SuccessTeam, &r.Special, &r.SpecialName, &r.SpecialTips, &r.SeasonTitle, &r.SeasonSubTitle); err != nil {
			log.Error("Contests:row.Scan() error(%v)", err)
			return
		}
		res = append(res, r)
	}
	if err = rows.Err(); err != nil {
		log.Error("rows.Err() error(%v)", err)
	}
	return
}

// Teams  teams by id.
func (d *Dao) Teams(c context.Context, homeID, awayID int64) (res []*mdlesp.Team, err error) {
	var (
		rows *xsql.Rows
	)
	if rows, err = d.db.Query(c, _teamSQL, homeID, awayID); err != nil {
		log.Error("Teams:d.db.Query homeID(%d) awayID(%d) error(%v)", homeID, awayID, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := new(mdlesp.Team)
		if err = rows.Scan(&r.ID, &r.Title, &r.SubTitle); err != nil {
			log.Error("Teams:row.Scan() error(%v)", err)
			return
		}
		res = append(res, r)
	}
	if err = rows.Err(); err != nil {
		log.Error("rows.Err() error(%v)", err)
	}
	return
}

// Arcs archives by ids.
func (d *Dao) Arcs(c context.Context, id int64, limit int) (res []*mdlesp.Arc, err error) {
	var (
		rows *xsql.Rows
	)
	if rows, err = d.db.Query(c, _arcSQL, id, limit); err != nil {
		log.Error("Arcs:d.db.Query id(%d) limit(%d) error(%v)", id, limit, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := new(mdlesp.Arc)
		if err = rows.Scan(&r.ID, &r.Aid, &r.Score, &r.IsDeleted); err != nil {
			log.Error("Arcs:row.Scan() error(%v)", err)
			return
		}
		res = append(res, r)
	}
	if err = rows.Err(); err != nil {
		log.Error("rows.Err() error(%v)", err)
	}
	return
}

// UpArcScore  update  archive score.
func (d *Dao) UpArcScore(c context.Context, partArcs []*mdlesp.Arc, arcs map[int64]*arcmdl.Arc) (err error) {
	var (
		caseStr string
		aids    []int64
		score   int64
	)
	for _, v := range partArcs {
		if arc, ok := arcs[v.Aid]; ok {
			score = d.score(arc)
		} else {
			continue
		}
		caseStr = fmt.Sprintf("%s WHEN aid = %d THEN %d", caseStr, v.Aid, score)
		aids = append(aids, v.Aid)
	}
	if len(aids) == 0 {
		return
	}
	if _, err = d.db.Exec(c, fmt.Sprintf(_arcEditSQL, caseStr, xstr.JoinInts(aids))); err != nil {
		err = errors.Wrapf(err, "UpArcScore  d.db.Exec")
	}
	return
}

func (d *Dao) score(arc *arcmdl.Arc) (res int64) {
	tmpRs := float64(arc.Stat.Coin)*d.c.Rule.CoinPercent +
		float64(arc.Stat.Fav)*d.c.Rule.FavPercent + float64(arc.Stat.Danmaku)*d.c.Rule.DmPercent +
		float64(arc.Stat.Reply)*d.c.Rule.ReplyPercent + float64(arc.Stat.View)*d.c.Rule.ViewPercent +
		float64(arc.Stat.Like)*d.c.Rule.LikePercent + float64(arc.Stat.Share)*d.c.Rule.SharePercent
	now := time.Now()
	hours := now.Sub(arc.PubDate.Time()).Hours()
	if hours/24 <= d.c.Rule.NewDay {
		tmpRs = tmpRs * 1.5
	}
	decimal, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", tmpRs), 64)
	res = int64(decimal * 100)
	return
}
