package tag

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"go-common/app/interface/main/app-show/conf"
	"go-common/app/interface/main/app-show/model/region"
	"go-common/app/interface/main/app-show/model/tag"
	tagm "go-common/app/interface/main/tag/model"
	"go-common/library/ecode"
	"go-common/library/log"
	httpx "go-common/library/net/http/blademaster"
)

const (
	_tagURL              = "/x/internal/tag/info"
	_tagHotURL           = "/tag/hot/%d/%d.json"
	_tagNewURL           = "/x/internal/tag/ranking/archives"
	_similarTagURL       = "/x/internal/tag/similar"
	_tagHotsIDURL        = "/x/internal/tag/hots"
	_similarTagChangeURL = "/x/internal/tag/change/similar"
	_tagDetailURL        = "/x/internal/tag/detail"
	_tagRankingURL       = "/x/internal/tag/detail/ranking"
	_tagArchiveURL       = "/x/internal/tag/archive/tags"
)

// Dao is tag dao.
type Dao struct {
	client      *httpx.Client
	clientParam *httpx.Client
	// url
	tagURL              string
	tagHotURL           string
	tagNewURL           string
	similarTagURL       string
	tagHotsIDURL        string
	similarTagChangeURL string
	tagDetailURL        string
	tagRankingURL       string
	tagArchiveURL       string
}

// New tag dao.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		client:              httpx.NewClient(conf.Conf.HTTPClient),
		clientParam:         httpx.NewClient(conf.Conf.HTTPClient),
		tagURL:              c.Host.ApiCo + _tagURL,
		tagHotURL:           c.Host.Hetongzi + _tagHotURL,
		tagNewURL:           c.Host.ApiCo + _tagNewURL,
		similarTagURL:       c.Host.ApiCo + _similarTagURL,
		tagHotsIDURL:        c.Host.ApiCo + _tagHotsIDURL,
		similarTagChangeURL: c.Host.ApiCo + _similarTagChangeURL,
		tagDetailURL:        c.Host.ApiCo + _tagDetailURL,
		tagRankingURL:       c.Host.ApiCo + _tagRankingURL,
		tagArchiveURL:       c.Host.ApiCo + _tagArchiveURL,
	}
	return
}

// TagInfo get tag info.
func (d *Dao) TagInfo(c context.Context, mid int64, tagId int, now time.Time) (data *tag.Tag, err error) {
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("tag_id", strconv.Itoa(tagId))
	var res struct {
		Code int      `json:"code"`
		Data *tag.Tag `json:"data"`
	}
	if err = d.client.Get(c, d.tagURL, "", params, &res); err != nil {
		log.Error("tagInfo url(%s) error(%v)", d.tagURL+"?"+params.Encode(), err)
		return
	}
	if res.Code != 0 {
		log.Error("tagInfo url(%s) error(%v)", d.tagURL+"?"+params.Encode(), res.Code)
		err = fmt.Errorf("tagInfo api response code(%v)", res)
		return
	}
	data = res.Data
	return
}

// Hots
func (d *Dao) Hots(c context.Context, rid, tagId, pn, ps int, now time.Time) (data []int64, err error) {
	var (
		uri = fmt.Sprintf(d.tagHotURL, rid, tagId)
		res struct {
			Code int     `json:"code"`
			Data []int64 `json:"data"`
		}
		count int
		start int
	)
	if err = d.clientParam.Get(c, uri, "", nil, &res); err != nil {
		log.Error("d.paramclient.Get(%s) error(%v)", uri, err)
		return
	}
	if res.Code != 0 {
		log.Error("tag region hots url(%s) code:%d", uri, res.Code)
		err = fmt.Errorf("tag region hots api response code(%v)", res)
		return
	}
	count = len(res.Data)
	if count == 0 {
		return
	}
	start = (pn - 1) * ps
	if start > count {
		return
	}
	data = res.Data
	return
}

// NewArcs
func (d *Dao) NewArcs(c context.Context, rid, tagId, pn, ps int, now time.Time) (arcids []int64, err error) {
	params := url.Values{}
	params.Set("rid", strconv.Itoa(rid))
	params.Set("tag_id", strconv.Itoa(tagId))
	params.Set("pn", strconv.Itoa(pn))
	params.Set("ps", strconv.Itoa(ps))
	var res struct {
		Code int `json:"code"`
		Data struct {
			Archives []struct {
				Aid int64 `json:"aid"`
			} `json:"archives"`
		} `json:"data"`
	}
	if err = d.client.Get(c, d.tagNewURL, "", params, &res); err != nil {
		log.Error("tag region news url(%s) error(%v)", d.tagNewURL+"?"+params.Encode(), err)
		return
	}
	if res.Code != 0 {
		log.Error("tag region news url(%s) error(%v)", d.tagNewURL+"?"+params.Encode(), res.Code)
		err = fmt.Errorf("tag region news api response code(%v)", res)
		return
	}
	for _, arcs := range res.Data.Archives {
		arcids = append(arcids, arcs.Aid)
	}
	return
}

// SimilarTag tag similar
func (d *Dao) SimilarTag(c context.Context, rid, tagId int, now time.Time) (data []*region.SimilarTag, err error) {
	params := url.Values{}
	params.Set("rid", strconv.Itoa(rid))
	params.Set("tid", strconv.Itoa(tagId))
	var res struct {
		Code int                  `json:"code"`
		Data []*region.SimilarTag `json:"data"`
	}
	if err = d.client.Get(c, d.similarTagURL, "", params, &res); err != nil {
		log.Error("tag similarTagURL url(%s) error(%v)", d.similarTagURL, err)
		return
	}
	if res.Code != 0 {
		log.Error("tag similarTagURL url(%s) Code(%v)", d.similarTagURL, res.Code)
		err = fmt.Errorf("tag api response code(%v)", res)
		return
	}
	data = res.Data
	return
}

// TagHotsId
func (d *Dao) TagHotsId(c context.Context, rid int, now time.Time) (tags []*tag.Tag, err error) {
	params := url.Values{}
	params.Set("rid", strconv.Itoa(rid))
	var res struct {
		Code int `json:"code"`
		Data []struct {
			Rid  int        `json:"rid"`
			Tags []*tag.Tag `json:"tags"`
		} `json:"data"`
	}
	if err = d.client.Get(c, d.tagHotsIDURL, "", params, &res); err != nil {
		log.Error("tag tagHotsIDURL url(%s) error(%v)", d.tagHotsIDURL, err)
		return
	}
	if res.Code != 0 {
		log.Error("tag tagHotsIDURL url(%s) Code(%v)", d.tagHotsIDURL, res.Code)
		err = fmt.Errorf("tag api response code(%v)", res)
		return
	}
	if len(res.Data) == 0 {
		return
	}
	tags = res.Data[0].Tags
	return
}

// SimilarTagChangetag tag similar no rid
func (d *Dao) SimilarTagChange(c context.Context, tagID int, now time.Time) (data []*region.SimilarTag, err error) {
	params := url.Values{}
	params.Set("tag_id", strconv.Itoa(tagID))
	var res struct {
		Code int                  `json:"code"`
		Data []*region.SimilarTag `json:"data"`
	}
	if err = d.client.Get(c, d.similarTagChangeURL, "", params, &res); err != nil {
		log.Error("tag similarTagChangeURL, url(%s) error(%v)", d.similarTagChangeURL+"?"+params.Encode(), err)
		return
	}
	if res.Code != 0 {
		log.Error("tag similarTagChangeURL url(%s) Code(%v)", d.similarTagChangeURL+"?"+params.Encode(), res.Code)
		err = fmt.Errorf("tag api response code(%v)", res)
		return
	}
	data = res.Data
	return
}

// Detail tag detail
func (d *Dao) Detail(c context.Context, tagID int, pn, ps int, now time.Time) (arcids []int64, err error) {
	params := url.Values{}
	params.Set("tag_id", strconv.Itoa(tagID))
	params.Set("pn", strconv.Itoa(pn))
	params.Set("ps", strconv.Itoa(ps))
	var res struct {
		Code int `json:"code"`
		Data struct {
			News struct {
				Archives []struct {
					Aid int64 `json:"aid"`
				} `json:"archives"`
			} `json:"news"`
		} `json:"data"`
	}
	if err = d.client.Get(c, d.tagDetailURL, "", params, &res); err != nil {
		log.Error("d.client.Get(%s) error(%v)", d.tagDetailURL+"?"+params.Encode(), err)
		return
	}
	if res.Code != 0 {
		log.Error("tag Detail url(%s) Code(%v)", d.tagDetailURL, res.Code)
		err = fmt.Errorf("tag api response code(%v)", res)
		return
	}
	for _, arcs := range res.Data.News.Archives {
		arcids = append(arcids, arcs.Aid)
	}
	return
}

// DetailRanking tag detail ranking
func (d *Dao) DetailRanking(c context.Context, reid, tagID int, pn, ps int, now time.Time) (arcids []int64, err error) {
	params := url.Values{}
	params.Set("prid", strconv.Itoa(reid))
	params.Set("tag_id", strconv.Itoa(tagID))
	params.Set("pn", strconv.Itoa(pn))
	params.Set("ps", strconv.Itoa(ps))
	var res struct {
		Code int `json:"code"`
		Data struct {
			Archives []struct {
				Aid int64 `json:"aid"`
			} `json:"archives"`
		} `json:"data"`
	}
	if err = d.client.Get(c, d.tagRankingURL, "", params, &res); err != nil {
		log.Error("d.client.Get(%s) error(%v)", d.tagRankingURL+"?"+params.Encode(), err)
		return
	}
	if res.Code != 0 {
		log.Error("tag DetailRanking url(%s) Code(%v)", d.tagRankingURL+"?"+params.Encode(), res.Code)
		err = fmt.Errorf("tag api response code(%v)", res)
		return
	}
	for _, arcs := range res.Data.Archives {
		arcids = append(arcids, arcs.Aid)
	}
	return
}

// TagArchive archive tags
func (d *Dao) TagArchive(c context.Context, aid int64) (data []*tagm.Tag, err error) {
	params := url.Values{}
	params.Set("aid", strconv.FormatInt(aid, 10))
	var res struct {
		Code int         `json:"code"`
		Data []*tagm.Tag `json:"data"`
	}
	if err = d.client.Get(c, d.tagArchiveURL, "", params, &res); err != nil {
		log.Error("TagArchive url(%s) error(%v)", d.tagArchiveURL+"?"+params.Encode(), err)
		return
	}
	if res.Code != 0 {
		log.Error("TagArchive url(%s) error(%v)", d.tagArchiveURL+"?"+params.Encode(), res.Code)
		err = ecode.Int(res.Code)
		return
	}
	data = res.Data
	return
}
