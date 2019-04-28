package service

import (
	"context"
	"time"

	"go-common/app/job/main/vip/model"
	"go-common/library/log"

	"github.com/pkg/errors"
)

// CheckUserData check vip_user_info data.
func (s *Service) CheckUserData(c context.Context) (diffs map[int64]string, err error) {
	var (
		maxID  int
		size   = s.c.Property.BatchSize
		ids    = []int64{}
		ousers = make(map[int64]*model.VipUserInfoOld, size)
		nusers = make(map[int64]*model.VipUserInfo, _ps)
		ou     *model.VipUserInfoOld
		nu     *model.VipUserInfo
		ok     bool
	)
	diffs = make(map[int64]string)
	if maxID, err = s.dao.SelOldUserInfoMaxID(context.TODO()); err != nil {
		err = errors.WithStack(err)
		return
	}
	page := maxID / size
	if maxID%size != 0 {
		page++
	}
	log.Info("check vip_user_info total(%d)", page)
	for i := 0; i < page; i++ {
		log.Info("check vip_user_info page index(%d) total(%d)", i, page)
		startID := i * size
		endID := (i + 1) * size
		if endID > maxID {
			endID = maxID
		}
		if ousers, err = s.dao.SelOldUserInfoMaps(context.TODO(), startID, endID); err != nil {
			return
		}
		j := 1
		for _, v := range ousers {
			ids = append(ids, v.Mid)
			if j%_ps == 0 || j == len(ousers) {
				if nusers, err = s.dao.SelVipByIds(context.TODO(), ids); err != nil {
					return
				}
				for _, mid := range ids {
					if ou, ok = ousers[mid]; !ok {
						diffs[mid] = "old not found"
						continue
					}
					if nu, ok = nusers[mid]; !ok {
						diffs[mid] = "new not found"
						continue
					}
					if nu.Type != ou.Type {
						diffs[mid] = "vip_type"
						continue
					}
					if nu.Status != ou.Status {
						diffs[mid] = "vip_status"
						continue
					}
					if !nu.OverdueTime.Time().Equal(ou.OverdueTime.Time()) {
						diffs[mid] = "vip_overdue_time"
						continue
					}
					if !nu.AnnualVipOverdueTime.Time().Equal(ou.AnnualVipOverdueTime.Time()) {
						diffs[mid] = "annual_vip_overdue_time"
						continue
					}
					if nu.PayType != ou.IsAutoRenew {
						diffs[mid] = "vip_pay_type"
						continue
					}
					if nu.PayChannelID != ou.PayChannelID {
						diffs[mid] = "pay_channel_id"
						continue
					}
					if !nu.IosOverdueTime.Time().Equal(ou.IosOverdueTime.Time()) {
						diffs[mid] = "ios_overdue_time"
						continue
					}
				}
				// reset
				ids = []int64{}
			}
			j++
		}
		log.Info("check index (%d) vip_user_info diff len (%d)", i, len(diffs))
		log.Info("check index (%d) vip_user_info diff data mids(%v)", i, diffs)
		time.Sleep(time.Millisecond * _defsleepmsec)
	}
	return
}

//CheckBcoinData check bcoin data
func (s *Service) CheckBcoinData(c context.Context) (mids []int64, err error) {
	var (
		maxID int
		size  = s.c.Property.BatchSize
	)
	if maxID, err = s.dao.SelMaxID(context.TODO()); err != nil {
		err = errors.WithStack(err)
		return
	}
	page := maxID / size
	if maxID%size != 0 {
		page++
	}

	for i := 0; i < page; i++ {
		startID := size * i
		endID := (i + 1) * size
		var res []*model.VipUserInfo
		if res, err = s.dao.SelUserInfos(context.TODO(), startID, endID); err != nil {
			err = errors.WithStack(err)
			return
		}
		var (
			tempMids    []int64
			bcoinMap    map[int64][]*model.VipBcoinSalary
			oldBcoinMap map[int64][]*model.VipBcoinSalary
		)
		for _, v := range res {
			tempMids = append(tempMids, v.Mid)
		}

		if bcoinMap, err = s.dao.SelBcoinSalaryDataMaps(context.TODO(), tempMids); err != nil {
			err = errors.WithStack(err)
			return
		}

		if oldBcoinMap, err = s.dao.SelOldBcoinSalaryDataMaps(context.TODO(), tempMids); err != nil {
			err = errors.WithStack(err)
			return
		}

		if len(bcoinMap) > len(oldBcoinMap) {
			for key, val := range bcoinMap {
				salaries := oldBcoinMap[key]
				if len(salaries) != len(val) {
					mids = append(mids, key)
				}
			}
		} else {
			for key, val := range oldBcoinMap {
				salaries := bcoinMap[key]
				if len(salaries) != len(val) {
					mids = append(mids, key)
				}
			}
		}
	}
	log.Info("cur not sync data mid is(%+v)", mids)
	return
}

//CheckChangeHistory check change history data
func (s *Service) CheckChangeHistory(c context.Context) (mids []int64, err error) {
	var (
		maxID int
		size  = 2000
	)
	if maxID, err = s.dao.SelMaxID(context.TODO()); err != nil {
		err = errors.WithStack(err)
		return
	}
	page := maxID / size
	if maxID%size != 0 {
		page++
	}

	for i := 0; i < page; i++ {
		startID := size * i
		endID := (i + 1) * size
		var res []*model.VipUserInfo
		if res, err = s.dao.SelUserInfos(context.TODO(), startID, endID); err != nil {
			err = errors.WithStack(err)
			return
		}

		var (
			tempMids      []int64
			historyMap    map[int64][]*model.VipChangeHistory
			oldHistoryMap map[int64][]*model.VipChangeHistory
		)

		for _, v := range res {
			tempMids = append(tempMids, v.Mid)
		}

		if historyMap, err = s.dao.SelChangeHistoryMaps(context.TODO(), tempMids); err != nil {
			err = errors.WithStack(err)
			return
		}
		if oldHistoryMap, err = s.dao.SelOldChangeHistoryMaps(context.TODO(), tempMids); err != nil {
			err = errors.WithStack(err)
			return
		}

		if len(historyMap) > len(oldHistoryMap) {
			for key, val := range historyMap {
				histories := oldHistoryMap[key]
				if len(histories) != len(val) {
					mids = append(mids, key)
				}
			}
		} else {
			for key, val := range oldHistoryMap {
				histories := historyMap[key]
				if len(histories) != len(val) {
					mids = append(mids, key)
				}
			}
		}

	}
	log.Info("cur not sync data mid is(%+v)", mids)
	return
}
