package service

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"go-common/app/service/main/member/model"
	"go-common/library/database/sql"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/metadata"
	ctime "go-common/library/time"
)

// Moral get user moral from db.
func (s *Service) Moral(c context.Context, mid int64) (moral *model.Moral, err error) {
	if moral, err = s.mbDao.Moral(c, mid); err != nil {
		log.Error("s.mbDao.Moral(%d) error(%v)", mid, err)
		return
	}
	return
}

// MoralLog get user moral log from hbase.
func (s *Service) MoralLog(c context.Context, mid int64) (logs []*model.UserLog, err error) {
	if logs, err = s.mbDao.MoralLog(c, mid); err != nil {
		log.Error("s.mbDao.MoralLog(%d) error(%v)", mid, err)
		return
	}
	return
}

// UndoMoral moralChange.
func (s *Service) UndoMoral(c context.Context, logID, remark, operator string) error {
	useLog, err := s.mbDao.MoralLogByID(c, logID)
	if err != nil {
		log.Error("s.mbDao.MoralLogByID(%v) error(%v)", logID, err)
		return err
	}
	oldContent := useLog.Content
	content := make(map[string]string, len(oldContent))
	for k, v := range oldContent {
		content[k] = v
	}
	content["status"] = strconv.FormatInt(model.RevokedMoralStatus, 10)
	if err = s.mbDao.DeleteMoralLog(c, useLog.LogID); err != nil {
		log.Error("Failed to delete moral log: log: %+v: %+v", useLog, err)
	}
	s.mbDao.AddMoralLogReport(c, &model.UserLog{
		Mid:     useLog.Mid,
		IP:      useLog.IP,
		TS:      useLog.TS,
		LogID:   useLog.LogID,
		Content: content,
	})

	// bcontent := make(map[string][]byte, len(content))
	// for k, v := range content {
	// 	bcontent[k] = []byte(v)
	// }
	// s.mbDao.AddMoralLog(c, useLog.Mid, useLog.TS, bcontent)

	arg := &model.ArgUpdateMoral{
		Mid:      useLog.Mid,
		Status:   model.IrrevocableMoralStatus,
		IP:       metadata.String(c, metadata.RemoteIP),
		Reason:   oldContent["reason"],
		Operator: operator,
		Remark:   remark,
	}
	if arg.Origin, err = strconv.ParseInt(oldContent["origin"], 10, 64); err != nil {
		log.Error("strconv.ParseInt(%v) error(%v)", oldContent["origin"], err)
		return err
	}
	fromMoral, err := strconv.ParseInt(oldContent["from_moral"], 10, 64)
	if err != nil {
		log.Error("strconv.ParseInt(%v) error(%v)", oldContent["from_moral"], err)
		return err
	}
	toMoral, err := strconv.ParseInt(oldContent["to_moral"], 10, 64)
	if err != nil {
		log.Error("strconv.ParseInt(%v) error(%v)", oldContent["to_moral"], err)
		return err
	}
	arg.Delta = fromMoral - toMoral
	arg.ReasonType = model.SysReasonType
	arg.IsNotify = false
	return s.UpdateMoral(c, arg)
}

// UpdateMorals batch update user moral .
func (s *Service) UpdateMorals(c context.Context, args *model.ArgUpdateMorals) (afterMorals map[int64]int64, err error) {
	var (
		before, after int64
		tx            *sql.Tx
		beforeMorals  map[int64]int64
		originType    *model.OriginType
		ok            bool
	)
	if originType, ok = model.OriginTypes[args.Origin]; !ok {
		err = ecode.RequestErr
		return
	}
	if originType.NeedReason && len(args.Reason) == 0 {
		err = ecode.RequestErr
		return
	}
	beforeMorals = make(map[int64]int64, len(args.Mids))
	afterMorals = make(map[int64]int64, len(args.Mids))
	if tx, err = s.mbDao.BeginTran(c); err != nil {
		return
	}
	for _, mid := range args.Mids {
		if before, after, err = s.updateMoral(tx, mid, args.Delta); err != nil {
			tx.Rollback()
			log.Error("s.updateMoral(%v) error(%v)", mid, err)
			return
		}
		beforeMorals[mid] = before
		afterMorals[mid] = after
	}
	if err = tx.Commit(); err != nil {
		log.Error("UpdateMoral tx.Commit(%v) error(%v)", args.Mids, err)
		return
	}
	ts := int64(time.Now().Unix())
	for mid, before := range beforeMorals {
		after = afterMorals[mid]

		// log to report
		l := &model.UserLog{
			Mid:   mid,
			IP:    args.IP,
			TS:    ts,
			LogID: model.UUID4(),
			Content: map[string]string{
				"from_moral": strconv.FormatInt(before, 10),
				"to_moral":   strconv.FormatInt(after, 10),
				"origin":     strconv.FormatInt(args.Origin, 10),
				"status":     strconv.FormatInt(args.Status, 10),
				"remark":     args.Remark,
				"operater":   args.Operator,
				"reason":     args.Reason,
			},
		}
		s.mbDao.AddMoralLogReport(c, l)

		// content := make(map[string][]byte, 10)
		// content["mid"] = []byte(strconv.FormatInt(mid, 10))
		// content["ip"] = []byte(args.IP)
		// content["from_moral"] = []byte(strconv.FormatInt(before, 10))
		// content["to_moral"] = []byte(strconv.FormatInt(after, 10))
		// content["origin"] = []byte(strconv.FormatInt(args.Origin, 10))
		// content["status"] = []byte(strconv.FormatInt(args.Status, 10))
		// content["remark"] = []byte(args.Remark)
		// content["operater"] = []byte(args.Operator)
		// content["reason"] = []byte(args.Reason)
		// s.mbDao.AddMoralLog(c, mid, ts, content)

		s.mbDao.DelMoralCache(c, mid)
		s.moralNotice(c, mid, before, after, args.ReasonType, args.Origin, args.Operator, args.IsNotify)
	}
	return
}

// UpdateMoral update user moral .
func (s *Service) UpdateMoral(c context.Context, arg *model.ArgUpdateMoral) (err error) {
	var (
		before, after int64
		tx            *sql.Tx
		originType    *model.OriginType
		ok            bool
	)
	if originType, ok = model.OriginTypes[arg.Origin]; !ok {
		err = ecode.RequestErr
		return
	}
	if originType.NeedReason && len(arg.Reason) == 0 {
		err = ecode.RequestErr
		return
	}
	ts := int64(time.Now().Unix())
	if tx, err = s.mbDao.BeginTran(c); err != nil {
		return
	}
	if before, after, err = s.updateMoral(tx, arg.Mid, arg.Delta); err != nil {
		tx.Rollback()
		return
	}
	if err = tx.Commit(); err != nil {
		log.Error("UpdateMoral tx.Commit(%v) error(%v)", arg.Mid, err)
		return
	}

	// log to report
	l := &model.UserLog{
		Mid:   arg.Mid,
		IP:    arg.IP,
		TS:    ts,
		LogID: model.UUID4(),
		Content: map[string]string{
			"from_moral": strconv.FormatInt(before, 10),
			"to_moral":   strconv.FormatInt(after, 10),
			"origin":     strconv.FormatInt(arg.Origin, 10),
			"status":     strconv.FormatInt(arg.Status, 10),
			"mid":        strconv.FormatInt(arg.Mid, 10),
			"remark":     arg.Remark,
			"operater":   arg.Operator,
			"reason":     arg.Reason,
		},
	}
	s.mbDao.AddMoralLogReport(c, l)

	// origin log
	// content := make(map[string][]byte, 10)
	// content["mid"] = []byte(strconv.FormatInt(arg.Mid, 10))
	// content["ip"] = []byte(arg.IP)
	// content["from_moral"] = []byte(strconv.FormatInt(before, 10))
	// content["to_moral"] = []byte(strconv.FormatInt(after, 10))
	// content["origin"] = []byte(strconv.FormatInt(arg.Origin, 10))
	// content["status"] = []byte(strconv.FormatInt(arg.Status, 10))
	// content["remark"] = []byte(arg.Remark)
	// content["operater"] = []byte(arg.Operator)
	// content["reason"] = []byte(arg.Reason)
	// s.mbDao.AddMoralLog(c, arg.Mid, ts, content)

	s.mbDao.DelMoralCache(c, arg.Mid)
	s.moralNotice(c, arg.Mid, before, after, arg.ReasonType, arg.Origin, arg.Operator, arg.IsNotify)
	return
}
func (s *Service) moralNotice(c context.Context, mid, before, after, reasonType, origin int64, operator string, notifyChange bool) {
	rt, ok := model.ReasonTypes[reasonType]
	if !ok || len(rt.NotifyType) == 0 {
		return
	}
	delta := abs(after - before)
	if delta == 0 {
		return
	}
	switch {
	case before >= 6000 && after <= 6000 && after >= 3000:
		s.mbDao.SendMessage(c, mid, model.Less6000Notice.Title, model.Less6000Notice.Message, model.Less6000Notice.NoticeType)
	case before >= 3000 && after <= 3000:
		s.mbDao.SendMessage(c, mid, model.Less3000Notice.Title, model.Less3000Notice.Message, model.Less3000Notice.NoticeType)
	case before < 6000 && after >= 6000:
		s.mbDao.SendMessage(c, mid, model.Greater6000Notice.Title, model.Greater6000Notice.Message, model.Greater6000Notice.NoticeType)
	}
	if !notifyChange {
		return
	}
	switch {
	case origin == model.PunishmentType && operator == "系统":
		s.mbDao.SendMessage(c, mid, fmt.Sprintf(model.SysPunishmentNotice.Title, moralStr(delta)), fmt.Sprintf(model.SysPunishmentNotice.Message, moralStr(delta)), rt.NotifyType)
	case origin == model.PunishmentType:
		s.mbDao.SendMessage(c, mid, fmt.Sprintf(model.PunishmentNotice.Title, moralStr(delta)), fmt.Sprintf(model.PunishmentNotice.Message, moralStr(delta)), rt.NotifyType)
	case origin == model.ReportRewardType:
		s.mbDao.SendMessage(c, mid, fmt.Sprintf(model.RewardNotice.Title, rt.Name), fmt.Sprintf(model.RewardNotice.Message, rt.Name, moralStr(delta)), rt.NotifyType)
	default:
	}
}

func (s *Service) updateMoral(tx *sql.Tx, mid, moralAdd int64) (before, after int64, err error) {
	var (
		moral *model.Moral
	)
	ts := int64(time.Now().Unix())
	if moral, err = s.mbDao.TxMoralDB(tx, mid); err != nil {
		return
	}
	if moral == nil {
		if err = s.mbDao.TxInitMoral(tx, mid, model.DefaultMoral, 0, 0, model.DefaultTime); err != nil {
			return
		}
		before = model.DefaultMoral
	} else {
		before = moral.Moral
	}
	if moralAdd > 0 {
		if after, err = s.incrMoral(tx, mid, abs(moralAdd), before, ts); err != nil {
			return
		}
		return
	}
	if after, err = s.decMoral(tx, mid, abs(moralAdd), before, ts); err != nil {
		return
	}
	return
}

func (s *Service) incrMoral(tx *sql.Tx, mid, delta int64, before, ts int64) (after int64, err error) {
	after = before + delta
	if after > model.MaxMoral {
		delta = model.MaxMoral - before
		after = model.MaxMoral
	}
	if err = s.mbDao.TxUpdateMoral(tx, mid, delta, delta, 0); err != nil {
		return
	}
	return
}

func (s *Service) decMoral(tx *sql.Tx, mid, delta, before, ts int64) (after int64, err error) {
	after = before - delta
	if after < 0 {
		delta = before
		after = 0
	}
	if err = s.mbDao.TxUpdateMoral(tx, mid, -delta, 0, delta); err != nil {
		return
	}
	if before >= model.DefaultMoral && after < model.DefaultMoral {
		err = s.mbDao.TxUpdateMoralRecoverDate(tx, mid, ctime.Time(ts))
	}
	return
}

func abs(v int64) int64 {
	if v > 0 {
		return v
	}
	return -v
}

func moralStr(delta int64) string {
	return fmt.Sprintf("%0.2f", float64(delta)/float64(100))
}
