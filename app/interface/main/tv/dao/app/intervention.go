package dao

import (
	"context"
	"database/sql"
	"fmt"

	"go-common/app/interface/main/tv/model"
	"go-common/app/interface/main/tv/model/search"
	xsql "go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/xstr"
)

const (
	_intervs      = "SELECT cont_id, cont_type FROM tv_rank WHERE is_deleted = 0"
	_parseInterv  = "SELECT cont_id, cont_type FROM tv_rank "
	_modIntervs   = "%s WHERE module_id = %d AND is_deleted = 0 ORDER BY position ASC LIMIT %d"
	_rankList     = "%s WHERE rank = %d AND category = %d AND is_deleted = 0 ORDER BY position ASC LIMIT %d"
	_rmPGCInterv  = "UPDATE tv_rank SET is_deleted = 1 WHERE cont_type != 2 AND is_deleted = 0 AND cont_id IN (%s)"
	_rmUGCInterv  = "UPDATE tv_rank SET is_deleted = 1 WHERE cont_type = 2 AND is_deleted = 0 AND cont_id IN (%s)"
	_idxIntervs   = "SELECT rank, category, cont_id FROM tv_rank WHERE category IN (?,?) AND rank > 0 AND is_deleted = 0 ORDER BY position ASC"
	_pgcIdxInterv = 6
	_ugcIdxInterv = 7
)

// ModIntervs get intervention data with a given mod ID
func (d *Dao) ModIntervs(c context.Context, modID int, capacity int) (resp *model.RespModInterv, err error) {
	var sql = fmt.Sprintf(_modIntervs, _parseInterv, modID, capacity)
	return d.rowsTreat(c, sql)
}

// ZoneIntervs get db data
func (d *Dao) ZoneIntervs(c context.Context, req *model.ReqZoneInterv) (resp *model.RespModInterv, err error) {
	var sql = fmt.Sprintf(_rankList, _parseInterv, req.RankType, req.Category, req.Limit)
	return d.rowsTreat(c, sql)
}

// IdxIntervs def.
func (d *Dao) IdxIntervs(c context.Context) (idxSave *search.IdxIntervSave, err error) {
	var (
		rows *xsql.Rows
	)
	idxSave = &search.IdxIntervSave{
		Pgc: make(map[int][]int64),
		Ugc: make(map[int][]int64),
	}
	if rows, err = d.db.Query(c, _idxIntervs, _pgcIdxInterv, _ugcIdxInterv); err != nil {
		log.Error("IdxIntervs d.db.Query (%s) error(%v)", _idxIntervs, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		li := &model.Rank{}
		if err = rows.Scan(&li.Rank, &li.Category, &li.ContID); err != nil {
			return
		}
		if li.Category == _pgcIdxInterv {
			idxSave.Pgc[li.Rank] = append(idxSave.Pgc[li.Rank], li.ContID)
		} else {
			idxSave.Ugc[li.Rank] = append(idxSave.Ugc[li.Rank], li.ContID)
		}
	}
	err = rows.Err()
	return
}

// AllIntervs picks all the active intervention data
func (d *Dao) AllIntervs(c context.Context) (sids []int64, aids []int64, err error) {
	var (
		pgcMap = make(map[int64]int)
		ugcMap = make(map[int64]int)
	)
	rows, err := d.db.Query(c, _intervs)
	if err != nil {
		log.Error("AllIntervs d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		li := &model.SimpleRank{}
		if err = rows.Scan(&li.ContID, &li.ContType); err != nil {
			log.Error("AllIntervs row.Scan error(%v)", err)
			return
		}
		if li.IsUGC() {
			ugcMap[li.ContID] = 1
		} else {
			pgcMap[li.ContID] = 1
		}
	}
	for avid := range ugcMap {
		aids = append(aids, avid)
	}
	for sid := range pgcMap {
		sids = append(sids, sid)
	}
	log.Info("AllIntervs Distinct UGC Intervs %d, PGC Intervs %d", len(aids), len(sids))
	return
}

func (d *Dao) rowsTreat(c context.Context, sql string) (resp *model.RespModInterv, err error) {
	var ranks []*model.SimpleRank
	resp = new(model.RespModInterv)
	rows, err := d.db.Query(c, sql)
	if err != nil {
		log.Error("ZoneIntervs d.db.Query (%s) error(%v)", sql, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		li := &model.SimpleRank{}
		if err = rows.Scan(&li.ContID, &li.ContType); err != nil {
			log.Error("ZoneIntervs row.Scan error(%v)", err)
			return
		}
		ranks = append(ranks, li)
	}
	resp = &model.RespModInterv{
		Ranks: ranks,
	}
	for _, v := range ranks {
		if v.IsUGC() {
			resp.AIDs = append(resp.AIDs, v.ContID)
		} else {
			resp.SIDs = append(resp.SIDs, v.ContID)
		}
	}
	return
}

// RmInterv removes invalids interventions
func (d *Dao) RmInterv(c context.Context, aids []int64, sids []int64) (err error) {
	var (
		resPGC, resUGC   sql.Result
		rowsPGC, rowsUGC int64
	)
	if len(sids) > 0 {
		if resPGC, err = d.db.Exec(c, fmt.Sprintf(_rmPGCInterv, xstr.JoinInts(sids))); err != nil {
			log.Error("RmInterv Sids %v, Err %v", sids, err)
			return
		}
		if rowsPGC, err = resPGC.RowsAffected(); err != nil {
			log.Error("RmInterv Sids %v, Err %v", sids, err)
			return
		}
	}
	if len(aids) > 0 {
		if resUGC, err = d.db.Exec(c, fmt.Sprintf(_rmUGCInterv, xstr.JoinInts(aids))); err != nil {
			log.Error("RmInterv Aids %v, Err %v", aids, err)
			return
		}
		if rowsUGC, err = resUGC.RowsAffected(); err != nil {
			log.Error("RmInterv Aids %v, Err %v", aids, err)
			return
		}
	}
	log.Warn("RmInterv Aids %v, Sids %v, UGCDel Rows %d, PGCDel Rows %d", aids, sids, rowsUGC, rowsPGC)
	return
}
