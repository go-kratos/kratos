package service

import (
	"encoding/json"
	"github.com/pkg/errors"
	"go-common/app/service/main/upcredit/dao/upcrmdao"
	"go-common/app/service/main/upcredit/model/canal"
	"go-common/app/service/main/upcredit/model/upcrmmodel"
	"go-common/library/log"
	"math/rand"
	"sync"
	"time"
)

// 必须是upcrmmodel.UpScoreHistoryTableCount
var batchDataList = make([][]*upcrmdao.UpQualityInfo, upcrmmodel.UpScoreHistoryTableCount)
var batchMutex = sync.Mutex{}

func (s *Service) onBusinessDatabus(data *canal.Msg) {
	if data == nil {
		return
	}
	var err error
	var upInfo = &upcrmdao.UpQualityInfo{}
	if err = json.Unmarshal(data.New, upInfo); err != nil {
		log.Error("unmarshal business canal message fail, err=%+v, value=%s", err, string(data.New))
		return
	}
	var index = upInfo.Mid % 100
	batchMutex.Lock()
	batchDataList[index] = append(batchDataList[index], upInfo)
	batchMutex.Unlock()
	s.businessScoreChan <- upInfo
}

func (s *Service) updateScoreProc() {
	defer func() {
		s.wg.Done()
		if r := recover(); r != nil {
			r = errors.WithStack(r.(error))
			log.Error("updateScoreProc Runtime error caught, try recover: %+v", r)
			s.wg.Add(1)
			go s.updateScoreProc()
		}
	}()
	var stop = false
	for s.running && !stop {
		<-s.limit.Token()
		select {
		case upInfo, ok := <-s.businessScoreChan:
			if !ok {
				log.Error("business score chan closed")
				stop = true
				break
			}
			if upInfo == nil {
				continue
			}
			// update up player's baseinfo
			affectedRow, err := s.upcrmdb.UpdateQualityAndPrScore(upInfo.PrValue, upInfo.QualityValue, upInfo.Mid)
			if err != nil {
				log.Error("error when update quality score, err=%s", err)
				break
			}
			log.Info("update up's score, info=%+v, affected=%d", upInfo, affectedRow)
		case <-s.closeChan:
			log.Info("server closing")
			stop = true
		}

	}

	// drain the channel
	for upInfo := range s.businessScoreChan {
		// update up player's baseinfo
		affectedRow, err := s.upcrmdb.UpdateQualityAndPrScore(upInfo.PrValue, upInfo.QualityValue, upInfo.Mid)
		if err != nil {
			log.Error("error when update quality score, err=%s", err)
			break
		}
		log.Info("update up's score, info=%+v, affected=%d", upInfo, affectedRow)
	}

}
func (s *Service) batchWriteProc() {
	defer func() {
		s.wg.Done()
		if r := recover(); r != nil {
			r = errors.WithStack(r.(error))
			log.Error("batchWriteProc Runtime error caught, try recover: %+v", r)
			s.wg.Add(1)
			go s.batchWriteProc()
		}
	}()
	var ticker = time.NewTicker(10 * time.Second)
	for s.running {
		// 每秒钟查一次
		select {
		case <-ticker.C:
		case <-s.closeChan:
			log.Info("server closing")
		}
		batchMutex.Lock()
		for i, list := range batchDataList {
			s.batchUpdate(list, i)
			batchDataList[i] = nil
		}
		batchMutex.Unlock()
	}
	// drain the list
	for i, list := range batchDataList {
		// batch update
		s.batchUpdate(list, i)
	}
}

// leastCount to flush
func (s *Service) batchUpdate(list []*upcrmdao.UpQualityInfo, tablenum int) {
	var leftCount = len(list)
	if leftCount == 0 {
		return
	}

	var affected, err = s.upcrmdb.InsertBatchScoreHistory(list, tablenum)
	if err != nil {
		log.Error("error insert score history, err=%s", err)
		return
	}
	log.Info("batch insert into score history, count=%d, affected=%d", leftCount, affected)
}

//Test test func
func (s *Service) Test() {
	var list = []*upcrmdao.UpQualityInfo{
		{Mid: 123, QualityValue: rand.Int() % 100, PrValue: rand.Int() % 100, Cdate: "2018-07-13"},
	}
	s.upcrmdb.InsertBatchScoreHistory(list, 23)
}
