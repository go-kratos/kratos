package dao

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"go-common/app/job/main/growup/model"

	"go-common/library/ecode"
	"go-common/library/log"
)

const (
	// insert
	_inBubbleMetaSQL   = "INSERT INTO lottery_av_info(av_id, date, b_type) VALUES %s ON DUPLICATE KEY UPDATE date=VALUES(date)"
	_inBubbleIncomeSQL = "INSERT INTO lottery_av_income(av_id,mid,tag_id,upload_time,total_income,income,tax_money,date,base_income,b_type) VALUES %s ON DUPLICATE KEY UPDATE av_id=VALUES(av_id),mid=VALUES(mid),tag_id=VALUES(tag_id),upload_time=VALUES(upload_time),total_income=VALUES(total_income),income=VALUES(income),tax_money=VALUES(tax_money),date=VALUES(date),base_income=VALUES(base_income)"
)

// InsertBubbleMeta insert lottery avs
func (d *Dao) InsertBubbleMeta(c context.Context, values string) (rows int64, err error) {
	if values == "" {
		return
	}
	res, err := d.db.Exec(c, fmt.Sprintf(_inBubbleMetaSQL, values))
	if err != nil {
		log.Error("dao.InsertBubbleMeta exec error(%v)", err)
		return
	}
	return res.RowsAffected()
}

// InsertBubbleIncome insert lottery income
func (d *Dao) InsertBubbleIncome(c context.Context, vals string) (rows int64, err error) {
	if vals == "" {
		return
	}
	res, err := d.db.Exec(c, fmt.Sprintf(_inBubbleIncomeSQL, vals))
	if err != nil {
		log.Error("dao.InsertBubbleIncome exec error(%v)", err)
		return
	}
	return res.RowsAffected()
}

// GetLotteryRIDs get lottery id
func (d *Dao) GetLotteryRIDs(c context.Context, start, end, offset int64) (info *model.LotteryRes, err error) {
	info = &model.LotteryRes{}
	params := url.Values{}
	params.Set("begin", strconv.FormatInt(start, 10))
	params.Set("end", strconv.FormatInt(end, 10))
	params.Set("offset", strconv.FormatInt(offset, 10))

	var res struct {
		Code    int               `json:"code"`
		Message string            `json:"message"`
		Data    *model.LotteryRes `json:"data"`
	}

	uri := d.avLotteryURL
	if err = d.client.Get(c, uri, "", params, &res); err != nil {
		log.Error("d.client.Get uri(%s) error(%v)", uri+"?"+params.Encode(), err)
		return
	}
	if res.Code != 0 {
		log.Error("GetLotteryRIDs code != 0. res.Code(%d) | uri(%s) error(%v)", res.Code, uri+"?"+params.Encode(), err)
		err = ecode.GrowupGetActivityError
		return
	}
	info = res.Data
	return
}

// VoteBIZArchive fetch avs in vote biz
func (d *Dao) VoteBIZArchive(c context.Context, start, end int64) (data []*model.VoteBIZArchive, err error) {
	var (
		uri    = d.avVoteURL
		params = url.Values{}
		res    struct {
			Code    int                     `json:"code"`
			Message string                  `json:"message"`
			Data    []*model.VoteBIZArchive `json:"data"`
		}
	)
	params.Set("start", strconv.FormatInt(start, 10))
	params.Set("end", strconv.FormatInt(end, 10))
	if err = d.client.Get(c, uri, "", params, &res); err != nil {
		log.Error("d.client.Get uri(%s) error(%v)", uri+"?"+params.Encode(), err)
		return
	}
	if res.Code != 0 {
		log.Error("VoteBIZArchive code != 0. res.Code(%d) | uri(%s) errorMsg(%s)", res.Code, uri+"?"+params.Encode(), res.Message)
		err = ecode.Errorf(ecode.ServerErr, "获取参与投票的视频失败")
		return
	}
	data = res.Data
	return
}
