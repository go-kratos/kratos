package service

import (
	"bytes"
	"context"
	"sort"
	"strconv"
	"sync"
	"time"

	"go-common/app/job/main/growup/model"
	"go-common/library/ecode"
	"go-common/library/log"

	"golang.org/x/sync/errgroup"
)

type bubbleType int

const (
	_bubbleLottery bubbleType = 1
	_bubbleVote    bubbleType = 2
)

type bubbleMeta struct {
	avid  int64
	aType bubbleType
}

// use var map[] + register map ???
func (s *Service) syncBubbleMetaFuncM() map[bubbleType]func(c context.Context, start, end time.Time) ([]*bubbleMeta, error) {
	return map[bubbleType]func(c context.Context, start, end time.Time) ([]*bubbleMeta, error){
		_bubbleLottery: s.syncLotteryAvs,
		_bubbleVote:    s.syncVoteBIZAvs,
	}
}

func (s *Service) syncLotteryAvs(c context.Context, start, end time.Time) (data []*bubbleMeta, err error) {
	var (
		offset  int64
		hasMore = 1
		avs     []int64
	)
	for hasMore == 1 {
		var info *model.LotteryRes
		info, err = s.dao.GetLotteryRIDs(c, start.Unix(), end.Unix()-1, offset)
		if err != nil {
			log.Error("s.dao.GetLotteryRIDs error(%v)", err)
			return
		}
		if len(info.RIDs) > 0 {
			avs = append(avs, info.RIDs...)
		}
		hasMore = info.HasMore
		offset = info.Offset
	}
	data = make([]*bubbleMeta, 0, len(avs))
	for _, av := range avs {
		data = append(data, &bubbleMeta{
			avid:  av,
			aType: _bubbleLottery,
		})
	}
	return
}

func (s *Service) syncVoteBIZAvs(c context.Context, start, end time.Time) (data []*bubbleMeta, err error) {
	if end.After(start.AddDate(0, 6, 0)) {
		err = ecode.Error(ecode.RequestErr, "同步投票视频时间区间不能超过半年")
		return
	}
	data = make([]*bubbleMeta, 0)
	var voteAVs []*model.VoteBIZArchive
	voteAVs, err = s.dao.VoteBIZArchive(c, start.Unix(), end.Unix())
	if err != nil {
		return
	}
	if len(voteAVs) == 0 {
		return
	}
	for _, av := range voteAVs {
		data = append(data, &bubbleMeta{
			avid:  av.Aid,
			aType: _bubbleVote,
		})
	}
	return
}

func (s *Service) saveBubbleMeta(c context.Context, metas []*bubbleMeta, date time.Time) (rows int64, err error) {
	return s.dao.InsertBubbleMeta(c, assembleBubbleMeta(metas, date))
}

func assembleBubbleMeta(data []*bubbleMeta, syncDate time.Time) (values string) {
	var buf bytes.Buffer
	for _, v := range data {
		buf.WriteString("(")
		buf.WriteString(strconv.FormatInt(v.avid, 10))
		buf.WriteByte(',')
		buf.WriteString("'" + syncDate.Format(_layout) + "'")
		buf.WriteByte(',')
		buf.WriteString(strconv.Itoa(int(v.aType)))
		buf.WriteString(")")
		buf.WriteByte(',')
	}
	if buf.Len() > 0 {
		buf.Truncate(buf.Len() - 1)
	}
	values = buf.String()
	buf.Reset()
	return
}

// SyncIncomeBubbleMeta .
func (s *Service) SyncIncomeBubbleMeta(c context.Context, start, end time.Time, tp int) (rows int64, err error) {
	syncFuncM := s.syncBubbleMetaFuncM()
	syncFunc, ok := syncFuncM[bubbleType(tp)]
	if !ok {
		err = ecode.Errorf(ecode.ReqParamErr, "illegal bubbleType(%d)", tp)
		return
	}
	metas, err := syncFunc(c, start, end)
	if err != nil {
		log.Error("s.syncVoteBIZAvs err(%v)", err)
		return
	}
	rows, err = s.saveBubbleMeta(c, metas, time.Now())
	if err != nil {
		log.Error("s.saveBubbleMeta err(%v)", err)
	}
	return
}

// SyncIncomeBubbleMetaTask sync meta data for income bubble to growup
func (s *Service) SyncIncomeBubbleMetaTask(c context.Context, date time.Time) (err error) {
	defer func() {
		GetTaskService().SetTaskStatus(c, TaskBubbleMeta, date.Format(_layout), err)
	}()

	var (
		lock          sync.Mutex
		group         errgroup.Group
		syncMetaFuncM = s.syncBubbleMetaFuncM()
		start, end    = date, date.AddDate(0, 0, 1)
		metas         []*bubbleMeta
	)
	for _, syncMetaFunc := range syncMetaFuncM {
		group.Go(func() error {
			var data []*bubbleMeta
			data, err = syncMetaFunc(c, start, end)
			if err != nil {
				return err
			}
			lock.Lock()
			metas = append(metas, data...)
			lock.Unlock()
			return nil
		})
	}
	err = group.Wait()
	if err != nil {
		return
	}
	_, err = s.saveBubbleMeta(c, metas, date)
	if err != nil {
		log.Error("s.saveBubbleMeta err(%v)", err)
	}
	return
}

func (s *Service) avToBType(bubbleMeta map[int64][]int) map[int64]int {
	var (
		res         = make(map[int64]int)
		typeToRatio = make(map[int]float64)
		chooseBType = func(bTypes []int) (bType int) {
			if len(bTypes) == 1 {
				return bTypes[0]
			}
			sort.Slice(bTypes, func(i, j int) bool {
				bti, btj := bTypes[i], bTypes[j]
				if typeToRatio[bti] == typeToRatio[btj] {
					return bti < btj
				}
				return typeToRatio[bti] < typeToRatio[btj]
			})
			return bTypes[0]
		}
	)

	for _, v := range s.conf.Bubble.BRatio {
		typeToRatio[v.BType] = v.Ratio
	}
	for avID, bTypes := range bubbleMeta {
		res[avID] = chooseBType(bTypes)
	}
	return res
}

// SnapshotBubbleIncomeTask get income lottery
func (s *Service) SnapshotBubbleIncomeTask(c context.Context, date time.Time) (err error) {
	defer func() {
		GetTaskService().SetTaskStatus(c, TaskSnapshotBubbleIncome, date.Format("2006-01-02"), err)
	}()
	err = GetTaskService().TaskReady(c, date.Format("2006-01-02"), TaskCreativeIncome)
	if err != nil {
		return
	}

	bubbleMeta, err := s.getBubbleMeta(c)
	if err != nil {
		log.Error("s.getBubbleAVs error(%v)", err)
		return
	}
	avToBType := s.avToBType(bubbleMeta)

	var (
		eg         errgroup.Group
		incomeCh   = make(chan []*model.IncomeInfo, 2000)
		snapshotCh = make(chan []*model.IncomeInfo, 2000)
	)
	// get av income by date
	eg.Go(func() (err error) {
		defer close(incomeCh)
		var from int64
		for {
			var infos []*model.IncomeInfo
			infos, err = s.dao.GetAvIncome(c, date, from, int64(_dbLimit))
			if err != nil {
				log.Error("dao.GetAvIncome error(%v)", err)
				return
			}
			if len(infos) == 0 {
				break
			}
			incomeCh <- infos
			from = infos[len(infos)-1].ID
		}
		return
	})

	eg.Go(func() (err error) {
		defer close(snapshotCh)
		for income := range incomeCh {
			lotteryIncome := make([]*model.IncomeInfo, 0)
			for _, av := range income {
				if bType, ok := avToBType[av.AVID]; ok {
					av.BType = bType
					lotteryIncome = append(lotteryIncome, av)
				}
			}
			if len(lotteryIncome) > 0 {
				snapshotCh <- lotteryIncome
			}
		}
		return
	})

	eg.Go(func() (err error) {
		for income := range snapshotCh {
			_, err = s.dao.InsertBubbleIncome(c, assembleLotteryIncome(income))
			if err != nil {
				return
			}
		}
		return
	})

	if err = eg.Wait(); err != nil {
		log.Error("eg.Wait error(%v)", err)
	}
	return
}

func (s *Service) getBubbleMeta(c context.Context) (data map[int64][]int, err error) {
	var id int64
	data = make(map[int64][]int)
	for {
		var meta map[int64][]int
		meta, id, err = s.income.GetBubbleMeta(c, id, int64(_dbLimit))
		if err != nil {
			return
		}
		if len(meta) == 0 {
			break
		}
		for avID, bTypes := range meta {
			data[avID] = append(data[avID], bTypes...)
		}
	}
	return
}

func assembleLotteryIncome(as []*model.IncomeInfo) (values string) {
	var buf bytes.Buffer
	for _, a := range as {
		buf.WriteString("(")
		buf.WriteString(strconv.FormatInt(a.AVID, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(a.MID, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(a.TagID, 10))
		buf.WriteByte(',')
		buf.WriteString("'" + a.UploadTime.Format("2006-01-02 15:04:05") + "'")
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(a.TotalIncome, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(a.Income, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(a.TaxMoney, 10))
		buf.WriteByte(',')
		buf.WriteString("'" + a.Date.Format(_layout) + "'")
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(a.BaseIncome, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.Itoa(a.BType))
		buf.WriteString(")")
		buf.WriteByte(',')
	}
	if buf.Len() > 0 {
		buf.Truncate(buf.Len() - 1)
	}
	values = buf.String()
	buf.Reset()
	return
}
