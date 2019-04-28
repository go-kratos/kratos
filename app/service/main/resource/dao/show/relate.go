package show

import (
	"context"
	"encoding/json"
	"time"

	"go-common/app/service/main/resource/model"
	"go-common/library/log"
)

const (
	_relateSQL = "SELECT `id`,`param`,`title`,`rec_reason`,`position`,`plat_ver`,`stime`,`etime`,`pgc_ids` FROM app_rcmd_pos WHERE `state`=1 AND `goto`='special' AND `pgc_relation`=1 AND `stime`<? AND `etime`>?"
)

// Relate get all relate rec.
func (d *Dao) Relate(c context.Context, now time.Time) (relates []*model.Relate, err error) {
	rows, err := d.db.Query(c, _relateSQL, now, now)
	if err != nil {
		log.Error("d.Relate.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := &model.Relate{}
		var verStr string
		if err = rows.Scan(&r.ID, &r.Param, &r.Title, &r.RecReason, &r.Position, &verStr, &r.STime, &r.ETime, &r.PgcIDs); err != nil {
			log.Error("d.Relate.rows.Scan error(%v)", err)
			return
		}
		if verStr != "" {
			var verStruct []*model.Version
			if err = json.Unmarshal([]byte(verStr), &verStruct); err != nil {
				log.Error("json.Unmarshal(%s) error(%v)", verStr, err)
				return
			}
			vm := make(map[int8][]*model.Version, len(verStruct))
			for _, v := range verStruct {
				vm[v.Plat] = append(vm[v.Plat], v)
			}
			r.Versions = vm
		}
		relates = append(relates, r)
	}
	err = rows.Err()
	return
}
