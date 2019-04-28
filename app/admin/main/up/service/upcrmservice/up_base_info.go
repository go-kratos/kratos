package upcrmservice

import (
	"context"
	"errors"
	"sync"
	"time"

	"go-common/app/admin/main/up/dao/global"
	"go-common/app/admin/main/up/model/datamodel"
	"go-common/app/admin/main/up/model/upcrmmodel"
	"go-common/app/admin/main/up/util"
	accgrpc "go-common/app/service/main/account/api"
	"go-common/library/database/elastic"
	"go-common/library/log"
	"go-common/library/net/metadata"
	"go-common/library/sync/errgroup"
)

const (
	maxSearchItemCount = 10000
	maxBatchCount      = 100
)

var (
	//ErrTooManySearchItem too many search items, 搜索只支持返回前10000条数据
	ErrTooManySearchItem = errors.New("筛选仅支持展示前1万条")
	//ErrNoMid no mid
	ErrNoMid = errors.New("Mid为空")
)

//UpBaseInfoQuery query
func (s *Service) UpBaseInfoQuery(context context.Context, arg *upcrmmodel.InfoQueryArgs) (result *upcrmmodel.InfoQueryResult, err error) {

	var data, e = s.upBaseInfoQueryBatch(s.crmdb.QueryUpBaseInfoBatchByMid, arg.Mid)
	err = e
	// 没找到按照错误处理
	if err != nil || len(data) == 0 {
		log.Error("get from db fail, req=%+v, err=%+v", arg, err)
		return
	}

	result = data[0]

	log.Info("query base info ok, req=%+v, err=%+v", arg, result)
	return
}

//QueryDbFunc query func type
type QueryDbFunc func(fields string, mid ...int64) (result []upcrmmodel.UpBaseInfo, err error)

func (s *Service) upBaseInfoQueryBatch(queryfunc QueryDbFunc, ids ...int64) (result []*upcrmmodel.InfoQueryResult, err error) {
	var data, e = queryfunc("*", ids...)
	err = e
	if err != nil {
		log.Error("get from db fail, err=%+v", err)
		return
	}

	for _, v := range data {
		var info = upcrmmodel.InfoQueryResult{}
		info.CopyFromBaseInfo(v)
		info.CalculateAttr()
		result = append(result, &info)
	}
	return
}

// UpAccountInfo get account info
func (s *Service) UpAccountInfo(c context.Context, arg *upcrmmodel.InfoAccountInfoArgs) (res []*accgrpc.Info, err error) {
	var (
		infosReply *accgrpc.InfosReply
		mids       = util.ExplodeInt64(arg.Mids, ",")
	)
	if infosReply, err = global.GetAccClient().Infos3(c, &accgrpc.MidsReq{Mids: mids, RealIp: metadata.String(c, metadata.RemoteIP)}); err != nil {
		return
	}
	if infosReply == nil || infosReply.Infos == nil {
		return
	}
	for _, v := range infosReply.Infos {
		res = append(res, v)
	}
	log.Info("query acount info ok, req=%+v, result=%+v", arg, res)
	return
}

//SearchResult  struct
type SearchResult struct {
	AccountState           int `json:"account_state"`
	Activity               int `json:"activity"`
	Attr                   int `json:"attr"`
	ArticleCountAccumulate int `json:"article_count_accumulate"`
	ID                     uint32
	Mid                    int64
}

func getAttrFormat(attrs upcrmmodel.UpAttr) (result []int) {
	// 什么要shift，因为es的位是从1开始的，而存储的位是从0开始的
	const shift = 1
	if attrs.AttrVideo != 0 {
		result = append(result, upcrmmodel.AttrBitVideo+shift)
	}
	if attrs.AttrAudio != 0 {
		result = append(result, upcrmmodel.AttrBitAudio+shift)
	}
	if attrs.AttrArticle != 0 {
		result = append(result, upcrmmodel.AttrBitArticle+shift)
	}
	if attrs.AttrPhoto != 0 {
		result = append(result, upcrmmodel.AttrBitPhoto+shift)
	}
	if attrs.AttrSign != 0 {
		result = append(result, upcrmmodel.AttrBitSign+shift)
	}

	if attrs.AttrGrowup != 0 {
		result = append(result, upcrmmodel.AttrBitGrowup+shift)
	}

	if attrs.AttrVerify != 0 {
		result = append(result, upcrmmodel.AttrBitVerify+shift)
	}

	return
}

func getEsCombo(attrs upcrmmodel.UpAttr) (combos []*elastic.Combo) {
	const shift = 1

	var attrs1, attrs2 []interface{}
	var attrFlagList = getAttrFormat(attrs)
	for _, v := range attrFlagList {
		if _, ok := upcrmmodel.AttrGroup1[v-shift]; ok {
			attrs1 = append(attrs1, v)
		} else if _, ok := upcrmmodel.AttrGroup2[v-shift]; ok {
			attrs2 = append(attrs2, v)
		}
	}
	if attrs1 != nil {
		var attrGroup = make(map[string][]interface{})
		attrGroup["attr_format"] = attrs1
		cmb := &elastic.Combo{}
		cmb.ComboIn([]map[string][]interface{}{
			attrGroup},
		).MinIn(1).MinAll(1)
		combos = append(combos, cmb)
	}

	if attrs2 != nil {
		var attrGroup = make(map[string][]interface{})
		attrGroup["attr_format"] = attrs2
		cmb := &elastic.Combo{}
		cmb.ComboIn([]map[string][]interface{}{
			attrGroup},
		).MinIn(1).MinAll(1)
		combos = append(combos, cmb)
	}
	return
}

//UpInfoSearch info search
func (s *Service) UpInfoSearch(c context.Context, arg *upcrmmodel.InfoSearchArgs) (result upcrmmodel.InfoSearchResult, err error) {
	//调用搜索的接口
	var searchData esResult
	searchData, err = s.searchFromEs(c, arg)
	if err != nil {
		log.Error("search arg=%+v, err=%+v", arg, err)
		return
	}

	if len(searchData.Result) == 0 {
		log.Info("no data return from search, just return")
		return
	}
	var ids []int64
	for _, v := range searchData.Result {
		ids = append(ids, int64(v.ID))
	}

	result.Result, err = s.queryUpBaseInfo(c, ids...)
	if err != nil {
		log.Error("query up base info fail, err=%+v", err)
		return
	}
	result.PageInfo = searchData.Page.ToPageInfo()

	log.Info("res=%+v, page=%+v", searchData, result.PageInfo)
	return
}

type esPage struct {
	Num   int `json:"num"`
	Size  int `json:"size"`
	Total int `json:"total"`
}

//ToPageInfo cast to page info
func (e *esPage) ToPageInfo() (pageInfo upcrmmodel.PageInfo) {
	if e == nil {
		return
	}

	pageInfo.TotalCount = e.Total
	pageInfo.Size = e.Size
	pageInfo.Page = e.Num

	return
}

type esResult struct {
	Page   *esPage         `json:"page"`
	Result []*SearchResult `json:"result"`
}

func (s *Service) searchFromEs(c context.Context, arg *upcrmmodel.InfoSearchArgs) (searchData esResult, err error) {
	if arg.Page*arg.Size > maxSearchItemCount {
		err = ErrTooManySearchItem
		return
	}

	var searchSdk = elastic.NewElastic(nil)
	var r = searchSdk.NewRequest("up_crm_info")
	r.Pn(arg.Page).Ps(arg.Size).Index("up_base_info").
		Fields("id")
	if arg.Mid != 0 {
		r.WhereEq("mid", arg.Mid)
	} else {
		if arg.AccountState != 0 {
			// 字段有0值，所以接口改为1，2表示字段的0，1，接口的0表示没有此条件
			var realArg = arg.AccountState - 1
			r.WhereEq("account_state", realArg)
		}
		if arg.Activity != 0 {
			r.WhereEq("activity", arg.Activity)
		}
		var startdate, _ = time.Parse(upcrmmodel.TimeFmtDate, arg.FirstDateBegin)
		var enddate, _ = time.Parse(upcrmmodel.TimeFmtDate, arg.FirstDateEnd)
		enddate = enddate.AddDate(0, 0, 1)
		var startStr = startdate.Format(upcrmmodel.TimeFmtMysql)
		var endStr = enddate.Format(upcrmmodel.TimeFmtMysql)
		if arg.FirstDateBegin != "" && arg.FirstDateEnd != "" {
			r.WhereRange("first_up_time", startStr, endStr, elastic.RangeScopeLcRc)
		}

		if arg.Order.Order == "" {
			arg.Order.Order = "desc"
		}

		if arg.Order.Field == "" {
			arg.Order.Field = "first_up_time"
			arg.Order.Order = "desc"
		}
		var combos = getEsCombo(arg.Attrs)
		r.WhereCombo(combos...)
		r.Order(arg.Order.Field, arg.Order.Order)
	}

	err = r.Scan(c, &searchData)
	return
}

func (s *Service) queryUpBaseInfo(c context.Context, ids ...int64) (result []*upcrmmodel.InfoQueryResult, err error) {

	var group, ctx = errgroup.WithContext(c)
	var infoData []*upcrmmodel.InfoQueryResult
	group.Go(func() error {
		var e error
		infoData, e = s.upBaseInfoQueryBatch(s.crmdb.QueryUpBaseInfoBatchByID, ids...)
		if e != nil {
			err = e
			log.Error("get base info fail, err=%+v", err)
		}
		return nil
	})

	var tidMap = make(map[int64]*datamodel.UpArchiveTypeData)
	var mapLock sync.Mutex
	for _, mid := range ids {
		group.Go(func() error {
			var arg = datamodel.GetUpArchiveTypeInfoArg{Mid: mid}
			var tidData, e = s.dataService.GetUpArchiveTypeInfo(ctx, &arg)
			if e != nil || tidData == nil {
				log.Error("get up type info err, err=%v", e)
				return nil
			}
			mapLock.Lock()
			tidMap[arg.Mid] = tidData
			mapLock.Unlock()
			return nil
		})
	}
	if err = group.Wait(); err != nil {
		log.Error("get data error, err=%v", err)
		return
	}

	var infoIDMap = make(map[uint32]*upcrmmodel.InfoQueryResult)
	for _, v := range infoData {
		infoIDMap[v.ID] = v
	}

	for _, v := range ids {
		var info, ok = infoIDMap[uint32(v)]
		if !ok {
			continue
		}
		if typeInfo, ok := tidMap[v]; ok {
			info.ActiveTid = typeInfo.Tid
			info.ActiveSubtid = typeInfo.SubTid
		}
		result = append(result, info)
	}
	return
}

//QueryUpInfoWithViewerData query with view data
func (s *Service) QueryUpInfoWithViewerData(c context.Context, arg *upcrmmodel.UpInfoWithViewerArg) (result upcrmmodel.UpInfoWithViewerResult, err error) {
	if arg.Page*arg.Size > maxSearchItemCount {
		err = ErrTooManySearchItem
		return
	}
	// 如果是0，则默认设置所有的tag
	if arg.Flag == 0 {
		arg.Flag = -1
	}
	var mids []int64
	if arg.Mids != "" {
		mids = util.ExplodeInt64(arg.Mids, ",")
		var midlen = len(mids)
		if midlen == 0 {
			err = ErrNoMid
			log.Error("no mid get from mids, arg=%+v", arg)
			return
		}
		if midlen > maxBatchCount {
			mids = mids[:maxBatchCount]
		}
	} else {
		var searchSdk = elastic.NewElastic(nil)
		var r = searchSdk.NewRequest("up_crm_info")
		if arg.Size > maxBatchCount {
			arg.Size = maxBatchCount
		}

		r.Pn(arg.Page).Ps(arg.Size).Index("up_base_info").
			Fields("mid").
			Order(arg.Sort, arg.Order)

		var searchData esResult
		err = r.Scan(c, &searchData)
		if err != nil {
			log.Error("fail to get from search, arg=%+v", arg)
			return
		}

		for _, v := range searchData.Result {
			mids = append(mids, int64(v.Mid))
		}

		result.PageInfo = searchData.Page.ToPageInfo()
	}

	var group, ctx = errgroup.WithContext(c)
	var infoData []*upcrmmodel.InfoQueryResult
	var playData []*upcrmmodel.UpPlayInfo
	group.Go(func() error {
		if arg.Flag&upcrmmodel.FlagUpBaseData != 0 {
			infoData, err = s.upBaseInfoQueryBatch(s.crmdb.QueryUpBaseInfoBatchByMid, mids...)
			if err != nil {
				log.Error("query up base error, err=%v", err)
				return err
			}
		}
		if arg.Flag&upcrmmodel.FlagUpPlayData != 0 {
			playData, err = s.crmdb.QueryPlayInfoBatch(mids, upcrmmodel.BusinessTypeVideo)
			if err != nil {
				log.Error("query play info err, err=%v", err)
				return err
			}
		}
		return nil
	})

	var dataMap = make(map[int64]*upcrmmodel.UpInfoWithViewerData)
	if arg.Flag&upcrmmodel.FlagViewData != 0 {
		for _, v := range mids {
			var mid = v // copy this v
			group.Go(func() error {
				var info, e = s.dataService.GetViewData(ctx, mid)
				if e != nil {
					err = e
					log.Error("query up view info from hbase error, err=%v", err)
					return err
				}
				var data = getOrCreateUpInfo(dataMap, mid)
				data.ViewerBase = info.Base
				data.ViewerTrend = info.Trend
				data.ViewerArea = info.Area
				return err
			})
		}
	}
	if err = group.Wait(); err != nil {
		log.Error("get data fail, err=%v", err)
		return
	}

	for _, baseInfo := range infoData {
		var data = getOrCreateUpInfo(dataMap, baseInfo.Mid)
		data.UpBaseInfo = baseInfo
	}

	for _, playInfo := range playData {
		var data = getOrCreateUpInfo(dataMap, playInfo.Mid)
		data.UpPlayInfo = playInfo
	}

	for _, mid := range mids {
		var data, ok = dataMap[mid]
		if !ok {
			log.Warn("up info not found, mid=%d", mid)
			continue
		}
		data.Mid = mid
		result.Result = append(result.Result, data)
	}
	log.Info("query up with view ok, arg=%+v, result count=%d", arg, len(mids))
	return
}

var dataMapMutex sync.Mutex

func getOrCreateUpInfo(dataMap map[int64]*upcrmmodel.UpInfoWithViewerData, mid int64) (result *upcrmmodel.UpInfoWithViewerData) {
	dataMapMutex.Lock()
	defer dataMapMutex.Unlock()

	var ok bool
	if result, ok = dataMap[mid]; !ok {
		result = &upcrmmodel.UpInfoWithViewerData{}
		dataMap[mid] = result
	}
	return result
}
