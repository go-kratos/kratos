package service

import (
	"context"
	"fmt"
	"go-common/library/log"
	"sync"
	"sync/atomic"
)

const (
	_stateActive = 0
	_defaultType = 23
)

//CommentReg 评论开关
func (s *Service) CommentReg(c context.Context, oid int64, state int16) (err error) {
	req := map[string]interface{}{
		"oid":   oid,
		"type":  _defaultType,
		"adid":  0,
		"mid":   0,
		"state": state,
	}
	err = s.dao.ReplyReg(c, req)
	return
}

// AutoRegAll 批量注册评论
func (s *Service) AutoRegAll(c context.Context) {
	var (
		IDs     []int64
		e       int
		err     error
		lastID  int64
		num     int64
		failNum int64
	)
	for {
		IDs, lastID, _ = s.dao.GetVideoByLastID(c, lastID)
		if len(IDs) == 0 {
			break
		}
		sum := len(IDs)
		partNum := 200
		rnum := int(sum / partNum)
		log.Info("rountine start num[%d]", rnum)
		var wg sync.WaitGroup
		for i := 0; i < rnum+1; i++ {
			wg.Add(1)
			if e = (i + 1) * partNum; e > sum {
				e = sum
			}
			r := IDs[i*partNum : e]
			go func(ids []int64, k int) {
				defer wg.Done()
				for _, svid := range ids {
					err = s.CommentReg(c, svid, _stateActive)
					if err != nil {
						atomic.AddInt64(&failNum, 1)
						fmt.Printf("Comment active fail [%v]\n", err)
						log.Errorv(c, log.KV("log", fmt.Sprintf("Comment active fail [%v]", err)))
					} else {
						atomic.AddInt64(&num, 1)
						fmt.Printf("Comment active success oid:[%d]\n", svid)
						log.Infov(c, log.KV("log", fmt.Sprintf("Comment active success oid:[%d]", svid)))
					}
				}
			}(r, i)
		}
		wg.Wait()
	}
	log.Info("comment reg complete succ[%d] fail[%d]", num, failNum)
	fmt.Printf("comment reg complete succ[%d] fail[%d]", num, failNum)

}
