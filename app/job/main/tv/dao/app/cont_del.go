package app

import (
	"context"
	dsql "database/sql"
	"fmt"
	"time"

	model "go-common/app/job/main/tv/model/pgc"
	"go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/xstr"
)

const (
	_deletedContSQL = "SELECT id, title, subtitle, `desc`, cover, epid, season_id, cid FROM tv_content WHERE state = ? AND is_deleted = 1 AND audit_time < UNIX_TIMESTAMP(now()) ORDER BY season_id, epid LIMIT 0,"
	_delaySyncSQL   = "UPDATE tv_content SET audit_time = ? WHERE epid IN (%s)"
	_delSyncContSQL = "UPDATE tv_content SET state = ? WHERE is_deleted = 1 AND epid = ?"
)

// DelCont picks the deleted content data to sync
func (d *Dao) DelCont(c context.Context) (res []*model.Content, err error) {
	var (
		rows   *sql.Rows
		nbData = d.conf.Sync.LConf.SizeMsg
	)
	if rows, err = d.DB.Query(c, _deletedContSQL+fmt.Sprintf("%d", nbData), EPToAudit); err != nil {
		log.Error("d._deletedEPSQL.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var r = &model.Content{}
		if err = rows.Scan(&r.ID, &r.Title, &r.Subtitle, &r.Desc, &r.Cover, &r.EPID, &r.SeasonID, &r.CID); err != nil {
			log.Error("DelCont row.Scan() error(%v)", err)
			return
		}
		res = append(res, r)
	}
	if err = rows.Err(); err != nil {
		log.Error("d.DelCont.Query error(%v)", err)
	}
	return
}

// SyncCont updates the delete content's state to not pass, avoid being selected again
func (d *Dao) SyncCont(c context.Context, cid int) (nbRows int64, err error) {
	var res dsql.Result
	if res, err = d.DB.Exec(c, _delSyncContSQL, EPNotPass, cid); err != nil {
		log.Error("_delSyncContSQL, failed to update to auditing: (%v), Error: %v", cid, err)
		return
	}
	return res.RowsAffected()
}

// DelaySync postpones the sync of the deleted content, in case of there is an error happenning in the interface
func (d *Dao) DelaySync(c context.Context, conts []*model.Content) (nbRows int64, err error) {
	var (
		res     dsql.Result
		xstrIds []int64
		delay   = time.Now().Unix() + int64(d.conf.Sync.Frequency.AuditDelay) // postpone the season's eps auditing to when
	)
	for _, v := range conts {
		xstrIds = append(xstrIds, int64(v.EPID))
	}
	if res, err = d.DB.Exec(c, fmt.Sprintf(_delaySyncSQL, xstr.JoinInts(xstrIds)), delay); err != nil {
		log.Error("_delSyncContSQL, failed to delay: (%v,%v), Error: %v", delay, xstrIds, err)
		return
	}
	return res.RowsAffected()
}
