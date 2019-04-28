package vip

import (
	"context"

	"go-common/app/service/main/vip/api"
	"go-common/library/log"
)

// WelfareList get welfare list
func (s *Service) WelfareList(c context.Context, tid, recommend, pn, ps int64) (res *v1.WelfareReply, err error) {
	welfareReq := &v1.WelfareReq{Tid: tid, Recommend: recommend, Ps: ps, Pn: pn}

	if res, err = s.vipgRPC.WelfareList(c, welfareReq); err != nil {
		log.Error("vipSvc.WelfareList(%+v) err(%+v)", welfareReq, err)
	}

	return
}

// WelfareTypeList get welfare type list
func (s *Service) WelfareTypeList(c context.Context) (res *v1.WelfareTypeReply, err error) {
	welfareTypeReq := &v1.WelfareTypeReq{}

	if res, err = s.vipgRPC.WelfareTypeList(c, welfareTypeReq); err != nil {
		log.Error("vipSvc.WelfareTypeList err(%+v)", err)
	}

	return
}

// WelfareInfo get welfare info
func (s *Service) WelfareInfo(c context.Context, wid, mid int64) (res *v1.WelfareInfoReply, err error) {
	welfareInfoReq := &v1.WelfareInfoReq{Id: wid, Mid: mid}

	if res, err = s.vipgRPC.WelfareInfo(c, welfareInfoReq); err != nil {
		log.Error("vipSvc.WelfareInfo(%+v) err(%+v)", wid, err)
	}

	return
}

// WelfareReceive receive welfare
func (s *Service) WelfareReceive(c context.Context, wid, mid int64) (res *v1.WelfareReceiveReply, err error) {
	welfareReceiveReq := &v1.WelfareReceiveReq{Wid: wid, Mid: mid}

	if res, err = s.vipgRPC.WelfareReceive(c, welfareReceiveReq); err != nil {
		log.Error("vipSvc.WelfareReceive(%+v) err(%+v)", wid, err)
	}

	return
}

// MyWelfare get my welfare
func (s *Service) MyWelfare(c context.Context, mid int64) (res *v1.MyWelfareReply, err error) {
	myWelfareReq := &v1.MyWelfareReq{Mid: mid}

	if res, err = s.vipgRPC.MyWelfare(c, myWelfareReq); err != nil {
		log.Error("vipSvc.MyWelfare err(%+v)", err)
	}

	return
}
