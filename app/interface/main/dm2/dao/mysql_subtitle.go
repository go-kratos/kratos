package dao

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"go-common/app/interface/main/dm2/model"
	xsql "go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/xstr"
)

const (
	_subtitleSharding = 100

	_getSubtitlePubIds = "SELECT subtitle_id FROM subtitle_pub WHERE oid=? AND type=? AND is_delete=0"
	_addSubtitlePub    = "INSERT INTO subtitle_pub(oid,type,lan,subtitle_id,is_delete) VALUES(?,?,?,?,?) ON DUPLICATE KEY UPDATE subtitle_id=?,is_delete=?"

	_addSubtitleSubject = "INSERT INTO subtitle_subject(aid,allow,lan_code) VALUES(?,?,?) ON DUPLICATE KEY UPDATE allow=?,lan_code=?"
	_getSubtitleSubject = "SELECT aid,allow,attr,lan_code from subtitle_subject WHERE aid=?"

	_getSubtitleOne   = "SELECT id,oid,type,lan,aid,mid,author_mid,up_mid,is_sign,is_lock,status,checksum,subtitle_url,pub_time,reject_comment FROM subtitle_%02d WHERE oid=? AND type=? AND lan=? AND status=5 ORDER BY pub_time DESC limit 1"
	_getSubtitles     = "SELECT id,oid,type,lan,aid,mid,author_mid,up_mid,is_sign,is_lock,status,checksum,subtitle_url,pub_time,reject_comment,mtime FROM subtitle_%02d WHERE id in (%s)"
	_getSubtitle      = "SELECT id,oid,type,lan,aid,mid,author_mid,up_mid,is_sign,is_lock,status,checksum,subtitle_url,pub_time,reject_comment,mtime FROM subtitle_%02d WHERE id = ? AND status!=4"
	_getSubtitleDraft = "SELECT id,oid,type,lan,aid,mid,author_mid,up_mid,is_sign,is_lock,status,checksum,subtitle_url,pub_time,reject_comment FROM subtitle_%02d WHERE oid=? AND type=? AND lan=? AND mid=? AND pub_time=0"
	_updateSubtitle   = "UPDATE subtitle_%02d SET aid=?,author_mid=?,up_mid=?,is_sign=?,is_lock=?,status=?,checksum=?,subtitle_url=?,pub_time=?,reject_comment=? WHERE id=?"
	_addSubtitle      = "INSERT INTO subtitle_%02d(id,oid,type,lan,aid,mid,author_mid,up_mid,is_sign,is_lock,status,checksum,subtitle_url,pub_time,reject_comment) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"

	_addWaveForm = "INSERT INTO subtitle_waveform(oid,type,state,wave_form_url) VALUES(?,?,?,?) ON DUPLICATE KEY UPDATE state=?,wave_form_url=?"
	_getWaveForm = "SELECT oid,type,state,wave_form_url,mtime FROM subtitle_waveform WHERE oid=? AND type=?"

	_getSubtitleLans = "SELECT code,lan,doc_zh FROM subtitle_language WHERE is_delete=0"
	_addSubtitleLan  = "INSERT INTO subtitle_language(code,lan,doc_zh,doc_en,is_delete) VALUES(?,?,?,?,?) ON DUPLICATE KEY UPDATE lan=?,doc_zh=?,doc_en=?,is_delete=?"
)

func (d *Dao) hitSubtitle(oid int64) int64 {
	return oid % _subtitleSharding
}

// AddSubtitleSubject .
func (d *Dao) AddSubtitleSubject(c context.Context, subtitleSubject *model.SubtitleSubject) (err error) {
	if _, err = d.dbDM.Exec(c, _addSubtitleSubject, subtitleSubject.Aid, subtitleSubject.Allow, subtitleSubject.Lan, subtitleSubject.Allow, subtitleSubject.Lan); err != nil {
		log.Error("params(subtitleSubject:%+v),error(%v)", subtitleSubject, err)
		return
	}
	return
}

// GetSubtitleSubject .
func (d *Dao) GetSubtitleSubject(c context.Context, aid int64) (subtitleSubject *model.SubtitleSubject, err error) {
	subtitleSubject = new(model.SubtitleSubject)
	row := d.dbDM.QueryRow(c, _getSubtitleSubject, aid)
	if err = row.Scan(&subtitleSubject.Aid, &subtitleSubject.Allow, &subtitleSubject.Attr, &subtitleSubject.Lan); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			subtitleSubject = nil
			return
		}
		log.Error("params(aid:%v),error(%v)", aid, err)
		return
	}
	return
}

// GetSubtitleIds .
func (d *Dao) GetSubtitleIds(c context.Context, oid int64, tp int32) (subtitlIds []int64, err error) {
	rows, err := d.dbDM.Query(c, _getSubtitlePubIds, oid, tp)
	if err != nil {
		log.Error("params(oid:%v,tp:%v),error(%v)", oid, tp, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var subtitleID int64
		if err = rows.Scan(&subtitleID); err != nil {
			log.Error("params(oid:%v,tp:%v),=error(%v)", oid, tp, err)
			return
		}
		subtitlIds = append(subtitlIds, subtitleID)
	}
	if err = rows.Err(); err != nil {
		log.Error("params(oid:%v,tp:%v),=error(%v)", oid, tp, err)
		return
	}
	return
}

// GetSubtitles .
func (d *Dao) GetSubtitles(c context.Context, oid int64, subtitleIds []int64) (subtitles []*model.Subtitle, err error) {
	rows, err := d.dbDM.Query(c, fmt.Sprintf(_getSubtitles, d.hitSubtitle(oid), xstr.JoinInts(subtitleIds)))
	if err != nil {
		log.Error("params(subtitleIds:%v),=error(%v)", subtitleIds, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		subtitle := &model.Subtitle{}
		var t time.Time
		if err = rows.Scan(&subtitle.ID, &subtitle.Oid, &subtitle.Type, &subtitle.Lan, &subtitle.Aid, &subtitle.Mid, &subtitle.AuthorID, &subtitle.UpMid, &subtitle.IsSign, &subtitle.IsLock, &subtitle.Status, &subtitle.CheckSum, &subtitle.SubtitleURL, &subtitle.PubTime, &subtitle.RejectComment, &t); err != nil {
			log.Error("params(subtitleIds:%v),=error(%v)", subtitleIds, err)
			return
		}
		subtitle.Mtime = t.Unix()
		subtitles = append(subtitles, subtitle)
	}
	if err = rows.Err(); err != nil {
		log.Error("params(subtitleIds:%v),=error(%v)", subtitleIds, err)
		return
	}
	return
}

// GetSubtitleDraft query a SubtitleDrfat
func (d *Dao) GetSubtitleDraft(c context.Context, oid int64, tp int32, mid int64, lan uint8) (subtitle *model.Subtitle, err error) {
	subtitle = &model.Subtitle{}
	row := d.dbDM.QueryRow(c, fmt.Sprintf(_getSubtitleDraft, d.hitSubtitle(oid)), oid, tp, lan, mid)
	if err = row.Scan(&subtitle.ID, &subtitle.Oid, &subtitle.Type, &subtitle.Lan, &subtitle.Aid, &subtitle.Mid, &subtitle.AuthorID, &subtitle.UpMid, &subtitle.IsSign, &subtitle.IsLock, &subtitle.Status, &subtitle.CheckSum, &subtitle.SubtitleURL, &subtitle.PubTime, &subtitle.RejectComment); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			subtitle = nil
			return
		}
		log.Error("params(oid:%v, tp:%v, mid:%v, lan:%v),=error(%v)", oid, tp, mid, lan, err)
		return
	}
	return
}

// GetSubtitle query a SubtitleDrfat
func (d *Dao) GetSubtitle(c context.Context, oid int64, subtitleID int64) (subtitle *model.Subtitle, err error) {
	var t time.Time
	subtitle = &model.Subtitle{}
	row := d.dbDM.QueryRow(c, fmt.Sprintf(_getSubtitle, d.hitSubtitle(oid)), subtitleID)
	if err = row.Scan(&subtitle.ID, &subtitle.Oid, &subtitle.Type, &subtitle.Lan, &subtitle.Aid, &subtitle.Mid, &subtitle.AuthorID, &subtitle.UpMid, &subtitle.IsSign, &subtitle.IsLock, &subtitle.Status, &subtitle.CheckSum, &subtitle.SubtitleURL, &subtitle.PubTime, &subtitle.RejectComment, &t); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			subtitle = nil
			return
		}
		log.Error("params(subtitleID:%v),error(%v)", subtitleID, err)
		return
	}
	subtitle.Mtime = t.Unix()
	return
}

// AddSubtitle .
func (d *Dao) AddSubtitle(c context.Context, subtitle *model.Subtitle) (insertID int64, err error) {
	var res sql.Result
	if res, err = d.dbDM.Exec(c, fmt.Sprintf(_addSubtitle, d.hitSubtitle(subtitle.Oid)),
		subtitle.ID, subtitle.Oid, subtitle.Type, subtitle.Lan, subtitle.Aid, subtitle.Mid, subtitle.AuthorID, subtitle.UpMid, subtitle.IsSign, subtitle.IsLock, subtitle.Status,
		subtitle.CheckSum, subtitle.SubtitleURL, subtitle.PubTime, subtitle.RejectComment); err != nil {
		log.Error("params(%+v),error(%v)", subtitle, err)
		return
	}
	if insertID, err = res.LastInsertId(); err != nil {
		log.Error("params(%+v),error(%v)", subtitle, err)
		return
	}
	return
}

// UpdateSubtitle .
func (d *Dao) UpdateSubtitle(c context.Context, subtitle *model.Subtitle) (err error) {
	if _, err = d.dbDM.Exec(c, fmt.Sprintf(_updateSubtitle, d.hitSubtitle(subtitle.Oid)),
		subtitle.Aid, subtitle.AuthorID, subtitle.UpMid, subtitle.IsSign, subtitle.IsLock, subtitle.Status, subtitle.CheckSum, subtitle.SubtitleURL, subtitle.PubTime, subtitle.RejectComment,
		subtitle.ID); err != nil {
		log.Error("params(%+v),error(%v)", subtitle, err)
		return
	}
	return
}

// TxUpdateSubtitle .
func (d *Dao) TxUpdateSubtitle(tx *xsql.Tx, subtitle *model.Subtitle) (err error) {
	if _, err = tx.Exec(fmt.Sprintf(_updateSubtitle, d.hitSubtitle(subtitle.Oid)),
		subtitle.Aid, subtitle.AuthorID, subtitle.UpMid, subtitle.IsSign, subtitle.IsLock, subtitle.Status, subtitle.CheckSum, subtitle.SubtitleURL, subtitle.PubTime, subtitle.RejectComment,
		subtitle.ID); err != nil {
		log.Error("params(%+v),error(%v)", subtitle, err)
		return
	}
	return
}

// TxAddSubtitlePub .
func (d *Dao) TxAddSubtitlePub(tx *xsql.Tx, subtitlePub *model.SubtitlePub) (err error) {
	if _, err = tx.Exec(_addSubtitlePub, subtitlePub.Oid, subtitlePub.Type, subtitlePub.Lan, subtitlePub.SubtitleID, subtitlePub.IsDelete, subtitlePub.SubtitleID, subtitlePub.IsDelete); err != nil {
		log.Error("params(%+v),error(%v)", subtitlePub, err)
		return
	}
	return
}

// TxGetSubtitleOne .
func (d *Dao) TxGetSubtitleOne(tx *xsql.Tx, oid int64, tp int32, lan uint8) (subtitle *model.Subtitle, err error) {
	subtitle = &model.Subtitle{}
	row := tx.QueryRow(fmt.Sprintf(_getSubtitleOne, d.hitSubtitle(oid)), oid, tp, lan)
	if err = row.Scan(&subtitle.ID, &subtitle.Oid, &subtitle.Type, &subtitle.Lan, &subtitle.Aid, &subtitle.Mid, &subtitle.AuthorID, &subtitle.UpMid, &subtitle.IsSign, &subtitle.IsLock, &subtitle.Status, &subtitle.CheckSum, &subtitle.SubtitleURL, &subtitle.PubTime, &subtitle.RejectComment); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			subtitle = nil
			return
		}
		log.Error("params(oid:%v, tp:%v, lan:%v),error(%v)", oid, tp, lan, err)
		return
	}
	return
}

// SubtitleLans .
func (d *Dao) SubtitleLans(c context.Context) (subtitleLans []*model.SubtitleLan, err error) {
	rows, err := d.dbDM.Query(c, _getSubtitleLans)
	if err != nil {
		log.Error("params(query:%v),error(%v)", _getSubtitleLans, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		subtitleLan := new(model.SubtitleLan)
		if err = rows.Scan(&subtitleLan.Code, &subtitleLan.Lan, &subtitleLan.DocZh); err != nil {
			log.Error("params.Sacn(query:%v),error(%v)", _getSubtitleLans, err)
			return
		}
		subtitleLans = append(subtitleLans, subtitleLan)
	}
	if err = rows.Err(); err != nil {
		log.Error("params.Err(query:%v),error(%v)", _getSubtitleLans, err)
		return
	}
	return
}

// SubtitleLanAdd .
func (d *Dao) SubtitleLanAdd(c context.Context, subtitleLan *model.SubtitleLan) (err error) {
	if _, err = d.dbDM.Exec(c, _addSubtitleLan, subtitleLan.Code, subtitleLan.Lan, subtitleLan.DocZh, subtitleLan.DocEn, subtitleLan.IsDelete, subtitleLan.Lan, subtitleLan.DocZh, subtitleLan.DocEn, subtitleLan.IsDelete); err != nil {
		log.Error("SubtitleLanAdd.params(subtitleLan:%+v),error(%v)", subtitleLan, err)
		return
	}
	return
}

// UpsertWaveFrom .
func (d *Dao) UpsertWaveFrom(c context.Context, waveForm *model.WaveForm) (err error) {
	if _, err = d.dbDM.Exec(c, _addWaveForm, waveForm.Oid, waveForm.Type, waveForm.State, waveForm.WaveFromURL, waveForm.State, waveForm.WaveFromURL); err != nil {
		log.Error("params(waveForm:%+v),error(%v)", waveForm, err)
		return
	}
	return
}

// GetWaveForm .
func (d *Dao) GetWaveForm(c context.Context, oid int64, tp int32) (waveForm *model.WaveForm, err error) {
	var t time.Time
	row := d.dbDM.QueryRow(c, _getWaveForm, oid, tp)
	waveForm = &model.WaveForm{}
	if err = row.Scan(&waveForm.Oid, &waveForm.Type, &waveForm.State, &waveForm.WaveFromURL, &t); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			waveForm = nil
			return
		}
		log.Error("params(oid:%v, tp:%v),error(%v)", oid, tp, err)
		return
	}
	waveForm.Mtime = t.Unix()
	return
}
