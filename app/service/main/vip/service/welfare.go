package service

import (
	"context"
	"fmt"
	"go-common/library/log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"go-common/app/service/main/vip/model"
	"go-common/library/ecode"
	xtime "go-common/library/time"

	"github.com/pkg/errors"
)

const (
	_onlyOnce     = 1
	_onceMonth    = 2
	_noTime       = -62135596800
	_zeroTime     = 0
	_expired      = 0
	_notVip       = 0
	_vip          = 1
	_yearVip      = 2
	_copyCode     = 2
	_redirectUrl  = 3
	_notRecommend = 0
	_recommend    = 1
	_both         = 2
)

// WelfareList get welfare list
func (s *Service) WelfareList(c context.Context, arg *model.ArgWelfareList) (res []*model.WelfareListResp, count int64, err error) {
	if arg.Pn == 0 {
		arg.Pn = 1
	}
	if arg.Ps == 0 {
		arg.Ps = 4
	}
	arg.NowTime = xtime.Time(time.Now().Unix())
	if _recommend == arg.Recommend {
		if res, err = s.dao.GetRecommendWelfare(c, arg); err != nil {
			return
		}
		if count, err = s.dao.CountRecommendWelfare(c, arg); err != nil {
			return
		}
	} else if _notRecommend == arg.Recommend {
		if res, err = s.dao.GetWelfareList(c, arg); err != nil {
			return
		}
		if count, err = s.dao.CountWelfare(c, arg); err != nil {
			return
		}
	} else if _both == arg.Recommend {
		arg.Recommend = 1
		if res, err = s.dao.GetWelfareList(c, arg); err != nil {
			return
		}
		if count, err = s.dao.CountWelfare(c, arg); err != nil {
			return
		}
	} else {
		log.Error("WelfareList param recommend error, recommend == (%+v)", arg.Recommend)
		return
	}
	randomBFSHost := fmt.Sprintf(s.c.Property.WelfareBgHost, rand.Intn(3))
	for _, r := range res {
		r.HomepageUri = fmt.Sprintf("%v%v", randomBFSHost, r.HomepageUri)
		r.BackdropUri = fmt.Sprintf("%v%v", randomBFSHost, r.BackdropUri)
	}

	return
}

// WelfareTypeList get welfare type list
func (s *Service) WelfareTypeList(c context.Context) (res []*model.WelfareTypeListResp, err error) {
	return s.dao.GetWelfareTypeList(c)
}

// WelfareInfo get welfare info
func (s *Service) WelfareInfo(c context.Context, arg *model.ArgWelfareInfo) (res *model.WelfareInfoResp, err error) {
	count := int64(0)
	if res, err = s.dao.GetWelfareInfo(c, arg.ID); err != nil {
		return
	}
	if res == nil {
		err = ecode.VipWelfareNotExist
		return
	}
	if arg.MID != 0 && res.UsageForm == _redirectUrl {
		if count, err = s.dao.CountReceiveRedirectWelfare(c, arg.ID, arg.MID); err != nil {
			return
		}
		if count > 0 {
			res.Received = true
		}
	}
	if arg.MID != 0 && res.UsageForm == _copyCode {
		var (
			batchResp     = []*model.WelfareBatchResp{}
			codeResp      = []*model.ReceivedCodeResp{}
			receivedCount int
			total         int
		)
		if batchResp, err = s.dao.GetWelfareBatch(c, arg.ID); err != nil {
			return
		}
		for _, br := range batchResp {
			if (br.Vtime.Time().Month() == res.Stime.Time().Month() && br.Vtime.Time().Year() == res.Stime.Time().Year()) ||
				(br.Vtime.Time().Month() == res.Etime.Time().Month() && br.Vtime.Time().Year() == res.Etime.Time().Year()) ||
				br.Vtime == _noTime || br.Vtime == _zeroTime {
				receivedCount += br.ReceivedCount
				total += br.Count
			}
		}
		if receivedCount == total {
			res.Finished = true
		}

		if codeResp, err = s.dao.GetReceivedCode(c, arg.ID, arg.MID); err != nil {
			return
		}
		if res.ReceiveRate == _onlyOnce && len(codeResp) != 0 {
			res.Received = true
		} else if res.ReceiveRate == _onceMonth {
			nowTime := time.Now()
			currentMonth := nowTime.Month()
			currentYear := nowTime.Year()
			for _, cr := range codeResp {
				if cr.Mtime.Time().Month() == currentMonth && cr.Mtime.Time().Year() == currentYear {
					res.Received = true
				}
			}
		}
	}
	randomBFSHost := fmt.Sprintf(s.c.Property.WelfareBgHost, rand.Intn(3))
	res.HomepageUri = fmt.Sprintf("%v%v", randomBFSHost, res.HomepageUri)
	res.BackdropUri = fmt.Sprintf("%v%v", randomBFSHost, res.BackdropUri)

	return
}

// WelfareReceive receive welfare
func (s *Service) WelfareReceive(c context.Context, arg *model.ArgWelfareReceive) (err error) {
	var (
		batchResp   []*model.WelfareBatchResp
		bids        []string
		welfareInfo *model.WelfareInfoResp
		codeResp    []*model.ReceivedCodeResp
		vipType     int32
		vipStatus   int32
		count       int64
	)
	if welfareInfo, err = s.dao.GetWelfareInfo(c, arg.Wid); err != nil {
		return
	}
	if welfareInfo == nil {
		err = ecode.VipWelfareNotExist
		return
	}
	//check vip_type
	if vipType, vipStatus, err = s.getVipType(c, arg.Mid); err != nil {
		return
	}
	if vipType == _notVip || vipStatus == _expired {
		err = ecode.VipWelfareVipOnly
		return
	}
	if vipType == _vip && welfareInfo.VipType == _yearVip {
		err = ecode.VipWelfareYearVipOnly
		return
	}

	//check time
	currentTime := time.Now()
	if currentTime.Unix() > welfareInfo.Etime.Time().Unix() {
		err = ecode.VipWelfareOffLine
		return
	}
	if currentTime.Unix() < welfareInfo.Stime.Time().Unix() {
		err = ecode.VipWelfareNotStart
		return
	}

	//check usage_form
	if welfareInfo.UsageForm == _redirectUrl {
		if count, err = s.dao.CountReceiveRedirectWelfare(c, arg.Wid, arg.Mid); err != nil {
			err = errors.WithStack(err)
			return
		}
		if count > 0 {
			err = ecode.VipWelfareAlreadyReceived
			return
		}
		if err = s.dao.AddReceiveRedirectWelfare(c, arg.Wid, arg.Mid); err != nil {
			err = errors.WithStack(err)
			return
		}
		return
	}

	//check receive rate
	currentMonth := currentTime.Month()
	currentYear := currentTime.Year()
	if codeResp, err = s.dao.GetReceivedCode(c, arg.Wid, arg.Mid); err != nil {
		return
	}
	if welfareInfo.ReceiveRate == _onlyOnce && len(codeResp) != 0 {
		err = ecode.VipWelfareAlreadyReceived
		return
	} else if welfareInfo.ReceiveRate == _onceMonth {
		for _, cr := range codeResp {
			if cr.Mtime.Time().Month() == currentMonth && cr.Mtime.Time().Year() == currentYear {
				err = ecode.VipWelfareAlreadyReceived
				return
			}
		}
	}

	if batchResp, err = s.dao.GetWelfareBatch(c, arg.Wid); err != nil {
		return
	}
	for _, br := range batchResp {
		if (br.Vtime.Time().Month() == currentMonth && br.Vtime.Time().Year() == currentYear) || (br.Vtime == _noTime || br.Vtime == _zeroTime) {
			bids = append(bids, strconv.Itoa(br.Id))
		}
	}
	if len(bids) == 0 {
		err = ecode.VipWelfareCodeRunOut
		return
	}
	for retry := 0; retry < 3; retry++ {
		if err = s.findCodeAndReceive(c, arg, bids, welfareInfo.ReceiveRate, currentTime); err != nil {
			if err == ecode.VipWelfareCodeUpdConflict {
				continue
			}
		}
		return
	}

	return
}

// MyWelfare get my welfare
func (s *Service) MyWelfare(c context.Context, mid int64) (res []*model.MyWelfareResp, err error) {
	if res, err = s.dao.GetMyWelfare(c, mid); err != nil {
		return
	}
	currentTime := time.Now().Unix()
	for _, r := range res {
		if currentTime > r.Etime.Time().Unix() {
			r.Expired = true
		}
	}
	return
}

func (s *Service) getVipType(c context.Context, mid int64) (vipType, vipStatus int32, err error) {
	var (
		v *model.VipInfo
	)
	if v, err = s.VipInfo(c, mid); err != nil {
		err = errors.WithStack(err)
		return
	}
	if v == nil {
		return
	}
	vipType = v.VipType
	vipStatus = v.VipStatus
	return
}

func (s *Service) findCodeAndReceive(c context.Context, arg *model.ArgWelfareReceive, bids []string, receiveRate int, currentTime time.Time) (err error) {
	var (
		codeNotUsedResp []*model.UnReceivedCodeResp
		random          int
		affectedRows    int64
		monthYear       int64
	)
	if codeNotUsedResp, err = s.dao.GetWelfareCodeUnReceived(c, arg.Wid, bids); err != nil {
		return
	}
	tx, err := s.dao.StartTx(c)
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	if len(codeNotUsedResp) == 0 {
		err = ecode.VipWelfareCodeRunOut
		return
	}
	random = rand.Intn(len(codeNotUsedResp))
	if affectedRows, err = s.dao.UpdateWelfareCodeUser(c, tx, codeNotUsedResp[random].Id, arg.Mid); err != nil {
		tx.Rollback()
		return
	}
	if affectedRows == 0 {
		tx.Rollback()
		err = ecode.VipWelfareCodeUpdConflict
		return
	}
	if err = s.dao.UpdateWelfareBatch(c, tx, codeNotUsedResp[random].Bid); err != nil {
		tx.Rollback()
		return
	}
	//insert record to prevent repeat receive
	if receiveRate == _onceMonth {
		monthYearStr := fmt.Sprintf("%v%v", currentTime.Year(), int(currentTime.Month()))
		monthYear, _ = strconv.ParseInt(monthYearStr, 10, 64)
	}
	if err = s.dao.InsertReceiveRecord(c, tx, arg.Mid, arg.Wid, monthYear); err != nil {
		tx.Rollback()
		if strings.Contains(err.Error(), "Error 1062: Duplicate entry") {
			err = ecode.VipWelfareAlreadyReceived
		}
		return
	}
	return tx.Commit()
}
