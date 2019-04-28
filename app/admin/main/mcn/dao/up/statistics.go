package up

import (
	"context"
	"fmt"
	"net/url"

	"go-common/app/admin/main/mcn/model"
	dtmdl "go-common/app/interface/main/mcn/model/datamodel"
	ifmdl "go-common/app/interface/main/mcn/model/mcnmodel"
	xsql "go-common/library/database/sql"
	"go-common/library/ecode"
	"go-common/library/net/metadata"
	xtime "go-common/library/time"

	"github.com/pkg/errors"
)

const (
	_mcnDataOverviewSQL             = `SELECT mcns, sign_ups, sign_ups_incr, fans_50, fans_10, fans_1, fans_incr_50, fans_incr_10, fans_incr_1 FROM mcn_data_overview WHERE generate_date = ?`
	_mcnRankFansOverviewSQL         = `SELECT id, sign_id, mid, data_view, data_type, rank, fans_incr, fans FROM mcn_rank_fans_overview WHERE data_view IN (1,2,3,4) AND data_type = ? AND generate_date = ? ORDER BY rank ASC limit ?`
	_mcnRankArchiveLikesOverviewSQL = `SELECT id, mcn_mid, up_mid, sign_id, avid, tid, rank, data_type, likes, plays FROM mcn_rank_archive_likes_overview WHERE data_type = ? AND generate_date = ? ORDER BY rank ASC limit ?`
	_mcnDataTypeSummarySQL          = `SELECT id, tid, data_view, data_type, amount FROM mcn_data_type_summary WHERE data_view IN (1,2,3,4) AND data_type IN (1,2) AND generate_date = ?`
	_arcTopURL                      = "/x/internal/mcn/rank/archive_likes"
	_dataFansURL                    = "/x/internal/mcn/data/fans"
	_dataFansBaseAttrURL            = "/x/internal/mcn/data/fans/base/attr"
	_dataFansAreaURL                = "/x/internal/mcn/data/fans/area"
	_dataFansTypeURL                = "/x/internal/mcn/data/fans/type"
	_dataFansTagURL                 = "/x/internal/mcn/data/fans/tag"
)

// McnDataOverview .
func (d *Dao) McnDataOverview(c context.Context, date xtime.Time) (m *model.McnDataOverview, err error) {
	row := d.db.QueryRow(c, _mcnDataOverviewSQL, date)
	m = new(model.McnDataOverview)
	if err = row.Scan(&m.Mcns, &m.SignUps, &m.SignUpsIncr, &m.Fans50, &m.Fans10, &m.Fans1, &m.FansIncr50, &m.FansIncr10, &m.FansIncr1); err != nil {
		if err == xsql.ErrNoRows {
			err = nil
			m = nil
			return
		}
	}
	return
}

// McnRankFansOverview .
func (d *Dao) McnRankFansOverview(c context.Context, dataType model.DataType, date xtime.Time, topLen int) (mrf map[int8][]*model.McnRankFansOverview, mids []int64, err error) {
	rows, err := d.db.Query(c, _mcnRankFansOverviewSQL, dataType, date, topLen*4)
	if err != nil {
		return
	}
	defer rows.Close()
	mrf = make(map[int8][]*model.McnRankFansOverview, topLen*4)
	for rows.Next() {
		rf := new(model.McnRankFansOverview)
		err = rows.Scan(&rf.ID, &rf.SignID, &rf.Mid, &rf.DataView, &rf.DataType, &rf.Rank, &rf.FansIncr, &rf.Fans)
		if err != nil {
			if err == xsql.ErrNoRows {
				err = nil
			}
			return
		}
		mids = append(mids, rf.Mid)
		mrf[rf.DataView] = append(mrf[rf.DataView], rf)
	}
	err = rows.Err()
	return
}

// McnRankArchiveLikesOverview .
func (d *Dao) McnRankArchiveLikesOverview(c context.Context, dataType model.DataType, date xtime.Time, topLen int) (ras []*model.McnRankArchiveLikesOverview, mids, avids, tids []int64, err error) {
	rows, err := d.db.Query(c, _mcnRankArchiveLikesOverviewSQL, dataType, date, topLen)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		ra := new(model.McnRankArchiveLikesOverview)
		err = rows.Scan(&ra.ID, &ra.McnMid, &ra.UpMid, &ra.SignID, &ra.Avid, &ra.Tid, &ra.Rank, &ra.DataType, &ra.Likes, &ra.Plays)
		if err != nil {
			if err == xsql.ErrNoRows {
				err = nil
			}
			return
		}
		ras = append(ras, ra)
		mids = append(mids, ra.McnMid)
		mids = append(mids, ra.UpMid)
		avids = append(avids, ra.Avid)
		tids = append(tids, int64(ra.Tid))
	}
	err = rows.Err()
	return
}

// McnDataTypeSummary .
func (d *Dao) McnDataTypeSummary(c context.Context, date xtime.Time) (mmd map[string][]*model.McnDataTypeSummary, tids []int64, err error) {
	rows, err := d.db.Query(c, _mcnDataTypeSummarySQL, date)
	if err != nil {
		return
	}
	defer rows.Close()
	mmd = make(map[string][]*model.McnDataTypeSummary)
	for rows.Next() {
		md := new(model.McnDataTypeSummary)
		err = rows.Scan(&md.ID, &md.Tid, &md.DataView, &md.DataType, &md.Amount)
		if err != nil {
			if err == xsql.ErrNoRows {
				err = nil
			}
			return
		}
		tids = append(tids, int64(md.Tid))
		mmd[fmt.Sprintf("%d-%d", md.DataView, md.DataType)] = append(mmd[fmt.Sprintf("%d-%d", md.DataView, md.DataType)], md)
	}
	err = rows.Err()
	return
}

// ArcTopDataStatistics .
func (d *Dao) ArcTopDataStatistics(c context.Context, arg *model.McnGetRankReq) (reply *model.McnGetRankUpFansReply, err error) {
	params := url.Values{}
	params.Set("sign_id", fmt.Sprintf("%d", arg.SignID))
	params.Set("tid", fmt.Sprintf("%d", arg.Tid))
	params.Set("data_type", fmt.Sprintf("%d", arg.DataType))
	var res struct {
		Code int                          `json:"code"`
		Data *model.McnGetRankUpFansReply `json:"data"`
	}
	if err = d.client.Get(c, d.arcTopURL, metadata.String(c, metadata.RemoteIP), params, &res); err != nil {
		return
	}
	if res.Code != 0 {
		err = errors.Wrapf(ecode.Int(res.Code), "arcRankFansTop d.client.Get(%s,%d)", d.arcTopURL+"?"+params.Encode(), res.Code)
	}
	reply = res.Data
	if reply != nil {
		for _, v := range reply.Result {
			v.PlayAccumulate = int64(v.Stat.View)
		}
	}
	return
}

// DataFans .
func (d *Dao) DataFans(c context.Context, arg *model.McnCommonReq) (reply *dtmdl.DmConMcnFansD, err error) {
	params := url.Values{}
	params.Set("sign_id", fmt.Sprintf("%d", arg.SignID))
	var res struct {
		Code int                       `json:"code"`
		Data *ifmdl.McnGetMcnFansReply `json:"data"`
	}
	if err = d.client.Get(c, d.dataFansURL, metadata.String(c, metadata.RemoteIP), params, &res); err != nil {
		return
	}
	if res.Code != 0 {
		err = errors.Wrapf(ecode.Int(res.Code), "dataFans d.client.Get(%s,%d)", d.dataFansURL+"?"+params.Encode(), res.Code)
	}
	if res.Data == nil {
		return
	}
	reply = &dtmdl.DmConMcnFansD{
		LogDate:    res.Data.LogDate,
		FansAll:    res.Data.FansAll,
		FansInc:    res.Data.FansInc,
		ActFans:    res.Data.ActFans,
		FansDecAll: res.Data.FansDecAll,
		FansDec:    res.Data.FansDec,
	}
	return
}

// DataFansBaseAttr .
func (d *Dao) DataFansBaseAttr(c context.Context, arg *model.McnCommonReq) (sex *dtmdl.DmConMcnFansSexW, age *dtmdl.DmConMcnFansAgeW, playWay *dtmdl.DmConMcnFansPlayWayW, err error) {
	params := url.Values{}
	params.Set("sign_id", fmt.Sprintf("%d", arg.SignID))
	params.Set("user_type", ifmdl.UserTypeFans)
	var res struct {
		Code int                            `json:"code"`
		Data *ifmdl.McnGetBaseFansAttrReply `json:"data"`
	}
	if err = d.client.Get(c, d.dataFansBaseAttrURL, metadata.String(c, metadata.RemoteIP), params, &res); err != nil {
		return
	}
	if res.Code != 0 {
		err = errors.Wrapf(ecode.Int(res.Code), "aataFansBaseAttr d.client.Get(%s,%d)", d.dataFansBaseAttrURL+"?"+params.Encode(), res.Code)
	}
	if res.Data == nil {
		return
	}
	sex = res.Data.FansSex
	age = res.Data.FansAge
	playWay = res.Data.FansPlayWay
	return
}

// DataFansArea .
func (d *Dao) DataFansArea(c context.Context, arg *model.McnCommonReq) (reply []*dtmdl.DmConMcnFansAreaW, err error) {
	params := url.Values{}
	params.Set("sign_id", fmt.Sprintf("%d", arg.SignID))
	params.Set("user_type", ifmdl.UserTypeFans)
	var res struct {
		Code int                        `json:"code"`
		Data *ifmdl.McnGetFansAreaReply `json:"data"`
	}
	if err = d.client.Get(c, d.dataFansAreaURL, metadata.String(c, metadata.RemoteIP), params, &res); err != nil {
		return
	}
	if res.Code != 0 {
		err = errors.Wrapf(ecode.Int(res.Code), "dataFansArea d.client.Get(%s,%d)", d.dataFansAreaURL+"?"+params.Encode(), res.Code)
	}
	if res.Data == nil {
		return
	}
	reply = res.Data.Result
	return
}

// DataFansType .
func (d *Dao) DataFansType(c context.Context, arg *model.McnCommonReq) (reply []*dtmdl.DmConMcnFansTypeW, err error) {
	params := url.Values{}
	params.Set("sign_id", fmt.Sprintf("%d", arg.SignID))
	params.Set("user_type", ifmdl.UserTypeFans)
	var res struct {
		Code int                        `json:"code"`
		Data *ifmdl.McnGetFansTypeReply `json:"data"`
	}
	if err = d.client.Get(c, d.dataFansTypeURL, metadata.String(c, metadata.RemoteIP), params, &res); err != nil {
		return
	}
	if res.Code != 0 {
		err = errors.Wrapf(ecode.Int(res.Code), "dataFansType d.client.Get(%s,%d)", d.dataFansTypeURL+"?"+params.Encode(), res.Code)
	}
	if res.Data == nil {
		return
	}
	reply = res.Data.Result
	return
}

// DataFansTag .
func (d *Dao) DataFansTag(c context.Context, arg *model.McnCommonReq) (reply []*dtmdl.DmConMcnFansTagW, err error) {
	params := url.Values{}
	params.Set("sign_id", fmt.Sprintf("%d", arg.SignID))
	params.Set("user_type", ifmdl.UserTypeFans)
	var res struct {
		Code int                       `json:"code"`
		Data *ifmdl.McnGetFansTagReply `json:"data"`
	}
	if err = d.client.Get(c, d.dataFansTagURL, metadata.String(c, metadata.RemoteIP), params, &res); err != nil {
		return
	}
	if res.Code != 0 {
		err = errors.Wrapf(ecode.Int(res.Code), "dataFansTag d.client.Get(%s,%d)", d.dataFansTagURL+"?"+params.Encode(), res.Code)
	}
	if res.Data == nil {
		return
	}
	reply = res.Data.Result
	return
}
