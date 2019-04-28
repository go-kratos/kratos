package workflow

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	"go-common/app/interface/main/reply/conf"
	model "go-common/app/interface/main/reply/model/reply"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/xstr"
)

const (
	workflowHost     = "api.bilibili.co"
	_apiTopics       = "http://matsuri.bilibili.co/activity/pages"
	_apiActivitySub  = "http://matsuri.bilibili.co//activity/subject/url"
	_apiLiveActivity = "http://api.live.bilibili.co/comment/v1/relation/get_by_id"
	_apiNotice       = "http://api.bilibili.co/x/internal/credit/publish/infos"
	_apiBan          = "http://api.bilibili.co/x/internal/credit/blocked/infos"
	_apiCredit       = "http://api.bilibili.co/x/internal/credit/blocked/cases"
	_apiLiveVideo    = "http://api.vc.bilibili.co/clip/v1/video/detail"
	_apiLiveNotice   = "http://api.vc.bilibili.co/news/v1/notice/info"
	_apiLivePicture  = "http://api.vc.bilibili.co/link_draw/v1/doc/detail"
	_apiTopic        = "http://matsuri.bilibili.co/activity/page/one/%d"
	_apiDynamic      = "http://api.vc.bilibili.co/dynamic_repost/v0/dynamic_repost/ftch_rp_cont?dynamic_ids[]=%d"
	_apiHuoniao      = "http://manga.bilibili.co/twirp/verify.v0.Verify/InfoC"
	// link
	_linkBan         = "https://www.bilibili.com/blackroom/ban/%d"
	_linkNotice      = "https://www.bilibili.com/blackroom/notice/%d"
	_linkCredit      = "https://www.bilibili.com/judgement/case/%d"
	_linkLiveVideo   = "http://vc.bilibili.com/video/%d"
	_linkLiveNotice  = "http://link.bilibili.com/p/eden/news#/newsdetail?id=%d"
	_linkLivePicture = "http://h.bilibili.com/ywh/%d"
	_linkDynamic     = "http://t.bilibili.com/%d"
	_urlLiveNotice   = "http://api.vc.bilibili.co/news/v1/notice/info"
	_linkHuoniao     = "http://manga.bilibili.com/m/detail/mc%d"
)

type result struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		ChallengeNo int64 `json:"challengeNo"`
	}
}

// Dao dao.
type Dao struct {
	httpClient *bm.Client
}

// New new a dao and return.
func New(c *conf.Config) *Dao {
	d := &Dao{
		httpClient: bm.NewClient(c.HTTPClient),
	}
	return d
}

type extra struct {
	Like  int    `json:"like"`
	Link  string `json:"link"`
	Title string `json:"title"`
}

func (dao *Dao) AddReport(c context.Context, oid int64, typ int8, typeid int32, rpid int64, score int, reason int8, reporter int64, reported int64, like int, content string, link string, title string) (err error) {
	var res result
	rid := "1"
	params := url.Values{}
	params.Add("business", "13")
	params.Add("fid", strconv.FormatInt(int64(typ), 10))
	if model.GetReportType(reason) == model.ReportStateNewTwo {
		rid = "2"
	}
	params.Add("rid", rid)
	params.Add("eid", strconv.FormatInt(int64(rpid), 10))
	params.Add("score", strconv.FormatInt(int64(score), 10))
	params.Add("tid", strconv.FormatInt(int64(reason), 10))
	params.Add("oid", strconv.FormatInt(int64(oid), 10))
	params.Add("mid", strconv.FormatInt(int64(reporter), 10))
	params.Add("business_typeid", strconv.FormatInt(int64(typeid), 10))
	params.Add("business_title", content)
	params.Add("business_mid", strconv.FormatInt(int64(reported), 10))
	var ex extra
	ex.Like = like
	ex.Link = link
	ex.Title = title
	businessExtra, _ := json.Marshal(&ex)
	params.Add("business_extra", string(businessExtra))

	url := fmt.Sprintf("http://%s/%s", workflowHost, "x/internal/workflow/appeal/v3/add")
	if err = dao.httpClient.Post(c, url, "", params, &res); err != nil {
		log.Error("AddReport url(%s,%s) error(%v)", url, params.Encode(), err)
		return
	}
	if res.Code != 0 {
		log.Error("AddReport url(%s,%s) error(%v)", url, params.Encode(), res.Code)
		err = ecode.Int(res.Code)
		return
	}
	return
}

// TopicsLink get topic info.
func (d *Dao) TopicsLink(c context.Context, links map[int64]string, isTopic bool) (err error) {
	if len(links) == 0 {
		return
	}
	var res struct {
		Code int `json:"code"`
		Data *struct {
			List []struct {
				ID    int64  `json:"id"`
				PCURL string `json:"pc_url"`
				H5URL string `json:"h5_url"`
			} `json:"list"`
		} `json:"data"`
	}
	var ids []int64
	for oid := range links {
		ids = append(ids, oid)
	}
	params := url.Values{}
	params.Set("pids", xstr.JoinInts(ids))
	params.Set("all", "isOne")
	if isTopic {
		params.Set("mold", "1")
	}
	if err = d.httpClient.Get(c, _apiTopics, "", params, &res); err != nil {
		log.Error("TopicsTitle(%s) error(%v)", _apiTopics, err)
		return
	}
	if res.Data == nil {
		err = fmt.Errorf("url:%s code:%d", _apiTopics, res.Code)
		return
	}
	for _, data := range res.Data.List {
		if data.PCURL != "" {
			links[data.ID] = data.PCURL
		} else {
			links[data.ID] = data.H5URL
		}
	}
	return
}

// ActivitySub return activity sub link.
func (d *Dao) ActivitySub(c context.Context, oid int64) (title, link string, err error) {
	params := url.Values{}
	params.Set("oid", strconv.FormatInt(oid, 10))
	var res struct {
		Code int `json:"code"`
		Data *struct {
			Title string `json:"name"`
			Link  string `json:"act_url"`
		} `json:"data"`
	}
	if err = d.httpClient.Get(c, _apiActivitySub, "", params, &res); err != nil {
		log.Error("d.httpClient.Get(%s?%s) error(%v)", _apiActivitySub, params.Encode(), err)
		return
	}
	if res.Data == nil {
		err = fmt.Errorf("url:%s code:%d", _apiActivitySub, res.Code)
		return
	}
	title = res.Data.Title
	link = res.Data.Link
	return
}

// LiveNotice return link.
func (d *Dao) LiveNotice(c context.Context, oid int64) (title string, err error) {
	params := url.Values{}
	params.Set("id", strconv.FormatInt(oid, 10))
	var res struct {
		Code int `json:"code"`
		Data *struct {
			Title string `json:"title"`
		} `json:"data"`
	}
	if err = d.httpClient.Get(c, _urlLiveNotice, "", params, &res); err != nil {
		log.Error("d.httpClient.Get(%s?%s) error(%v)", _urlLiveNotice, params.Encode(), err)
		return
	}
	if res.Code != 0 || res.Data == nil {
		err = fmt.Errorf("url:%s?%s code:%d", _urlLiveNotice, params.Encode(), res.Code)
		return
	}
	title = res.Data.Title
	return
}

// LiveActivityTitle get live activity info.
func (d *Dao) LiveActivityTitle(c context.Context, oid int64) (title, link string, err error) {
	params := url.Values{}
	params.Set("id", strconv.FormatInt(oid, 10))
	var res struct {
		Code int `json:"code"`
		Data *struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"data"`
	}
	if err = d.httpClient.Get(c, _apiLiveActivity, "", params, &res); err != nil {
		log.Error("httpLiveActivityTitle(%s?%s) error(%v)", _apiLiveActivity, params.Encode(), err)
		return
	}
	if res.Code != 0 || res.Data == nil {
		err = fmt.Errorf("url:%s?%s code:%d", _apiLiveActivity, params.Encode(), res.Code)
		return
	}
	title = res.Data.Name
	link = res.Data.URL
	return
}

type notice struct {
	Title string `json:"title"`
}

type ban struct {
	Title string `json:"punishTitle"`
}
type credit struct {
	Title string `json:"punishTitle"`
}

// NoticeTitle get blackromm notice info.
func (d *Dao) NoticeTitle(c context.Context, oid int64) (title, link string, err error) {
	params := url.Values{}
	params.Set("ids", strconv.FormatInt(oid, 10))
	var res struct {
		Code int               `json:"code"`
		Data map[int64]*notice `json:"data"`
	}
	if err = d.httpClient.Get(c, _apiNotice, "", params, &res); err != nil {
		log.Error("httpNotice(%s) error(%v)", _apiNotice, err)
		return
	}
	if r := res.Data[oid]; r != nil {
		title = r.Title
	}
	link = fmt.Sprintf(_linkNotice, oid)
	return
}

// BanTitle get ban info.
func (d *Dao) BanTitle(c context.Context, oid int64) (title, link string, err error) {
	params := url.Values{}
	params.Set("ids", strconv.FormatInt(oid, 10))
	var res struct {
		Code int            `json:"code"`
		Data map[int64]*ban `json:"data"`
	}
	if err = d.httpClient.Get(c, _apiBan, "", params, &res); err != nil {
		log.Error("httpBan(%s) error(%v)", _apiBan, err)
		return
	}
	if r := res.Data[oid]; r != nil {
		title = r.Title
	}
	link = fmt.Sprintf(_linkBan, oid)
	return
}

// CreditTitle return link.
func (d *Dao) CreditTitle(c context.Context, oid int64) (title, link string, err error) {
	params := url.Values{}
	params.Set("ids", strconv.FormatInt(oid, 10))
	var res struct {
		Code int               `json:"code"`
		Data map[int64]*credit `json:"data"`
	}
	if err = d.httpClient.Get(c, _apiCredit, "", params, &res); err != nil {
		log.Error("d.httpClient.Get(%s?%s) error(%v)", _apiCredit, params.Encode(), err)
		return
	}
	if res.Code != 0 || res.Data == nil {
		err = fmt.Errorf("url:%s?%s code:%d", _apiCredit, params.Encode(), res.Code)
		return
	}
	if r := res.Data[oid]; r != nil {
		title = r.Title
	}
	link = fmt.Sprintf(_linkCredit, oid)
	return
}

// LiveVideoTitle get live video title.
func (d *Dao) LiveVideoTitle(c context.Context, oid int64) (title, link string, err error) {
	params := url.Values{}
	params.Set("video_id", strconv.FormatInt(oid, 10))
	var res struct {
		Code int `json:"code"`
		Data *struct {
			Item *struct {
				Description string `json:"description"`
			} `json:"item"`
		} `json:"data"`
	}
	if err = d.httpClient.Get(c, _apiLiveVideo, "", params, &res); err != nil {
		log.Error("httpLiveVideoTitle(%s?%s) error(%v)", _apiLiveVideo, params.Encode(), err)
		return
	}
	if res.Code != 0 || res.Data == nil || res.Data.Item == nil {
		err = fmt.Errorf("url:%s?%s code:%d", _apiLiveVideo, params.Encode(), res.Code)
		return
	}
	title = res.Data.Item.Description
	link = fmt.Sprintf(_linkLiveVideo, oid)
	return
}

// LiveNoticeTitle get live notice info.
func (d *Dao) LiveNoticeTitle(c context.Context, oid int64) (title, link string, err error) {
	params := url.Values{}
	params.Set("id", strconv.FormatInt(oid, 10))
	var res struct {
		Code int `json:"code"`
		Data *struct {
			Title string `json:"title"`
		} `json:"data"`
	}
	if err = d.httpClient.Get(c, _apiLiveNotice, "", params, &res); err != nil {
		log.Error("LiveNoticeTitle(%s?%s) error(%v)", _apiLiveNotice, params.Encode(), err)
		return
	}
	if res.Code != 0 || res.Data == nil {
		err = fmt.Errorf("url:%s?%s code:%d", _apiLiveNotice, params.Encode(), res.Code)
		return
	}
	title = res.Data.Title
	link = fmt.Sprintf(_linkLiveNotice, oid)
	return
}

// LivePictureTitle get live picture info.
func (d *Dao) LivePictureTitle(c context.Context, oid int64) (title, link string, err error) {
	params := url.Values{}
	params.Set("doc_id", strconv.FormatInt(oid, 10))
	var res struct {
		Code int `json:"code"`
		Data *struct {
			Item *struct {
				Title string `json:"title"`
				Desc  string `json:"description"`
			} `json:"item"`
		} `json:"data"`
	}
	if err = d.httpClient.Get(c, _apiLivePicture, "", params, &res); err != nil {
		log.Error("LivePictureTitle(%s?%s) error(%v)", _apiLivePicture, params.Encode(), err)
		return
	}
	if res.Code != 0 || res.Data == nil || res.Data.Item == nil {
		err = fmt.Errorf("url:%s?%s code:%d", _apiLivePicture, params.Encode(), res.Code)
		return
	}
	title = res.Data.Item.Title
	if title == "" {
		title = res.Data.Item.Desc
	}
	link = fmt.Sprintf(_linkLivePicture, oid)
	return
}

// TopicTitle get topic info.
func (d *Dao) TopicTitle(c context.Context, oid int64) (title, link string, err error) {
	var res struct {
		Code int `json:"code"`
		Data *struct {
			Title  string `json:"name"`
			PCLink string `json:"pc_url"`
			H5Link string `json:"h5_url"`
		} `json:"data"`
	}
	if err = d.httpClient.Get(c, fmt.Sprintf(_apiTopic, oid), "", nil, &res); err != nil {
		log.Error("TopicTitle(%s) error(%v)", fmt.Sprintf(_apiTopic, oid), err)
		return
	}
	if res.Data == nil {
		err = fmt.Errorf("url:%s code:%d", fmt.Sprintf(_apiTopic, oid), res.Code)
		return
	}
	title = res.Data.Title
	link = res.Data.PCLink
	if link == "" {
		link = res.Data.H5Link
	}
	return
}

// DynamicTitle return link and title.
func (d *Dao) DynamicTitle(c context.Context, oid int64) (title, link string, err error) {
	params := url.Values{}
	uri := fmt.Sprintf(_apiDynamic, oid)
	var res struct {
		Code int `json:"code"`
		Data *struct {
			Pairs []struct {
				DynamicID int64  `json:"dynamic_id"`
				Content   string `json:"rp_cont"`
				Type      int32  `json:"type"`
			} `json:"pairs"`
			TotalCount int64 `json:"total_count"`
		} `json:"data,omitempty"`
		Message string `json:"message"`
	}
	if err = d.httpClient.Get(c, uri, "", params, &res); err != nil {
		log.Error("d.httpClient.Get(%s?%s) error(%v)", uri, params.Encode(), err)
		return
	}

	if res.Code != 0 || res.Data == nil || len(res.Data.Pairs) == 0 {
		err = fmt.Errorf("get dynamic failed!url:%s?%s code:%d message:%s pairs:%v", uri, params.Encode(), res.Code, res.Message, res.Data.Pairs)
		return
	}
	title = res.Data.Pairs[0].Content
	link = fmt.Sprintf(_linkDynamic, oid)
	return
}

// HuoniaoTitle title and link
func (d *Dao) HuoniaoTitle(c context.Context, oid int64) (title, link string, err error) {
	params := url.Values{}
	params.Set("id", strconv.FormatInt(oid, 10))

	var res struct {
		Code int `json:"code"`
		Data *struct {
			Title   string `json:"title"`
			JumpURL string `json:"jump_url"`
		} `json:"data,omitempty"`
		Message string `json:"msg"`
	}

	if err = d.httpClient.Post(c, _apiHuoniao, "", params, &res); err != nil {
		log.Error("httpClient.Post(%s) error(%v)", params.Encode(), err)
		return
	}
	if res.Code != ecode.OK.Code() {
		err = ecode.Int(res.Code)
		return
	}
	if res.Data != nil {
		title = res.Data.Title
		link = res.Data.JumpURL
	}
	return
}
