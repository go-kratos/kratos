package dao

import (
	"context"
	"net/url"
	"strconv"

	"go-common/app/interface/main/dm2/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/metadata"
)

const (
	_workFlowTagList      = "/x/internal/workflow/tag/v3/list"
	_workFlowAppealAdd    = "/x/internal/workflow/appeal/v3/add"
	_workFlowAppealDelete = "/x/internal/workflow/appeal/v3/delete"
)

// WorkFlowTagList get tag list from workflow
func (d *Dao) WorkFlowTagList(c context.Context, bid, rid int64) (data []*model.WorkFlowTag, err error) {
	var (
		res    *model.WorkFlowTagListResp
		params = url.Values{}
		uri    = d.conf.Host.API + _workFlowTagList
		ip     = metadata.String(c, metadata.RemoteIP)
	)
	params.Set("bid", strconv.FormatInt(bid, 10))
	params.Set("rid", strconv.FormatInt(rid, 10))
	if err = d.httpCli.Get(c, uri, ip, params, &res); err != nil {
		log.Error("WorkFlowTagList(params:%v),error(%v)", params, err)
		return
	}
	if err = ecode.Int(res.Code); err != ecode.OK {
		log.Error("WorkFlowTagList(params:%v),error(%v)", params, err)
		return
	}
	data = res.Data
	err = nil
	return
}

// WorkFlowAppealAdd add a record to workflow
func (d *Dao) WorkFlowAppealAdd(c context.Context, req *model.WorkFlowAppealAddReq) (err error) {
	var (
		res *model.CommonResponse
		uri = d.conf.Host.API + _workFlowAppealAdd
		ip  = metadata.String(c, metadata.RemoteIP)
	)
	if err = d.httpCli.Post(c, uri, ip, req.Params(), &res); err != nil {
		log.Error("WorkFlowAppealAdd(req:%+v),error(%v)", req.Params(), err)
		return
	}
	if err = ecode.Int(res.Code); err != ecode.OK {
		log.Error("WorkFlowAppealAdd(req:%+v),error(%v)", req, err)
		return
	}
	err = nil
	return
}

// WorkFlowAppealDelete .
func (d *Dao) WorkFlowAppealDelete(c context.Context, bid, oid, subtitleID int64) (err error) {
	var (
		res    *model.CommonResponse
		params = url.Values{}
		uri    = d.conf.Host.API + _workFlowAppealDelete
		ip     = metadata.String(c, metadata.RemoteIP)
	)
	params.Set("business", strconv.FormatInt(bid, 10))
	params.Set("oid", strconv.FormatInt(oid, 10))
	params.Set("eid", strconv.FormatInt(subtitleID, 10))
	if err = d.httpCli.Post(c, uri, ip, params, &res); err != nil {
		log.Error("WorkFlowAppealDelete(params:%v),error(%v)", params, err)
		return
	}
	if err = ecode.Int(res.Code); err != ecode.OK {
		log.Error("WorkFlowAppealDelete(params:%v),error(%v)", params, err)
		return
	}
	err = nil
	return
}
