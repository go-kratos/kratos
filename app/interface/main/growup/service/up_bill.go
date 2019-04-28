package service

import (
	"context"
	"fmt"
	"time"

	"go-common/app/interface/main/growup/model"

	"go-common/library/database/sql"
	"go-common/library/log"
	xtime "go-common/library/time"
)

// UpBill get up bill
func (s *Service) UpBill(c context.Context, mid int64) (data interface{}, err error) {
	redisKey := fmt.Sprintf("growup-up-bill-v1:%d", mid)
	res, err := s.dao.GetIncomeCache(c, redisKey)
	if err != nil {
		log.Error("s.dao.GetIncomeCache error(%v)", err)
		return
	}
	if res != nil {
		data = res["data"]
		return
	}

	data, err = s.upBill(c, mid)
	if err != nil {
		log.Error("s.upBill error(%v)", err)
		return
	}

	err = s.dao.SetIncomeCache(c, redisKey, map[string]interface{}{"data": data})
	if err != nil {
		log.Error("s.dao.SetIncomeCache error(%v)", err)
	}
	return
}

func (s *Service) upBill(c context.Context, mid int64) (data interface{}, err error) {
	up := &model.UpBill{}
	// 判断up主是否在创作激励
	signedAt, err := s.dao.GetUpSignedAt(c, "up_info_video", mid)
	if err != nil {
		log.Error("s.dao.GetUpSignedAt error(%v)", err)
		return
	}
	if signedAt == 0 {
		up.Join = false
		data = up
		return
	}
	up.Join = true
	endAt := time.Date(2018, 10, 31, 0, 0, 0, 0, time.Local)
	if signedAt >= xtime.Time(endAt.AddDate(0, 0, 1).Unix()) {
		up.SignedAt = signedAt
		up.EndAt = xtime.Time(endAt.Unix())
		data = up
		return
	}

	up, err = s.dao.GetUpBill(c, mid)
	up.Join = true
	if err == sql.ErrNoRows {
		err = nil
		up.SignedAt = signedAt
		up.EndAt = xtime.Time(endAt.Unix())
		data = up
		return
	}
	if err != nil {
		log.Error("s.dao.GetUpBill error(%v)", err)
		return
	}
	title, err := s.getAvTitle(c, []int64{up.AvID})
	if err != nil {
		log.Error("s.getAvTitle error(%v)", err)
		return
	}
	up.AvTitle = title[up.AvID]

	upsMap, err := s.dao.AccountInfos(c, []int64{mid})
	if err != nil {
		log.Error("s.dao.AccountInfos error(%v)", err)
		return
	}
	if up.Title == "流量王" {
		up.Title = "激励101"
	}

	if info, ok := upsMap[mid]; ok {
		up.Nickname = info.Nickname
		up.Face = info.Face
	}
	data = up
	return
}
