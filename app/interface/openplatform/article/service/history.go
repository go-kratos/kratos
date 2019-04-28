package service

import (
	"context"
	"fmt"
	"time"

	hismdl "go-common/app/interface/main/history/model"
	"go-common/app/interface/openplatform/article/dao"
	"go-common/app/interface/openplatform/article/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

const _historyArtType = 5
const _historyListType = 7

// AddHistory .
func (s *Service) AddHistory(c context.Context, mid, cvid int64, listID int64, ip string, plat int8, from string) (err error) {
	if from == "articleSlide" {
		return
	}
	var client, typ int8
	client = model.HistoryClient(plat)
	var aid, cid int64
	if listID != 0 {
		typ = _historyListType
		aid = listID
		cid = cvid
	} else {
		typ = _historyArtType
		aid = cvid
	}
	arg := &hismdl.ArgHistory{
		Mid:    mid,
		RealIP: ip,
		History: &hismdl.History{
			Mid:  mid,
			Aid:  aid,
			TP:   typ,
			Cid:  cid,
			DT:   client,
			Unix: time.Now().Unix(),
		},
	}
	if err = s.hisRPC.Add(c, arg); err != nil {
		dao.PromError("rpc:添加历史记录")
		log.Error("s.historyRPC.Add(%+v +v) error(%v)", arg, arg.History, err)
	}
	return
}

func (s *Service) historyPosition(c context.Context, mid, listID int64) (cvid int64, err error) {
	arg := &hismdl.ArgPos{
		Mid: mid,
		Aid: listID,
		TP:  _historyListType,
	}
	history, err := s.hisRPC.Position(c, arg)
	if err != nil {
		if ecode.NothingFound.Equal(err) {
			dao.PromError("history:获取历史记录")
			log.Warnv(c, log.KV("log", fmt.Sprintf("s.historyRPC.Position(%+v) error(%v)", arg, err)))
		} else {
			dao.PromError("history:获取历史记录")
			log.Errorv(c, log.KV("log", fmt.Sprintf("s.historyRPC.Position(%+v) error(%v)", arg, err)))
		}
	}
	if history != nil {
		cvid = history.Cid
	}
	return
}
