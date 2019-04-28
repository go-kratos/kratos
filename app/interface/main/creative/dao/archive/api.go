package archive

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"go-common/app/interface/main/creative/conf"
	"go-common/app/interface/main/creative/dao/tool"
	"go-common/app/interface/main/creative/model/archive"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/xstr"
)

const (
	_view          = "/videoup/view"
	_views         = "/videoup/views"
	_del           = "/videoup/del"
	_video         = "/videoup/cid"
	_archives      = "/videoup/up/archives"
	_descFormat    = "/videoup/desc/format"
	_simpleVideos  = "/videoup/simplevideos"
	_simpleArchive = "/videoup/simplearchive"
	_videoJam      = "/videoup/video/jam"
	_upSpecial     = "/x/internal/uper/special"
	_flowjudge     = "/videoup/flow/list/judge"
	_staffApplies  = "/videoup/staff/apply/filter"
	_staffApply    = "/videoup/staff/apply/submit"
	_staffCheck    = "/videoup/staff/mid/applys"
)

// SimpleArchive fn
func (d *Dao) SimpleArchive(c context.Context, aid int64, ip string) (sa *archive.SpArchive, err error) {
	params := url.Values{}
	params.Set("aid", strconv.FormatInt(aid, 10))
	var res struct {
		Code int                `json:"code"`
		Data *archive.SpArchive `json:"data"`
	}
	if err = d.client.Get(c, d.simpleArchive, ip, params, &res); err != nil {
		log.Error("archive.simpleArchive url(%s) mid(%d) error(%v)", d.simpleArchive+"?"+params.Encode(), aid, err)
		err = ecode.CreativeArchiveAPIErr
		return
	}
	if res.Code != 0 {
		log.Error("archive.simpleArchive url(%s) mid(%d) res(%v)", d.simpleArchive+"?"+params.Encode(), aid, res)
		err = ecode.Int(res.Code)
		return
	}
	if res.Data == nil {
		log.Error("archive.simpleArchive url(%s) aid(%d) res(%v)", d.simpleVideos+"?"+params.Encode(), aid, res)
		return
	}
	sa = res.Data
	return
}

// SimpleVideos fn
func (d *Dao) SimpleVideos(c context.Context, aid int64, ip string) (vs []*archive.SpVideo, err error) {
	params := url.Values{}
	params.Set("aid", strconv.FormatInt(aid, 10))
	var res struct {
		Code int                `json:"code"`
		Data []*archive.SpVideo `json:"data"`
	}
	if err = d.client.Get(c, d.simpleVideos, ip, params, &res); err != nil {
		log.Error("archive.simpleVideos url(%s) mid(%d) error(%v)", d.simpleVideos+"?"+params.Encode(), aid, err)
		err = ecode.CreativeArchiveAPIErr
		return
	}
	if res.Code != 0 {
		log.Error("archive.simpleVideos url(%s) aid(%d) res(%v)", d.simpleVideos+"?"+params.Encode(), aid, res)
		err = ecode.Int(res.Code)
		return
	}
	if res.Data == nil {
		log.Error("archive.simpleVideos url(%s) aid(%d) res(%v)", d.simpleVideos+"?"+params.Encode(), aid, res)
		return
	}
	vs = res.Data
	sort.Slice(vs, func(i, j int) bool {
		return vs[i].Index <= vs[j].Index
	})
	return
}

// View get archive
func (d *Dao) View(c context.Context, mid, aid int64, ip string, needPOI, needVote int) (av *archive.ArcVideo, err error) {
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("aid", strconv.FormatInt(aid, 10))
	params.Set("need_poi", strconv.Itoa(needPOI))
	params.Set("need_vote", strconv.Itoa(needVote))
	var res struct {
		Code int               `json:"code"`
		Data *archive.ArcVideo `json:"data"`
	}
	if err = d.client.Get(c, d.view, ip, params, &res); err != nil {
		log.Error("archive.view url(%s) mid(%d) error(%v)", d.view+"?"+params.Encode(), mid, err)
		err = ecode.CreativeArchiveAPIErr
		return
	}
	if res.Code != 0 {
		log.Error("archive.view url(%s) mid(%d) res(%v)", d.view+"?"+params.Encode(), mid, res)
		err = ecode.Int(res.Code)
		return
	}
	if res.Data.Archive == nil {
		log.Error("archive.view url(%s) mid(%d) res(%v)", d.view+"?"+params.Encode(), mid, res)
		return
	}
	if res.Data.Archive.Staffs == nil {
		res.Data.Archive.Staffs = []*archive.StaffView{}
	}
	res.Data.Archive.StateDesc = d.c.StatDesc(int(res.Data.Archive.State))
	res.Data.Archive.StatePanel = archive.StatePanel(res.Data.Archive.State)
	av = res.Data
	return
}

// Views get archives
func (d *Dao) Views(c context.Context, mid int64, aids []int64, ip string) (avm map[int64]*archive.ArcVideo, err error) {
	params := url.Values{}
	params.Set("aids", xstr.JoinInts(aids))
	params.Set("mid", strconv.FormatInt(mid, 10))
	var res struct {
		Code int                          `json:"code"`
		Data map[string]*archive.ArcVideo `json:"data"`
	}
	if err = d.client.Get(c, d.views, ip, params, &res); err != nil {
		log.Error("archive.views url(%s) mid(%d) error(%v)", d.views+"?"+params.Encode(), mid, err)
		err = ecode.CreativeArchiveAPIErr
		return
	}
	if res.Code != 0 {
		log.Error("archive.views url(%s) mid(%d) res(%v)", d.views+"?"+params.Encode(), mid, res)
		err = ecode.Int(res.Code)
		return
	}
	avm = map[int64]*archive.ArcVideo{}
	for aidStr, av := range res.Data {
		var err error
		aid, err := strconv.ParseInt(aidStr, 10, 64)
		if err != nil {
			log.Error("strconv.ParseInt(%s) error(%v)", aidStr, err)
			continue
		}
		av.Archive.StateDesc = d.c.StatDesc(int(av.Archive.State))
		av.Archive.StatePanel = archive.StatePanel(av.Archive.State)
		avm[aid] = av
	}
	return
}

// Del delete archive.
func (d *Dao) Del(c context.Context, mid, aid int64, ip string) (err error) {
	params := url.Values{}
	params.Set("appkey", conf.Conf.App.Key)
	params.Set("appsecret", conf.Conf.App.Secret)
	params.Set("ts", strconv.FormatInt(time.Now().Unix(), 10))
	// uri
	var (
		query, _ = tool.Sign(params)
		uri      = d.del + "?" + query
	)
	// new request
	req, err := http.NewRequest("POST", uri, strings.NewReader(fmt.Sprintf(`{"mid":%d,"aid":%d}`, mid, aid)))
	if err != nil {
		log.Error("http.NewRequest(%s) error(%v); mid(%d), aid(%d), ip(%s)", d.del, err, mid, aid, ip)
		err = ecode.CreativeArchiveAPIErr
		return
	}
	req.Header.Set("X-BACKEND-BILI-REAL-IP", ip)
	var res struct {
		Code int `json:"code"`
	}
	if err = d.client.Do(c, req, &res); err != nil {
		log.Error("d.client.Do(%s) error(%v); mid(%d), aid(%d), ip(%s)", d.del, err, mid, aid, ip)
		err = ecode.CreativeArchiveAPIErr
		return
	}
	if res.Code != 0 {
		err = ecode.Int(res.Code)
		log.Error("del archive url(%s) res(%v); mid(%d), aid(%d), ip(%s)", d.del, res, mid, aid, ip)
		return
	}
	return
}

// VideoByCid get videos by cids.
func (d *Dao) VideoByCid(c context.Context, cid int64, ip string) (v *archive.Video, err error) {
	params := url.Values{}
	params.Set("cid", strconv.FormatInt(cid, 10))
	var res struct {
		Code int            `json:"code"`
		Data *archive.Video `json:"data"`
	}
	if err = d.client.Get(c, d.video, ip, params, &res); err != nil {
		log.Error("VideoByCid cidURI(%s)  error(%v)", d.video+"?"+params.Encode(), err)
		return
	}
	if res.Code != 0 {
		err = ecode.Int(res.Code)
		log.Error("VideoByCid cidURI(%s) Code=(%d)", d.video+"?"+params.Encode(), res.Code)
		return
	}
	v = res.Data
	return
}

// UpArchives get archives by mid.
func (d *Dao) UpArchives(c context.Context, mid, pn, ps, group int64, ip string) (aids []int64, count int64, err error) {
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("pn", strconv.FormatInt(pn, 10))
	params.Set("ps", strconv.FormatInt(ps, 10))
	params.Set("group", strconv.FormatInt(mid, 10)) //0:全部稿件；1:开放稿件；2:未开放稿件
	var res struct {
		Code int `json:"code"`
		Data struct {
			Aids  []int64 `json:"aids"`
			Count int64   `json:"count"`
		} `json:"data"`
	}
	if err = d.client.Get(c, d.upArchives, ip, params, &res); err != nil {
		log.Error("upArchives URI(%s)  error(%v)", d.upArchives+"?"+params.Encode(), err)
		err = ecode.CreativeArchiveAPIErr
		return
	}
	if res.Code != 0 {
		log.Error("upArchives URI(%s) Code=(%d)", d.upArchives+"?"+params.Encode(), res.Code)
		err = ecode.Int(res.Code)
		return
	}
	aids = res.Data.Aids
	count = res.Data.Count
	return
}

// DescFormat get desc format by typeid and copyright
func (d *Dao) DescFormat(c context.Context) (descs []*archive.DescFormat, err error) {
	params := url.Values{}
	var res struct {
		Code    int                   `json:"code"`
		Message string                `json:"message"`
		Data    []*archive.DescFormat `json:"data"`
	}
	descs = []*archive.DescFormat{}
	if err = d.client.Get(c, d.descFormat, "", params, &res); err != nil {
		log.Error("videoup descFormat error(%v) | descFormat(%s) params(%v)", err, d.descFormat+"?"+params.Encode(), params)
		err = ecode.CreativeArchiveAPIErr
		return
	}
	if res.Code != 0 {
		log.Error("videoup descFormat res.Code(%d) | descFormat(%s)params(%v) res(%v)", res.Code, d.descFormat+"?"+params.Encode(), params, res)
		err = ecode.Int(res.Code)
		return
	}
	descs = res.Data
	log.Info("res.Code(%d) descFormat(%s) params(%v)", res.Code, d.descFormat+"?"+params.Encode(), params)
	return
}

// VideoJam get video-check traffic jam level
func (d *Dao) VideoJam(c context.Context, ip string) (level int8, err error) {
	params := url.Values{}
	var res struct {
		Code int `json:"code"`
		Data *struct {
			Level int8 `json:"level"`
		} `json:"data"`
	}
	if err = d.client.Get(c, d.videoJam, ip, params, &res); err != nil {
		log.Error("archive.VideoJam url(%s) error(%v)", d.videoJam+"?"+params.Encode(), err)
		err = ecode.CreativeArchiveAPIErr
		return
	}
	if res.Code != 0 {
		log.Error("archive.VideoJam url(%s) res(%v)", d.videoJam+"?"+params.Encode(), res)
		err = ecode.Int(res.Code)
		return
	}
	if res.Data == nil {
		log.Error("archive.VideoJam url(%s) res(%v)", d.videoJam+"?"+params.Encode(), res)
		return
	}
	level = res.Data.Level
	return
}

// FlowJudge fn
func (d *Dao) FlowJudge(c context.Context, business, groupID int64, oids []int64) (hitOids []int64, err error) {
	params := url.Values{}
	params.Set("business", strconv.FormatInt(business, 10))
	params.Set("gid", strconv.FormatInt(groupID, 10))
	params.Set("oids", xstr.JoinInts(oids))
	var res struct {
		Code    int     `json:"code"`
		Message string  `json:"message"`
		Data    []int64 `json:"data"`
	}
	if err = d.client.Get(c, d.flowJudge, "", params, &res); err != nil {
		log.Error("archive.FlowJudge url(%s) error(%v)", d.upSpecialURL+"?"+params.Encode(), err)
		return
	}
	if res.Code != 0 {
		log.Error("FlowJudge api url(%s) res(%v) code(%d)", d.upSpecialURL, res, res.Code)
		err = ecode.CreativeArchiveAPIErr
		return
	}
	hitOids = res.Data
	return
}

// StaffApplies fn
func (d *Dao) StaffApplies(c context.Context, staffMid int64, aids []int64) (apply []*archive.StaffApply, err error) {
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(staffMid, 10))
	params.Set("aids", xstr.JoinInts(aids))
	var res struct {
		Code    int                   `json:"code"`
		Message string                `json:"message"`
		Data    []*archive.StaffApply `json:"data"`
	}
	if err = d.client.Get(c, d.staffApplies, "", params, &res); err != nil {
		log.Error("archive.staffApplies url(%s) error(%v)", d.staffApplies+"?"+params.Encode(), err)
		return
	}
	if res.Code != 0 {
		log.Error("staffApplies api url(%s) res(%v) code(%d)", d.staffApplies, res, res.Code)
		err = ecode.CreativeArchiveAPIErr
		return
	}
	apply = res.Data
	return
}

// StaffMidValidate fn
func (d *Dao) StaffMidValidate(c context.Context, mid int64) (data int, err error) {
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	var res struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    int    `json:"data"`
	}
	if err = d.client.Get(c, d.staffCheck, "", params, &res); err != nil {
		log.Error("archive.StaffMidValidate url(%s) error(%v)", d.staffCheck+"?"+params.Encode(), err)
		return
	}
	if res.Code != 0 {
		log.Error("StaffMidValidate api url(%s) res(%v) code(%d)", d.staffCheck, res, res.Code)
		err = ecode.Int(res.Code)
		return
	}
	data = res.Data
	return
}

// StaffApplySubmit fn
func (d *Dao) StaffApplySubmit(c context.Context, id, aid, mid, state, atype int64, flagAddBlack, flagRefuse int) (err error) {
	params := url.Values{}
	params.Set("id", strconv.FormatInt(id, 10))
	params.Set("state", strconv.FormatInt(state, 10))
	params.Set("type", strconv.FormatInt(atype, 10))
	params.Set("apply_aid", strconv.FormatInt(aid, 10))
	params.Set("apply_staff_mid", strconv.FormatInt(mid, 10))
	params.Set("flag_add_black", strconv.Itoa(flagAddBlack))
	params.Set("flag_refuse", strconv.Itoa(flagRefuse))
	var res struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
	if err = d.client.Post(c, d.staffApply, "", params, &res); err != nil {
		log.Error("archive.StaffApplySubmit url(%s) error(%v)", d.staffApply+"?"+params.Encode(), err)
		return
	}
	if res.Code != 0 {
		log.Error("StaffApplySubmit api url(%s) res(%v) code(%d)", d.staffApply, res, res.Code)
		err = ecode.Int(res.Code)
		return
	}
	return
}
