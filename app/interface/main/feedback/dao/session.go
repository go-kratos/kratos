package dao

import (
	"context"
	"fmt"
	"time"

	"go-common/app/interface/main/feedback/model"
	"go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/xstr"
	"strings"
)

const (
	_selSsn               = "SELECT id,mid,content,img_url,log_url,state,ctime FROM session WHERE buvid=? AND system=? AND version=? AND mid=?"
	_selSsnByMid          = "SELECT id,mid,content,img_url,log_url,state,ctime FROM session WHERE mid=? AND platform IN (%s)"
	_selSSnCntByMid       = `SELECT COUNT(id) AS count FROM session WHERE mid=? AND state IN (0,1,2) AND platform IN ("ugc","article")`
	_inSsn                = "INSERT INTO session (buvid,system,version,mid,aid,content,img_url,log_url,device,channel,ip,net_state,net_operator,agency_area,platform,browser,qq,email,state,laster_time,ctime,mtime) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
	_upSsn                = "UPDATE session SET device=?,channel=?,ip=?,net_state=?,net_operator=?,agency_area=?,platform=?,browser=?,qq=?,email=?,state=?,laster_time=?,mtime=? WHERE id=?"
	_upSsnMtime           = "UPDATE session SET mtime=? WHERE id=?"
	_upSsnState           = "UPDATE session SET state=? where id=?"
	_selTagByPlat         = "SELECT id,name,platform,type FROM tag where type=?  AND platform IN (%s)"
	_selSsnIDbByTagID     = "SELECT session_id FROM session_tag WHERE tag_id IN (%s)"
	_selSSnBySsnID        = "SELECT id,content,ctime,state FROM session WHERE id IN (%s) AND state IN (%s) AND session.ctime>? AND session.ctime<? ORDER BY id DESC"
	_selSSnBySsnIDAllSate = "SELECT id,content,ctime,state FROM session WHERE id IN (%s) AND session.ctime>? AND session.ctime<? ORDER BY id DESC"
	_selSSnID             = "select id from session where id=? limit 1"
)

// JudgeSsnRecord judge session is exist or not .
func (d *Dao) JudgeSsnRecord(c context.Context, sid int64) (cnt int, err error) {
	res, err := d.selSSnID.Exec(c, sid)
	if err != nil {
		log.Error("d.upSsnSta.Exec error(%v)", err)
		return
	}
	res.RowsAffected()
	return
}

// Session select feedback session
func (d *Dao) Session(c context.Context, buvid, system, version string, mid int64) (ssn *model.Session, err error) {
	row := d.selSsn.QueryRow(c, buvid, system, version, mid)
	ssn = &model.Session{}
	if err = row.Scan(&ssn.ID, &ssn.Mid, &ssn.Content, &ssn.ImgURL, &ssn.LogURL, &ssn.State, &ssn.CTime); err != nil {
		if err == sql.ErrNoRows {
			err, ssn = nil, nil
		} else {
			log.Error("row.Scan() error(%v)", err)
		}
	}
	return
}

//SessionCount session count.
func (d *Dao) SessionCount(c context.Context, mid int64) (cnt int, err error) {
	row := d.selSSnCntByMid.QueryRow(c, mid)
	if err = row.Scan(&cnt); err != nil {
		log.Error("row.Scan error(%v)", err)
	}
	return
}

// UpSsnMtime up ssn mtime.
func (d *Dao) UpSsnMtime(c context.Context, now time.Time, id int64) (err error) {
	_, err = d.upSsnMtime.Exec(c, now, id)
	if err != nil {
		log.Error("d.upSsnMtime error(%v)", err)
	}
	return
}

// TxUpSsnMtime up ssn mtime.
func (d *Dao) TxUpSsnMtime(tx *sql.Tx, now time.Time, id int64) (err error) {
	_, err = tx.Exec(_upSsnMtime, now, id)
	if err != nil {
		log.Error("d.upSsnMtime error(%v)", err)
	}
	return
}

// SessionIDByTagID session find by time state and mid .
func (d *Dao) SessionIDByTagID(c context.Context, tagID []int64) (sid []int64, err error) {
	rows, err := d.dbMs.Query(c, fmt.Sprintf(_selSsnIDbByTagID, xstr.JoinInts(tagID)))
	if err != nil {
		log.Error("d.dbMs.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var id int64
		if err = rows.Scan(&id); err != nil {
			log.Error("row.Scan() error(%v)", err)
			return
		}
		sid = append(sid, id)
	}
	return
}

// SessionBySsnID get session by ssnID
func (d *Dao) SessionBySsnID(c context.Context, sid []int64, state string, start, end time.Time) (ssns []*model.Session, err error) {
	rows, err := d.dbMs.Query(c, fmt.Sprintf(_selSSnBySsnID, xstr.JoinInts(sid), state), start, end)
	if err != nil {
		log.Error("d.selSSnBySsnID.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		ssn := &model.Session{}
		if err = rows.Scan(&ssn.ID, &ssn.Content, &ssn.CTime, &ssn.State); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		ssns = append(ssns, ssn)
	}
	return
}

// SSnBySsnIDAllSate get all SSn by ssnID all sate.
func (d *Dao) SSnBySsnIDAllSate(c context.Context, sid []int64, start, end time.Time) (ssns []*model.Session, err error) {
	rows, err := d.dbMs.Query(c, fmt.Sprintf(_selSSnBySsnIDAllSate, xstr.JoinInts(sid)), start, end)
	if err != nil {
		log.Error("d.selSSnBySsnID.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		ssn := &model.Session{}
		if err = rows.Scan(&ssn.ID, &ssn.Content, &ssn.CTime, &ssn.State); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		ssns = append(ssns, ssn)
	}
	return
}

// SessionByMid select feedback session by mid
func (d *Dao) SessionByMid(c context.Context, mid int64, platform string) (ssns []*model.Session, err error) {
	rows, err := d.dbMs.Query(c, fmt.Sprintf(_selSsnByMid, platConvert(platform)), mid)
	if err != nil {
		log.Error("d.dbMs.Query error(%v)", err)
		return
	}
	defer rows.Close()
	ssns = []*model.Session{}
	for rows.Next() {
		ssn := &model.Session{}
		if err = rows.Scan(&ssn.ID, &ssn.Mid, &ssn.Content, &ssn.ImgURL, &ssn.LogURL, &ssn.State, &ssn.CTime); err != nil {
			log.Error("row.Scan() error(%v)", err)
			return
		}
		ssns = append(ssns, ssn)
	}
	return
}

// TxUpdateSessionState up date session state.
func (d *Dao) TxUpdateSessionState(tx *sql.Tx, state int, sid int64) (err error) {
	res, err := tx.Exec(_upSsnState, state, sid)
	if err != nil {
		log.Error("d.upSsnSta.Exec error(%v)", err)
		return
	}
	_, err = res.RowsAffected()
	return
}

// UpdateSessionState update session state.
func (d *Dao) UpdateSessionState(c context.Context, state int, sid int64) (err error) {
	res, err := d.upSsnSta.Exec(c, state, sid)
	if err != nil {
		log.Error("d.upSsnSta.Exec error(%v)", err)
		return
	}
	_, err = res.RowsAffected()
	return
}

// Tags get tags.
func (d *Dao) Tags(c context.Context, mold int, platform string) (tMap map[string][]*model.Tag, err error) {
	rows, err := d.dbMs.Query(c, fmt.Sprintf(_selTagByPlat, platConvert(platform)), mold)
	if err != nil {
		log.Error("d.selTagByPlat.Query error(%v)", err)
		return
	}
	defer rows.Close()
	tMap = make(map[string][]*model.Tag)
	for rows.Next() {
		tag := &model.Tag{}
		if err = rows.Scan(&tag.ID, &tag.Name, &tag.Platform, &tag.Type); err != nil {
			log.Error("row.Scan() error(%v)", err)
			return
		}
		if tag.Type == 0 {
			tMap[tag.Platform] = append(tMap[tag.Platform], tag)
		}
	}
	return
}

// AddSession add feedback session
func (d *Dao) AddSession(c context.Context, s *model.Session) (id int64, err error) {
	res, err := d.inSsn.Exec(c, s.Buvid, s.System, s.Version, s.Mid, s.Aid, s.Content, s.ImgURL, s.LogURL, s.Device, s.Channel, s.IP, s.NetState, s.NetOperator, s.AgencyArea, s.Platform, s.Browser, s.QQ, s.Email, s.State, s.LasterTime, s.CTime, s.MTime)
	if err != nil {
		log.Error("AddSession tx.Exec() error(%v)", err)
		return
	}
	return res.LastInsertId()
}

// TxAddSession add session
func (d *Dao) TxAddSession(tx *sql.Tx, s *model.Session) (id int64, err error) {
	res, err := tx.Exec(_inSsn, s.Buvid, s.System, s.Version, s.Mid, s.Aid, s.Content, s.ImgURL, s.LogURL, s.Device, s.Channel, s.IP, s.NetState, s.NetOperator, s.AgencyArea, s.Platform, s.Browser, s.QQ, s.Email, s.State, s.LasterTime, s.CTime, s.MTime)
	if err != nil {
		log.Error("AddSession tx.Exec() error(%v)", err)
		return
	}
	return res.LastInsertId()
}

// AddSessionTag add feedback session and tag.
func (d *Dao) AddSessionTag(c context.Context, sessionID, tagID int64, now time.Time) (id int64, err error) {
	res, err := d.inSsnTag.Exec(c, sessionID, tagID, now)
	if err != nil {
		log.Error("AddSession tx.Exec() error(%v)", err)
		return
	}
	return res.LastInsertId()
}

// TxAddSessionTag add session tag.
func (d *Dao) TxAddSessionTag(tx *sql.Tx, sessionID, tagID int64, now time.Time) (id int64, err error) {
	res, err := tx.Exec(_inSsnTag, sessionID, tagID, now)
	if err != nil {
		log.Error("TxAddSessionTag tx.Exec() error(%v)", err)
		return
	}
	return res.LastInsertId()
}

// UpdateSession update feedback session
func (d *Dao) UpdateSession(c context.Context, s *model.Session) (affected int64, err error) {
	res, err := d.upSsn.Exec(c, s.Device, s.Channel, s.IP, s.NetState, s.NetOperator, s.AgencyArea, s.Platform, s.Browser, s.QQ, s.Email, s.State, s.LasterTime, s.MTime, s.ID)
	if err != nil {
		log.Error("UpdateSession tx.Exec() error(%v)", err)
		return
	}
	return res.RowsAffected()
}

// TagIDBySid get tagid by sid.
func (d *Dao) TagIDBySid(c context.Context, sids []int64) (tsMap map[int64]int64, err error) {
	rows, err := d.dbMs.Query(c, fmt.Sprintf(_selTagID, xstr.JoinInts(sids)))
	if err != nil {
		log.Error("d.dbMs.Query() error(%v)", err)
		return
	}
	defer rows.Close()
	tsMap = make(map[int64]int64)
	for rows.Next() {
		var (
			tid int64
			sid int64
		)
		if err = rows.Scan(&tid, &sid); err != nil {
			log.Error("row.Scan() error(%v)", err)
			return
		}
		tsMap[tid] = sid
	}
	return
}

// platConvert plat covert.
func platConvert(platform string) (s string) {
	s = `"` + strings.Replace(platform, ",", `","`, -1) + `"`
	return
}
