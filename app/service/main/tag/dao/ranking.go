package dao

import (
	"context"

	"go-common/app/service/main/tag/model"
	"go-common/library/log"
)

var (
	_rankHotSQL = "SELECT tag_id,tag_name FROM ranking_tag"
)

// RankHots .
func (d *Dao) RankHots(c context.Context) (ts []*model.Tag, err error) {
	rows, err := d.db.Query(c, _rankHotSQL)
	if err != nil {
		log.Error("d.rkiHotsStmt.Query() error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		t := &model.Tag{}
		if err = rows.Scan(&t.ID, &t.Name); err != nil {
			log.Error("rows.Scan() error(%v)", err)
			return
		}
		ts = append(ts, t)
	}
	return
}

var (
	_bangumiSQL = "SELECT id,season_id,cover,name,color FROM ranking_bangumi"
)

// Bangumis .
func (d *Dao) Bangumis(c context.Context) (seasonIds []int64, bgms []*model.RankingBangumi, err error) {
	rows, err := d.db.Query(c, _bangumiSQL)
	if err != nil {
		log.Error("d.rkiBgmsStmt.Query() error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		bgm := &model.RankingBangumi{}
		if err = rows.Scan(&bgm.ID, &bgm.SeasonID, &bgm.Cover, &bgm.Name, &bgm.Color); err != nil {
			log.Error("rows.Scan() error(%v)", err)
			return
		}
		seasonIds = append(seasonIds, bgm.SeasonID)
		bgms = append(bgms, bgm)
	}
	return
}

var (
	_regionsSQL = "SELECT target_id,target_name,is_tag FROM ranking_part WHERE type=? ORDER BY sort"
)

// Regions .
func (d *Dao) Regions(c context.Context, rid int64) (ps []*model.RankingRegion, err error) {
	rows, err := d.db.Query(c, _regionsSQL, rid)
	if err != nil {
		log.Error("d.rkiRgnsStmt.Query(%d) error(%v)", rid, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		p := &model.RankingRegion{}
		if err = rows.Scan(&p.Tid, &p.Tname, &p.IsTag); err != nil {
			log.Error("rows.Scan() error(%v)", err)
			ps = nil
			return
		}
		ps = append(ps, p)
	}
	return
}
