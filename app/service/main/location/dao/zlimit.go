package dao

import (
	"context"

	"go-common/app/service/main/location/model"
	"go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/xstr"

	"github.com/pkg/errors"
)

const (
	_getPolicySQL       = "SELECT id, play_auth, down_auth, zone_id FROM policy_item WHERE zone_id <> '' AND state=1"
	_getRelationSQL     = "SELECT policy_id FROM archive_relation WHERE aid=?"
	_getGolbalPolicySQL = "SELECT group_id,group_concat(id) FROM policy_item WHERE zone_id <> '' AND state=1 GROUP BY group_id"
	_getGroupZone       = "SELECT a.group_id,a.play_auth,a.zone_id FROM policy_item AS a,policy_group AS b WHERE a.zone_id <> '' AND a.group_id=b.id AND b.type=2 AND a.state=1 AND b.state=1"
)

// Policies get policy data from db
func (d *Dao) Policies(c context.Context) (res map[int64]map[int64]int64, err error) {
	var (
		tmpres map[int64]int64
		ok     bool
	)
	rows, err := d.db.Query(c, _getPolicySQL)
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	defer rows.Close()
	res = make(map[int64]map[int64]int64)
	for rows.Next() {
		var (
			pid, playAuth, downAuth int64
			zoneID                  string
			zoneIDs                 []int64
		)
		if err = rows.Scan(&pid, &playAuth, &downAuth, &zoneID); err != nil {
			err = errors.WithStack(err)
			return
		}
		if zoneIDs, err = xstr.SplitInts(zoneID); err != nil {
			log.Error("xstr.SplitInts(%s) error(%v)", zoneID, err)
			continue
		}
		for _, zoneid := range zoneIDs {
			if tmpres, ok = res[pid]; !ok {
				tmpres = make(map[int64]int64)
				res[pid] = tmpres
			}
			resCode := playAuth<<8 | downAuth
			tmpres[zoneid] = resCode
		}
	}
	err = errors.WithStack(err)
	return
}

// GroupPolicies get policy data from db group by group_id
func (d *Dao) GroupPolicies(c context.Context) (res map[int64][]int64, err error) {
	rows, err := d.db.Query(c, _getGolbalPolicySQL)
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	defer rows.Close()
	res = make(map[int64][]int64)
	for rows.Next() {
		var (
			groupID int64
			pids    string
			zoneIDs []int64
		)
		if err = rows.Scan(&groupID, &pids); err != nil {
			err = errors.WithStack(err)
			return
		}
		if zoneIDs, err = xstr.SplitInts(pids); err != nil {
			log.Error("xstr.SplitInts(%s) error(%v)", pids, err)
			continue
		}
		res[groupID] = zoneIDs
	}
	err = errors.WithStack(err)
	return
}

// Groupid get gid from db by aid
func (d *Dao) Groupid(c context.Context, aid int64) (gid int64, err error) {
	row := d.db.QueryRow(c, _getRelationSQL, aid)
	if err = row.Scan(&gid); err != nil {
		if err == sql.ErrNoRows {
			gid = 0
			err = nil
		} else {
			err = errors.WithStack(err)
		}
	}
	return
}

// GroupAuthZone zone_id by group_id.
func (d *Dao) GroupAuthZone(c context.Context) (res map[int64]map[int64]map[int64]int64, err error) {
	var (
		tmpAres map[int64]map[int64]int64
		tmpZres map[int64]int64
		ok      bool
	)
	rows, err := d.db.Query(c, _getGroupZone)
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	defer rows.Close()
	res = make(map[int64]map[int64]map[int64]int64)
	for rows.Next() {
		var (
			gid, playAuth int64
			zoneID        string
			zoneIDs       []int64
		)
		if err = rows.Scan(&gid, &playAuth, &zoneID); err != nil {
			err = errors.WithStack(err)
			return
		}
		if playAuth != model.Forbidden && playAuth != model.Allow {
			playAuth = model.Allow
		}
		if zoneIDs, err = xstr.SplitInts(zoneID); err != nil {
			log.Error("xstr.SplitInts(%s) error(%v)", zoneID, err)
			continue
		}
		for _, zoneid := range zoneIDs {
			if tmpAres, ok = res[gid]; !ok {
				tmpAres = make(map[int64]map[int64]int64)
				res[gid] = tmpAres
			}
			if tmpZres, ok = tmpAres[playAuth]; !ok {
				tmpZres = make(map[int64]int64)
				tmpAres[playAuth] = tmpZres
			}
			if _, ok = tmpZres[zoneid]; !ok {
				tmpZres[zoneid] = zoneid
			}
		}
	}
	err = errors.WithStack(err)
	return
}
