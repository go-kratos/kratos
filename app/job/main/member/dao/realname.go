package dao

import (
	"context"
	"net/http"
	"net/url"
	"time"

	"go-common/app/job/main/member/conf"
	"go-common/app/job/main/member/model"
	xsql "go-common/library/database/sql"
	"go-common/library/log"

	"github.com/pkg/errors"
)

var (
	_selRealnameAlipayApplyList = "SELECT id,mid,realname,card,img,status,reason,bizno,ctime,mtime FROM realname_alipay_apply WHERE id>? AND status=? AND mtime>=? AND mtime<=? ORDER BY id ASC LIMIT ?"

	_selRealnameInfo        = `SELECT id,mid,channel,realname,country,card_type,card,card_md5,status,reason,ctime,mtime FROM realname_info WHERE mid = ? LIMIT 1`
	_upsertRealnameInfo     = `INSERT INTO realname_info (mid,channel,realname,country,card_type,card,card_md5,status,reason) VALUES (?,?,?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE channel=?,realname=?,country=?,card_type=?,card=?,card_md5=?,status=?,reason=?`
	_upsertRealnameApply    = `INSERT INTO realname_apply (id,mid,realname,country,card_type,card_num,card_md5,hand_img,front_img,back_img,status,operator,operator_id,operator_time,remark,remark_status,ctime,mtime) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE realname=?,country=?,card_type=?,card_num=?,card_md5=?,hand_img=?,front_img=?,back_img=?,status=?,operator=?,operator_time=?,remark=?,remark_status=?,mtime=?`
	_upsertRealnameApplyIMG = `INSERT INTO realname_apply_img (id,img_data,ctime,mtime) VALUES (?,?,?,?) ON DUPLICATE KEY UPDATE img_data=?,mtime=?`
)

// UpdateRealnameFromMSG is
func (d *Dao) UpdateRealnameFromMSG(c context.Context, ms *model.RealnameApplyMessage) (err error) {
	var (
		tx           *xsql.Tx
		cardEncrpted = ms.CardData()
		cardMD5      = ms.CardMD5()
	)
	if tx, err = d.db.Begin(c); err != nil {
		return
	}
	if _, err = tx.Exec(_upsertRealnameApply, ms.ID, ms.MID, ms.Realname, ms.Country(), ms.CardType(), cardEncrpted, cardMD5, ms.FrontIMG, ms.FrontIMG2, ms.BackIMG, ms.Status, ms.Operater, 0, ms.OperaterTime(), ms.Remark, ms.RemarkStatus, ms.ApplyTime(), ms.ApplyTime(), ms.Realname, ms.Country(), ms.CardType(), cardEncrpted, cardMD5, ms.FrontIMG, ms.FrontIMG2, ms.BackIMG, ms.Status, ms.Operater, ms.OperaterTime(), ms.Remark, ms.RemarkStatus, time.Now()); err != nil {
		err = errors.WithStack(err)
		tx.Rollback()
		return
	}
	if _, err = tx.Exec(_upsertRealnameInfo, ms.MID, 0, ms.Realname, ms.Country(), ms.CardType(), cardEncrpted, cardMD5, ms.Status, ms.Remark, 0, ms.Realname, ms.Country(), ms.CardType(), cardEncrpted, cardMD5, ms.Status, ms.Remark); err != nil {
		err = errors.WithStack(err)
		tx.Rollback()
		return
	}
	if err = tx.Commit(); err != nil {
		return
	}
	return
}

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

// UpsertRealnameInfo is
func (d *Dao) UpsertRealnameInfo(c context.Context, ms *model.RealnameInfo) (err error) {
	if _, err = d.db.Exec(c, _upsertRealnameInfo, ms.MID, 0, ms.Realname, ms.Country, ms.CardType, ms.Card, ms.CardMD5, ms.Status, ms.Reason, 0, ms.Realname, ms.Country, ms.CardType, ms.Card, ms.CardMD5, ms.Status, ms.Reason); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

// UpsertRealnameApplyImg is
func (d *Dao) UpsertRealnameApplyImg(c context.Context, ms *model.RealnameApplyImgMessage) (err error) {
	if _, err = d.db.Exec(c, _upsertRealnameApplyIMG, ms.ID, ms.IMGData, ms.AddTime(), ms.AddTime(), ms.IMGData, time.Now()); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

// RealnameAlipayApplyList is
func (d *Dao) RealnameAlipayApplyList(c context.Context, startID int64, status model.RealnameApplyStatus, fromTime, toTime time.Time, limit int) (maxID int64, list []*model.RealnameAlipayApply, err error) {
	var (
		rows *xsql.Rows
	)
	if rows, err = d.db.Query(c, _selRealnameAlipayApplyList, startID, status, fromTime, toTime, limit); err != nil {
		err = errors.WithStack(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var (
			apply = &model.RealnameAlipayApply{}
		)
		if err = rows.Scan(&apply.ID, &apply.MID, &apply.Realname, &apply.Card, &apply.IMG, &apply.Status, &apply.Reason, &apply.Bizno, &apply.CTime, &apply.MTime); err != nil {
			err = errors.WithStack(err)
			return
		}
		if maxID < apply.ID {
			maxID = apply.ID
		}
		list = append(list, apply)
	}

	if err = rows.Err(); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

// AlipayQuery .
func (d *Dao) AlipayQuery(c context.Context, param url.Values) (pass bool, reason string, err error) {
	var (
		req *http.Request
	)
	url := conf.Conf.Biz.RealnameAlipayGateway + "?" + param.Encode()
	if req, err = http.NewRequest("GET", url, nil); err != nil {
		err = errors.Wrapf(err, "http.NewRequest(GET,%s)", url)
		return
	}
	var resp struct {
		Resp struct {
			respAlipay
			Passed          string `json:"passed"`
			FailedReason    string `json:"failed_reason"`
			IdentityInfo    string `json:"identity_info"`
			AttributeInfo   string `json:"attribute_info"`
			ChannelStatuses string `json:"channel_statuses"`
		} `json:"zhima_customer_certification_query_response"`
		Sign string `json:"sign"`
	}
	if err = d.client.Do(c, req, &resp); err != nil {
		return
	}
	log.Info("Realname alipay query \n\tparam : %+v \n\tresp : %+v", param, resp)
	if err = resp.Resp.Error(); err != nil {
		return
	}
	if resp.Resp.Passed == "true" {
		pass = true
	} else {
		pass = false
	}
	reason = resp.Resp.FailedReason
	return
}

type respAlipay struct {
	Code    string `json:"code"`
	Msg     string `json:"msg"`
	SubCode string `json:"sub_code"`
	SubMsg  string `json:"sub_msg"`
}

func (r *respAlipay) Error() error {
	if r.Code == "10000" {
		return nil
	}
	return errors.Errorf("alipay response failed , code : %s, msg : %s, sub_code : %s, sub_msg : %s", r.Code, r.Msg, r.SubCode, r.SubMsg)
}
