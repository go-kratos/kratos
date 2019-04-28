package tag

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"go-common/app/interface/main/app-tag/model/tag"
	"go-common/library/log"
)

const (
	_detail              = "/x/internal/tag/detail"
	_tagHotsIDURL        = "/x/internal/tag/hots"
	_similarTagChangeURL = "/x/internal/tag/change/similar"
	_tagURL              = "/x/internal/tag/info"
	_tagRankingURL       = "/x/internal/tag/detail/ranking"
)

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

func (d *Dao) Detail(c context.Context, tagID int64, pn, ps int, now time.Time) (arcids []int64, err error) {
	params := url.Values{}
	params.Set("tag_id", strconv.FormatInt(tagID, 10))
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
	if err = d.client.Get(c, d.detailURL, "", params, &res); err != nil {
		log.Error("d.client.Get(%s) error(%v)", d.detailURL+"?"+params.Encode(), err)
		return
	}
	if res.Code != 0 {
		log.Error("DetailTag url(%s) code:%d", d.detailURL+"?"+params.Encode(), res.Code)
		err = fmt.Errorf("DetailTag code error(%d)", res.Code)
		return
	}
	for _, arcs := range res.Data.News.Archives {
		arcids = append(arcids, arcs.Aid)
	}
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
func (d *Dao) SimilarTagChange(c context.Context, tagID int64, now time.Time) (data []*tag.SimilarTag, err error) {
	params := url.Values{}
	params.Set("tag_id", strconv.FormatInt(tagID, 10))
	var res struct {
		Code int               `json:"code"`
		Data []*tag.SimilarTag `json:"data"`
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

// DetailRanking tag detail ranking
func (d *Dao) DetailRanking(c context.Context, reid int, tagID int64, pn, ps int, now time.Time) (arcids []int64, err error) {
	params := url.Values{}
	params.Set("prid", strconv.Itoa(reid))
	params.Set("tag_id", strconv.FormatInt(tagID, 10))
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
