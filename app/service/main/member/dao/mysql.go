package dao

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"go-common/app/service/main/member/model"
	xsql "go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/time"

	"github.com/pkg/errors"
)

const (
	_shard = 100
	// base
	_selBaseInfo = "SELECT mid,name,sex,face,sign,rank,birthday FROM user_base_%02d WHERE mid=? "
	_setSex      = "INSERT INTO user_base_%02d(mid,sex) VALUES (?,?) ON DUPLICATE KEY UPDATE sex=?"
	_setName     = "INSERT INTO user_base_%02d(mid,name) VALUES (?,?) ON DUPLICATE KEY UPDATE name=?"
	_setRank     = "INSERT INTO user_base_%02d(mid,rank) VALUES (?,?) ON DUPLICATE KEY UPDATE rank=?"
	_setSign     = "INSERT INTO user_base_%02d(mid,sign) VALUES (?,?) ON DUPLICATE KEY UPDATE sign=?"
	_setBirthday = "INSERT INTO user_base_%02d(mid,birthday) VALUES (?,?) ON DUPLICATE KEY UPDATE birthday=?"
	_setFace     = "INSERT INTO user_base_%02d(mid,face) VALUES (?,?) ON DUPLICATE KEY UPDATE face=?"
	_setBase     = `INSERT INTO user_base_%02d(mid,name,sex,face,sign,rank,birthday) VALUES (?,?,?,?,?,?,?)
		ON DUPLICATE KEY UPDATE name=?,sex=?,face=?,sign=?,rank=?,birthday=?`
	// _exp
	_setExp    = "INSERT INTO user_exp_%02d (mid,exp) VALUES(?,?) ON DUPLICATE KEY UPDATE exp=VALUES(exp)"
	_updateExp = "UPDATE user_exp_%02d SET exp=exp+? WHERE mid=?"
	_selExp    = "SELECT exp FROM user_exp_%02d where mid=?"
	//user_flag
	_selUserAttr    = "SELECT flag FROM user_flag where mid =? and flag & ?= ?"
	_updateUserAttr = "INSERT INTO user_flag(mid,flag) VALUES(?,?) ON DUPLICATE KEY UPDATE flag=flag|?"
	// user_official
	_selOfficials = "SELECT mid,role,title,description FROM user_official WHERE role>0"
	// moral
	_selMoral = "SELECT moral,added,deducted,last_recover_date from user_moral where mid=?"

	_initMoral              = `INSERT IGNORE INTO user_moral (mid,moral,added,deducted,last_recover_date) VALUES(?,?,?,?,?)`
	_updateMoral            = `update user_moral set moral=moral+?,added=added+?, deducted=deducted + ? where mid = ?`
	_updateMoralRecoverDate = `update user_moral set last_recover_date = ? where mid = ?`

	// official
	_setOfficialDoc = `INSERT INTO user_official_doc (mid,name,state,role,title,description,extra,submit_source,submit_time) VALUES (?,?,?,?,?,?,?,?,?) 
	ON DUPLICATE KEY UPDATE name=VALUES(name), state=VALUES(state), role=VALUES(role), title=VALUES(title), description=VALUES(description), extra=VALUES(extra), submit_source=VALUES(submit_source), submit_time=VALUES(submit_time)`
	_selOfficialDoc = "SELECT mid,name,state,role,title,description,reject_reason,extra FROM user_official_doc WHERE mid=?"
	_selOfficial    = "SELECT role,title,description FROM user_official WHERE mid=? AND role>0"

	//official addit
	_setOfficialDocAddit = `INSERT INTO user_official_doc_addit (mid,property,vstring) VALUES (?,?,?) ON DUPLICATE KEY UPDATE vstring=VALUES(vstring)`

	// realname
	_selRealnameInfo           = `SELECT id,mid,channel,realname,country,card_type,card,card_md5,status,reason,ctime,mtime FROM realname_info WHERE mid = ? LIMIT 1`
	_selRealnameInfoByCard     = `SELECT id,mid,channel,realname,country,card_type,card,card_md5,status,reason,ctime,mtime FROM realname_info WHERE card_md5=? AND status in (0,1) LIMIT 1`
	_selRealnameInfoMidByCards = `SELECT mid,card_md5 FROM realname_info WHERE card_md5 IN (%s) AND status in (0,1)`
	_upsertRealnameInfo        = `INSERT INTO realname_info (mid,channel,realname,country,card_type,card,card_md5,status,reason) VALUES (?,?,?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE channel=?,realname=?,country=?,card_type=?,card=?,card_md5=?,status=?,reason=?`
	_updateRealnameInfo        = `UPDATE realname_info SET status=?,reason=? WHERE mid=?`
	_selRealnameApply          = "SELECT id,mid,realname,country,card_type,card_num,card_md5,hand_img,front_img,back_img,status,operator,operator_id,operator_time,remark,remark_status,ctime,mtime FROM realname_apply WHERE mid = ? ORDER BY ID DESC LIMIT 1"
	_insertRealnameApply       = `INSERT INTO realname_apply (mid,realname,country,card_type,card_num,card_md5,hand_img,front_img,back_img,status) VALUES (?,?,?,?,?,?,?,?,?,?)`
	_selRealnameApplyImg       = "SELECT id,img_data,ctime,mtime FROM realname_apply_img WHERE id = ? LIMIT 1"
	_insertRealnameApplyImg    = `INSERT INTO realname_apply_img (img_data) VALUES (?)`
	// alipay
	_selRealnameAlipayApply    = "SELECT id,mid,realname,card,img,status,reason,bizno,ctime,mtime FROM realname_alipay_apply WHERE mid = ? ORDER BY ID DESC LIMIT 1"
	_insertRealnameAlipayApply = `INSERT INTO realname_alipay_apply (mid,realname,card,img,status,reason,bizno) VALUES (?,?,?,?,?,?,?)`
	_updateRealnameAlipayApply = `UPDATE realname_alipay_apply SET status=?,reason=? WHERE id=?`

	// realname old
	_insertOldRealnameApply    = `INSERT INTO dede_identification_card_apply (mid,realname,type,card_data,card_for_search,front_img,front_img2,back_img,apply_time,status) VALUES (?,?,?,?,?,?,?,?,?,?)`
	_insertOldRealnameApplyImg = `INSERT INTO dede_identification_card_apply_img (img_data,add_time) VALUES (?,?)`
)

func hit(id int64) int64 {
	return id % _shard
}

// BeginTran begin transcation.
func (d *Dao) BeginTran(c context.Context) (tx *xsql.Tx, err error) {
	return d.db.Begin(c)
}

// BaseInfo base info of user.
func (d *Dao) BaseInfo(c context.Context, mid int64) (r *model.BaseInfo, err error) {
	r = &model.BaseInfo{}
	row := d.db.Master().QueryRow(c, fmt.Sprintf(_selBaseInfo, hit(mid)), mid)
	if err = row.Scan(&r.Mid, &r.Name, &r.Sex, &r.Face, &r.Sign, &r.Rank, &r.Birthday); err != nil {
		if err == xsql.ErrNoRows {
			r = nil
			err = nil
			return
		}
		err = errors.Wrapf(err, "dao base mid(%d)", mid)
	}
	r.RandFaceURL()
	return
}

// SetBase set base.
func (d *Dao) SetBase(c context.Context, base *model.BaseInfo) (err error) {
	if _, err = d.db.Exec(c, fmt.Sprintf(_setBase, hit(base.Mid)), base.Mid, base.Name, base.Sex, base.Face, base.Sign,
		base.Rank, base.Birthday, base.Name, base.Sex, base.Face, base.Sign, base.Rank, base.Birthday); err != nil {
		err = errors.Wrapf(err, "dao set base(%v)", base)
	}
	return
}

// SetSign set user sign.
func (d *Dao) SetSign(c context.Context, mid int64, sign string) (err error) {
	if _, err = d.db.Exec(c, fmt.Sprintf(_setSign, hit(mid)), mid, sign, sign); err != nil {
		err = errors.Wrapf(err, "dao set sign mid(%d)", mid)
	}
	return
}

// SetName set user name.
func (d *Dao) SetName(c context.Context, mid int64, name string) (err error) {
	if _, err = d.db.Exec(c, fmt.Sprintf(_setName, hit(mid)), mid, name, name); err != nil {
		err = errors.Wrapf(err, "dao set name mid(%d) name(%s)", mid, name)
	}
	return
}

// SetRank set user rank.
func (d *Dao) SetRank(c context.Context, mid, rank int64) (err error) {
	if _, err = d.db.Exec(c, fmt.Sprintf(_setRank, hit(mid)), mid, rank, rank); err != nil {
		err = errors.Wrapf(err, "dao set rank mid(%d) rank(%d)", mid, rank)
	}
	return
}

// SetSex set sex.
func (d *Dao) SetSex(c context.Context, mid, sex int64) (err error) {
	if _, err = d.db.Exec(c, fmt.Sprintf(_setSex, hit(mid)), mid, sex, sex); err != nil {
		err = errors.Wrapf(err, "dao set sex mid(%d) sex(%d)", mid, sex)
	}
	return
}

// SetBirthday set birthday.
func (d *Dao) SetBirthday(c context.Context, mid int64, birthday time.Time) (err error) {
	if _, err = d.db.Exec(c, fmt.Sprintf(_setBirthday, hit(mid)), mid, birthday, birthday); err != nil {
		err = errors.Wrapf(err, "dao set birthday mid(%d) birthday(%d)", mid, birthday)
	}
	return
}

// SetFace set face.
func (d *Dao) SetFace(c context.Context, mid int64, face string) (err error) {
	if _, err = d.db.Exec(c, fmt.Sprintf(_setFace, hit(mid)), mid, face, face); err != nil {
		err = errors.Wrapf(err, "dao set face mid(%d) face(%v)", mid, face)
	}
	return
}

// ExpDB get user exp from db.
func (d *Dao) ExpDB(c context.Context, mid int64) (count int64, err error) {
	row := d.db.QueryRow(c, fmt.Sprintf(_selExp, hit(mid)), mid)
	if err = row.Scan(&count); err != nil {
		if err == xsql.ErrNoRows {
			err = nil
			return
		}
		err = errors.Wrapf(err, "dao exp mid(%d)", mid)
	}
	return
}

// SetExp set user exp to count.
func (d *Dao) SetExp(c context.Context, mid, count int64) (affect int64, err error) {
	row, err := d.db.Exec(c, fmt.Sprintf(_setExp, hit(mid)), mid, count)
	if err != nil {
		err = errors.Wrapf(err, "dao set exp mid(%d) exp(%d)", mid, count)
		return
	}
	return row.RowsAffected()
}

// UpdateExp incr user exp by delta.
func (d *Dao) UpdateExp(c context.Context, mid, delta int64) (affect int64, err error) {
	row, err := d.db.Exec(c, fmt.Sprintf(_updateExp, hit(mid)), delta, mid)
	if err != nil {
		err = errors.Wrapf(err, "dao update exp mid(%d) exp(%d)", mid, delta)
		return
	}
	return row.RowsAffected()
}

// ------------- exp ---------------- //

// UserAttrDB get attr.
func (d *Dao) UserAttrDB(c context.Context, mid int64, attr uint) (hasAttr bool, err error) {
	var flag int8
	row := d.db.QueryRow(c, _selUserAttr, mid, attr, attr)
	if err = row.Scan(&flag); err != nil {
		if err == xsql.ErrNoRows {
			err = nil
			return
		}
		err = errors.Wrapf(err, "dao attr mid(%d) attr(%d)", mid, attr)
	}
	hasAttr = true
	return
}

// SetUserAttr update attr .
func (d *Dao) SetUserAttr(c context.Context, mid int64, attr uint) (err error) {
	if _, err = d.db.Exec(c, _updateUserAttr, mid, attr, attr); err != nil {
		err = errors.Wrapf(err, "dao attr mid(%d) attr(%d)", mid, attr)
		return
	}
	return
}

// Officials all officials info of user.
func (d *Dao) Officials(c context.Context) (om map[int64]*model.OfficialInfo, err error) {
	om = make(map[int64]*model.OfficialInfo)
	rows, err := d.db.Query(c, _selOfficials)
	if err != nil {
		err = errors.Wrap(err, "dao officials")
		return
	}
	defer rows.Close()
	for rows.Next() {
		var mid int64
		o := &model.OfficialInfo{}
		if err = rows.Scan(&mid, &o.Role, &o.Title, &o.Desc); err != nil {
			err = errors.Wrap(err, "dao officials scan")
			return
		}
		om[mid] = o
	}
	return
}

// Official is.
func (d *Dao) Official(c context.Context, mid int64) (*model.OfficialInfo, error) {
	row := d.db.QueryRow(c, _selOfficial, mid)
	o := &model.OfficialInfo{}
	if err := row.Scan(&o.Role, &o.Title, &o.Desc); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.Wrap(err, "dao official scan")
	}
	return o, nil
}

// MoralDB get user moral from db.
func (d *Dao) MoralDB(c context.Context, mid int64) (moral *model.Moral, err error) {
	moral = &model.Moral{}
	row := d.db.QueryRow(c, _selMoral, mid)
	if err = row.Scan(&moral.Moral, &moral.Added, &moral.Deducted, &moral.LastRecoverDate); err != nil {
		if err == sql.ErrNoRows {
			moral = nil
			err = nil
			return
		}
		log.Error(" SelMoral row.Scan() error(%v) mid(%v)", err, mid)
		err = errors.Wrapf(err, "dao moral mid(%d)", mid)
		return
	}
	return
}

// TxMoralDB get user moral from db.
func (d *Dao) TxMoralDB(tx *xsql.Tx, mid int64) (moral *model.Moral, err error) {
	moral = &model.Moral{}
	row := tx.QueryRow(_selMoral, mid)
	if err = row.Scan(&moral.Moral, &moral.Added, &moral.Deducted, &moral.LastRecoverDate); err != nil {
		if err == sql.ErrNoRows {
			moral = nil
			err = nil
			return
		}
		log.Error(" SelMoral row.Scan() error(%v) mid(%v)", err, mid)
		err = errors.Wrapf(err, "dao moral mid(%d)", mid)
		return
	}
	return
}

// TxUpdateMoral set user moral.
func (d *Dao) TxUpdateMoral(tx *xsql.Tx, mid, moral, added, deducted int64) (err error) {
	if _, err = tx.Exec(_updateMoral, moral, added, deducted, mid); err != nil {
		err = errors.Wrapf(err, "TxUpdateMoral mid(%d) moralAdded(%v)", mid, moral)
		return
	}
	return
}

// TxUpdateMoralRecoverDate update moral recover date.
func (d *Dao) TxUpdateMoralRecoverDate(tx *xsql.Tx, mid int64, recoverDate time.Time) (err error) {
	if _, err = tx.Exec(_updateMoralRecoverDate, recoverDate, mid); err != nil {
		err = errors.Wrapf(err, "TxUpdateMoralRecoverDate mid(%d) recoverDate(%v)", mid, recoverDate)
		return
	}
	return
}

// TxInitMoral set user moral.
func (d *Dao) TxInitMoral(tx *xsql.Tx, mid, moral, added, deducted int64, lastRecoverDate time.Time) (err error) {
	if _, err = tx.Exec(_initMoral, mid, moral, added, deducted, lastRecoverDate); err != nil {
		err = errors.Wrapf(err, "TxInitMoral mid(%d) moral(%v)", mid, moral)
		return
	}
	return
}

// SetOfficialDoc add official doc.
func (d *Dao) SetOfficialDoc(c context.Context, od *model.OfficialDoc) (err error) {
	_, err = d.db.Exec(c, _setOfficialDoc, od.Mid, od.Name, model.OfficialStateWait, od.Role,
		od.Title, od.Desc, od.OfficialExtra.String(), od.SubmitSource, od.SubmitTime)
	if err != nil {
		err = errors.Wrapf(err, "dao add official doc")
		return
	}

	return
}

// SetOfficialDocAddit add official doc addit.
func (d *Dao) SetOfficialDocAddit(c context.Context, mid int64, property, vstring string) error {
	_, err := d.db.Exec(c, _setOfficialDocAddit, mid, property, vstring)
	if err != nil {
		err = errors.Wrapf(err, "dao add official doc addit")
		return err
	}
	return nil
}

// OfficialDoc get official doc.
func (d *Dao) OfficialDoc(c context.Context, mid int64) (*model.OfficialDoc, error) {
	od := new(model.OfficialDoc)
	row := d.db.QueryRow(c, _selOfficialDoc, mid)
	if err := row.Scan(&od.Mid, &od.Name, &od.State, &od.Role, &od.Title, &od.Desc, &od.RejectReason, &od.Extra); err != nil {
		err = errors.Wrapf(err, "official doc")
		return nil, err
	}
	od.ParseExtra()
	return od, nil
}

// Realname

// RealnameInfo is.
func (d *Dao) RealnameInfo(c context.Context, mid int64) (info *model.RealnameInfo, err error) {
	row := d.db.QueryRow(c, _selRealnameInfo, mid)
	info = &model.RealnameInfo{}
	if err = row.Scan(&info.ID, &info.MID, &info.Channel, &info.Realname, &info.Country, &info.CardType, &info.Card, &info.CardMD5, &info.Status, &info.Reason, &info.CTime, &info.MTime); err != nil {
		if err == xsql.ErrNoRows {
			err = nil
			info = nil
			return
		}
		err = errors.Wrapf(err, "dao RealnameInfo mid(%d)", mid)
		return
	}
	return
}

// RealnameInfoByCard is.
func (d *Dao) RealnameInfoByCard(c context.Context, cardMD5 string) (info *model.RealnameInfo, err error) {
	row := d.db.QueryRow(c, _selRealnameInfoByCard, cardMD5)
	info = &model.RealnameInfo{}
	if err = row.Scan(&info.ID, &info.MID, &info.Channel, &info.Realname, &info.Country, &info.CardType, &info.Card, &info.CardMD5, &info.Status, &info.Reason, &info.CTime, &info.MTime); err != nil {
		if err == xsql.ErrNoRows {
			err = nil
			info = nil
			return
		}
		err = errors.Wrapf(err, "dao RealnameInfo cardMD5(%s)", cardMD5)
		return
	}
	return
}

// UpsertTxRealnameInfo is.
func (d *Dao) UpsertTxRealnameInfo(c context.Context, tx *xsql.Tx, info *model.RealnameInfo) (err error) {
	if _, err = tx.Exec(_upsertRealnameInfo, info.MID, info.Channel, info.Realname, info.Country, info.CardType, info.Card, info.CardMD5, info.Status, info.Reason, info.Channel, info.Realname, info.Country, info.CardType, info.Card, info.CardMD5, info.Status, info.Reason); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

// UpdateTxRealnameInfo is.
func (d *Dao) UpdateTxRealnameInfo(c context.Context, tx *xsql.Tx, mid int64, status model.RealnameApplyStatus, reason string) (err error) {
	if _, err = tx.Exec(_updateRealnameInfo, status, reason, mid); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

// RealnameApply realname
func (d *Dao) RealnameApply(c context.Context, mid int64) (apply *model.RealnameApply, err error) {
	row := d.db.Master().QueryRow(c, _selRealnameApply, mid)
	apply = &model.RealnameApply{}
	if err = row.Scan(&apply.ID, &apply.MID, &apply.Realname, &apply.Country, &apply.CardType, &apply.CardNum, &apply.CardMD5, &apply.HandIMG, &apply.FrontIMG, &apply.BackIMG, &apply.Status, &apply.Operator, &apply.OperatorID, &apply.OperatorTime, &apply.Remark, &apply.RemarkStatus, &apply.CTime, &apply.MTime); err != nil {
		if err == xsql.ErrNoRows {
			err = nil
			apply = nil
			return
		}
		err = errors.Wrapf(err, "dao RealnameApply mid(%d)", mid)
		return
	}
	return
}

// InsertRealnameApply is
func (d *Dao) InsertRealnameApply(c context.Context, data *model.RealnameApply) (err error) {
	if _, err = d.db.Exec(c, _insertRealnameApply, data.MID, data.Realname, data.Country, data.CardType, data.CardNum, data.CardMD5, data.HandIMG, data.FrontIMG, data.BackIMG, data.Status); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

// RealnameApplyIMG is
func (d *Dao) RealnameApplyIMG(c context.Context, id int) (img *model.RealnameApplyImage, err error) {
	row := d.db.Master().QueryRow(c, _selRealnameApplyImg, id)
	img = &model.RealnameApplyImage{}
	if err = row.Scan(&img.ID, &img.IMGData, &img.CTime, &img.MTime); err != nil {
		if err == xsql.ErrNoRows {
			err = nil
			img = nil
			return
		}
		err = errors.Wrapf(err, "dao RealnameApplyIMG mid(%d)", id)
		return
	}
	return
}

// InsertRealnameApplyImg is
func (d *Dao) InsertRealnameApplyImg(c context.Context, data *model.RealnameApplyImage) (id int64, err error) {
	var res sql.Result
	if res, err = d.db.Exec(c, _insertRealnameApplyImg, data.IMGData); err != nil {
		err = errors.WithStack(err)
		return
	}
	if id, err = res.LastInsertId(); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

// InsertOldRealnameApply is
func (d *Dao) InsertOldRealnameApply(c context.Context, data *model.RealnameApply) (id int64, err error) {
	var (
		res      sql.Result
		cardType = int16(data.CardType)
	)
	if data.Country > 0 {
		cardType = data.Country
	}
	if res, err = d.accdb.Exec(c, _insertOldRealnameApply, data.MID, data.Realname, cardType, data.CardNum, data.CardMD5, data.FrontIMG, data.HandIMG, data.BackIMG, data.CTime.Unix(), data.Status); err != nil {
		err = errors.WithStack(err)
		return
	}
	if id, err = res.LastInsertId(); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

// InsertOldRealnameApplyImg is
func (d *Dao) InsertOldRealnameApplyImg(c context.Context, data *model.RealnameApplyImage) (id int64, err error) {
	var res sql.Result
	if res, err = d.accdb.Exec(c, _insertOldRealnameApplyImg, data.IMGData, data.CTime); err != nil {
		err = errors.WithStack(err)
		return
	}
	if id, err = res.LastInsertId(); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

// InsertTxRealnameAlipayApply .
func (d *Dao) InsertTxRealnameAlipayApply(c context.Context, tx *xsql.Tx, data *model.RealnameAlipayApply) (err error) {
	if _, err = tx.Exec(_insertRealnameAlipayApply, data.MID, data.Realname, data.Card, data.IMG, data.Status, data.Reason, data.Bizno); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

// UpdateTxRealnameAlipayApply is.
func (d *Dao) UpdateTxRealnameAlipayApply(c context.Context, tx *xsql.Tx, id int64, status model.RealnameApplyStatus, reason string) (err error) {
	if _, err = tx.Exec(_updateRealnameAlipayApply, status, reason, id); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

// RealnameAlipayApply .
func (d *Dao) RealnameAlipayApply(c context.Context, mid int64) (apply *model.RealnameAlipayApply, err error) {
	row := d.db.QueryRow(c, _selRealnameAlipayApply, mid)
	apply = &model.RealnameAlipayApply{}
	if err = row.Scan(&apply.ID, &apply.MID, &apply.Realname, &apply.Card, &apply.IMG, &apply.Status, &apply.Reason, &apply.Bizno, &apply.CTime, &apply.MTime); err != nil {
		if err == xsql.ErrNoRows {
			err = nil
			apply = nil
			return
		}
		err = errors.Wrapf(err, "dao RealnameAlipayApply mid(%d)", mid)
		return
	}
	return
}

func prepareStringArray(in []string) string {
	surrounded := make([]string, 0, len(in))
	for _, s := range in {
		surrounded = append(surrounded, fmt.Sprintf(`'%s'`, s))
	}
	return strings.Join(surrounded, ",")
}

// MidByRealnameCards is
func (d *Dao) MidByRealnameCards(ctx context.Context, cardMD5s []string) (map[string]int64, error) {
	rows, err := d.db.Query(ctx, fmt.Sprintf(_selRealnameInfoMidByCards, prepareStringArray(cardMD5s)))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := make(map[string]int64, len(cardMD5s))
	for rows.Next() {
		mid, cardMD5 := int64(0), ""
		if err := rows.Scan(&mid, &cardMD5); err != nil {
			log.Warn("Failed to scan realname info: %+v", err)
			continue
		}
		result[cardMD5] = mid
	}
	return result, nil
}
