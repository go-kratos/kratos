package datadao

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"go-common/app/interface/main/mcn/dao/cache"
	"go-common/app/interface/main/mcn/dao/global"
	"go-common/app/interface/main/mcn/model/datamodel"
	"go-common/app/interface/main/mcn/model/mcnmodel"
	"go-common/app/interface/main/mcn/tool/datacenter"
	tagmdl "go-common/app/interface/main/tag/api"
	"go-common/library/log"
	"go-common/library/sync/errgroup"
)

// const url for api
const (
	APIMcnSummary          = "http://berserker.bilibili.co/avenger/api/155/query" // 7   see doc http://info.bilibili.co/pages/viewpage.action?pageId=11545690
	APIIndexInc            = "http://berserker.bilibili.co/avenger/api/156/query" // 3.1
	APIIndexSource         = "http://berserker.bilibili.co/avenger/api/159/query" // 3.2
	APIPlaySource          = "http://berserker.bilibili.co/avenger/api/161/query" // 3.3
	APIMcnFans             = "http://berserker.bilibili.co/avenger/api/168/query" // 3.4
	APIMcnFansInc          = "http://berserker.bilibili.co/avenger/api/171/query" // 3.5
	APIMcnFansDec          = "http://berserker.bilibili.co/avenger/api/169/query" // 3.6
	APIMcnFansAttentionWay = "http://berserker.bilibili.co/avenger/api/170/query" // 3.7
	APIMcnFansSex          = "http://berserker.bilibili.co/avenger/api/162/query" // 3.8
	APIMcnFansAge          = "http://berserker.bilibili.co/avenger/api/163/query" // 3.8
	APIMcnFansPlayWay      = "http://berserker.bilibili.co/avenger/api/164/query" // 3.8
	APIMcnFansArea         = "http://berserker.bilibili.co/avenger/api/165/query" // 3.9
	APIMcnFansType         = "http://berserker.bilibili.co/avenger/api/166/query" // 3.10
	APIMcnFansTag          = "http://berserker.bilibili.co/avenger/api/167/query" // 3.11
)

func (d *Dao) callDataAPI(c context.Context, api string, query *datacenter.Query, res interface{}) (err error) {
	var response = &datacenter.Response{
		Result: res,
	}
	if query.Error() != nil {
		err = query.Error()
		log.Error("query error, err=%s", err)
		return
	}
	var params = url.Values{}
	params.Add("query", query.String())

	if err = d.Client.Get(c, api, params, response); err != nil {
		log.Error("fail to get response, err=%+v", err)
		return
	}

	if response.Code != http.StatusOK {
		err = fmt.Errorf("code:%d, msg:%s", response.Code, response.Msg)
		return
	}
	return
}

// GetMcnSummary 7
// see doc http://info.bilibili.co/pages/viewpage.action?pageId=11545690#id-对外接口文档-7.mcn获取概要数据
func (d *Dao) GetMcnSummary(c context.Context, signID int64, date time.Time) (res *mcnmodel.McnGetDataSummaryReply, err error) {
	res = new(mcnmodel.McnGetDataSummaryReply)
	var tmp []*datamodel.DmConMcnArchiveD
	var q = &datacenter.Query{}
	q.Select("*").Where(
		datacenter.ConditionMapType{"log_date": datacenter.ConditionLte(date)},
		datacenter.ConditionMapType{"sign_id": datacenter.ConditionIn(signID)},
	).Limit(1, 0).Order("log_date desc")
	var api = APIMcnSummary
	if err = d.callDataAPI(c, api, q, &tmp); err != nil {
		log.Error("call data api fail, api=%s, err=%s", api, err)
		return
	}
	if len(tmp) > 0 {
		res.CopyFromDmConMcnArchiveD(tmp[0])
	}
	//log.Info("%s query arg(%d,%+v) res(%+v)", api, signID, date, tmp[0])
	return
}

//GetMcnSummaryCache GetMcnSummary with cache
func (d *Dao) GetMcnSummaryCache(c context.Context, signID int64, date time.Time) (res *mcnmodel.McnGetDataSummaryReply, err error) {
	res = new(mcnmodel.McnGetDataSummaryReply)
	var cache = NewCacheMcnDataSignID(signID, date, res, "McnGetDataSummaryReply", func(c context.Context, signID int64, date time.Time) (res interface{}, err error) {
		return d.GetMcnSummary(c, signID, date)
	})
	if err = d.McWrapper.GetOrLoad(c, cache); err != nil {
		log.Error("cache get err, err=%v", err)
		return
	}

	return
}

// GetIndexInc 3.1
// see doc http://info.bilibili.co/pages/viewpage.action?pageId=11545690#id-对外接口文档-3.1.查询MCN增量趋势
func (d *Dao) GetIndexInc(c context.Context, signID int64, date time.Time, tp string) (res *mcnmodel.McnGetIndexIncReply, err error) {
	res = new(mcnmodel.McnGetIndexIncReply)
	var q = &datacenter.Query{}
	q.Select("*").Where(
		datacenter.ConditionMapType{"log_date": datacenter.ConditionLte(date)},
		datacenter.ConditionMapType{"sign_id": datacenter.ConditionIn(signID)},
		datacenter.ConditionMapType{"type": datacenter.ConditionIn(tp)},
	).Limit(30, 0).Order("log_date desc")
	var api = APIIndexInc
	if err = d.callDataAPI(c, api, q, &res.Result); err != nil {
		log.Error("call data api fail, api=%s, err=%s", api, err)
		return
	}
	//log.Info("%s query arg(%d,%+v,%s) res(%+v)", api, signID, date, tp, res.Result[0])
	return
}

//GetIndexIncCache GetIndexInc with cache
func (d *Dao) GetIndexIncCache(c context.Context, signID int64, date time.Time, tp string) (res *mcnmodel.McnGetIndexIncReply, err error) {
	res = new(mcnmodel.McnGetIndexIncReply)
	var cache = NewCacheMcnDataWithTp(signID, date, tp, res, "McnGetIndexIncReply", func(c context.Context, signID int64, date time.Time, tp string) (res interface{}, err error) {
		return d.GetIndexInc(c, signID, date, tp)
	})
	if err = d.McWrapper.GetOrLoad(c, cache); err != nil {
		log.Error("cache get err, err=%v", err)
		return
	}

	return
}

// GetIndexSource 3.2
// see doc http://info.bilibili.co/pages/viewpage.action?pageId=11545690#id-对外接口文档-3.2.查询MCN下播放稿件来源所在分区
func (d *Dao) GetIndexSource(c context.Context, signID int64, date time.Time, tp string) (res *mcnmodel.McnGetIndexSourceReply, err error) {
	res = new(mcnmodel.McnGetIndexSourceReply)

	var q = &datacenter.Query{}
	q.Select("*").Where(
		datacenter.ConditionMapType{"log_date": datacenter.ConditionIn(date)},
		datacenter.ConditionMapType{"sign_id": datacenter.ConditionIn(signID)},
		datacenter.ConditionMapType{"type": datacenter.ConditionIn(tp)},
	)
	var api = APIIndexSource
	if err = d.callDataAPI(c, api, q, &res.Result); err != nil {
		log.Error("call data api fail, api=%s, err=%s", api, err)
		return
	}

	var tids []int64
	for _, v := range res.Result {
		tids = append(tids, v.TypeID)
	}

	tpNames := cache.GetTidNames(tids)
	for _, v := range res.Result {
		if tpName, ok := tpNames[v.TypeID]; ok {
			v.TypeName = tpName
		}
	}
	//log.Info("%s query arg(%d,%+v,%s) res(%+v)", api, signID, date, tp, res.Result[0])
	return
}

//GetIndexSourceCache GetIndexSource with cache
func (d *Dao) GetIndexSourceCache(c context.Context, signID int64, date time.Time, tp string) (res *mcnmodel.McnGetIndexSourceReply, err error) {
	res = new(mcnmodel.McnGetIndexSourceReply)
	var cache = NewCacheMcnDataWithTp(signID, date, tp, res, "McnGetIndexSourceReply", func(c context.Context, signID int64, date time.Time, tp string) (res interface{}, err error) {
		return d.GetIndexSource(c, signID, date, tp)
	})
	if err = d.McWrapper.GetOrLoad(c, cache); err != nil {
		log.Error("cache get err, err=%v", err)
		return
	}

	return
}

// GetPlaySource 3.3
// see doc http://info.bilibili.co/pages/viewpage.action?pageId=11545690#id-对外接口文档-3.3.查询MCN播放设备占比
func (d *Dao) GetPlaySource(c context.Context, signID int64, date time.Time) (res *mcnmodel.McnGetPlaySourceReply, err error) {
	res = new(mcnmodel.McnGetPlaySourceReply)

	var q = &datacenter.Query{}
	q.Select("*").Where(
		datacenter.ConditionMapType{"log_date": datacenter.ConditionLte(date)},
		datacenter.ConditionMapType{"sign_id": datacenter.ConditionIn(signID)},
	).Limit(1, 0).Order("log_date desc")
	var api = APIPlaySource
	var tmp []*mcnmodel.McnGetPlaySourceReply
	if err = d.callDataAPI(c, api, q, &tmp); err != nil {
		log.Error("call data api fail, api=%s, err=%s", api, err)
		return
	}
	if len(tmp) > 0 {
		res = tmp[0]
	}
	//log.Info("%s query arg(%d,%+v) res(%+v)", api, signID, date, tmp)
	return
}

//GetPlaySourceCache GetPlaySource with cache
func (d *Dao) GetPlaySourceCache(c context.Context, signID int64, date time.Time) (res *mcnmodel.McnGetPlaySourceReply, err error) {
	res = new(mcnmodel.McnGetPlaySourceReply)
	var cache = NewCacheMcnDataSignID(signID, date, res, "McnGetPlaySourceReply", func(c context.Context, signID int64, date time.Time) (res interface{}, err error) {
		return d.GetPlaySource(c, signID, date)
	})
	if err = d.McWrapper.GetOrLoad(c, cache); err != nil {
		log.Error("cache get err, err=%v", err)
		return
	}

	return
}

// GetMcnFans 3.4
// see doc http://info.bilibili.co/pages/viewpage.action?pageId=11545690#id-对外接口文档-3.4.查询MCN粉丝数与活跃度
func (d *Dao) GetMcnFans(c context.Context, signID int64, date time.Time) (res *mcnmodel.McnGetMcnFansReply, err error) {
	res = new(mcnmodel.McnGetMcnFansReply)

	var q = &datacenter.Query{}
	q.Select("*").Where(
		datacenter.ConditionMapType{"log_date": datacenter.ConditionLte(date)},
		datacenter.ConditionMapType{"sign_id": datacenter.ConditionIn(signID)},
	).Limit(1, 0).Order("log_date desc")
	var tmp []*mcnmodel.McnGetMcnFansReply
	var api = APIMcnFans
	if err = d.callDataAPI(c, api, q, &tmp); err != nil {
		log.Error("call data api fail, api=%s, err=%s", api, err)
		return
	}
	if len(tmp) > 0 {
		res = tmp[0]
	}
	//log.Info("%s query arg(%d,%+v) res(%+v)", api, signID, date, tmp[0])
	return
}

//GetMcnFansCache GetMcnFans with cache
func (d *Dao) GetMcnFansCache(c context.Context, signID int64, date time.Time) (res *mcnmodel.McnGetMcnFansReply, err error) {
	res = new(mcnmodel.McnGetMcnFansReply)
	var cache = NewCacheMcnDataSignID(signID, date, res, "McnGetMcnFansReply", func(c context.Context, signID int64, date time.Time) (res interface{}, err error) {
		return d.GetMcnFans(c, signID, date)
	})
	if err = d.McWrapper.GetOrLoad(c, cache); err != nil {
		log.Error("cache get err, err=%v", err)
		return
	}

	return
}

// GetMcnFansInc 3.5
// see doc http://info.bilibili.co/pages/viewpage.action?pageId=11545690#id-对外接口文档-3.5.查询MCN粉丝按天增量
func (d *Dao) GetMcnFansInc(c context.Context, signID int64, date time.Time) (res *mcnmodel.McnGetMcnFansIncReply, err error) {
	res = new(mcnmodel.McnGetMcnFansIncReply)

	var q = &datacenter.Query{}
	q.Select("*").Where(
		datacenter.ConditionMapType{"log_date": datacenter.ConditionLte(date)},
		datacenter.ConditionMapType{"sign_id": datacenter.ConditionIn(signID)},
	).Limit(30, 0).Order("log_date desc")

	var api = APIMcnFansInc
	if err = d.callDataAPI(c, api, q, &res.Result); err != nil {
		log.Error("call data api fail, api=%s, err=%s", api, err)
		return
	}
	//log.Info("%s query arg(%d,%+v) res(%+v)", api, signID, date, res.Result[0])
	return
}

//GetMcnFansIncCache GetMcnFansInc with cache
func (d *Dao) GetMcnFansIncCache(c context.Context, signID int64, date time.Time) (res *mcnmodel.McnGetMcnFansIncReply, err error) {
	res = new(mcnmodel.McnGetMcnFansIncReply)
	var cache = NewCacheMcnDataSignID(signID, date, res, "McnGetMcnFansIncReply", func(c context.Context, signID int64, date time.Time) (res interface{}, err error) {
		return d.GetMcnFansInc(c, signID, date)
	})
	if err = d.McWrapper.GetOrLoad(c, cache); err != nil {
		log.Error("cache get err, err=%v", err)
		return
	}

	return
}

// GetMcnFansDec 3.6
// see doc http://info.bilibili.co/pages/viewpage.action?pageId=11545690#id-对外接口文档-3.6.查询MCN粉丝取关数按天
func (d *Dao) GetMcnFansDec(c context.Context, signID int64, date time.Time) (res *mcnmodel.McnGetMcnFansDecReply, err error) {
	res = new(mcnmodel.McnGetMcnFansDecReply)

	var q = &datacenter.Query{}
	q.Select("*").Where(
		datacenter.ConditionMapType{"log_date": datacenter.ConditionLte(date)},
		datacenter.ConditionMapType{"sign_id": datacenter.ConditionIn(signID)},
	).Limit(30, 0).Order("log_date desc")

	var api = APIMcnFansDec
	if err = d.callDataAPI(c, api, q, &res.Result); err != nil {
		log.Error("call data api fail, api=%s, err=%s", api, err)
		return
	}
	//log.Info("%s query arg(%d,%+v) res(%+v)", api, signID, date, res.Result[0])
	return
}

//GetMcnFansDecCache GetMcnFansDec with cache
func (d *Dao) GetMcnFansDecCache(c context.Context, signID int64, date time.Time) (res *mcnmodel.McnGetMcnFansDecReply, err error) {
	res = new(mcnmodel.McnGetMcnFansDecReply)
	var cache = NewCacheMcnDataSignID(signID, date, res, "McnGetMcnFansDecReply", func(c context.Context, signID int64, date time.Time) (res interface{}, err error) {
		return d.GetMcnFansDec(c, signID, date)
	})
	if err = d.McWrapper.GetOrLoad(c, cache); err != nil {
		log.Error("cache get err, err=%v", err)
		return
	}

	return
}

// GetMcnFansAttentionWay 3.7
// see doc http://info.bilibili.co/pages/viewpage.action?pageId=11545690#id-对外接口文档-3.7.查询MCN粉丝关注渠道
func (d *Dao) GetMcnFansAttentionWay(c context.Context, signID int64, date time.Time) (res *mcnmodel.McnGetMcnFansAttentionWayReply, err error) {
	res = new(mcnmodel.McnGetMcnFansAttentionWayReply)

	var q = &datacenter.Query{}
	q.Select("*").Where(
		datacenter.ConditionMapType{"log_date": datacenter.ConditionLte(date)},
		datacenter.ConditionMapType{"sign_id": datacenter.ConditionIn(signID)},
	).Limit(1, 0).Order("log_date desc")
	var tmp []*mcnmodel.McnGetMcnFansAttentionWayReply
	var api = APIMcnFansAttentionWay
	if err = d.callDataAPI(c, api, q, &tmp); err != nil {
		log.Error("call data api fail, api=%s, err=%s", api, err)
		return
	}
	if len(tmp) > 0 {
		res = tmp[0]
	}
	//log.Info("%s query arg(%d,%+v) res(%+v)", api, signID, date, tmp[0])
	return
}

//GetMcnFansAttentionWayCache GetMcnFansAttentionWay with cache
func (d *Dao) GetMcnFansAttentionWayCache(c context.Context, signID int64, date time.Time) (res *mcnmodel.McnGetMcnFansAttentionWayReply, err error) {
	res = new(mcnmodel.McnGetMcnFansAttentionWayReply)
	var cache = NewCacheMcnDataSignID(signID, date, res, "McnGetMcnFansAttentionWayReply", func(c context.Context, signID int64, date time.Time) (res interface{}, err error) {
		return d.GetMcnFansAttentionWay(c, signID, date)
	})
	if err = d.McWrapper.GetOrLoad(c, cache); err != nil {
		log.Error("cache get err, err=%v", err)
		return
	}

	return
}

// GetFansBaseFansAttr 3.8
// see doc  http://info.bilibili.co/pages/viewpage.action?pageId=11545690#id-对外接口文档-3.8.查询MCN粉丝/游客基本属性分析(性别占比 观众年龄 观看途径)
func (d *Dao) GetFansBaseFansAttr(c context.Context, signID int64, date time.Time, tp string) (res *mcnmodel.McnGetBaseFansAttrReply, err error) {
	res = new(mcnmodel.McnGetBaseFansAttrReply)

	var group, _ = errgroup.WithContext(c)
	group.Go(func() (err error) {
		var q = &datacenter.Query{}
		q.Select("*").Where(
			datacenter.ConditionMapType{"log_date": datacenter.ConditionLte(date)},
			datacenter.ConditionMapType{"sign_id": datacenter.ConditionIn(signID)},
			datacenter.ConditionMapType{"type": datacenter.ConditionIn(tp)},
		).Limit(1, 0).Order("log_date desc")
		var api = APIMcnFansSex
		var tmp []*datamodel.DmConMcnFansSexW
		if err = d.callDataAPI(context.Background(), api, q, &tmp); err != nil {
			log.Error("call data api fail, api=%s, err=%s", api, err)
			return
		}
		if len(tmp) > 0 {
			res.FansSex = tmp[0]
		}
		//log.Info("%s query arg(%d,%+v,%s) res(%+v)", api, signID, date, tp, tmp[0])
		return
	})
	group.Go(func() (err error) {
		var q = &datacenter.Query{}
		q.Select("*").Where(
			datacenter.ConditionMapType{"log_date": datacenter.ConditionLte(date)},
			datacenter.ConditionMapType{"sign_id": datacenter.ConditionIn(signID)},
			datacenter.ConditionMapType{"type": datacenter.ConditionIn(tp)},
		).Limit(1, 0).Order("log_date desc")
		var tmp []*datamodel.DmConMcnFansAgeW
		var api = APIMcnFansAge
		if err = d.callDataAPI(context.Background(), api, q, &tmp); err != nil {
			log.Error("call data api fail, api=%s, err=%s", api, err)
			return
		}
		if len(tmp) > 0 {
			res.FansAge = tmp[0]
		}
		//log.Info("%s query arg(%d,%+v,%s) res(%+v)", api, signID, date, tp, tmp[0])
		return
	})
	group.Go(func() (err error) {
		var q = &datacenter.Query{}
		q.Select("*").Where(
			datacenter.ConditionMapType{"log_date": datacenter.ConditionLte(date)},
			datacenter.ConditionMapType{"sign_id": datacenter.ConditionIn(signID)},
			datacenter.ConditionMapType{"type": datacenter.ConditionIn(tp)},
		).Limit(1, 0).Order("log_date desc")
		var tmp []*datamodel.DmConMcnFansPlayWayW
		var api = APIMcnFansPlayWay
		if err = d.callDataAPI(context.Background(), api, q, &tmp); err != nil {
			log.Error("call data api fail, api=%s, err=%s", api, err)
			return
		}
		if len(tmp) > 0 {
			res.FansPlayWay = tmp[0]
		}
		//log.Info("%s query arg(%d,%+v,%s) res(%+v)", api, signID, date, tp, tmp[0])
		return
	})
	err = group.Wait()
	if err != nil {
		log.Error("fail to get data, err=%v", err)
		return
	}

	return
}

//GetFansBaseFansAttrCache GetFansBaseFansAttr with cache
func (d *Dao) GetFansBaseFansAttrCache(c context.Context, signID int64, date time.Time, tp string) (res *mcnmodel.McnGetBaseFansAttrReply, err error) {
	res = new(mcnmodel.McnGetBaseFansAttrReply)
	var cache = NewCacheMcnDataWithTp(signID, date, tp, res, "McnGetBaseFansAttrReply", func(c context.Context, signID int64, date time.Time, tp string) (res interface{}, err error) {
		return d.GetFansBaseFansAttr(c, signID, date, tp)
	})
	if err = d.McWrapper.GetOrLoad(c, cache); err != nil {
		log.Error("cache get err, err=%v", err)
		return
	}

	return
}

// GetFansArea 3.9
// see doc http://info.bilibili.co/pages/viewpage.action?pageId=11545690#id-对外接口文档-3.9.查询MCN粉丝/游客地区分布分析
func (d *Dao) GetFansArea(c context.Context, signID int64, date time.Time, tp string) (res *mcnmodel.McnGetFansAreaReply, err error) {
	res = new(mcnmodel.McnGetFansAreaReply)

	var q = &datacenter.Query{}
	q.Select("*").Where(
		datacenter.ConditionMapType{"log_date": datacenter.ConditionIn(date)},
		datacenter.ConditionMapType{"sign_id": datacenter.ConditionIn(signID)},
		datacenter.ConditionMapType{"type": datacenter.ConditionIn(tp)},
	)

	var api = APIMcnFansArea
	if err = d.callDataAPI(c, api, q, &res.Result); err != nil {
		log.Error("call data api fail, api=%s, err=%s", api, err)
		return
	}
	//log.Info("%s query arg(%d,%+v,%s) res(%+v)", api, signID, date, tp, res.Result[0])
	return
}

//GetFansAreaCache GetFansArea with cache
func (d *Dao) GetFansAreaCache(c context.Context, signID int64, date time.Time, tp string) (res *mcnmodel.McnGetFansAreaReply, err error) {
	res = new(mcnmodel.McnGetFansAreaReply)
	var cache = NewCacheMcnDataWithTp(signID, date, tp, res, "McnGetFansAreaReply", func(c context.Context, signID int64, date time.Time, tp string) (res interface{}, err error) {
		return d.GetFansArea(c, signID, date, tp)
	})
	if err = d.McWrapper.GetOrLoad(c, cache); err != nil {
		log.Error("cache get err, err=%v", err)
		return
	}

	return
}

// GetFansType 3.10
// see doc http://info.bilibili.co/pages/viewpage.action?pageId=11545690#id-对外接口文档-3.10.查询MCN粉丝/游客内容倾向分析
func (d *Dao) GetFansType(c context.Context, signID int64, date time.Time, tp string) (res *mcnmodel.McnGetFansTypeReply, err error) {
	res = new(mcnmodel.McnGetFansTypeReply)

	var q = &datacenter.Query{}
	q.Select("*").Where(
		datacenter.ConditionMapType{"log_date": datacenter.ConditionIn(date)},
		datacenter.ConditionMapType{"sign_id": datacenter.ConditionIn(signID)},
		datacenter.ConditionMapType{"type": datacenter.ConditionIn(tp)},
	)

	var api = APIMcnFansType
	if err = d.callDataAPI(c, api, q, &res.Result); err != nil {
		log.Error("call data api fail, api=%s, err=%s", api, err)
		return
	}

	var tids []int64
	for _, v := range res.Result {
		tids = append(tids, v.TypeID)
	}

	tpNames := cache.GetTidNames(tids)
	for _, v := range res.Result {
		if tpName, ok := tpNames[v.TypeID]; ok {
			v.TypeName = tpName
		}
	}
	//log.Info("%s query arg(%d,%+v,%s) res(%+v)", api, signID, date, tp, res.Result[0])
	return
}

//GetFansTypeCache GetFansType with cache
func (d *Dao) GetFansTypeCache(c context.Context, signID int64, date time.Time, tp string) (res *mcnmodel.McnGetFansTypeReply, err error) {
	res = new(mcnmodel.McnGetFansTypeReply)
	var cache = NewCacheMcnDataWithTp(signID, date, tp, res, "McnGetFansTypeReply", func(c context.Context, signID int64, date time.Time, tp string) (res interface{}, err error) {
		return d.GetFansType(c, signID, date, tp)
	})
	if err = d.McWrapper.GetOrLoad(c, cache); err != nil {
		log.Error("cache get err, err=%v", err)
		return
	}

	return
}

// GetFansTag 3.11
// see doc http://info.bilibili.co/pages/viewpage.action?pageId=11545690#id-对外接口文档-3.11.查询MCN粉丝/游客标签地图分析
func (d *Dao) GetFansTag(c context.Context, signID int64, date time.Time, tp string) (res *mcnmodel.McnGetFansTagReply, err error) {
	res = new(mcnmodel.McnGetFansTagReply)

	var q = &datacenter.Query{}
	q.Select("*").Where(
		datacenter.ConditionMapType{"log_date": datacenter.ConditionIn(date)},
		datacenter.ConditionMapType{"sign_id": datacenter.ConditionIn(signID)},
		datacenter.ConditionMapType{"type": datacenter.ConditionIn(tp)},
	)

	var api = APIMcnFansTag
	if err = d.callDataAPI(c, api, q, &res.Result); err != nil {
		log.Error("call data api fail, api=%s, err=%s", api, err)
		return
	}

	var tagIDs []int64
	for _, v := range res.Result {
		tagIDs = append(tagIDs, v.TagID)
	}

	var tagsReply *tagmdl.TagsReply
	if tagsReply, err = global.GetTagGRPC().Tags(c, &tagmdl.TagsReq{Tids: tagIDs}); err != nil {
		log.Error("tag(%+v) grpc client fail, err=%s", tagIDs, err)
		err = nil
	}

	for _, v := range res.Result {
		if tagsReply == nil {
			continue
		}
		if tag, ok := tagsReply.Tags[v.TagID]; ok {
			v.TagName = tag.Name
		}
	}
	//log.Info("%s query arg(%d,%+v,%s) res(%+v)", api, signID, date, tp, res.Result[0])
	return
}

//GetFansTagCache GetFansTag with cache
func (d *Dao) GetFansTagCache(c context.Context, signID int64, date time.Time, tp string) (res *mcnmodel.McnGetFansTagReply, err error) {
	res = new(mcnmodel.McnGetFansTagReply)
	var cache = NewCacheMcnDataWithTp(signID, date, tp, res, "McnGetFansTagReply", func(c context.Context, signID int64, date time.Time, tp string) (res interface{}, err error) {
		return d.GetFansTag(c, signID, date, tp)
	})
	if err = d.McWrapper.GetOrLoad(c, cache); err != nil {
		log.Error("cache get err, err=%v", err)
		return
	}

	return
}
