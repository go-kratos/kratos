package service

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"go-common/app/interface/main/tag/conf"
	"go-common/app/interface/main/tag/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/metadata"
)

const (
	_workflowAppeaListURI   = "/x/internal/workflow/appeal/v3/list"
	_workflowAppealStateURI = "/x/internal/workflow/appeal/v3/state"
	_addWorkflowAppeaURI    = "/x/internal/workflow/appeal/v3/add"
)

var (
	errFilterFaile      = fmt.Errorf("filter failed")
	errVideoTagUpFaile  = fmt.Errorf("videoTagUp failed")
	errAIRecommandFaile = fmt.Errorf("ai channel recommand code != 0")
	errappeaListFaile   = fmt.Errorf("workflow appeal list code != 0")
	errappealStateFaile = fmt.Errorf("workflow appeal state code != 0")
	errAddAppealFaile   = fmt.Errorf("add workflow appeal code != 0")
)

// setArcTag .
func (s *Service) setArcTag(c context.Context, aid int64, tag string) (err error) {
	params := url.Values{}
	params.Set("aid", strconv.FormatInt(aid, 10))
	params.Set("tag", tag)
	var res struct {
		Code int `json:"code"`
	}
	if err = s.client.Post(c, s.videoTagUpURL, metadata.String(c, metadata.RemoteIP), params, &res); err != nil {
		conf.PromError("videoup服务更新tag接口", "d.client.Post(%s) error(%v)", s.videoTagUpURL+"?"+params.Encode(), err)
		log.Error("s.client.Post(%s) error(%v)", s.videoTagUpURL+"?"+params.Encode(), err)
		return
	}
	if res.Code != 0 {
		conf.PromError("videoup服务更新tag接口", "d.client.Post(%s) res code(%d)", s.videoTagUpURL+"?"+params.Encode(), res.Code)
		log.Error("s.client.Post(%s) res code(%d)", s.videoTagUpURL+"?"+params.Encode(), res.Code)
		err = errVideoTagUpFaile
	}
	return
}

// similarsTids .
func (s *Service) similarsTids(c context.Context, tid int64) (tids []int64, err error) {
	var res struct {
		Code int     `json:"code"`
		Data []int64 `json:"data"`
	}
	params := url.Values{}
	params.Set("tid", strconv.FormatInt(tid, 10))
	if err = s.siClient.Get(c, s.similarURL, metadata.String(c, metadata.RemoteIP), params, &res); err != nil {
		conf.PromError("大数据不带分区相关tag接口", "d.siClient.Get(%s) error(%v)", s.similarURL+"?"+params.Encode(), err)
		log.Error("s.siClient.Get(%s) error(%v)", s.similarURL+"?"+params.Encode(), err)
		return
	}
	if res.Code != 0 {
		conf.PromError("大数据不带分区相关tag接口", "d.siClient.Get(%s) res code(%d) or res.data(%v)", s.similarURL+"?"+params.Encode(), res.Code, res.Data)
		log.Error("s.siClient.Get(%s) res code(%d) or res.data(%v)", s.similarURL+"?"+params.Encode(), res.Code, res.Data)
		return
	}
	return res.Data, nil
}

// similars .
func (s *Service) similars(c context.Context, rid, tid int64) (sis []*model.SimilarTag, err error) {
	var res struct {
		Code int                 `json:"code"`
		Data []*model.SimilarTag `json:"data"`
	}
	params := url.Values{}
	params.Set("rid", strconv.FormatInt(rid, 10))
	params.Set("tid", strconv.FormatInt(tid, 10))
	if err = s.siClient.Get(c, s.similarURL, metadata.String(c, metadata.RemoteIP), params, &res); err != nil {
		conf.PromError("大数据带分区相关tag接口", "d.siClient.Get(%s) error(%v)", s.similarURL+"?"+params.Encode(), err)
		log.Error("s.siClient.Get(%s) error(%v)", s.similarURL+"?"+params.Encode(), err)
		return
	}
	if res.Code != 0 {
		conf.PromError("大数据带分区相关tag接口", "d.siClient.Get(%s) res code(%d) or res.data(%v)", s.similarURL+"?"+params.Encode(), res.Code, res.Data)
		log.Error("s.siClient.Get(%s) res code(%d) or res.data(%v)", s.similarURL+"?"+params.Encode(), res.Code, res.Data)
		sis = nil
		return
	}
	return res.Data, nil
}

// filter .
func (s *Service) filter(c context.Context, msg string, now time.Time) (err error) {
	params := url.Values{}
	params.Set("area", "tag")
	params.Set("msg", msg)
	var res struct {
		Code int `json:"code"`
		Data struct {
			Level int `json:"level"`
		}
	}
	if err = s.client.Get(c, s.filterURL, metadata.String(c, metadata.RemoteIP), params, &res); err != nil {
		conf.PromError("过滤服务单个接口", "d.client.Get(%s) error(%v)", s.filterURL+"?"+params.Encode(), err)
		log.Error("s.client.Get(%s) error(%v)", s.filterURL+"?"+params.Encode(), err)
		return
	}
	if res.Code != 0 {
		conf.PromError("过滤服务单个接口", "d.client.Get(%s) res code(%d) or res.data(%v)", s.filterURL+"?"+params.Encode(), res.Code, res.Data)
		log.Error("s.client.Get(%s) res code(%d) or res.data(%v)", s.filterURL+"?"+params.Encode(), res.Code, res.Data)
		err = errFilterFaile
		return
	}
	if res.Data.Level > 10 {
		log.Error("s.filter() tag name:%+v,Level:%+v", msg, res.Data.Level)
		err = ecode.TagIsSealing
	}
	return
}

// mfilter .
func (s *Service) mfilter(c context.Context, msgs []string, now time.Time) (checked []string, err error) {
	params := url.Values{}
	params.Set("area", "tag")
	for _, msg := range msgs {
		params.Add("msg", msg)
	}
	var res struct {
		Code int             `json:"code"`
		Data []*model.Filter `json:"data"`
	}
	if err = s.client.Get(c, s.mFilterURL, metadata.String(c, metadata.RemoteIP), params, &res); err != nil {
		conf.PromError("过滤服务多个接口", "d.client.Get(%s) error(%v)", s.mFilterURL+"?"+params.Encode(), err)
		log.Error("s.client.Get(%s) error(%v)", s.mFilterURL+"?"+params.Encode(), err)
		return
	}
	if res.Code != 0 {
		conf.PromError("过滤服务多个接口", "d.client.Get(%s) res code(%d) or res.data(%v)", s.mFilterURL+"?"+params.Encode(), res.Code, res.Data)
		log.Error("s.client.Get(%s) res code(%d) or res.data(%v)", s.mFilterURL+"?"+params.Encode(), res.Code, res.Data)
		err = errFilterFaile
		return
	}
	for _, f := range res.Data {
		if f.Level <= 10 {
			checked = append(checked, f.Msg)
		}
	}
	return
}

// TODO
// AI接口的业务请求来源针对频道详情页只有一个字段：9，待AI接口修改后另行修改.
func (s *Service) aiRecommand(c context.Context, arg *model.ArgChannelResource) (oids []int64, err error) {
	params := url.Values{}
	params.Set("cmd", "video")
	params.Set("timeout", strconv.Itoa(s.c.Tag.AITimeout))
	params.Set("mid", strconv.FormatInt(arg.Mid, 10))
	params.Set("buvid", arg.Buvid)
	params.Set("build", strconv.FormatInt(int64(arg.Build), 10))
	params.Set("plat", strconv.FormatInt(int64(arg.Plat), 10))
	params.Set("login_event", strconv.FormatInt(int64(arg.LoginEvent), 10))
	params.Set("request_cnt", strconv.FormatInt(int64(arg.RequestCNT), 10))
	params.Set("display_id", strconv.FormatInt(int64(arg.DisplayID), 10))
	params.Set("page_type", strconv.FormatInt(int64(arg.From), 10))
	switch arg.Channel {
	case model.TagChannelYes:
		params.Set("chn_id", strconv.FormatInt(arg.Tid, 10))
		params.Set("from", strconv.FormatInt(int64(model.AIRecommandChannel), 10))
		params.Set("top_channel", strconv.FormatInt(int64(model.TagChannelYes), 10))
	default:
		params.Set("tag", strconv.FormatInt(arg.Tid, 10))
		params.Set("from", strconv.FormatInt(int64(model.AIRecommandTag), 10))
		params.Set("top_channel", strconv.FormatInt(int64(model.TagChannelNo), 10))
	}
	res := &struct {
		Code int                         `json:"code"`
		Data []*model.AIChannelRecommand `json:"data"`
	}{}
	if err = s.client.Get(c, s.aiRecommandlURL, arg.RealIP, params, &res); err != nil {
		conf.PromError("AI channel recommand", "d.client.Get(%s) error(%v)", s.aiRecommandlURL+"?"+params.Encode(), err)
		log.Error("s.client.Get(%s) error(%v)", s.aiRecommandlURL+"?"+params.Encode(), err)
		return
	}
	switch res.Code {
	case 0, -4:
		for _, data := range res.Data {
			oids = append(oids, data.Oid)
		}
	case -3:
		err = ecode.ChannelAINoData
	case -2:
		err = ecode.ChannelAITimeout
	default:
		err = errAIRecommandFaile
	}
	if err != nil {
		conf.PromError("AI channel recommand", "d.client.Get(%s) error(%v)", s.aiRecommandlURL+"?"+params.Encode(), err)
		log.Error("s.client.Get(%s) res code(%d) or res.data(%v)", s.aiRecommandlURL+"?"+params.Encode(), res.Code, res.Data)
	}
	return
}

func (s *Service) workflowAppeaList(c context.Context, appeal *model.WorkflowAppealInfo) (appeals []*model.WorkflowAppeal, err error) {
	uri := s.c.Host.APICo + _workflowAppeaListURI
	params := url.Values{}
	params.Set("business", fmt.Sprintf("%d", appeal.Business))
	params.Set("mid", fmt.Sprintf("%d", appeal.RptMid))
	params.Set("oid", fmt.Sprintf("%d", appeal.Oid))
	res := &struct {
		Code int                     `json:"code"`
		Data []*model.WorkflowAppeal `json:"data"`
	}{}
	if err = s.client.Get(c, uri, appeal.RealIP, params, &res); err != nil {
		log.Error("s.client.Get(%s) error(%v)", uri+"?"+params.Encode(), err)
		return
	}
	if res.Code != ecode.OK.Code() {
		log.Error("s.client.Get(%s) res code(%d) or res.data(%v)", uri+"?"+params.Encode(), res.Code, res.Data)
		err = errappeaListFaile
		return
	}
	appeals = res.Data
	return
}

func (s *Service) workflowAppealState(c context.Context, appeal *model.WorkflowAppealInfo) (state int32, err error) {
	uri := s.c.Host.APICo + _workflowAppealStateURI
	params := url.Values{}
	params.Set("business", fmt.Sprintf("%d", appeal.Business))
	params.Set("eid", fmt.Sprintf("%d", appeal.EID))
	params.Set("oid", fmt.Sprintf("%d", appeal.Oid))
	res := &struct {
		Code int `json:"code"`
		Data struct {
			State int32 `json:"state"`
		} `json:"data"`
	}{}
	if err = s.client.Get(c, uri, appeal.RealIP, params, &res); err != nil {
		log.Error("s.client.Get(%s) error(%v)", uri+"?"+params.Encode(), err)
		return
	}
	if res.Code != ecode.OK.Code() {
		log.Error("s.client.Get(%s) res code(%d) or res.data(%v)", uri+"?"+params.Encode(), res.Code, res.Data)
		err = errappealStateFaile
		return
	}
	state = res.Data.State
	return
}

func (s *Service) addWorkflowAppeal(c context.Context, appeal *model.WorkflowAppealInfo) (err error) {
	uri := s.c.Host.APICo + _addWorkflowAppeaURI
	params := url.Values{}
	params.Set("business", fmt.Sprintf("%d", appeal.Business))
	params.Set("fid", fmt.Sprintf("%d", appeal.FID))
	params.Set("rid", fmt.Sprintf("%d", appeal.RID))
	params.Set("eid", fmt.Sprintf("%d", appeal.EID))
	params.Set("score", fmt.Sprintf("%d", appeal.Score))
	params.Set("tid", fmt.Sprintf("%d", appeal.ReasonID))
	params.Set("oid", fmt.Sprintf("%d", appeal.Oid))
	params.Set("mid", fmt.Sprintf("%d", appeal.RptMid))
	params.Set("business_typeid", fmt.Sprintf("%d", appeal.RegionID))
	params.Set("business_mid", fmt.Sprintf("%d", appeal.Mid))
	params.Set("business_title", appeal.TName)
	res := &struct {
		Code int `json:"code"`
	}{}
	if err = s.client.Post(c, uri, appeal.RealIP, params, &res); err != nil {
		log.Error("s.client.Get(%s) error(%v)", uri+"?"+params.Encode(), err)
		return
	}
	if res.Code != ecode.OK.Code() {
		log.Error("s.client.Get(%s) res code(%d)", uri+"?"+params.Encode(), res.Code)
		err = errAddAppealFaile
	}
	return
}
