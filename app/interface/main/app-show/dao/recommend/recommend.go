package recommend

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"go-common/app/interface/main/app-show/conf"
	"go-common/app/interface/main/app-show/model/recommend"
	"go-common/library/ecode"
	"go-common/library/log"
	httpx "go-common/library/net/http/blademaster"
	xtime "go-common/library/time"

	"github.com/pkg/errors"
)

const (
	// _hotUrl           = "/y3kflg2k/ranking-m.json"
	_hotUrl    = "/data/rank/reco-tmzb.json"
	_regionUrl = "/8669rank/mobile_random/%s/1.json" // %s must be replaced to concrete tid
	// _regionHotUrl     = "/y3kflg2k/catalogy/%d-recommend-m.json"
	_regionListUrl = "/list"
	// _regionChildHotUrl = "/y3kflg2k/catalogy/catalogy-%d-3-m.json"
	_regionChildHotUrl   = "/data/rank/recent_region-%d-3.json"
	_regionArcListUrl    = "/x/v2/archive/rank"
	_rankRegionUrl       = "/y3kflg2k/rank/%s-03-%d.json"
	_rankOriginalUrl     = "/y3kflg2k/rank/%s-03.json"
	_rankBangumiUrl      = "/y3kflg2k/rank/all-3-33.json"
	_feedDynamicUrl      = "/feed/tag/top"
	_rankAllAppUrl       = "/data/rank/recent_all-app.json"
	_rankOriginAppUrl    = "/data/rank/recent_origin-app.json"
	_rankRegionAppUrl    = "/data/rank/recent_region-%d-app.json"
	_rankBangumiAppUrl   = "/data/rank/all_region-33-app.json"
	_hottabURL           = "/data/rank/reco-app-remen.json"
	_hotHeTongtabURL     = "/data/rank/reco-app-remen-%d.json"
	_hotHeTongtabcardURL = "/data/rank/reco-app-remen-card-%d.json"
)

// Dao is recommend dao.
type Dao struct {
	client              *httpx.Client
	clientAsyn          *httpx.Client
	clientParam         *httpx.Client
	hotUrl              string
	regionUrl           string
	regionChildHotUrl   string
	regionListUrl       string
	regionArcListUrl    string
	rankRegionUrl       string
	rankOriginalUrl     string
	rankBangumilUrl     string
	feedDynamicUrl      string
	rankAllAppUrl       string
	rankOriginAppUrl    string
	rankRegionAppUrl    string
	rankBangumiAppUrl   string
	hottabURL           string
	hotHetongURL        string
	hotHeTongtabcardURL string
}

//New recommend dao.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		client:      httpx.NewClient(conf.Conf.HTTPClient),
		clientAsyn:  httpx.NewClient(c.HTTPClientAsyn),
		clientParam: httpx.NewClient(conf.Conf.HTTPClient),
		// hotUrl:       c.Host.Hetongzi + _hotUrl,
		hotUrl:    c.Host.HetongziRank + _hotUrl,
		regionUrl: c.Host.Hetongzi + _regionUrl,
		// regionHotUrl: c.Host.Hetongzi + _regionHotUrl,
		// regionChildHotUrl: c.Host.Hetongzi + _regionChildHotUrl,
		regionChildHotUrl:   c.Host.HetongziRank + _regionChildHotUrl,
		regionListUrl:       c.Host.ApiCo + _regionListUrl,
		regionArcListUrl:    c.Host.ApiCoX + _regionArcListUrl,
		rankRegionUrl:       c.Host.Hetongzi + _rankRegionUrl,
		rankOriginalUrl:     c.Host.Hetongzi + _rankOriginalUrl,
		rankBangumilUrl:     c.Host.Hetongzi + _rankBangumiUrl,
		feedDynamicUrl:      c.Host.Data + _feedDynamicUrl,
		rankAllAppUrl:       c.Host.HetongziRank + _rankAllAppUrl,
		rankOriginAppUrl:    c.Host.HetongziRank + _rankOriginAppUrl,
		rankRegionAppUrl:    c.Host.HetongziRank + _rankRegionAppUrl,
		rankBangumiAppUrl:   c.Host.HetongziRank + _rankBangumiAppUrl,
		hottabURL:           c.Host.Data + _hottabURL,
		hotHetongURL:        c.Host.Data + _hotHeTongtabURL,
		hotHeTongtabcardURL: c.Host.Data + _hotHeTongtabcardURL,
	}
	return
}

// Hots get recommends.
func (d *Dao) Hots(c context.Context) (arcids []int64, err error) {
	var res struct {
		Code int `json:"code"`
		List []struct {
			Aid int64 `json:"aid"`
		} `json:"list"`
	}
	if err = d.clientAsyn.Get(c, d.hotUrl, "", nil, &res); err != nil {
		log.Error("recommend hots url(%s) error(%v)", d.hotUrl, err)
		return
	}
	b, _ := json.Marshal(&res)
	log.Info("recommend hots url(%s) json(%s)", d.hotUrl, b)
	if res.Code != 0 {
		log.Error("recommend hots url(%s) error(%v)", d.hotUrl, res.Code)
		err = fmt.Errorf("recommend api response code(%v)", res)
		return
	}
	for _, arcs := range res.List {
		arcids = append(arcids, arcs.Aid)
	}
	return
}

// Region get region recommend.
func (d *Dao) Region(c context.Context, tid string) (arcids []int64, err error) {
	var res struct {
		Code int `json:"code"`
		Data []struct {
			Aid string `json:"aid"`
		} `json:"list"`
	}
	api := fmt.Sprintf(d.regionUrl, tid)
	if err = d.clientAsyn.Get(c, api, "", nil, &res); err != nil {
		log.Error("recommend region url(%s) error(%v)", api, err)
		return
	}
	if res.Code != 0 {
		log.Error("recommend region url(%s) error(%v)", api, res.Code)
		err = fmt.Errorf("recommend region api response code(%v)", res)
		return
	}
	for _, arcs := range res.Data {
		arcids = append(arcids, aidToInt(arcs.Aid))
	}
	return
}

// RegionHots get hots recommend
func (d *Dao) RegionHots(c context.Context, tid int) (arcids []int64, err error) {
	var res struct {
		Code int `json:"code"`
		List []struct {
			Aid int64 `json:"aid"`
		} `json:"list"`
	}
	api := fmt.Sprintf(d.rankRegionAppUrl, tid)
	if err = d.clientAsyn.Get(c, api, "", nil, &res); err != nil {
		log.Error("recommend region hots url(%s) error(%v)", api, err)
		return
	}
	if res.Code != 0 {
		log.Error("recommend region hots url(%s) error(%v)", api, res.Code)
		err = fmt.Errorf("recommend region hots api response code(%v)", res)
		return
	}
	for _, arcs := range res.List {
		arcids = append(arcids, arcs.Aid)
	}
	return
}

// RegionList
func (d *Dao) RegionList(c context.Context, rid, tid, audit, pn, ps int, order string) (arcids []int64, err error) {
	params := url.Values{}
	params.Set("order", order)
	params.Set("filtered", strconv.Itoa(audit))
	params.Set("page", strconv.Itoa(pn))
	params.Set("pagesize", strconv.Itoa(ps))
	params.Set("tid", strconv.Itoa(rid))
	if tid > 0 {
		params.Set("tag_id", strconv.Itoa(tid))
	}
	params.Set("apiver", "2")
	params.Set("ver", "2")
	var res struct {
		Code int `json:"code"`
		List []struct {
			Aid interface{} `json:"aid"`
		} `json:"list"`
	}
	if err = d.client.Get(c, d.regionListUrl, "", params, &res); err != nil {
		log.Error("recommend region news url(%s) error(%v)", d.regionListUrl+"?"+params.Encode(), err)
		return
	}
	if res.Code != 0 && res.Code != -1 {
		log.Error("recommend region news url(%s) error(%v)", d.regionListUrl+"?"+params.Encode(), res.Code)
		err = fmt.Errorf("recommend region news api response code(%v)", res)
		return
	}
	for _, arcs := range res.List {
		var aidInt int64
		switch aid := arcs.Aid.(type) {
		case string:
			aidInt = aidToInt(aid)
		case float64:
			aidInt = int64(aid)
		}
		arcids = append(arcids, aidInt)
	}
	return
}

// TwoRegionHots
func (d *Dao) RegionChildHots(c context.Context, rid int) (arcids []int64, err error) {
	var res struct {
		Code int `json:"code"`
		List []struct {
			Aid int64 `json:"aid"`
		} `json:"list"`
	}
	api := fmt.Sprintf(d.regionChildHotUrl, rid)
	if err = d.clientAsyn.Get(c, api, "", nil, &res); err != nil {
		log.Error("recommend region child hots url(%s) error(%v)", api, err)
		return
	}
	if res.Code != 0 {
		log.Error("recommend region child hots url(%s) error(%v)", api, res.Code)
		err = fmt.Errorf("recommend region child hots api response code(%v)", res)
		return
	}
	for _, arcs := range res.List {
		arcids = append(arcids, arcs.Aid)
	}
	return
}

func (d *Dao) RegionArcList(c context.Context, rid, pn, ps int, now time.Time) (arcids []int64, err error) {
	params := url.Values{}
	params.Set("rid", strconv.Itoa(rid))
	params.Set("pn", strconv.Itoa(pn))
	params.Set("ps", strconv.Itoa(ps))
	var res struct {
		Code int `json:"code"`
		Data struct {
			List []struct {
				Aid int64 `json:"aid"`
			} `json:"archives"`
		} `json:"data"`
	}
	if err = d.client.Get(c, d.regionArcListUrl, "", params, &res); err != nil {
		log.Error("recommend regionArc news url(%s) error(%v)", d.regionArcListUrl+"?"+params.Encode(), err)
		return
	}
	if res.Code != 0 && res.Code != -1 {
		log.Error("recommend regionArc news url(%s) error(%v)", d.regionArcListUrl+"?"+params.Encode(), res.Code)
		err = fmt.Errorf("recommend regionArc news api response code(%v)", res)
		return
	}
	for _, arcs := range res.Data.List {
		arcids = append(arcids, arcs.Aid)
	}
	return
}

// RegionRank
func (d *Dao) RankRegion(c context.Context, rid int, order string) (data []*recommend.Arc, err error) {
	var res struct {
		Data struct {
			Code int              `json:"code"`
			List []*recommend.Arc `json:"list"`
		} `json:"rank"`
	}
	api := fmt.Sprintf(d.rankRegionUrl, order, rid)
	if err = d.clientAsyn.Get(c, api, "", nil, &res); err != nil {
		log.Error("recommend region rank hots url(%s) error(%v)", api, err)
		return
	}
	if res.Data.Code != 0 {
		log.Error("recommend region rank hots url(%s) error(%v)", api, res.Data.Code)
		err = fmt.Errorf("recommend region rank hots api response code(%v)", res)
		return
	}
	data = res.Data.List
	return
}

// RankAll
func (d *Dao) RankAll(c context.Context, order string) (data []*recommend.Arc, err error) {
	var res struct {
		Data struct {
			Code int              `json:"code"`
			List []*recommend.Arc `json:"list"`
		} `json:"rank"`
	}
	api := fmt.Sprintf(d.rankOriginalUrl, order)
	if err = d.clientAsyn.Get(c, api, "", nil, &res); err != nil {
		log.Error("recommend region rank hots url(%s) error(%v)", api, err)
		return
	}
	if res.Data.Code != 0 {
		log.Error("recommend region rank hots url(%s) error(%v)", api, res.Data.Code)
		err = fmt.Errorf("recommend region rank hots api response code(%v)", res)
		return
	}
	data = res.Data.List
	return
}

// RankAll
func (d *Dao) RankBangumi(c context.Context) (data []*recommend.Arc, err error) {
	var res struct {
		Data struct {
			Code int              `json:"code"`
			List []*recommend.Arc `json:"list"`
		} `json:"rank"`
	}
	if err = d.clientAsyn.Get(c, d.rankBangumilUrl, "", nil, &res); err != nil {
		log.Error("recommend region rank hots url(%s) error(%v)", d.rankBangumilUrl, err)
		return
	}
	if res.Data.Code != 0 {
		log.Error("recommend region rank hots url(%s) error(%v)", d.rankBangumilUrl, res.Data.Code)
		err = fmt.Errorf("recommend region rank hots api response code(%v)", res)
		return
	}
	data = res.Data.List
	return
}

// FeedDynamic
func (d *Dao) FeedDynamic(c context.Context, pull bool, rid, tid int, ctime, mid int64, now time.Time) (hotAids, newAids []int64, ctop, cbottom xtime.Time, err error) {
	var pn string
	if pull {
		pn = "1"
	} else {
		pn = "2"
	}
	params := url.Values{}
	params.Set("src", "2")
	params.Set("pn", pn)
	params.Set("mid", strconv.FormatInt(mid, 10))
	if ctime != 0 {
		params.Set("ctime", strconv.FormatInt(ctime, 10))
	}
	if rid != 0 {
		params.Set("rid", strconv.Itoa(rid))
	}
	if tid != 0 {
		params.Set("tag", strconv.Itoa(tid))
	}
	var res struct {
		Code    int        `json:"code"`
		Data    []int64    `json:"data"`
		Hot     []int64    `json:"hot"`
		CTop    xtime.Time `json:"ctop"`
		CBottom xtime.Time `json:"cbottom"`
	}
	if err = d.client.Get(c, d.feedDynamicUrl, "", params, &res); err != nil {
		log.Error("region feed dynamic d.client.Get(%s) error(%v)", d.feedDynamicUrl+"?"+params.Encode(), err)
		return
	}
	b, _ := json.Marshal(&res)
	log.Info("region feed dynamic url(%s) response(%s)", d.feedDynamicUrl+"?"+params.Encode(), b)
	if res.Code != 0 {
		log.Error("region feed dynamic d.client.Get(%s) error(%v)", d.regionArcListUrl+"?"+params.Encode(), res.Code)
		err = fmt.Errorf("region feed dynamicapi response code(%v)", res)
		return
	}
	hotAids = res.Hot
	newAids = res.Data
	ctop = res.CTop
	cbottom = res.CBottom
	return
}

func (d *Dao) RankAppRegion(c context.Context, rid int) (aids []int64, others, scores map[int64]int64, err error) {
	var res struct {
		Code int `json:"code"`
		List []struct {
			Aid    int64 `json:"aid"`
			Score  int64 `json:"score"`
			Others []struct {
				Aid   int64 `json:"aid"`
				Score int64 `json:"score"`
			} `json:"others"`
		} `json:"list"`
	}
	api := fmt.Sprintf(d.rankRegionAppUrl, rid)
	if err = d.client.Get(c, api, "", nil, &res); err != nil {
		log.Error("recommend region rank hots url(%s) error(%v)", api, err)
		return
	}
	if res.Code != 0 && res.Code != -1 {
		log.Error("recommend region rank hots url(%s) error(%v)", api, res.Code)
		err = fmt.Errorf("recommend region rank hots api response code(%v)", res)
		return
	}
	scores = map[int64]int64{}
	others = map[int64]int64{}
	for _, arcs := range res.List {
		aids = append(aids, arcs.Aid)
		scores[arcs.Aid] = arcs.Score
		for _, o := range arcs.Others {
			aids = append(aids, o.Aid)
			scores[o.Aid] = o.Score
			others[o.Aid] = arcs.Aid
		}
	}
	return
}

func (d *Dao) RankAppOrigin(c context.Context) (aids []int64, others, scores map[int64]int64, err error) {
	var res struct {
		Code int `json:"code"`
		List []struct {
			Aid    int64 `json:"aid"`
			Score  int64 `json:"score"`
			Others []struct {
				Aid   int64 `json:"aid"`
				Score int64 `json:"score"`
			} `json:"others"`
		} `json:"list"`
	}
	if err = d.client.Get(c, d.rankOriginAppUrl, "", nil, &res); err != nil {
		log.Error("recommend Origin rank hots url(%s) error(%v)", d.rankOriginAppUrl, err)
		return
	}
	if res.Code != 0 && res.Code != -1 {
		log.Error("recommend Origin rank hots url(%s) error(%v)", d.rankOriginAppUrl, res.Code)
		err = fmt.Errorf("recommend Origin rank hots api response code(%v)", res)
		return
	}
	scores = map[int64]int64{}
	others = map[int64]int64{}
	for _, arcs := range res.List {
		aids = append(aids, arcs.Aid)
		scores[arcs.Aid] = arcs.Score
		for _, o := range arcs.Others {
			aids = append(aids, o.Aid)
			scores[o.Aid] = o.Score
			others[o.Aid] = arcs.Aid
		}
	}
	return
}

func (d *Dao) RankAppAll(c context.Context) (aids []int64, others, scores map[int64]int64, err error) {
	var res struct {
		Code int `json:"code"`
		List []struct {
			Aid    int64 `json:"aid"`
			Score  int64 `json:"score"`
			Others []struct {
				Aid   int64 `json:"aid"`
				Score int64 `json:"score"`
			} `json:"others"`
		} `json:"list"`
	}
	if err = d.client.Get(c, d.rankAllAppUrl, "", nil, &res); err != nil {
		log.Error("recommend All rank hots url(%s) error(%v)", d.rankAllAppUrl, err)
		return
	}
	if res.Code != 0 && res.Code != -1 {
		log.Error("recommend All rank hots url(%s) error(%v)", d.rankAllAppUrl, res.Code)
		err = fmt.Errorf("recommend All rank hots api response code(%v)", res)
		return
	}
	scores = map[int64]int64{}
	others = map[int64]int64{}
	for _, arcs := range res.List {
		aids = append(aids, arcs.Aid)
		scores[arcs.Aid] = arcs.Score
		for _, o := range arcs.Others {
			aids = append(aids, o.Aid)
			scores[o.Aid] = o.Score
			others[o.Aid] = arcs.Aid
		}
	}
	return
}

func (d *Dao) RankAppBangumi(c context.Context) (aids []int64, others, scores map[int64]int64, err error) {
	var res struct {
		Code int `json:"code"`
		List []struct {
			Aid    int64 `json:"aid"`
			Score  int64 `json:"score"`
			Others []struct {
				Aid   int64 `json:"aid"`
				Score int64 `json:"score"`
			} `json:"others"`
		} `json:"list"`
	}
	if err = d.client.Get(c, d.rankBangumiAppUrl, "", nil, &res); err != nil {
		log.Error("recommend bangumi rank hots url(%s) error(%v)", d.rankBangumiAppUrl, err)
		return
	}
	if res.Code != 0 && res.Code != -1 {
		log.Error("recommend bangumi rank hots url(%s) error(%v)", d.rankBangumiAppUrl, res.Code)
		err = fmt.Errorf("recommend bangumi rank hots api response code(%v)", res)
		return
	}
	scores = map[int64]int64{}
	others = map[int64]int64{}
	for _, arcs := range res.List {
		aids = append(aids, arcs.Aid)
		scores[arcs.Aid] = arcs.Score
		for _, o := range arcs.Others {
			aids = append(aids, o.Aid)
			scores[o.Aid] = o.Score
			others[o.Aid] = arcs.Aid
		}
	}
	return
}

func aidToInt(aidstr string) (aid int64) {
	aid, _ = strconv.ParseInt(aidstr, 10, 64)
	return
}

// HotTab hot tab
func (d *Dao) HotTab(c context.Context) (list []*recommend.List, err error) {
	var res struct {
		Code int               `json:"code"`
		List []*recommend.List `json:"list"`
	}
	if err = d.client.Get(c, d.hottabURL, "", nil, &res); err != nil {
		log.Error("hottab hots url(%s) error(%v)", d.hottabURL, err)
		return
	}
	b, _ := json.Marshal(&res)
	log.Info("hottab list url(%s) response(%s)", d.hottabURL, b)
	if res.Code != 0 {
		err = ecode.Int(res.Code)
		log.Error("hottab hots url(%s) code(%d)", d.hottabURL, res.Code)
		return
	}
	list = res.List
	return
}

// HotTenTab hot tab
func (d *Dao) HotTenTab(c context.Context, i int) (list []*recommend.List, err error) {
	var res struct {
		Code int               `json:"code"`
		List []*recommend.List `json:"list"`
	}
	if err = d.client.Get(c, fmt.Sprintf(d.hotHetongURL, i), "", nil, &res); err != nil {
		err = errors.Wrap(err, fmt.Sprintf(d.hotHetongURL, i))
		return
	}
	if res.Code != 0 {
		err = errors.Wrap(err, fmt.Sprintf("code(%d)", res.Code))
		return
	}
	list = res.List
	return
}

// HotHeTongTabCard hot tab card
func (d *Dao) HotHeTongTabCard(c context.Context, i int) (list []*recommend.CardList, err error) {
	var res struct {
		Code int                   `json:"code"`
		List []*recommend.CardList `json:"list"`
	}
	if err = d.client.Get(c, fmt.Sprintf(d.hotHeTongtabcardURL, i), "", nil, &res); err != nil {
		err = errors.Wrap(err, fmt.Sprintf(d.hotHeTongtabcardURL, i))
		return
	}
	if res.Code != 0 {
		err = errors.Wrap(err, fmt.Sprintf("code(%d)", res.Code))
		return
	}
	list = res.List
	return
}
