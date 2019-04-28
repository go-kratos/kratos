package dao

import (
	"context"
	"fmt"
	"net/url"

	"go-common/app/interface/main/web/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/metadata"
)

const (
	_rankURI          = "%s.json"
	_rankAllURI       = "all-%d"
	_rankAllRidURI    = "all_region-%d-%d"
	_rankAllRecURI    = "recent_all-%d"
	_rankAllRecRidURI = "recent_region-%d-%d"
	_rankOriAllURI    = "all_origin-%d"
	_rankOriAllRidURI = "all_region_origin-%d-%d"
	_rankOriRecURI    = "recent_origin-%d"
	_rankOriRecRidURI = "recent_region_origin-%d-%d"
	_rankAllNewURI    = "all_rookie-%d"
	_rankAllNewRidURI = "all_region_rookie-%d-%d"
	_rankRegionURI    = "recent_region%s-%d-%d.json"
	_rankRecURI       = "reco_region-%d.json"
	_rankTagURI       = "/tag/hot/web/%d/%d.json"
	_rankIndexURI     = "reco-%d.json"
	_customURI        = "game_custom_2.json"
)

// Ranking get ranking data from new api
func (d *Dao) Ranking(c context.Context, rid int16, rankType, day, arcType int) (res *model.RankNew, err error) {
	var (
		params = url.Values{}
		ip     = metadata.String(c, metadata.RemoteIP)
	)
	suffix := rankURI(rid, model.RankType[rankType], day, arcType)
	var rs struct {
		Code int                     `json:"code"`
		Note string                  `json:"note"`
		List []*model.RankNewArchive `json:"list"`
	}
	if err = d.httpBigData.RESTfulGet(c, d.rankURL, ip, params, &rs, suffix); err != nil {
		log.Error("d.httpBigData.RESTfulGet(%s) error(%v)", suffix, err)
		return
	}
	if rs.Code != ecode.OK.Code() {
		log.Error("d.httpBigData.RESTfulGet(%s) error code(%d)", suffix, rs.Code)
		err = ecode.Int(rs.Code)
		return
	}
	res = &model.RankNew{Note: rs.Note, List: rs.List}
	return
}

// RankingIndex get rank index data from bigdata
func (d *Dao) RankingIndex(c context.Context, day int) (res []*model.NewArchive, err error) {
	var (
		params   = url.Values{}
		remoteIP = metadata.String(c, metadata.RemoteIP)
	)
	var rs struct {
		Code int                 `json:"code"`
		Num  int                 `json:"num"`
		List []*model.NewArchive `json:"list"`
	}
	if err = d.httpBigData.RESTfulGet(c, d.rankIndexURL, remoteIP, params, &rs, day); err != nil {
		log.Error("d.httpBigData.RESTfulGet(%d) error(%v)", day, err)
		return
	}
	if rs.Code != ecode.OK.Code() {
		log.Error("d.httpBigData.RESTfulGet(%d) error(%v)", day, rs.Code)
		err = ecode.Int(rs.Code)
		return
	}
	res = rs.List
	return
}

// RankingRegion get rank region data from bigdata
func (d *Dao) RankingRegion(c context.Context, rid int16, day, original int) (res []*model.NewArchive, err error) {
	var (
		params   = url.Values{}
		remoteIP = metadata.String(c, metadata.RemoteIP)
	)
	var rs struct {
		Code int                 `json:"code"`
		List []*model.NewArchive `json:"list"`
	}
	if err = d.httpBigData.RESTfulGet(c, d.rankRegionURL, remoteIP, params, &rs, model.OriType[original], rid, day); err != nil {
		log.Error("d.httpBigData.RESTfulGet(%d,%d,%d) error(%v)", original, rid, day, err)
		return
	}
	if rs.Code != ecode.OK.Code() {
		log.Error("d.httpBigData.RESTfulGet(%d,%d,%d) error code(%d)", original, rid, day, rs.Code)
		err = ecode.Int(rs.Code)
		return
	}
	res = rs.List
	return
}

// RankingRecommend get rank recommend data from bigdata.
func (d *Dao) RankingRecommend(c context.Context, rid int16) (res []*model.NewArchive, err error) {
	var (
		params   = url.Values{}
		remoteIP = metadata.String(c, metadata.RemoteIP)
	)
	var rs struct {
		Code int                 `json:"code"`
		Num  int                 `json:"num"`
		List []*model.NewArchive `json:"list"`
	}
	if err = d.httpBigData.RESTfulGet(c, d.rankRecURL, remoteIP, params, &rs, rid); err != nil {
		log.Error("d.httpBigData.RESTfulGet(%d) error(%v)", rid, err)
		return
	}
	if rs.Code != ecode.OK.Code() {
		log.Error("d.httpBigData.RESTfulGet(%d) error(%v)", rid, rs.Code)
		err = ecode.Int(rs.Code)
		return
	}
	res = rs.List
	return
}

// RankingTag get rank tag data from bigdata
func (d *Dao) RankingTag(c context.Context, rid int16, tagID int64) (rs []*model.NewArchive, err error) {
	var (
		params   = url.Values{}
		remoteIP = metadata.String(c, metadata.RemoteIP)
	)
	var res struct {
		Code int                 `json:"code"`
		List []*model.NewArchive `json:"list"`
	}
	if err = d.httpBigData.RESTfulGet(c, d.rankTagURL, remoteIP, params, &res, rid, tagID); err != nil {
		log.Error("d.httpBigData.RESTfulGet(%d,%d) error(%v)", rid, tagID, err)
		return
	}
	if res.Code != ecode.OK.Code() {
		log.Error("d.httpBigData.RESTfulGet(%d,%d) code(%d)", rid, tagID, res.Code)
		err = ecode.Int(res.Code)
		return
	}
	rs = res.List
	return
}

// RegionCustom get region(game) custom data from big data
func (d *Dao) RegionCustom(c context.Context) (rs []*model.Custom, err error) {
	var (
		params = url.Values{}
		ip     = metadata.String(c, metadata.RemoteIP)
	)
	var res struct {
		Code int
		List []*model.Custom
	}
	if err = d.httpBigData.Get(c, d.customURL, ip, params, &res); err != nil {
		log.Error("d.httpBigData.Get(%s) error(%v)", d.customURL, err)
		return
	}
	if res.Code != ecode.OK.Code() {
		log.Error("d.httpBigData.Get(%s) error(%v)", d.customURL, err)
		err = ecode.Int(res.Code)
		return
	}
	rs = res.List
	return
}

func rankURI(rid int16, rankType string, day, arcType int) string {
	if rankType == model.RankType[1] {
		if arcType == 1 {
			if rid > 0 {
				return fmt.Sprintf(_rankAllRecRidURI, rid, day)
			}
			return fmt.Sprintf(_rankAllRecURI, day)
		}
		if rid > 0 {
			return fmt.Sprintf(_rankAllRidURI, rid, day)
		}
		return fmt.Sprintf(_rankAllURI, day)
	} else if rankType == model.RankType[2] {
		if arcType == 1 {
			if rid > 0 {
				return fmt.Sprintf(_rankOriRecRidURI, rid, day)
			}
			return fmt.Sprintf(_rankOriRecURI, day)
		}
		if rid > 0 {
			return fmt.Sprintf(_rankOriAllRidURI, rid, day)
		}
		return fmt.Sprintf(_rankOriAllURI, day)
	}
	if rid > 0 {
		return fmt.Sprintf(_rankAllNewRidURI, rid, day)
	}
	return fmt.Sprintf(_rankAllNewURI, day)
}
