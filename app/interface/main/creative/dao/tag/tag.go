package tag

import (
	"context"
	"net/url"
	"strconv"

	"go-common/app/interface/main/creative/model/tag"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/xstr"
)

// TagList get tag list.
func (d *Dao) TagList(c context.Context, ids []int64) (tgs []*tag.Meta, err error) {
	params := url.Values{}
	params.Set("tag_id", xstr.JoinInts(ids))
	var res struct {
		Code int         `json:"code"`
		Data []*tag.Meta `json:"data"`
	}
	if err = d.client.Get(c, d.tagList, "", params, &res); err != nil {
		log.Error("TagList url(%s) response(%v) error(%v)", d.tagList, res, err)
		err = ecode.CreativeTagErr
		return
	}
	if res.Code != 0 {
		log.Error("TagList url(%s) res(%v)", d.tagList, res)
		err = ecode.CreativeTagErr
		return
	}
	tgs = res.Data
	return
}

// TagCheck tag check
func (d *Dao) TagCheck(c context.Context, mid int64, tagName string) (t *tag.Tag, err error) {
	params := url.Values{}
	params.Set("tag_name", tagName)
	params.Set("mid", strconv.FormatInt(mid, 10))
	var res struct {
		Code int      `json:"code"`
		Data *tag.Tag `json:"data"`
	}
	if err = d.client.Get(c, d.tagCheck, "", params, &res); err != nil {
		log.Error("TagCheck url(%s) p(%+v)  response(%v) error(%v)", d.tagCheck, params.Encode(), res, err)
		err = ecode.CreativeTagErr
		return
	}
	log.Info("TagCheck url(%s) | p(%+v) |res(%v)", d.tagCheck, params.Encode(), res)
	if res.Code != 0 {
		log.Error("TagCheck url(%s) res(%v)", d.tagCheck, res)
		err = ecode.Int(res.Code)
		return
	}
	t = res.Data
	return
}

// AppealTag appeal tag from videoup.
func (d *Dao) AppealTag(c context.Context, aid int64, ip string) (tid int64, err error) {
	params := url.Values{}
	params.Set("aid", strconv.FormatInt(aid, 10))
	var res struct {
		Code int `json:"code"`
		Data *struct {
			TagID int64 `json:"tag_id"`
		} `json:"data"`
	}
	if err = d.client.Get(c, d.appealTag, ip, params, &res); err != nil {
		log.Error("appeal tag error(%v)", err)
		err = ecode.CreativeTagErr
		return
	}
	if res.Code != 0 {
		log.Error("appeal tag url(%s) res(%v)", d.appealTag, res)
		err = ecode.Int(res.Code)
		return
	}
	if res.Data != nil {
		tid = res.Data.TagID
	}
	return
}

// StaffTitleList 获取联合投稿职能列表
func (d *Dao) StaffTitleList(c context.Context) (staffTitles []*tag.StaffTitle, err error) {
	var res struct {
		Code int `json:"code"`
		Data *struct {
			Data []*tag.StaffTitle `json:"data"`
		} `json:"data"`
	}
	staffTitles = make([]*tag.StaffTitle, 0)
	params := url.Values{}
	params.Set("bid", strconv.FormatInt(StaffTagBid, 10))
	params.Set("state", "1")
	params.Set("ps", "50")
	if err = d.client.Get(c, d.mngTagListURI, "", params, &res); err != nil {
		log.Error("StaffTitleList error(%v)", err)
		err = ecode.CreativeTagErr
		return
	}
	if res.Code != 0 {
		log.Error("StaffTitleList url(%s) res(%v)", d.appealTag, res)
		err = ecode.Int(res.Code)
		return
	}
	if res.Data != nil && res.Data.Data != nil {
		staffTitles = res.Data.Data
	}
	return
}
