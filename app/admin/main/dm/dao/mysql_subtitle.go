package dao

import (
	"context"
	"fmt"
	"time"

	"go-common/app/admin/main/dm/model"
	"go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/xstr"
)

const (
	_subtitleSharding = 100

	_getSubtitles       = "SELECT id,oid,type,lan,aid,mid,author_mid,up_mid,is_sign,is_lock,status,subtitle_url,pub_time,mtime FROM subtitle_%02d WHERE id in (%s)"
	_getSubtitle        = "SELECT id,oid,type,lan,aid,mid,author_mid,up_mid,is_sign,is_lock,status,subtitle_url,pub_time,mtime FROM subtitle_%02d WHERE id=?"
	_countSubtitleDraft = "SELECT COUNT(id) FROM subtitle_%02d WHERE oid=? AND type=? AND lan=? AND mid=? AND pub_time=0"

	_updateSubtitle = "UPDATE subtitle_%02d set status=?,pub_time=? where id=?"

	_getSubtitlePubID = "SELECT id FROM subtitle_%02d WHERE oid=? AND type=? AND lan=? AND status=5 ORDER BY pub_time DESC limit 1"
	_addSubtitlePub   = "INSERT INTO subtitle_pub(oid,type,lan,subtitle_id,is_delete) VALUES(?,?,?,?,?) ON DUPLICATE KEY UPDATE subtitle_id=?,is_delete=?"

	_getSubtitleLans = "SELECT code,lan,doc_zh FROM subtitle_language WHERE is_delete=0"

	_addSubtitleSubject = "INSERT INTO subtitle_subject(aid,allow,attr,lan_code) VALUES(?,?,?,?) ON DUPLICATE KEY UPDATE allow=?,attr=?,lan_code=?"
	_getSubtitleSubject = "SELECT aid,allow,attr,lan_code from subtitle_subject WHERE aid=?"
)

func (d *Dao) hitSubtitle(oid int64) int64 {
	return oid % _subtitleSharding
}

// GetSubtitles .
func (d *Dao) GetSubtitles(c context.Context, oid int64, subtitleIds []int64) (subtitles []*model.Subtitle, err error) {
	rows, err := d.biliDM.Query(c, fmt.Sprintf(_getSubtitles, d.hitSubtitle(oid), xstr.JoinInts(subtitleIds)))
	if err != nil {
		log.Error("params(subtitleIds:%v),error(%v)", subtitleIds, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		subtitle := &model.Subtitle{}
		var t time.Time
		if err = rows.Scan(&subtitle.ID, &subtitle.Oid, &subtitle.Type, &subtitle.Lan, &subtitle.Aid, &subtitle.Mid, &subtitle.AuthorID, &subtitle.UpMid, &subtitle.IsSign, &subtitle.IsLock, &subtitle.Status, &subtitle.SubtitleURL, &subtitle.PubTime, &t); err != nil {
			log.Error("params(subtitleIds:%v),error(%v)", subtitleIds, err)
			return
		}
		subtitle.Mtime = t.Unix()
		subtitles = append(subtitles, subtitle)
	}
	if err = rows.Err(); err != nil {
		log.Error("params(subtitleIds:%v),error(%v)", subtitleIds, err)
		return
	}
	return
}

// GetSubtitle .
func (d *Dao) GetSubtitle(c context.Context, oid, subtitleID int64) (subtitle *model.Subtitle, err error) {
	row := d.biliDM.QueryRow(c, fmt.Sprintf(_getSubtitle, d.hitSubtitle(oid)), subtitleID)
	subtitle = &model.Subtitle{}
	var t time.Time
	if err = row.Scan(&subtitle.ID, &subtitle.Oid, &subtitle.Type, &subtitle.Lan, &subtitle.Aid, &subtitle.Mid, &subtitle.AuthorID, &subtitle.UpMid, &subtitle.IsSign, &subtitle.IsLock, &subtitle.Status, &subtitle.SubtitleURL, &subtitle.PubTime, &t); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			subtitle = nil
			return
		}
		log.Error("params(subtitleIds:%v),error(%v)", subtitleID, err)
		return
	}
	subtitle.Mtime = t.Unix()
	return
}

// UpdateSubtitle .
func (d *Dao) UpdateSubtitle(c context.Context, subtitle *model.Subtitle) (err error) {
	if _, err = d.biliDM.Exec(c, fmt.Sprintf(_updateSubtitle, d.hitSubtitle(subtitle.Oid)), subtitle.Status, subtitle.PubTime, subtitle.ID); err != nil {
		log.Error("UpdateSubtitle.params(subtitle:%+v),error(%v)", subtitle, err)
		return
	}
	return
}

// CountSubtitleDraft .
func (d *Dao) CountSubtitleDraft(c context.Context, oid int64, mid int64, lan uint8, tp int32) (count int64, err error) {
	row := d.biliDM.QueryRow(c, fmt.Sprintf(_countSubtitleDraft, d.hitSubtitle(oid)), oid, tp, lan, mid)
	if err = row.Scan(&count); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			count = 0
			return
		}
		log.Error("params(oid:%v, tp:%v, mid:%v, lan:%v),error(%v)", oid, tp, mid, lan, err)
		return
	}
	return
}

// TxUpdateSubtitle .
func (d *Dao) TxUpdateSubtitle(tx *sql.Tx, subtitle *model.Subtitle) (err error) {
	if _, err = tx.Exec(fmt.Sprintf(_updateSubtitle, d.hitSubtitle(subtitle.Oid)), subtitle.Status, subtitle.PubTime, subtitle.ID); err != nil {
		log.Error("TxUpdateSubtitle.params(subtitle:%+v),error(%v)", subtitle, err)
		return
	}
	return
}

// TxGetSubtitleID .
func (d *Dao) TxGetSubtitleID(tx *sql.Tx, oid int64, tp int32, lan uint8) (subtitleID int64, err error) {
	row := tx.QueryRow(fmt.Sprintf(_getSubtitlePubID, d.hitSubtitle(oid)), oid, tp, lan)
	if err = row.Scan(&subtitleID); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			subtitleID = 0
			return
		}
		log.Error("TxGetSubtitleID.Scac(err:%v)", err)
		return
	}
	return
}

// TxUpdateSubtitlePub .
func (d *Dao) TxUpdateSubtitlePub(tx *sql.Tx, subtitlePub *model.SubtitlePub) (err error) {
	if _, err = tx.Exec(_addSubtitlePub, subtitlePub.Oid, subtitlePub.Type, subtitlePub.Lan, subtitlePub.SubtitleID, subtitlePub.IsDelete, subtitlePub.SubtitleID, subtitlePub.IsDelete); err != nil {
		log.Error("params(%+v),error(%v)", subtitlePub, err)
		return
	}
	return
}

// SubtitleLans .
func (d *Dao) SubtitleLans(c context.Context) (subtitleLans []*model.SubtitleLan, err error) {
	rows, err := d.biliDM.Query(c, _getSubtitleLans)
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

// AddSubtitleSubject .
func (d *Dao) AddSubtitleSubject(c context.Context, subtitleSubject *model.SubtitleSubject) (err error) {
	if _, err = d.biliDM.Exec(c, _addSubtitleSubject, subtitleSubject.Aid, subtitleSubject.Allow, subtitleSubject.Attr, subtitleSubject.Lan, subtitleSubject.Allow, subtitleSubject.Attr, subtitleSubject.Lan); err != nil {
		log.Error("params(subtitleSubject:%+v),error(%v)", subtitleSubject, err)
		return
	}
	return
}

// GetSubtitleSubject .
func (d *Dao) GetSubtitleSubject(c context.Context, aid int64) (subtitleSubject *model.SubtitleSubject, err error) {
	subtitleSubject = new(model.SubtitleSubject)
	row := d.biliDM.QueryRow(c, _getSubtitleSubject, aid)
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
