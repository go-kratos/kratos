package dao

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"
	"time"

	"go-common/app/job/main/figure-timer/model"
	"go-common/library/log"

	"github.com/pkg/errors"
	"github.com/tsuna/gohbase/hrpc"
)

var (
	// record
	_hbaseRecordTable           = "ugc:figurecalcrecord"
	_hbaseRecordFC              = "x"
	_hbaseRecordQLawfulPosX     = "pos_lawful"
	_hbaseRecordQLawfulNegX     = "neg_lawful"
	_hbaseRecordQWidePosX       = "pos_wide"
	_hbaseRecordQWideNegX       = "neg_wide"
	_hbaseRecordQFriendlyPosX   = "pos_friendly"
	_hbaseRecordQFriendlyNegX   = "neg_friendly"
	_hbaseRecordQCreativityPosX = "pos_creativity"
	_hbaseRecordQCreativityNegX = "neg_creativity"
	_hbaseRecordQBountyPosX     = "pos_bounty"
	_hbaseRecordQBountyNegX     = "neg_bounty"
	// UserInfo
	_hbaseUserTable                = "ugc:figureuserstatus"
	_hbaseUserFC                   = "user"
	_hbaseUserQExp                 = "exp"
	_hbaseUserQSpy                 = "spy_score"
	_hbaseUserQArchiveViews        = "archive_views"
	_hbaseUserQVip                 = "vip_status"
	_hbaseUserQDisciplineCommittee = "discipline_committee"
	// ActionCounter
	_hbaseActionTable          = "ugc:figureactioncounter"
	_hbaseActionFC             = "user"
	_hbaseActionQCoinCount     = "coins"
	_hbaseActionQReplyCount    = "replies"
	_hbaseActionQDanmakuCount  = "danmaku"
	_hbaseActionQCoinLowRisk   = "coin_low_risk"
	_hbaseActionQCoinHighRisk  = "coin_high_risk"
	_hbaseActionQReplyLowRisk  = "reply_low_risk"
	_hbaseActionQReplyHighRisk = "reply_high_risk"
	_hbaseActionQReplyLiked    = "reply_liked"
	_hbaseActionQReplyUnLiked  = "reply_hate"

	_hbaseActionQReportReplyPassed     = "report_reply_passed"
	_hbaseActionQReportDanmakuPassed   = "report_danmaku_passed"
	_hbaseActionQPublishReplyDeleted   = "publish_reply_deleted"
	_hbaseActionQPublishDanmakuDeleted = "publish_danmaku_deleted"
	_hbaseActionQPayMoney              = "pay_money"
	_hbaseActionQPayLiveMoney          = "pay_live_money"
)

func rowKeyFigureRecord(mid int64, weekTS int64) (key string) {
	return fmt.Sprintf("%d_%d", mid, weekTS)
}

func rowKeyUserInfo(mid int64) (key string) {
	return fmt.Sprintf("%d", mid)
}

func rowKeyActionCounter(mid int64, dayTS int64) (key string) {
	return fmt.Sprintf("%d_%d", mid, dayTS)
}

// UserInfo get it from hbase
func (d *Dao) UserInfo(c context.Context, mid int64, weekVer int64) (userInfo *model.UserInfo, err error) {
	var (
		result      *hrpc.Result
		key         = rowKeyUserInfo(mid)
		ctx, cancel = context.WithTimeout(c, time.Duration(d.c.Hbase.ReadTimeout))
	)
	defer cancel()
	if result, err = d.hbase.GetStr(ctx, _hbaseUserTable, key); err != nil {
		err = errors.Wrapf(err, "hbase.GetStr(%s,%s)", _hbaseUserTable, key)
		return
	}
	userInfo = &model.UserInfo{Mid: mid, SpyScore: math.MaxUint64}
	for _, c := range result.Cells {
		if c == nil {
			continue
		}
		if bytes.Equal([]byte(_hbaseUserFC), c.Family) {
			switch string(c.Qualifier) {
			case _hbaseUserQExp:
				userInfo.Exp = binary.BigEndian.Uint64(c.Value)
			case _hbaseUserQSpy:
				userInfo.SpyScore = binary.BigEndian.Uint64(c.Value)
			case _hbaseUserQArchiveViews:
				userInfo.ArchiveViews = binary.BigEndian.Uint64(c.Value)
			case _hbaseUserQVip:
				userInfo.VIPStatus = binary.BigEndian.Uint64(c.Value)
			case _hbaseUserQDisciplineCommittee:
				userInfo.DisciplineCommittee = binary.BigEndian.Uint16(c.Value)
			}
		}
	}
	// 未获取到spy分值特殊处理
	if userInfo.SpyScore == math.MaxUint64 {
		userInfo.SpyScore = 100
	}
	log.Info("User Info [%+v]", userInfo)
	return
}

// PutCalcRecord put record all the params x.
func (d *Dao) PutCalcRecord(c context.Context, record *model.FigureRecord, weekTS int64) (err error) {
	var (
		key            = rowKeyFigureRecord(record.Mid, weekTS)
		lawfulPosX     = make([]byte, 8)
		lawfulNegX     = make([]byte, 8)
		widePosX       = make([]byte, 8)
		wideNegX       = make([]byte, 8)
		friendlyPosX   = make([]byte, 8)
		friendlyNegX   = make([]byte, 8)
		creativityPosX = make([]byte, 8)
		creativityNegX = make([]byte, 8)
		bountyPosX     = make([]byte, 8)
		bountyNegX     = make([]byte, 8)
		ctx, cancel    = context.WithTimeout(c, time.Duration(d.c.Hbase.WriteTimeout))
	)
	defer cancel()
	binary.BigEndian.PutUint64(lawfulPosX, uint64(record.XPosLawful))
	binary.BigEndian.PutUint64(lawfulNegX, uint64(record.XNegLawful))
	binary.BigEndian.PutUint64(widePosX, uint64(record.XPosWide))
	binary.BigEndian.PutUint64(wideNegX, uint64(record.XNegWide))
	binary.BigEndian.PutUint64(friendlyPosX, uint64(record.XPosFriendly))
	binary.BigEndian.PutUint64(friendlyNegX, uint64(record.XNegFriendly))
	binary.BigEndian.PutUint64(creativityPosX, uint64(record.XPosCreativity))
	binary.BigEndian.PutUint64(creativityNegX, uint64(record.XNegCreativity))
	binary.BigEndian.PutUint64(bountyPosX, uint64(record.XPosBounty))
	binary.BigEndian.PutUint64(bountyNegX, uint64(record.XNegBounty))
	values := map[string]map[string][]byte{_hbaseRecordFC: {
		_hbaseRecordQLawfulPosX:     lawfulPosX,
		_hbaseRecordQLawfulNegX:     lawfulNegX,
		_hbaseRecordQWideNegX:       wideNegX,
		_hbaseRecordQFriendlyPosX:   friendlyPosX,
		_hbaseRecordQFriendlyNegX:   friendlyNegX,
		_hbaseRecordQCreativityPosX: creativityPosX,
		_hbaseRecordQCreativityNegX: creativityNegX,
		_hbaseRecordQBountyPosX:     bountyPosX,
		_hbaseRecordQBountyNegX:     bountyNegX,
	}}
	if _, err = d.hbase.PutStr(ctx, _hbaseRecordTable, key, values); err != nil {
		err = errors.Wrapf(err, "hbase.Put(%s,%s,%+v) error(%v)", _hbaseRecordTable, key, values, err)
	}
	return
}

// CalcRecords get it from hbase
func (d *Dao) CalcRecords(c context.Context, mid int64, weekTSFrom, weekTSTo int64) (figureRecords []*model.FigureRecord, err error) {
	var (
		scanner     hrpc.Scanner
		keyFrom     = rowKeyFigureRecord(mid, weekTSFrom)
		keyTo       = rowKeyFigureRecord(mid, weekTSTo)
		ctx, cancel = context.WithTimeout(c, time.Duration(d.c.Hbase.ReadTimeout))
	)
	defer cancel()
	if scanner, err = d.hbase.ScanRangeStr(ctx, _hbaseRecordTable, keyFrom, keyTo); err != nil {
		err = errors.Wrapf(err, "hbase.ScanRangeStr(%s,%s,%s)", _hbaseRecordTable, keyFrom, keyTo)
		return
	}
	for {
		res, err := scanner.Next()
		if err != nil {
			if err != io.EOF {
				err = errors.WithStack(err)
				return nil, err
			}
			break
		}
		figureRecord := &model.FigureRecord{Mid: mid}
		for i, c := range res.Cells {
			if c == nil {
				continue
			}
			if bytes.Equal([]byte(_hbaseRecordFC), c.Family) {
				switch string(c.Qualifier) {
				case _hbaseRecordQLawfulPosX:
					figureRecord.XPosLawful = int64(binary.BigEndian.Uint64(c.Value))
				case _hbaseRecordQLawfulNegX:
					figureRecord.XNegLawful = int64(binary.BigEndian.Uint64(c.Value))
				case _hbaseRecordQWidePosX:
					figureRecord.XPosWide = int64(binary.BigEndian.Uint64(c.Value))
				case _hbaseRecordQWideNegX:
					figureRecord.XNegWide = int64(binary.BigEndian.Uint64(c.Value))
				case _hbaseRecordQFriendlyPosX:
					figureRecord.XPosFriendly = int64(binary.BigEndian.Uint64(c.Value))
				case _hbaseRecordQFriendlyNegX:
					figureRecord.XNegFriendly = int64(binary.BigEndian.Uint64(c.Value))
				case _hbaseRecordQCreativityPosX:
					figureRecord.XPosCreativity = int64(binary.BigEndian.Uint64(c.Value))
				case _hbaseRecordQCreativityNegX:
					figureRecord.XNegCreativity = int64(binary.BigEndian.Uint64(c.Value))
				case _hbaseRecordQBountyPosX:
					figureRecord.XPosBounty = int64(binary.BigEndian.Uint64(c.Value))
				case _hbaseRecordQBountyNegX:
					figureRecord.XNegBounty = int64(binary.BigEndian.Uint64(c.Value))
				}
			}
			if i == 0 {
				figureRecord.Version = time.Unix(parseVersion(string(c.Row), c.Timestamp), 0)
			}
		}
		log.Info("User figure record [%+v]", figureRecord)
		figureRecords = append(figureRecords, figureRecord)
	}
	return
}

// ActionCounter .
func (d *Dao) ActionCounter(c context.Context, mid int64, ts int64) (actionCounter *model.ActionCounter, err error) {
	var (
		result      *hrpc.Result
		key         = rowKeyActionCounter(mid, ts)
		ctx, cancel = context.WithTimeout(c, time.Duration(d.c.Hbase.ReadTimeout))
	)
	defer cancel()
	log.Info("Get mid [%d] key [%s]", mid, key)
	if result, err = d.hbase.GetStr(ctx, _hbaseActionTable, key); err != nil {
		err = errors.Wrapf(err, "hbase.GetStr(%s,%s)", _hbaseActionTable, key)
		return
	}
	actionCounter = &model.ActionCounter{Mid: mid}
	for i, c := range result.Cells {
		if c == nil {
			continue
		}
		if bytes.Equal([]byte(_hbaseActionFC), c.Family) {
			switch string(c.Qualifier) {
			case _hbaseActionQCoinCount:
				actionCounter.CoinCount = binary.BigEndian.Uint64(c.Value)
			case _hbaseActionQReplyCount:
				actionCounter.ReplyCount = int64(binary.BigEndian.Uint64(c.Value))
			case _hbaseActionQDanmakuCount:
				actionCounter.DanmakuCount = int64(binary.BigEndian.Uint64(c.Value))
			case _hbaseActionQCoinLowRisk:
				actionCounter.CoinLowRisk = binary.BigEndian.Uint64(c.Value)
			case _hbaseActionQCoinHighRisk:
				actionCounter.CoinHighRisk = binary.BigEndian.Uint64(c.Value)
			case _hbaseActionQReplyLowRisk:
				actionCounter.ReplyLowRisk = binary.BigEndian.Uint64(c.Value)
			case _hbaseActionQReplyHighRisk:
				actionCounter.ReplyHighRisk = binary.BigEndian.Uint64(c.Value)
			case _hbaseActionQReplyLiked:
				actionCounter.ReplyLiked = int64(binary.BigEndian.Uint64(c.Value))
			case _hbaseActionQReplyUnLiked:
				actionCounter.ReplyUnliked = int64(binary.BigEndian.Uint64(c.Value))
			case _hbaseActionQReportReplyPassed:
				actionCounter.ReportReplyPassed = int64(binary.BigEndian.Uint64(c.Value))
			case _hbaseActionQReportDanmakuPassed:
				actionCounter.ReportDanmakuPassed = int64(binary.BigEndian.Uint64(c.Value))
			case _hbaseActionQPublishReplyDeleted:
				actionCounter.PublishReplyDeleted = int64(binary.BigEndian.Uint64(c.Value))
			case _hbaseActionQPublishDanmakuDeleted:
				actionCounter.PublishDanmakuDeleted = int64(binary.BigEndian.Uint64(c.Value))
			case _hbaseActionQPayMoney:
				actionCounter.PayMoney = int64(binary.BigEndian.Uint64(c.Value))
			case _hbaseActionQPayLiveMoney:
				actionCounter.PayLiveMoney = int64(binary.BigEndian.Uint64(c.Value))

			}
		}
		if i == 0 {
			actionCounter.Version = time.Unix(parseVersion(string(c.Row), c.Timestamp), 0)
		}
	}
	log.Info("User action counter [%+v]", actionCounter)
	return
}

// // ActionCounters .
// func (d *Dao) ActionCounters(c context.Context, mid int64, tsfrom int64, tsto int64) (actionCounters []*model.ActionCounter, err error) {
// 	var (
// 		scan        *hrpc.Scan
// 		results     []*hrpc.Result
// 		keyFrom     = rowKeyActionCounter(mid, tsfrom)
// 		keyTo       = rowKeyActionCounter(mid, tsto)
// 		ctx, cancel = context.WithTimeout(c, time.Duration(d.c.Hbase.ReadTimeout))
// 	)
// 	defer cancel()
// 	log.Info("Scan mid [%d] keyFrom [%s] keyTo [%s]", mid, keyFrom, keyTo)
// 	if scan, err = hrpc.NewScanRangeStr(ctx, _hbaseActionTable, keyFrom, keyTo); err != nil {
// 		err = errors.Wrapf(err, "hrcp.NewScanRangeStr(%s,%s,%s)", _hbaseActionTable, keyFrom, keyTo)
// 		return
// 	}
// 	if results, err = d.hbase.Scan(ctx, scan); err != nil {
// 		err = errors.Wrapf(err, "hbase.Scan(%s,%s,%s)", _hbaseActionTable, keyFrom, keyTo)
// 		return
// 	}
// 	if len(results) == 0 {
// 		return
// 	}
// 	for _, res := range results {
// 		if res == nil {
// 			continue
// 		}
// 		actionCounter := &model.ActionCounter{Mid: mid}
// 		for i, c := range res.Cells {
// 			if c == nil {
// 				continue
// 			}
// 			if bytes.Equal([]byte(_hbaseActionFC), c.Family) {
// 				switch string(c.Qualifier) {
// 				case _hbaseActionQCoinCount:
// 					actionCounter.CoinCount = binary.BigEndian.Uint64(c.Value)
// 				case _hbaseActionQReplyCount:
// 					actionCounter.ReplyCount = int64(binary.BigEndian.Uint64(c.Value))
// 				case _hbaseActionQCoinLowRisk:
// 					actionCounter.CoinLowRisk = binary.BigEndian.Uint64(c.Value)
// 				case _hbaseActionQCoinHighRisk:
// 					actionCounter.CoinHighRisk = binary.BigEndian.Uint64(c.Value)
// 				case _hbaseActionQReplyLowRisk:
// 					actionCounter.ReplyLowRisk = binary.BigEndian.Uint64(c.Value)
// 				case _hbaseActionQReplyHighRisk:
// 					actionCounter.ReplyHighRisk = binary.BigEndian.Uint64(c.Value)
// 				case _hbaseActionQReplyLiked:
// 					actionCounter.ReplyLiked = int64(binary.BigEndian.Uint64(c.Value))
// 				case _hbaseActionQReplyUnLiked:
// 					actionCounter.ReplyUnliked = int64(binary.BigEndian.Uint64(c.Value))
// 				}
// 			}
// 			if i == 0 {
// 				actionCounter.Version = time.Unix(parseVersion(string(c.Row), c.Timestamp), 0)
// 			}
// 		}
// 		log.Info("User action counter [%+v]", actionCounter)
// 		actionCounters = append(actionCounters, actionCounter)
// 	}
// 	return
// }

func parseVersion(rowKey string, timestamp *uint64) (ts int64) {
	strs := strings.Split(rowKey, "_")
	if len(strs) < 2 {
		return int64(*timestamp)
	}
	var err error
	if ts, err = strconv.ParseInt(strs[1], 10, 64); err != nil {
		return int64(*timestamp)
	}
	return ts
}
