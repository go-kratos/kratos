package service

import (
	"context"
	"runtime/debug"
	"time"

	"go-common/app/job/main/figure-timer/model"
	"go-common/library/log"
)

func (s *Service) fixproc() {
	defer func() {
		if x := recover(); x != nil {
			log.Error("fixproc panic(%+v)", x)
			log.Error("%+v", string(debug.Stack()))
		}
	}()
	var (
		ctx                = context.TODO()
		err                error
		weekVerRecordsFrom = time.Now().AddDate(0, 0, -7*52).Unix()
		weekVerRecordsTo   = time.Now().Unix() + 1
	)
	for shard := s.c.Property.PendingMidStart; shard < 100; shard++ {
		log.Info("start fix: %d", shard)
		var (
			figures []*model.Figure
			fromMid = int64(shard)
			end     bool
		)
		for !end {
			if figures, end, err = s.dao.Figures(ctx, fromMid, 100); err != nil {
				log.Error("%+v", err)
				break
			}
			for _, figure := range figures {
				var (
					records []*model.FigureRecord
				)
				if fromMid < figure.Mid {
					fromMid = figure.Mid
				}
				if records, err = s.dao.CalcRecords(ctx, figure.Mid, weekVerRecordsFrom, weekVerRecordsTo); err != nil {
					log.Error("%+v", err)
					continue
				}
				s.fixRecord(ctx, figure.Mid, records)
			}
		}
		log.Info("end fix: %d", shard)
	}
}

// 全量清洗用户mid所有的records
func (s *Service) fixRecord(c context.Context, mid int64, records []*model.FigureRecord) {
	var (
		err    error
		action *model.ActionCounter
		x      float64
	)
	for _, record := range records {
		time.Sleep(time.Millisecond)
		beforeRecord := *record
		// 获得本次record 记录对应的 action 记录
		if action, err = s.dao.ActionCounter(c, mid, record.Version.Unix()); err != nil {
			log.Error("%+v", err)
			continue
		}
		actions := []*model.ActionCounter{action}
		// lawful
		x, _ = s.calcActionX(s.c.Property.Calc.LawfulPosL, bizTypePosLawful, actions, nil, record.Version.Unix())
		record.XPosLawful = int64(x)
		x, _ = s.calcActionX(s.c.Property.Calc.LawfulNegL, bizTypeNegLawful, actions, nil, record.Version.Unix())
		record.XNegLawful = int64(x)
		// wide
		record.XPosWide = 0
		record.XNegWide = 0
		// friendly
		x, _ = s.calcActionX(s.c.Property.Calc.FriendlyPosL, bizTypePosFriendly, actions, nil, record.Version.Unix())
		record.XPosFriendly = int64(x)
		x, _ = s.calcActionX(s.c.Property.Calc.FriendlyNegL, bizTypeNegFriendly, actions, nil, record.Version.Unix())
		record.XNegFriendly = int64(x)
		// bounty
		x, _ = s.calcActionX(s.c.Property.Calc.BountyPosL, bizTypePosBounty, actions, nil, record.Version.Unix())
		record.XPosBounty = int64(x)
		record.XNegBounty = 0
		// creativity
		x, _ = s.calcActionX(s.c.Property.Calc.CreativityPosL1, bizTypePosCreativity, actions, nil, record.Version.Unix())
		record.XPosCreativity = int64(x)
		record.XNegCreativity = 0
		// 更新本次的fix
		if err = s.dao.PutCalcRecord(c, record, record.Version.Unix()); err != nil {
			log.Error("%+v", err)
		} else {
			log.Info("fix figure record before [%+v] --> now [%+v]", beforeRecord, record)
		}
	}
}
