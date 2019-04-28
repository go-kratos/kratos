package zlimit

import (
	"context"

	"go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/xstr"
)

const (
	_getPolicySQL       = "SELECT id, play_auth, down_auth, zone_id FROM policy_item WHERE zone_id <> '' AND state=1"
	_getRelationSQL     = "SELECT policy_id FROM archive_relation WHERE aid=?"
	_getGolbalPolicySQL = "select group_id,group_concat(id) from policy_item WHERE zone_id <> '' AND state=1 GROUP BY group_id"
)

// policies get policy data from db
func (s *Service) policies(c context.Context) (res map[int64]map[int64]int64, err error) {
	var (
		tmpres map[int64]int64
		ok     bool
	)
	rows, err := s.db.Query(c, _getPolicySQL)
	if err != nil {
		log.Error("db.Query error(%v)", err)
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
			log.Error("rows.Scan error(%v)", err)
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
	return
}

// groupPolicies get policy data from db group by group_id
func (s *Service) groupPolicies(c context.Context) (res map[int64][]int64, err error) {
	rows, err := s.db.Query(c, _getGolbalPolicySQL)
	if err != nil {
		log.Error("db.Query error(%v)", err)
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
			log.Error("rows.Scan error(%v)", err)
			return
		}
		if zoneIDs, err = xstr.SplitInts(pids); err != nil {
			log.Error("xstr.SplitInts(%s) error(%v)", pids, err)
			continue
		}
		res[groupID] = zoneIDs
	}
	return
}

// policy get pids from db by aid
func (s *Service) groupid(c context.Context, aid int64) (gid int64, err error) {
	row := s.getRelationStmt.QueryRow(c, aid)
	if err = row.Scan(&gid); err != nil {
		if err == sql.ErrNoRows {
			gid = 0
			err = nil
		} else {
			log.Error("rows.Scan error(%v)", err)
		}
	}
	return
}
