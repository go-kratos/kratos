package service

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"go-common/app/job/main/member/model"
	smodel "go-common/app/service/main/member/model"
	"go-common/library/log"
)

func (s *Service) recoverMoral(c context.Context, mid int64) (err error) {
	var (
		moral        *smodel.Moral
		rowsAffected int64
	)
	if moral, err = s.dao.SelMoral(c, mid); err != nil {
		return
	}
	if moral == nil || moral.Moral >= smodel.DefaultMoral {
		log.Info("recoverMoral ignore, moral(%v)", moral)
		return
	}
	now := time.Now()
	deltaDays := deltaDays(now, moral.LastRecoverDate.Time())
	if deltaDays <= 0 {
		return
	}
	deltaMoral := deltaDays * 100
	if moral.Moral+deltaMoral > smodel.DefaultMoral {
		deltaMoral = smodel.DefaultMoral - moral.Moral
	}
	if rowsAffected, err = s.dao.RecoverMoral(c, mid, deltaMoral, deltaMoral, now.Format("2006-01-02")); err != nil {
		return
	}
	if rowsAffected == 0 {
		return
	}

	// report log
	ul := &model.UserLog{
		Mid:   mid,
		IP:    "127.0.0.1",
		TS:    now.Unix(),
		LogID: model.UUID4(),
		Content: map[string]string{
			"from_moral": strconv.FormatInt(moral.Moral, 10),
			"to_moral":   strconv.FormatInt(deltaMoral+moral.Moral, 10),
			"origin":     strconv.FormatInt(smodel.ManualRecoveryType, 10),
			"status":     strconv.FormatInt(smodel.RevocableMoralStatus, 10),
			"remark":     fmt.Sprintf("自动恢复(%d)天", deltaDays),
			"operater":   "系统",
			"reason":     fmt.Sprintf("时间：%d天", deltaDays),
		},
	}
	s.dao.AddMoralLog(c, ul)

	// origin log
	// content := make(map[string][]byte, 10)
	// content["mid"] = []byte(strconv.FormatInt(mid, 10))
	// content["ip"] = []byte("127.0.0.1")
	// content["from_moral"] = []byte(strconv.FormatInt(moral.Moral, 10))
	// content["to_moral"] = []byte(strconv.FormatInt(deltaMoral+moral.Moral, 10))
	// content["origin"] = []byte(strconv.FormatInt(smodel.ManualRecoveryType, 10))
	// content["status"] = []byte(strconv.FormatInt(smodel.RevocableMoralStatus, 10))
	// content["remark"] = []byte(fmt.Sprintf("自动恢复(%d)天", deltaDays))
	// content["operater"] = []byte("系统")
	// content["reason"] = []byte(fmt.Sprintf("时间：%d天", deltaDays))
	// if err = s.dao.AddLog(c, mid, now.Unix(), content, model.TableMoralLog); err != nil {
	// 	log.Error("recoverMoral mid: %v from_moral:%v,to_moral:%v  error(%v)", mid, content["from_moral"], content["to_moral"], err)
	// } else {
	// 	log.Info("recoverMoral mid: %v from_moral:%v,to_moral:%v  error(%v)", mid, content["from_moral"], content["to_moral"], err)
	// }
	s.dao.DelMoralCache(c, mid)
	return
}

func deltaDays(now, old time.Time) int64 {
	return int64(now.Sub(old).Hours() / 24)
}
