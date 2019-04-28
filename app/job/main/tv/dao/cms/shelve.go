package cms

import (
	"context"
	"fmt"

	"go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/xstr"
)

const (
	_validSns = "SELECT DISTINCT a.id FROM tv_ep_season a LEFT JOIN tv_content b ON a.id = b.season_id " +
		"WHERE a.is_deleted = 0 AND a.`check` = ? AND b.is_deleted = 0 AND b.valid = ? AND b.state = ?"
	_allPassedSns = "SELECT id, valid FROM tv_ep_season WHERE is_deleted = 0 AND `check` = 1"
	_actSns       = "UPDATE tv_ep_season SET valid = ? WHERE id IN (%s)"
	_offArcs      = "SELECT aid FROM ugc_archive WHERE aid IN (%s) AND valid = 0 AND deleted = 0 AND result = 1 "
	_reshelfArcs  = "UPDATE ugc_archive SET valid = 1 WHERE aid IN (%s)"
	_cmsOnline    = 1
	_cmsOffline   = 0
	_epPassed     = 3
)

// ValidSns gets all the seasons that should be on the shelves, which includes free and audited episodes.
func (d *Dao) ValidSns(ctx context.Context, onlyfree bool) (res map[int64]int, err error) {
	var (
		rows     *sql.Rows
		validSql = _validSns
	)
	res = make(map[int64]int)
	if onlyfree {
		validSql = validSql + " AND b.pay_status = 2" // free episode
	}
	if rows, err = d.DB.Query(ctx, validSql, _cmsOnline, _cmsOnline, _epPassed); err != nil {
		log.Error("d.ValidSns.Query: %s error(%v)", validSql, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var sid int64
		if err = rows.Scan(&sid); err != nil {
			log.Error("ValidSns row.Scan() error(%v)", err)
			return
		}
		res[sid] = 1
	}
	if err = rows.Err(); err != nil {
		log.Error("d.PgcCont.Query error(%v)", err)
	}
	return
}

// ShelveOp gets the status of all audited seasons on and off shelves, and compare the results with the "ValidSns" method above to determine which episodes need to be on or off shelves.
func (d *Dao) ShelveOp(ctx context.Context, validSns map[int64]int) (onIDs, offIDs []int64, err error) {
	var rows *sql.Rows
	if rows, err = d.DB.Query(ctx, _allPassedSns); err != nil {
		log.Error("d.ShelveOp.Query: %s error(%v)", _allPassedSns, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var sid, valid int64
		if err = rows.Scan(&sid, &valid); err != nil {
			log.Error("ValidSns row.Scan() error(%v)", err)
			return
		}
		_, ok := validSns[sid]
		if ok && valid == 0 {
			onIDs = append(onIDs, sid)
		}
		if !ok && valid == 1 {
			offIDs = append(offIDs, sid)
		}
	}
	if err = rows.Err(); err != nil {
		log.Error("d.PgcCont.Query error(%v)", err)
	}
	return
}

// ActOps carries out the action on the season which need to be on/off shelves
func (d *Dao) ActOps(ctx context.Context, ids []int64, on bool) (err error) {
	var action int
	if on {
		action = _cmsOnline
	} else {
		action = _cmsOffline
	}
	if _, err = d.DB.Exec(ctx, fmt.Sprintf(_actSns, xstr.JoinInts(ids)), action); err != nil {
		log.Error("ActOps, Ids %v, Err %v", ids, err)
	}
	return
}

// OffArcs takes the archives that passed but cms invalid archives
func (d *Dao) OffArcs(ctx context.Context, aids []int64) (offAids []int64, err error) {
	var rows *sql.Rows
	if rows, err = d.DB.Query(ctx, fmt.Sprintf(_offArcs, xstr.JoinInts(aids))); err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var aid int64
		if err = rows.Scan(&aid); err != nil {
			return
		}
		offAids = append(offAids, aid)
	}
	err = rows.Err()
	return
}

// ReshelfArcs re-put the arcs onshelf ( CMS valid )
func (d *Dao) ReshelfArcs(ctx context.Context, aids []int64) (err error) {
	_, err = d.DB.Exec(ctx, fmt.Sprintf(_reshelfArcs, xstr.JoinInts(aids)))
	return
}
