package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go-common/app/admin/main/aegis/model/common"
	"go-common/app/admin/main/aegis/model/net"
	"go-common/library/log"
	"go-common/library/xstr"
)

const netCacheNull = "{null}"

//ICacheObjDest 对象缓存接口
type ICacheObjDest interface {
	//GetKey 获取对象缓存的key
	GetKey(id int64, opt ...interface{}) (k []string)
	//AppendString--v is from cache, to decode & append to dest, && check null
	AppendString(id int64, v string) (ok bool, err error)
	//AppendRaw--get from db, to encode & append to dest & to convert to string
	AppendRaw(c context.Context, s *Service, miss []int64, kopt []interface{}, opt ...interface{}) (missCache map[string]string, err error)
}

//ICacheDecoder 解析字符串
type ICacheDecoder interface {
	//--get from cache, to decode & append to dest, && check null
	AppendString(id int64, v string) (ok bool, err error)
}

//FCacheKeyPro 根据id获取key
type FCacheKeyPro func(id int64, opt ...interface{}) (k []string)

//FCacheAggregateRelation db获取聚合字段和对象id的关联
type FCacheAggregateRelation func(c context.Context, miss []int64, kopt []interface{}, opt ...interface{}) (rela map[int64][]int64, err error)

//FCacheRawByAggr db根据聚合值，获取对象集合和关联
type FCacheRawByAggr func(c context.Context, miss []int64, objDest ICacheObjDest, kopt []interface{}, opt ...interface{}) (objCache map[string]string, aggrRelation map[int64][]int64, err error)

//AggrCacheDest 聚合缓存
type AggrCacheDest struct {
	Dest        map[int64][]int64
	GetKey      FCacheKeyPro
	AggrRelaRaw FCacheAggregateRelation
}

//AppendString 实现ICacheDecoder接口
func (a *AggrCacheDest) AppendString(id int64, v string) (ok bool, err error) {
	var (
		res []int64
	)
	if a.Dest == nil {
		a.Dest = map[int64][]int64{}
	}
	if res, err = xstr.SplitInts(v); err != nil {
		log.Error("AggrCacheDest AppendString(%s) error(%v)", v, err)
		return
	}

	ok = true
	a.Dest[id] = res
	return
}

//AppendRelationRaw 根据聚合字段值，从db获取聚合字段与对象id的关联
func (a *AggrCacheDest) AppendRelationRaw(c context.Context, miss []int64, kopt []interface{}, opt ...interface{}) (missCache map[string]string, err error) {
	var (
		res map[int64][]int64
	)

	if res, err = a.AggrRelaRaw(c, miss, kopt, opt...); err != nil {
		log.Error("AggrCacheDest AppendRelationRaw error(%v) miss(%+v)", err, miss)
		return
	}

	if len(res) == 0 {
		return
	}
	a.Dest = common.CopyMap(res, a.Dest, true)
	missCache = make(map[string]string, len(res))
	for aggr, ids := range res {
		v := xstr.JoinInts(ids)
		for _, k := range a.GetKey(aggr, kopt...) {
			missCache[k] = v
		}
	}
	return
}

//CacheWrap 以net对象为例子，可以根据business_id字段聚合
type CacheWrap struct {
	ObjDest       ICacheObjDest   //对象集合
	AggrObjRaw    FCacheRawByAggr //根据聚合字段获取对象集合from db
	AggregateDest *AggrCacheDest  //对象的某个聚合字段与对象id的集合
}

//AppendAggregateRaw 从db获取聚合字段映射的对象集合&关联
func (w *CacheWrap) AppendAggregateRaw(c context.Context, miss []int64, kopt []interface{}, opt ...interface{}) (missCache map[string]string, err error) {
	var (
		aggrRela map[int64][]int64
	)

	if missCache, aggrRela, err = w.AggrObjRaw(c, miss, w.ObjDest, kopt, opt...); err != nil {
		return
	}

	if missCache == nil {
		missCache = map[string]string{}
	}
	for aggrFieldVal, objIDList := range aggrRela {
		v := xstr.JoinInts(objIDList)
		for _, k := range w.AggregateDest.GetKey(aggrFieldVal, kopt...) {
			missCache[k] = v
		}
	}
	w.AggregateDest.Dest = common.CopyMap(aggrRela, w.AggregateDest.Dest, true)
	return
}

//getFromCache 获取缓存并解析
func (s *Service) getFromCache(c context.Context, ids []int64, keyPro FCacheKeyPro, dest ICacheDecoder, kopt []interface{}) (miss []int64) {
	var (
		caches []string
		err    error
		ok     bool
	)

	if len(ids) == 0 {
		return
	}

	//get from cache, chech whether miss
	ids = common.Unique(ids, true)
	keys := []string{}
	keymap := map[string]int64{}
	for _, id := range ids {
		for _, k := range keyPro(id, kopt...) {
			if keymap[k] != 0 {
				continue
			}

			keymap[k] = id
			keys = append(keys, k)
		}
	}

	if caches, err = s.redis.MGet(c, keys...); err != nil {
		miss = ids
		err = nil
		return
	}
	log.Info("getFromCache caches(%+v) keys(%+v)", caches, keys)
	missmap := map[int64]int64{}
	for i, item := range caches {
		id := keymap[keys[i]]
		if item != "" && item != netCacheNull {
			if ok, err = dest.AppendString(id, item); err == nil && ok {
				continue
			}
		}
		if missmap[id] > 0 || item == netCacheNull {
			continue
		}

		miss = append(miss, id)
		missmap[id] = id
	}
	log.Info("getFromCache miss(%+v) keys(%+v)", miss, keys)
	return
}

func (s *Service) setNetCache() {
	for {
		missCache, open := <-s.netCacheCh
		if !open {
			log.Info("s.netCacheCh closed!")
			break
		}

		s.redis.SetMulti(context.TODO(), missCache)
		time.Sleep(time.Millisecond * 1)
	}
}

//objCache 根据对象id，获取对象
func (s *Service) objCache(c context.Context, ids []int64, dest ICacheObjDest, kopt []interface{}, opt ...interface{}) (err error) {
	var (
		miss      []int64
		missCache map[string]string
	)

	//get from cache
	if miss = s.getFromCache(c, ids, dest.GetKey, dest, kopt); len(miss) == 0 {
		return
	}

	//get miss from db
	if missCache, err = dest.AppendRaw(c, s, miss, kopt, opt...); err != nil {
		return
	}

	//set miss to cache
	if missCache == nil {
		missCache = map[string]string{}
	}
	for _, missid := range miss {
		missks := dest.GetKey(missid, kopt...)
		for _, missk := range missks {
			if _, exist := missCache[missk]; exist {
				continue
			}
			missCache[missk] = netCacheNull
		}
	}
	log.Info("objCache missCache(%+v)", missCache)
	s.netCacheCh <- missCache
	return
}

//aggregateRelationCache 根据聚合值，获取聚合字段与对象的映射关系
func (s *Service) aggregateRelationCache(c context.Context, aggrID []int64, dest *AggrCacheDest, kopt []interface{}, aggrOpt ...interface{}) (err error) {
	var (
		miss      []int64
		missCache map[string]string
	)

	//get from cache
	if miss = s.getFromCache(c, aggrID, dest.GetKey, dest, kopt); len(miss) == 0 {
		return
	}

	//get miss from db
	if missCache, err = dest.AppendRelationRaw(c, miss, kopt, aggrOpt...); err != nil {
		return
	}

	//set miss to cache
	if missCache == nil {
		missCache = map[string]string{}
	}
	for _, missid := range miss {
		missks := dest.GetKey(missid, kopt...)
		for _, missk := range missks {
			if _, exist := missCache[missk]; exist {
				continue
			}
			missCache[missk] = netCacheNull
		}
	}
	log.Info("aggregateRelationCache missCache(%+v)", missCache)
	s.netCacheCh <- missCache
	return
}

//aggregateCache 根据聚合字段的值，获取对应的对象,比如：根据flowid获取获取其绑定关系binds
func (s *Service) aggregateCache(c context.Context, aggrID []int64, wrap *CacheWrap, aggrkopt, objkopt, aggrOpt, objOpt []interface{}) (err error) {
	var (
		ptrAggrMiss, ptrAggrAll, ptrCache, ptrObjMiss string = "aggrmiss", "aggrall", "cache", "objmiss"
		current                                       string
		needObjRelat                                  map[int64][]int64
		aggrMiss, objMiss                             []int64
		missCache, missObjCache                       map[string]string
	)

	//获取flowid对应哪些binds的映射关系
	if aggrMiss = s.getFromCache(c, aggrID, wrap.AggregateDest.GetKey, wrap.AggregateDest, aggrkopt); len(aggrMiss) == 0 {
		current = ptrAggrAll
	} else {
		current = ptrAggrMiss
	}
	needObjRelat = common.CopyMap(wrap.AggregateDest.Dest, needObjRelat, true)

	if current == ptrAggrMiss {
		if missCache, err = wrap.AppendAggregateRaw(c, aggrMiss, aggrkopt, aggrOpt...); err != nil {
			return
		}

		if len(needObjRelat) <= 0 { //全部从db获取
			current = ptrCache
		} else { //部分缓存获取聚合字段与对象id的关联，部分db根据聚合字段获取对象体
			current = ptrAggrAll
		}
	}

	if current == ptrAggrAll {
		objID := []int64{}
		for _, obji := range needObjRelat {
			objID = append(objID, obji...)
		}

		if objMiss = s.getFromCache(c, objID, wrap.ObjDest.GetKey, wrap.ObjDest, objkopt); len(objMiss) == 0 {
			return
		}
		current = ptrObjMiss
	}

	if missCache == nil {
		missCache = map[string]string{}
	}
	if current == ptrObjMiss {
		if missObjCache, err = wrap.ObjDest.AppendRaw(c, s, objMiss, objkopt, objOpt...); err != nil {
			return
		}
		for k, v := range missObjCache {
			missCache[k] = v
		}
		current = ptrCache
	}

	if current == ptrCache {
		for _, missid := range aggrMiss {
			missks := wrap.AggregateDest.GetKey(missid, aggrkopt...)
			for _, missk := range missks {
				if _, exist := missCache[missk]; exist {
					continue
				}
				missCache[missk] = netCacheNull
			}
		}
		for _, missid := range objMiss {
			missks := wrap.ObjDest.GetKey(missid, objkopt...)
			for _, missk := range missks {
				if _, exist := missCache[missk]; exist {
					continue
				}
				missCache[missk] = netCacheNull
			}
		}

		log.Info("aggregateCache missCache(%+v)", missCache)
		s.netCacheCh <- missCache
	}
	return
}

//delCache 删除缓存
func (s *Service) delCache(c context.Context, keys []string) (err error) {
	if len(keys) == 0 {
		return
	}

	if err = s.redis.DelMulti(c, keys...); err != nil {
		log.Error("delCache s.redis.DelMulti error(%v) keys(%+v)", err, keys)
	}
	return
}

/**
 * 业务应用
 * 需要缓存的对象
 ******************************************
 * one2many
 * key       type          val
 * flow  int64+avail   []*tokenbind,err--ok
 * flow  int64+dir     []*dir,err---ok
 * netid int64         []*flow,err---ok
 * netid []int64,avail,disp       []int64,err--transitionid---ok
 * tranid []int64,isbatch  []*tokenbind,err---ok
 * tranid int64,avail  []*dir.err--ok
 * biz    int64        []int64,err---netid ---ok
 *
 * one2one
 * ids    []int64,xxx,yy []*struct,err --- ok
 ******************************************
 */
//net的聚合字段business_id
type netArr struct {
	dest []*net.Net
}

func (n *netArr) GetKey(id int64, opt ...interface{}) []string {
	return []string{fmt.Sprintf("net:%d", id)}
}

func (n *netArr) AppendString(id int64, v string) (ok bool, err error) {
	one := &net.Net{}
	if err = json.Unmarshal([]byte(v), one); err != nil {
		log.Error("netArr AppendString json.Unmarshal error(%v) v(%s)", err, v)
		return
	}

	ok = true
	n.dest = append(n.dest, one)
	return
}

func (n *netArr) AppendRaw(c context.Context, s *Service, miss []int64, kopt []interface{}, opt ...interface{}) (missCache map[string]string, err error) {
	var (
		res []*net.Net
		bs  []byte
	)
	if res, err = s.gorm.Nets(c, miss); err != nil {
		log.Error("netArr AppendRaw error(%+v) miss(%+v)", err, miss)
		return
	}

	n.dest = append(n.dest, res...)
	missCache = map[string]string{}
	for _, item := range res {
		if bs, err = json.Marshal(item); err != nil {
			log.Error("netArr AppendRaw json.Marshal error(%v) item(%+v)", err, item)
			continue
		}
		v := string(bs)
		for _, k := range n.GetKey(item.ID, kopt...) {
			missCache[k] = v
		}
	}
	return
}

func (s *Service) bizNetRelation(c context.Context, miss []int64, kopt []interface{}, opt ...interface{}) (rela map[int64][]int64, err error) {
	rela, err = s.gorm.NetIDByBusiness(c, miss)
	return
}

func (s *Service) bizNetKey(id int64, opt ...interface{}) []string {
	return []string{fmt.Sprintf("business_net:%d", id)}
}

func (s *Service) netIDByBusiness(c context.Context, business int64) (res []int64, err error) {
	aggr := &AggrCacheDest{
		GetKey:      s.bizNetKey,
		AggrRelaRaw: s.bizNetRelation,
	}

	if err = s.aggregateRelationCache(c, []int64{business}, aggr, nil); err != nil {
		return
	}

	res = aggr.Dest[business]
	return
}

func (s *Service) netByID(c context.Context, id int64) (res *net.Net, err error) {
	dest := &netArr{}
	if err = s.objCache(c, []int64{id}, dest, nil); err != nil {
		return
	}
	if len(dest.dest) == 0 {
		return
	}

	res = dest.dest[0]
	return
}

func (s *Service) delNetCache(c context.Context, n *net.Net) (err error) {
	objdest := &netArr{}
	keys := objdest.GetKey(n.ID)
	keys = append(keys, s.bizNetKey(n.BusinessID)...)
	err = s.delCache(c, keys)
	return
}

//tokenArr 根据id查询对象
type tokenArr struct {
	dest []*net.Token
}

func (n *tokenArr) GetKey(id int64, opt ...interface{}) []string {
	return []string{fmt.Sprintf("token:%d", id)}
}

func (n *tokenArr) AppendString(id int64, v string) (ok bool, err error) {
	one := &net.Token{}
	if err = json.Unmarshal([]byte(v), one); err != nil {
		log.Error("tokenArr AppendString json.Unmarshal error(%v) v(%s)", err, v)
		return
	}

	ok = true
	n.dest = append(n.dest, one)
	return
}

func (n *tokenArr) AppendRaw(c context.Context, s *Service, miss []int64, kopt []interface{}, opt ...interface{}) (missCache map[string]string, err error) {
	var (
		res []*net.Token
		bs  []byte
	)
	if res, err = s.gorm.Tokens(c, miss); err != nil {
		log.Error("tokenArr AppendRaw error(%+v) miss(%+v)", err, miss)
		return
	}

	n.dest = append(n.dest, res...)
	missCache = map[string]string{}
	for _, item := range res {
		if bs, err = json.Marshal(item); err != nil {
			log.Error("tokenArr AppendRaw json.Marshal error(%v) item(%+v)", err, item)
			continue
		}
		v := string(bs)
		for _, k := range n.GetKey(item.ID, kopt...) {
			missCache[k] = v
		}
	}
	return
}

func (s *Service) tokens(c context.Context, ids []int64) (res []*net.Token, err error) {
	objdest := &tokenArr{}
	if err = s.objCache(c, ids, objdest, nil); err != nil {
		return
	}

	res = objdest.dest
	return
}

//token_bind的聚合字段element_id
type bindArr struct {
	dest []*net.TokenBind
}

func (n *bindArr) GetKey(id int64, opt ...interface{}) []string {
	return []string{fmt.Sprintf("bind:%d", id)}
}

func (n *bindArr) AppendString(id int64, v string) (ok bool, err error) {
	one := &net.TokenBind{}
	if err = json.Unmarshal([]byte(v), one); err != nil {
		log.Error("bindArr AppendString json.Unmarshal error(%v) v(%s)", err, v)
		return
	}

	ok = true
	n.dest = append(n.dest, one)
	return
}

func (n *bindArr) AppendRaw(c context.Context, s *Service, miss []int64, kopt []interface{}, opt ...interface{}) (missCache map[string]string, err error) {
	var (
		res []*net.TokenBind
		bs  []byte
	)
	if res, err = s.gorm.TokenBinds(c, miss); err != nil {
		log.Error("bindArr AppendRaw error(%+v) miss(%+v)", err, miss)
		return
	}

	n.dest = append(n.dest, res...)
	missCache = map[string]string{}
	for _, item := range res {
		if bs, err = json.Marshal(item); err != nil {
			log.Error("bindArr AppendRaw json.Marshal error(%v) item(%+v)", err, item)
			continue
		}
		v := string(bs)
		for _, k := range n.GetKey(item.ID, kopt...) {
			missCache[k] = v
		}
	}
	return
}

func (s *Service) elementBindKey(id int64, opt ...interface{}) []string {
	return []string{fmt.Sprintf("element_bind:%d", id)}
}

//全部类型的正常状态的绑定
func (s *Service) elementBindAppendRaw(c context.Context, miss []int64, objdest ICacheObjDest, kopt []interface{}, opt ...interface{}) (objCache map[string]string, aggrRelation map[int64][]int64, err error) {
	var (
		res map[int64][]*net.TokenBind
		bs  []byte
	)

	objdester := objdest.(*bindArr)
	if res, err = s.gorm.TokenBindByElement(c, miss, net.BindTypes, true); err != nil {
		return
	}

	aggrRelation = map[int64][]int64{}
	objCache = map[string]string{}
	for eid, tbs := range res {
		for _, item := range tbs {
			if bs, err = json.Marshal(item); err != nil {
				return
			}

			v := string(bs)
			for _, k := range objdest.GetKey(item.ID, kopt...) {
				objCache[k] = v
			}
			objdester.dest = append(objdester.dest, item)
			aggrRelation[eid] = append(aggrRelation[eid], item.ID)
		}
	}
	return
}

func (s *Service) tokenBindByElement(c context.Context, eid []int64, tp []int8) (res []*net.TokenBind, err error) {
	objdest := &bindArr{}
	w := &CacheWrap{
		ObjDest:    objdest,
		AggrObjRaw: s.elementBindAppendRaw,
		AggregateDest: &AggrCacheDest{
			GetKey: s.elementBindKey,
		},
	}

	//获取eid对应的所有类型的bind对象
	if err = s.aggregateCache(c, eid, w, nil, nil, nil, nil); err != nil {
		return
	}

	res = []*net.TokenBind{}
	for _, item := range objdest.dest {
		included := false
		for _, bindtp := range tp {
			if bindtp == item.Type {
				included = true
				break
			}
		}
		if included {
			res = append(res, item)
		}
	}
	return
}

func (s *Service) tokenBinds(c context.Context, ids []int64, onlyAvailable bool) (res []*net.TokenBind, err error) {
	objdest := &bindArr{}
	if err = s.objCache(c, ids, objdest, nil); err != nil {
		return
	}

	if !onlyAvailable {
		res = objdest.dest
		return
	}
	for _, item := range objdest.dest {
		if item.IsAvailable() {
			res = append(res, item)
		}
	}
	return
}

//flow的聚合字段是net_id
type flowArr struct {
	dest []*net.Flow
}

func (n *flowArr) GetKey(id int64, opt ...interface{}) []string {
	return []string{fmt.Sprintf("flow:%d", id)}
}

func (n *flowArr) AppendString(id int64, v string) (ok bool, err error) {
	one := &net.Flow{}
	if err = json.Unmarshal([]byte(v), one); err != nil {
		log.Error("flowArr AppendString json.Unmarshal error(%v) v(%s)", err, v)
		return
	}

	ok = true
	n.dest = append(n.dest, one)
	return
}

func (n *flowArr) AppendRaw(c context.Context, s *Service, miss []int64, kopt []interface{}, opt ...interface{}) (missCache map[string]string, err error) {
	var (
		res []*net.Flow
		bs  []byte
	)
	if res, err = s.gorm.Flows(c, miss); err != nil {
		log.Error("flowArr AppendRaw error(%+v) miss(%+v)", err, miss)
		return
	}

	n.dest = append(n.dest, res...)
	missCache = map[string]string{}
	for _, item := range res {
		if bs, err = json.Marshal(item); err != nil {
			log.Error("flowArr AppendRaw json.Marshal error(%v) item(%+v)", err, item)
			continue
		}
		v := string(bs)
		for _, k := range n.GetKey(item.ID, kopt...) {
			missCache[k] = v
		}
	}
	return
}

func (s *Service) netFlowKey(id int64, opt ...interface{}) []string {
	return []string{fmt.Sprintf("net_flow:%d", id)}
}

func (s *Service) netFlowAppendRaw(c context.Context, miss []int64, objdest ICacheObjDest, kopt []interface{}, opt ...interface{}) (objCache map[string]string, aggrRelation map[int64][]int64, err error) {
	var (
		res []*net.Flow
		bs  []byte
	)

	objdester := objdest.(*flowArr)
	if res, err = s.gorm.FlowsByNet(c, miss); err != nil {
		return
	}

	aggrRelation = map[int64][]int64{}
	objCache = map[string]string{}
	for _, item := range res {
		if bs, err = json.Marshal(item); err != nil {
			return
		}
		v := string(bs)
		for _, k := range objdest.GetKey(item.ID, kopt...) {
			objCache[k] = v
		}
		objdester.dest = append(objdester.dest, item)
		aggrRelation[item.NetID] = append(aggrRelation[item.NetID], item.ID)
	}
	return
}

func (s *Service) netFlowRelation(c context.Context, miss []int64, kopt []interface{}, opt ...interface{}) (rela map[int64][]int64, err error) {
	rela, err = s.gorm.FlowIDByNet(c, miss)
	return
}

func (s *Service) flowIDByNet(c context.Context, nid int64) (res []int64, err error) {
	dest := &AggrCacheDest{
		GetKey:      s.netFlowKey,
		AggrRelaRaw: s.netFlowRelation,
	}

	if err = s.aggregateRelationCache(c, []int64{nid}, dest, nil); err != nil {
		return
	}

	res = dest.Dest[nid]
	return
}

func (s *Service) flowsByNet(c context.Context, nid int64) (res []*net.Flow, err error) {
	objdest := &flowArr{}
	w := &CacheWrap{
		ObjDest:    objdest,
		AggrObjRaw: s.netFlowAppendRaw,
		AggregateDest: &AggrCacheDest{
			GetKey: s.netFlowKey,
		},
	}

	if err = s.aggregateCache(c, []int64{nid}, w, nil, nil, nil, nil); err != nil {
		return
	}

	res = objdest.dest
	return
}

func (s *Service) flows(c context.Context, ids []int64, onlyAvailable bool) (list []*net.Flow, err error) {
	objdest := &flowArr{}
	if err = s.objCache(c, ids, objdest, nil); err != nil {
		return
	}
	if !onlyAvailable {
		list = objdest.dest
		return
	}

	for _, item := range objdest.dest {
		if item.IsAvailable() {
			list = append(list, item)
		}
	}
	return
}

func (s *Service) flowByID(c context.Context, id int64) (f *net.Flow, err error) {
	var list []*net.Flow
	if list, err = s.flows(c, []int64{id}, false); err != nil {
		return
	}
	if len(list) == 0 {
		return
	}

	f = list[0]
	return
}

func (s *Service) delFlowCache(c context.Context, f *net.Flow, changedBind []int64) (err error) {
	objdest := &flowArr{}
	keys := objdest.GetKey(f.ID)
	keys = append(keys, s.netFlowKey(f.NetID)...)
	keys = append(keys, s.elementBindKey(f.ID)...)
	binds := &bindArr{}
	for _, bind := range changedBind {
		keys = append(keys, binds.GetKey(bind)...)
	}

	err = s.delCache(c, keys)
	return
}

//transition的聚合字段是net_id
type tranArr struct {
	dest []*net.Transition
}

func (n *tranArr) GetKey(id int64, opt ...interface{}) []string {
	return []string{fmt.Sprintf("transition:%d", id)}
}

func (n *tranArr) AppendString(id int64, v string) (ok bool, err error) {
	one := &net.Transition{}
	if err = json.Unmarshal([]byte(v), one); err != nil {
		log.Error("tranArr AppendString json.Unmarshal error(%v) v(%s)", err, v)
		return
	}

	ok = true
	n.dest = append(n.dest, one)
	return
}

func (n *tranArr) AppendRaw(c context.Context, s *Service, miss []int64, kopt []interface{}, opt ...interface{}) (missCache map[string]string, err error) {
	var (
		res []*net.Transition
		bs  []byte
	)
	if res, err = s.gorm.Transitions(c, miss); err != nil {
		log.Error("tranArr AppendRaw error(%+v) miss(%+v)", err, miss)
		return
	}

	n.dest = append(n.dest, res...)
	missCache = map[string]string{}
	for _, item := range res {
		if bs, err = json.Marshal(item); err != nil {
			log.Error("tranArr AppendRaw json.Marshal error(%v) item(%+v)", err, item)
			continue
		}
		v := string(bs)
		for _, k := range n.GetKey(item.ID, kopt...) {
			missCache[k] = v
		}
	}
	return
}

func (s *Service) netTranKey(id int64, opt ...interface{}) []string {
	return []string{fmt.Sprintf("net_transition_dispatch%v_avail%v:%d", opt[0].(bool), opt[1].(bool), id)}
}

func (s *Service) netTranRelation(c context.Context, miss []int64, kopt []interface{}, opt ...interface{}) (rela map[int64][]int64, err error) {
	rela, err = s.gorm.TransitionIDByNet(c, miss, opt[0].(bool), opt[1].(bool))
	return
}

func (s *Service) tranIDByNet(c context.Context, nid []int64, onlyDispatch bool, onlyAvailable bool) (res []int64, err error) {
	dest := &AggrCacheDest{
		GetKey:      s.netTranKey,
		AggrRelaRaw: s.netTranRelation,
	}

	if err = s.aggregateRelationCache(c, nid, dest, []interface{}{onlyDispatch, onlyAvailable}, onlyDispatch, onlyAvailable); err != nil {
		return
	}

	for _, ids := range dest.Dest {
		res = append(res, ids...)
	}
	return
}

func (s *Service) transitions(c context.Context, ids []int64, onlyAvailable bool) (list []*net.Transition, err error) {
	objdest := &tranArr{}
	if err = s.objCache(c, ids, objdest, nil); err != nil {
		return
	}
	if !onlyAvailable {
		list = objdest.dest
		return
	}

	for _, item := range objdest.dest {
		if item.IsAvailable() {
			list = append(list, item)
		}
	}
	return
}

func (s *Service) transitionByID(c context.Context, id int64) (res *net.Transition, err error) {
	var list []*net.Transition
	if list, err = s.transitions(c, []int64{id}, false); err != nil {
		return
	}
	if len(list) == 0 {
		return
	}

	res = list[0]
	return
}

func (s *Service) delTranCache(c context.Context, f *net.Transition, changedBind []int64) (err error) {
	objdest := &tranArr{}
	keys := objdest.GetKey(f.ID)
	keys = append(keys, s.netTranKey(f.NetID, true, true)...)
	keys = append(keys, s.netTranKey(f.NetID, true, false)...)
	keys = append(keys, s.netTranKey(f.NetID, false, true)...)
	keys = append(keys, s.netTranKey(f.NetID, false, false)...)
	keys = append(keys, s.elementBindKey(f.ID)...)
	binds := &bindArr{}
	for _, bind := range changedBind {
		keys = append(keys, binds.GetKey(bind)...)
	}

	err = s.delCache(c, keys)
	return
}

//direction的聚合字段是flow_id/transition_id
type dirArr struct {
	dest []*net.Direction
}

func (n *dirArr) GetKey(id int64, opt ...interface{}) []string {
	return []string{fmt.Sprintf("direction:%d", id)}
}

func (n *dirArr) AppendString(id int64, v string) (ok bool, err error) {
	one := &net.Direction{}
	if err = json.Unmarshal([]byte(v), one); err != nil {
		log.Error("dirArr AppendString json.Unmarshal error(%v) v(%s)", err, v)
		return
	}

	ok = true
	n.dest = append(n.dest, one)
	return
}

func (n *dirArr) AppendRaw(c context.Context, s *Service, miss []int64, kopt []interface{}, opt ...interface{}) (missCache map[string]string, err error) {
	var (
		res []*net.Direction
		bs  []byte
	)
	if res, err = s.gorm.Directions(c, miss); err != nil {
		log.Error("dirArr AppendRaw error(%+v) miss(%+v)", err, miss)
		return
	}

	n.dest = append(n.dest, res...)
	missCache = map[string]string{}
	for _, item := range res {
		if bs, err = json.Marshal(item); err != nil {
			log.Error("dirArr AppendRaw json.Marshal error(%v) item(%+v)", err, item)
			continue
		}
		v := string(bs)
		for _, k := range n.GetKey(item.ID, kopt...) {
			missCache[k] = v
		}
	}
	return
}

func (s *Service) flowDirKey(id int64, opt ...interface{}) []string {
	di := opt[0].([]int8)
	res := make([]string, len(di))
	for i, d := range di {
		res[i] = fmt.Sprintf("flow_direction_%d:%d", d, id)
	}
	return res
}
func (s *Service) tranDirKey(id int64, opt ...interface{}) []string {
	di := opt[0].([]int8)
	avail := opt[1].(bool)
	res := make([]string, len(di))
	for i, d := range di {
		res[i] = fmt.Sprintf("transition_direction_avail%v_%d:%d", avail, d, id)
	}
	return res

}

func (s *Service) flowDirAppendRaw(c context.Context, miss []int64, objdest ICacheObjDest, kopt []interface{}, opt ...interface{}) (objCache map[string]string, aggrRelation map[int64][]int64, err error) {
	var (
		res []*net.Direction
		bs  []byte
	)

	objdester := objdest.(*dirArr)
	if res, err = s.gorm.DirectionByFlowID(c, miss, opt[0].(int8)); err != nil {
		return
	}

	aggrRelation = map[int64][]int64{}
	objCache = map[string]string{}
	for _, item := range res {
		if bs, err = json.Marshal(item); err != nil {
			return
		}

		v := string(bs)
		for _, k := range objdest.GetKey(item.ID, kopt...) {
			objCache[k] = v
		}
		objdester.dest = append(objdester.dest, item)
		aggrRelation[item.FlowID] = append(aggrRelation[item.FlowID], item.ID)
	}
	return
}

func (s *Service) tranDirAppendRaw(c context.Context, miss []int64, objdest ICacheObjDest, kopt []interface{}, opt ...interface{}) (objCache map[string]string, aggrRelation map[int64][]int64, err error) {
	var (
		res []*net.Direction
		bs  []byte
	)

	objdester := objdest.(*dirArr)
	if res, err = s.gorm.DirectionByTransitionID(c, miss, opt[0].(int8), opt[1].(bool)); err != nil {
		return
	}

	aggrRelation = map[int64][]int64{}
	objCache = map[string]string{}
	for _, item := range res {
		if bs, err = json.Marshal(item); err != nil {
			return
		}

		v := string(bs)
		for _, k := range objdest.GetKey(item.ID, kopt...) {
			objCache[k] = v
		}
		objdester.dest = append(objdester.dest, item)
		aggrRelation[item.TransitionID] = append(aggrRelation[item.TransitionID], item.ID)
	}
	return
}

func (s *Service) dirByFlow(c context.Context, nid []int64, dir int8) (res []*net.Direction, err error) {
	objdest := &dirArr{}
	w := &CacheWrap{
		ObjDest:    objdest,
		AggrObjRaw: s.flowDirAppendRaw,
		AggregateDest: &AggrCacheDest{
			GetKey: s.flowDirKey,
		},
	}

	if err = s.aggregateCache(c, nid, w, []interface{}{[]int8{dir}}, nil, []interface{}{dir}, nil); err != nil {
		return
	}

	res = objdest.dest
	return
}

func (s *Service) dirByTran(c context.Context, nid []int64, dir int8, onlyAvailable bool) (res []*net.Direction, err error) {
	objdest := &dirArr{}
	w := &CacheWrap{
		ObjDest:    objdest,
		AggrObjRaw: s.tranDirAppendRaw,
		AggregateDest: &AggrCacheDest{
			GetKey: s.tranDirKey,
		},
	}

	if err = s.aggregateCache(c, nid, w, []interface{}{[]int8{dir}, onlyAvailable}, nil, []interface{}{dir, onlyAvailable}, nil); err != nil {
		return
	}

	res = objdest.dest
	return
}
func (s *Service) directions(c context.Context, ids []int64) (list []*net.Direction, err error) {
	objdest := &dirArr{}
	if err = s.objCache(c, ids, objdest, nil); err != nil {
		return
	}

	list = objdest.dest
	return
}

func (s *Service) delDirCache(c context.Context, dir *net.Direction) (err error) {
	objdest := &dirArr{}
	keys := objdest.GetKey(dir.ID)
	keys = append(keys, s.flowDirKey(dir.FlowID, []int8{net.DirInput, net.DirOutput})...)
	keys = append(keys, s.tranDirKey(dir.TransitionID, []int8{net.DirInput, net.DirOutput}, true)...)
	keys = append(keys, s.tranDirKey(dir.TransitionID, []int8{net.DirInput, net.DirOutput}, false)...)

	err = s.delCache(c, keys)
	return
}
