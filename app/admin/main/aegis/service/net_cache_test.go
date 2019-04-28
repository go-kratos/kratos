package service

/**
 * 需要缓存的对象
 ******************************************
 * one2many
 * key       type          val
 * flow  int64+avail   []*tokenbind,err---ok
 * flow  int64+dir     []*dir,err
 * netid int64         []*flow,err --- ok
 * netid []int64,avail,disp       []int64,err--transitionid
 * tranid []int64,isbatch  []*tokenbind,err---ok
 * tranid int64,avail  []*dir.err
 * biz    int64        []int64,err---netid ---ok
 *
 * one2one
 * ids    []int64,xxx,yy []*struct,err --- net/bind/flow/=ok
 ******************************************
 *  testing:
 *  1. get objs by ids
 *  2. get objs by obj.aggregate_field
 *  3. get obj ids by obj.aggregate_field
 */

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/smartystreets/goconvey/convey"
	"go-common/app/admin/main/aegis/model/net"
	"go-common/library/database/orm"
)

var (
	testaggrmap = map[int64]struct {
		Rela []int64
		JSON string
	}{
		1: {Rela: []int64{1, 2}, JSON: "1,2"},
		2: {Rela: []int64{3, 4}, JSON: "3,4"},
	}
	testaggrv     = []int64{1, 2}
	testaggrcache = map[string]string{
		"test_aggr:1": "1,2",
		"test_aggr:2": "3,4",
	}
	testaggrdest = &AggrCacheDest{
		GetKey: func(id int64, opt ...interface{}) (k []string) {
			k = []string{fmt.Sprintf("test_aggr:%d", id)}
			return
		},
		AggrRelaRaw: func(c context.Context, miss []int64, kopt []interface{}, opt ...interface{}) (rela map[int64][]int64, err error) {
			rela = map[int64][]int64{}
			for _, i := range miss {
				rela[i] = testaggrmap[i].Rela
			}
			return
		},
	}

	testflows = []*net.Flow{
		{ID: 50, Name: "init", NetID: 99},
		{ID: 51, Name: "first", NetID: 99},
		{ID: 52, Name: "init", NetID: 100},
	}
	testnets = []*net.Net{
		{ID: 8, ChName: "测试网1", BusinessID: 10},
		{ID: 9, ChName: "测试网1", BusinessID: 10},
		{ID: 10, ChName: "测试网1", BusinessID: 11},
	}
	testtokens = []*net.Token{
		{ID: 50, NetID: 10, Name: "test1", Compare: 0, Value: "0"},
		{ID: 51, NetID: 10, Name: "test2", Compare: 0, Value: "0"},
		{ID: 52, NetID: 11, Name: "test2", Compare: 0, Value: "0"},
	}
	testbinds = []*net.TokenBind{
		{ID: 50, Type: 1, ElementID: 10, TokenID: "8"},
		{ID: 51, Type: 1, ElementID: 10, TokenID: "9"},
		{ID: 52, Type: 2, ElementID: 11, TokenID: "8"},
	}
	testtrans = []*net.Transition{
		{ID: 50, NetID: 10, Trigger: 1, Limit: 1, Name: "test1"},
		{ID: 51, NetID: 10, Trigger: 1, Limit: 1, Name: "test2"},
		{ID: 52, NetID: 11, Trigger: 1, Limit: 1, Name: "test3"},
	}
	testdirs = []*net.Direction{
		{ID: 50, NetID: 10, FlowID: 10, TransitionID: 20, Direction: 1},
		{ID: 51, NetID: 10, FlowID: 11, TransitionID: 20, Direction: 2},
		{ID: 52, NetID: 10, FlowID: 11, TransitionID: 21, Direction: 1},
		{ID: 53, NetID: 10, FlowID: 12, TransitionID: 21, Direction: 2},
	}
	testdata = map[string]struct {
		RowObj        interface{}
		RowObjLen     int
		RowRela       map[int64][]int64
		Caches        map[string]string
		AggrID, ObjID []int64
	}{
		"flow": {
			RowObj:    testflows,
			RowObjLen: 3,
			RowRela:   map[int64][]int64{99: {50, 51}, 100: {52}},
			Caches: map[string]string{
				"net_flow:99":  "50,51",
				"net_flow:100": "52",
				"flow:50":      getjsonobj(testflows[0]),
				"flow:51":      getjsonobj(testflows[1]),
				"flow:52":      getjsonobj(testflows[2]),
			},
			AggrID: []int64{99, 100, 101},
			ObjID:  []int64{50, 51, 52},
		},
		"net": {
			RowObj:    testnets,
			RowObjLen: 3,
			RowRela:   map[int64][]int64{10: {8, 9}, 11: {10}},
			Caches: map[string]string{
				"business_net:10": "8,9",
				"business_net:11": "10",
				"net:8":           getjsonobj(testnets[0]),
				"net:9":           getjsonobj(testnets[1]),
				"net:10":          getjsonobj(testnets[2]),
			},
			AggrID: []int64{10, 11, 1200},
			ObjID:  []int64{8, 9, 10},
		},
		"token": {
			RowObj:    testtokens,
			RowObjLen: 3,
			Caches: map[string]string{
				"token:50": getjsonobj(testtokens[0]),
				"token:51": getjsonobj(testtokens[1]),
				"token:52": getjsonobj(testtokens[2]),
			},
			ObjID: []int64{50, 51, 52},
		},
		"bind": {
			RowObj:    testbinds,
			RowObjLen: 3,
			RowRela:   map[int64][]int64{10: {50, 51}, 11: {52}},
			Caches: map[string]string{
				"flow_bind:10":       "50,51",
				"transition_bind:11": "52",
				"bind:50":            getjsonobj(testbinds[0]),
				"bind:51":            getjsonobj(testbinds[1]),
				"bind:52":            getjsonobj(testbinds[2]),
			},
			AggrID: []int64{10, 11, 1200},
			ObjID:  []int64{50, 51, 52},
		},
		"tran": {
			RowObj:    testtrans,
			RowObjLen: 3,
			RowRela:   map[int64][]int64{10: {50, 51}, 11: {52}},
			Caches: map[string]string{
				"net_transition:10": "50,51",
				"net_transition:11": "52",
				"transition:50":     getjsonobj(testtrans[0]),
				"transition:51":     getjsonobj(testtrans[1]),
				"transition:52":     getjsonobj(testtrans[2]),
			},
			AggrID: []int64{10, 11, 12000},
			ObjID:  []int64{50, 51, 52},
		},
	}
	testdata2 = map[string]struct {
		RowObj    interface{}
		RowObjLen int
		RowRela   map[int64]map[int8][]int64
		Caches    map[string]string
		ObjID     []int64
		AggrID    []int64
	}{
		"dir": {
			RowObj:    testdirs,
			RowObjLen: 4,
			RowRela: map[int64]map[int8][]int64{
				10: {1: {50}},
				11: {1: {52}, 2: {51}},
				12: {2: {53}},
				20: {1: {50}, 2: {51}},
				21: {1: {52}, 2: {53}},
			},
			Caches: map[string]string{
				"flow_direction_1:10":       "50",
				"flow_direction_1:11":       "52",
				"flow_direction_2:11":       "51",
				"flow_direction_2:12":       "53",
				"transition_direction_1:20": "50",
				"transition_direction_2:20": "51",
				"transition_direction_1:21": "52",
				"transition_direction_2:21": "53",
				"direction:50":              getjsonobj(testdirs[0]),
				"direction:51":              getjsonobj(testdirs[1]),
				"direction:52":              getjsonobj(testdirs[2]),
				"direction:53":              getjsonobj(testdirs[3]),
			},
			AggrID: []int64{10, 11, 12, 20, 21, 12000},
			ObjID:  []int64{50, 51, 52, 53},
		},
	}

	getjsonobj = func(v interface{}) string {
		bs, _ := json.Marshal(v)
		return string(bs)
	}
)

func TestServiceAppendString(t *testing.T) {
	convey.Convey("AppendString", t, func(ctx convey.C) {
		testaggrdest.Dest = map[int64][]int64{}
		ctx.Convey("err=nil, ok=true, aggrdest expected", func(ctx convey.C) {
			for id, item := range testaggrmap {
				ok, err := testaggrdest.AppendString(id, item.JSON)
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(ok, convey.ShouldEqual, true)
				ctx.So(fmt.Sprintf("%v", testaggrdest.Dest[id]), convey.ShouldEqual, fmt.Sprintf("%v", item.Rela))
			}
		})
	})
}

func TestServiceAppendRelationRaw(t *testing.T) {
	convey.Convey("AppendRelationRaw", t, func(ctx convey.C) {
		missCache, err := testaggrdest.AppendRelationRaw(cntx, testaggrv, nil)
		ctx.Convey("Then err should be nil.missCache should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(missCache, convey.ShouldNotBeNil)
			for k, v := range missCache {
				ctx.So(v, convey.ShouldEqual, testaggrcache[k])
			}
		})
	})
}

func testdataprepare(k string, v interface{}) {
	testdao := orm.NewMySQL(s.c.ORM)
	objs := reflect.ValueOf(v)
	for i := 0; i < objs.Len(); i++ {
		testdao.Save(objs.Index(i).Interface())
	}
}

func TestServiceAppendAggregateRaw(t *testing.T) {
	convey.Convey("AppendAggregateRaw", t, func(ctx convey.C) {
		objdest := &flowArr{}
		w := &CacheWrap{
			ObjDest:    objdest,
			AggrObjRaw: s.netFlowAppendRaw,
			AggregateDest: &AggrCacheDest{
				GetKey:      s.netFlowKey,
				AggrRelaRaw: s.netFlowRelation,
			},
		}

		k := "flow"
		testdataprepare(k, testdata[k].RowObj)
		missCache, err := w.AppendAggregateRaw(cntx, testdata[k].AggrID, nil)
		ctx.Convey("Then err should be nil.missCache should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(missCache, convey.ShouldNotBeNil)
			//对象存在 & 聚合关系存在 & 待缓存值存在
			ctx.So(len(objdest.dest), convey.ShouldEqual, testdata[k].RowObjLen)
			for ck, v := range testdata[k].RowRela {
				ctx.So(fmt.Sprintf("%v", w.AggregateDest.Dest[ck]), convey.ShouldEqual, fmt.Sprintf("%v", v))
			}
			for ck, v := range testdata[k].Caches {
				pos := strings.Index(missCache[ck], "ctime") //json化的
				if pos == -1 {
					ctx.So(missCache[ck], convey.ShouldEqual, v)
				} else {
					ctx.So(missCache[ck][:pos], convey.ShouldEqual, v[:pos])
				}
			}
		})
	})
}

func TestServiceGetFromCache(t *testing.T) {
	convey.Convey("getFromCache", t, func(ctx convey.C) {
		k := "flow"
		//delete cache
		keypro := s.netFlowKey
		llen := len(testdata[k].AggrID) + 2
		keys := []string{}
		ids := append(testdata[k].AggrID, 1000, 2000)
		for _, id := range ids {
			keys = append(keys, keypro(id)...)
		}
		s.redis.DelMulti(cntx, keys...)
		//miss all
		miss := s.getFromCache(cntx, ids, keypro, &flowArr{}, nil)

		//set last 2 caches & get from cache
		s.redis.SetMulti(cntx, map[string]string{keys[llen-1]: "10000", keys[llen-2]: "20000"})
		miss2 := s.getFromCache(cntx, ids, keypro, &AggrCacheDest{}, nil)

		ctx.Convey("Then miss should not be nil.", func(ctx convey.C) {
			ctx.So(miss, convey.ShouldNotBeNil)
			ctx.So(fmt.Sprintf("%v", miss), convey.ShouldEqual, fmt.Sprintf("%v", ids))
			ctx.So(fmt.Sprintf("%v", miss2), convey.ShouldEqual, fmt.Sprintf("%v", testdata[k].AggrID))
		})
	})
}

func TestServiceobjCache(t *testing.T) {
	convey.Convey("objCache", t, func(ctx convey.C) {
		k := "flow"
		ids := testdata[k].ObjID
		objdest := &flowArr{}
		err := s.objCache(cntx, ids, objdest, nil)
		objm := map[int64]*net.Flow{}
		for _, item := range objdest.dest {
			objm[item.ID] = item
		}
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			for _, id := range ids {
				ctx.So(objm[id], convey.ShouldNotBeNil)
			}
		})
	})
}

func TestServiceaggregateRelationCache(t *testing.T) {
	convey.Convey("aggregateRelationCache", t, func(ctx convey.C) {
		k := "flow"
		aggrdest := &AggrCacheDest{
			GetKey:      s.netFlowKey,
			AggrRelaRaw: s.netFlowRelation,
		}
		err := s.aggregateRelationCache(cntx, testdata[k].AggrID, aggrdest, nil)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			for nid, v := range testdata[k].RowRela {
				ctx.So(len(aggrdest.Dest[nid]), convey.ShouldEqual, len(v))
				mmp := map[int64]int64{}
				for _, i := range v {
					mmp[i] = i
				}
				for _, i := range aggrdest.Dest[nid] {
					ctx.So(i, convey.ShouldEqual, mmp[i])
				}
			}
		})
	})
}

func TestServiceaggregateCache(t *testing.T) {
	convey.Convey("aggregateCache", t, func(ctx convey.C) {
		k := "flow"
		objdest := &flowArr{}
		w := &CacheWrap{
			ObjDest:    objdest,
			AggrObjRaw: s.netFlowAppendRaw,
			AggregateDest: &AggrCacheDest{
				GetKey: s.netFlowKey,
			},
		}
		err := s.aggregateCache(cntx, testdata[k].AggrID, w, nil, nil, nil, nil)
		objm := map[int64]*net.Flow{}
		for _, item := range objdest.dest {
			objm[item.ID] = item
		}

		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			for id, v := range testdata[k].RowRela {
				ctx.So(len(w.AggregateDest.Dest[id]), convey.ShouldEqual, len(v))
				mmp := map[int64]int64{}
				for _, i := range v {
					mmp[i] = i
				}
				for _, i := range w.AggregateDest.Dest[id] {
					ctx.So(i, convey.ShouldEqual, mmp[i])
				}
			}
			for _, id := range testdata[k].ObjID {
				ctx.So(objm[id], convey.ShouldNotBeNil)
			}
		})
	})
}

func TestServicedelCache(t *testing.T) {
	convey.Convey("delCache", t, func(ctx convey.C) {
		err := s.delCache(cntx, []string{})
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestNetArrGetKey(t *testing.T) {
	convey.Convey("netArrGetKey", t, func(ctx convey.C) {
		objdest := &netArr{}
		p1 := objdest.GetKey(0)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestNetArrAppendString(t *testing.T) {
	convey.Convey("netArrAppendString", t, func(ctx convey.C) {
		objdest := &netArr{}
		k := "net"
		testdataprepare(k, testdata[k].RowObj)
		v := testdata[k].Caches[objdest.GetKey(testdata[k].ObjID[0])[0]]
		ok, err := objdest.AppendString(testdata[k].ObjID[0], v)
		ctx.Convey("Then err should be nil.ok should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(ok, convey.ShouldEqual, true)
			ctx.So(objdest.dest[0].ID, convey.ShouldEqual, testdata[k].ObjID[0])
		})
	})
}

func TestNetArrAppendRaw(t *testing.T) {
	convey.Convey("netArrAppendRaw", t, func(ctx convey.C) {
		objdest := &netArr{}
		k := "net"
		missCache, err := objdest.AppendRaw(cntx, s, testdata[k].ObjID, nil)
		objm := map[int64]*net.Net{}
		for _, item := range objdest.dest {
			objm[item.ID] = item
		}
		ctx.Convey("Then err should be nil.missCache should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			for ck, v := range testdata[k].Caches {
				if strings.Contains(ck, "business_net") {
					continue
				}
				pos := strings.Index(v, "ctime")
				if pos == -1 {
					ctx.So(fmt.Sprintf("%v", missCache[ck]), convey.ShouldEqual, fmt.Sprintf("%v", v))
				} else {
					ctx.So(fmt.Sprintf("%v", missCache[ck][:pos]), convey.ShouldEqual, fmt.Sprintf("%v", v[:pos]))
				}
			}
			for _, id := range testdata[k].ObjID {
				ctx.So(objm[id], convey.ShouldNotBeNil)
			}
		})
	})
}

func TestNetCache(t *testing.T) {
	convey.Convey("NetCache", t, func(ctx convey.C) {
		k := "net"
		biz := testdata[k].AggrID[0]
		res, err := s.netIDByBusiness(cntx, biz)
		ctx.So(err, convey.ShouldBeNil)
		ctx.So(fmt.Sprintf("%v", res), convey.ShouldEqual, fmt.Sprintf("%v", testdata[k].RowRela[biz]))

		id := testdata[k].ObjID[0]
		res2, err := s.netByID(cntx, id)
		ctx.So(err, convey.ShouldBeNil)
		ctx.So(res2.ID, convey.ShouldEqual, id)

		err = s.delNetCache(cntx, res2)
		ctx.So(err, convey.ShouldBeNil)

		res3, err := s.netIDByBusiness(cntx, biz)
		ctx.So(err, convey.ShouldBeNil)
		ctx.So(fmt.Sprintf("%v", res3), convey.ShouldEqual, fmt.Sprintf("%v", testdata[k].RowRela[biz]))

		res4, err := s.netByID(cntx, id)
		ctx.So(err, convey.ShouldBeNil)
		ctx.So(res4.ID, convey.ShouldEqual, id)
	})
}

func TestTokenArrGetKey(t *testing.T) {
	convey.Convey("TokenArrGetKey", t, func(ctx convey.C) {
		objdest := &tokenArr{}
		p1 := objdest.GetKey(0)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestTokenArrAppendString(t *testing.T) {
	convey.Convey("TokenArrAppendString", t, func(ctx convey.C) {
		objdest := &tokenArr{}
		k := "token"
		testdataprepare(k, testdata[k].RowObj)
		v := testdata[k].Caches[objdest.GetKey(testdata[k].ObjID[0])[0]]
		ok, err := objdest.AppendString(testdata[k].ObjID[0], v)
		ctx.Convey("Then err should be nil.ok should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(ok, convey.ShouldEqual, true)
			ctx.So(objdest.dest[0].ID, convey.ShouldEqual, testdata[k].ObjID[0])
		})
	})
}

func TestTokenArrAppendRaw(t *testing.T) {
	convey.Convey("TokenArrAppendRaw", t, func(ctx convey.C) {
		objdest := &tokenArr{}
		k := "token"
		missCache, err := objdest.AppendRaw(cntx, s, testdata[k].ObjID, nil)
		objm := map[int64]*net.Token{}
		for _, item := range objdest.dest {
			objm[item.ID] = item
		}
		ctx.Convey("Then err should be nil.missCache should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			for ck, v := range testdata[k].Caches {
				pos := strings.Index(v, "ctime")
				if pos == -1 {
					ctx.So(fmt.Sprintf("%v", missCache[ck]), convey.ShouldEqual, fmt.Sprintf("%v", v))
				} else {
					ctx.So(fmt.Sprintf("%v", missCache[ck][:pos]), convey.ShouldEqual, fmt.Sprintf("%v", v[:pos]))
				}
			}
			for _, id := range testdata[k].ObjID {
				ctx.So(objm[id], convey.ShouldNotBeNil)
			}
		})
	})
}

func TestTokenCache(t *testing.T) {
	convey.Convey("TokenCache", t, func(ctx convey.C) {
		k := "token"
		biz := testdata[k].ObjID
		res, err := s.tokens(cntx, biz)
		ctx.So(err, convey.ShouldBeNil)
		resm := map[int64]*net.Token{}
		for _, item := range res {
			resm[item.ID] = item
		}
		for _, i := range biz {
			ctx.So(resm[i], convey.ShouldNotBeNil)
		}

		keys := []string{}
		for ck := range testdata[k].Caches {
			keys = append(keys, ck)
		}
		err = s.redis.DelMulti(cntx, keys...)
		ctx.So(err, convey.ShouldBeNil)

		res2, err := s.tokens(cntx, biz)
		ctx.So(err, convey.ShouldBeNil)
		resm2 := map[int64]*net.Token{}
		for _, item := range res2 {
			resm2[item.ID] = item
		}
		for _, i := range biz {
			ctx.So(resm2[i], convey.ShouldNotBeNil)
		}
	})
}

func TestBindArrGetKey(t *testing.T) {
	convey.Convey("BindArrGetKey", t, func(ctx convey.C) {
		objdest := &bindArr{}
		p1 := objdest.GetKey(0)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestBindArrAppendString(t *testing.T) {
	convey.Convey("BindArrAppendString", t, func(ctx convey.C) {
		objdest := &bindArr{}
		k := "bind"
		testdataprepare(k, testdata[k].RowObj)
		v := testdata[k].Caches[objdest.GetKey(testdata[k].ObjID[0])[0]]
		ok, err := objdest.AppendString(testdata[k].ObjID[0], v)
		ctx.Convey("Then err should be nil.ok should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(ok, convey.ShouldEqual, true)
			ctx.So(objdest.dest[0].ID, convey.ShouldEqual, testdata[k].ObjID[0])
		})
	})
}

func TestBindArrAppendRaw(t *testing.T) {
	convey.Convey("BindArrAppendRaw", t, func(ctx convey.C) {
		objdest := &bindArr{}
		k := "bind"
		missCache, err := objdest.AppendRaw(cntx, s, testdata[k].ObjID, nil)
		objm := map[int64]*net.TokenBind{}
		for _, item := range objdest.dest {
			objm[item.ID] = item
		}
		ctx.Convey("Then err should be nil.missCache should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			for ck, v := range testdata[k].Caches {
				if strings.Contains(ck, "flow_bind") || strings.Contains(ck, "transition_bind") {
					continue
				}
				pos := strings.Index(v, "ctime")
				if pos == -1 {
					ctx.So(fmt.Sprintf("%v", missCache[ck]), convey.ShouldEqual, fmt.Sprintf("%v", v))
				} else {
					ctx.So(fmt.Sprintf("%v", missCache[ck][:pos]), convey.ShouldEqual, fmt.Sprintf("%v", v[:pos]))
				}
			}
			for _, id := range testdata[k].ObjID {
				ctx.So(objm[id], convey.ShouldNotBeNil)
			}
		})
	})
}

func TestBindCache(t *testing.T) {
	convey.Convey("BindCache", t, func(ctx convey.C) {
		k := "bind"
		biz := testdata[k].AggrID
		res, err := s.tokenBindByElement(cntx, biz, []int8{net.BindTypeFlow})
		resm := map[int64]*net.TokenBind{}
		for _, item := range res {
			t.Logf("flow bind tokens in biz(%+v)  item(%+v)", biz, item)
			resm[item.ID] = item
		}
		ctx.So(err, convey.ShouldBeNil)
		for _, i := range testdata[k].ObjID[:2] {
			ctx.So(resm[i], convey.ShouldNotBeNil)
		}

		res, err = s.tokenBindByElement(cntx, biz, net.BindTranType)
		for _, item := range res {
			resm[item.ID] = item
		}
		ctx.So(err, convey.ShouldBeNil)

		res2, err := s.tokenBinds(cntx, testdata[k].ObjID, true)
		resm2 := map[int64]*net.TokenBind{}
		for _, item := range res2 {
			resm2[item.ID] = item
		}
		ctx.So(err, convey.ShouldBeNil)
		for _, i := range testdata[k].ObjID {
			ctx.So(resm2[i], convey.ShouldNotBeNil)
		}
	})
}

func TestFlowArrGetKey(t *testing.T) {
	convey.Convey("FlowArrGetKey", t, func(ctx convey.C) {
		objdest := &flowArr{}
		p1 := objdest.GetKey(0)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestFlowArrAppendString(t *testing.T) {
	convey.Convey("FlowArrAppendString", t, func(ctx convey.C) {
		objdest := &flowArr{}
		k := "flow"
		testdataprepare(k, testdata[k].RowObj)
		v := testdata[k].Caches[objdest.GetKey(testdata[k].ObjID[0])[0]]
		ok, err := objdest.AppendString(testdata[k].ObjID[0], v)
		ctx.Convey("Then err should be nil.ok should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(ok, convey.ShouldEqual, true)
			ctx.So(objdest.dest[0].ID, convey.ShouldEqual, testdata[k].ObjID[0])
		})
	})
}

func TestFlowArrAppendRaw(t *testing.T) {
	convey.Convey("FlowArrAppendRaw", t, func(ctx convey.C) {
		objdest := &flowArr{}
		k := "flow"
		missCache, err := objdest.AppendRaw(cntx, s, testdata[k].ObjID, nil)
		objm := map[int64]*net.Flow{}
		for _, item := range objdest.dest {
			objm[item.ID] = item
		}
		ctx.Convey("Then err should be nil.missCache should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			for ck, v := range testdata[k].Caches {
				if strings.Contains(ck, "net_flow") {
					continue
				}
				pos := strings.Index(v, "ctime")
				if pos == -1 {
					ctx.So(fmt.Sprintf("%v", missCache[ck]), convey.ShouldEqual, fmt.Sprintf("%v", v))
				} else {
					ctx.So(fmt.Sprintf("%v", missCache[ck][:pos]), convey.ShouldEqual, fmt.Sprintf("%v", v[:pos]))
				}
			}
			for _, id := range testdata[k].ObjID {
				ctx.So(objm[id], convey.ShouldNotBeNil)
			}
		})
	})
}

func TestFlowCache(t *testing.T) {
	convey.Convey("FlowCache", t, func(ctx convey.C) {
		k := "flow"
		biz := testdata[k].AggrID[0]
		res, err := s.flowsByNet(cntx, biz)
		resm := map[int64]*net.Flow{}
		for _, item := range res {
			resm[item.ID] = item
		}
		ctx.So(err, convey.ShouldBeNil)
		for _, i := range testdata[k].RowRela[biz] {
			ctx.So(resm[i], convey.ShouldNotBeNil)
		}

		res2, err := s.flowIDByNet(cntx, biz)
		ctx.So(err, convey.ShouldBeNil)
		mmp := map[int64]int64{}
		for _, i := range testdata[k].RowRela[biz] {
			mmp[i] = i
		}
		ctx.So(len(res2), convey.ShouldEqual, len(testdata[k].RowRela[biz]))
		for _, i := range res2 {
			ctx.So(i, convey.ShouldEqual, mmp[i])
		}

		res3, err := s.flows(cntx, testdata[k].ObjID, true)
		resm2 := map[int64]*net.Flow{}
		for _, item := range res3 {
			resm2[item.ID] = item
		}
		ctx.So(err, convey.ShouldBeNil)
		for _, i := range testdata[k].ObjID {
			ctx.So(resm2[i], convey.ShouldNotBeNil)
		}
	})
}

func TestTranArrGetKey(t *testing.T) {
	convey.Convey("TranArrGetKey", t, func(ctx convey.C) {
		objdest := &tranArr{}
		p1 := objdest.GetKey(0)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestTranArrAppendString(t *testing.T) {
	convey.Convey("TranArrAppendString", t, func(ctx convey.C) {
		objdest := &tranArr{}
		k := "tran"
		testdataprepare(k, testdata[k].RowObj)
		v := testdata[k].Caches[objdest.GetKey(testdata[k].ObjID[0])[0]]
		ok, err := objdest.AppendString(testdata[k].ObjID[0], v)
		ctx.Convey("Then err should be nil.ok should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(ok, convey.ShouldEqual, true)
			ctx.So(objdest.dest[0].ID, convey.ShouldEqual, testdata[k].ObjID[0])
		})
	})
}

func TestTranArrAppendRaw(t *testing.T) {
	convey.Convey("TranArrAppendRaw", t, func(ctx convey.C) {
		objdest := &tranArr{}
		k := "tran"
		_, err := objdest.AppendRaw(cntx, s, testdata[k].ObjID, nil)
		objm := map[int64]*net.Transition{}
		for _, item := range objdest.dest {
			objm[item.ID] = item
		}
		ctx.Convey("Then err should be nil.missCache should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			for _, id := range testdata[k].ObjID {
				ctx.So(objm[id], convey.ShouldNotBeNil)
			}
		})
	})
}

func TestTranCache(t *testing.T) {
	convey.Convey("TranCache", t, func(ctx convey.C) {
		k := "tran"
		biz := testdata[k].AggrID[0]
		res2, err := s.tranIDByNet(cntx, []int64{biz}, true, true)
		ctx.So(err, convey.ShouldBeNil)
		ctx.So(fmt.Sprintf("%v", res2), convey.ShouldEqual, fmt.Sprintf("%v", testdata[k].RowRela[biz]))

		res3, err := s.transitions(cntx, testdata[k].ObjID, true)
		resm2 := map[int64]*net.Transition{}
		for _, item := range res3 {
			resm2[item.ID] = item
		}
		ctx.So(err, convey.ShouldBeNil)
		for _, i := range testdata[k].ObjID {
			ctx.So(resm2[i], convey.ShouldNotBeNil)
		}
	})
}

func TestDirArrGetKey(t *testing.T) {
	convey.Convey("dirArrGetKey", t, func(ctx convey.C) {
		objdest := &dirArr{}
		p1 := objdest.GetKey(0)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDirArrAppendString(t *testing.T) {
	convey.Convey("dirArrAppendString", t, func(ctx convey.C) {
		objdest := &dirArr{}
		k := "dir"
		testdataprepare(k, testdata2[k].RowObj)
		v := testdata2[k].Caches[objdest.GetKey(testdata2[k].ObjID[0])[0]]
		ok, err := objdest.AppendString(testdata2[k].ObjID[0], v)
		ctx.Convey("Then err should be nil.ok should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(ok, convey.ShouldEqual, true)
			ctx.So(objdest.dest[0].ID, convey.ShouldEqual, testdata2[k].ObjID[0])
		})
	})
}

func TestDirArrAppendRaw(t *testing.T) {
	convey.Convey("dirArrAppendRaw", t, func(ctx convey.C) {
		objdest := &dirArr{}
		k := "dir"
		missCache, err := objdest.AppendRaw(cntx, s, testdata2[k].ObjID, nil)
		objm := map[int64]*net.Direction{}
		for _, item := range objdest.dest {
			objm[item.ID] = item
		}
		ctx.Convey("Then err should be nil.missCache should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			for ck, v := range testdata2[k].Caches {
				if strings.Contains(ck, "flow_direction") || strings.Contains(ck, "transition_direction") {
					continue
				}
				pos := strings.Index(v, "ctime")
				if pos == -1 {
					ctx.So(fmt.Sprintf("%v", missCache[ck]), convey.ShouldEqual, fmt.Sprintf("%v", v))
				} else {
					ctx.So(fmt.Sprintf("%v", missCache[ck][:pos]), convey.ShouldEqual, fmt.Sprintf("%v", v[:pos]))
				}
			}
			for _, id := range testdata2[k].ObjID {
				ctx.So(objm[id], convey.ShouldNotBeNil)
			}
		})
	})
}

func TestDirCache(t *testing.T) {
	convey.Convey("DirCache", t, func(ctx convey.C) {
		k := "dir"
		biz := testdata2[k].AggrID
		tp := net.DirInput
		res, err := s.dirByFlow(cntx, biz, tp)
		resm := map[int64]*net.Direction{}
		for _, item := range res {
			resm[item.ID] = item
		}
		ctx.So(err, convey.ShouldBeNil)
		expected := []int64{}
		for _, i := range biz {
			expected = append(expected, testdata2[k].RowRela[i][tp]...)
		}
		for _, i := range expected {
			ctx.So(resm[i], convey.ShouldNotBeNil)
		}

		res2, err := s.dirByTran(cntx, biz, net.DirOutput, true)
		for _, item := range res2 {
			resm[item.ID] = item
		}
		ctx.So(err, convey.ShouldBeNil)
		expected2 := []int64{}
		for _, i := range biz {
			expected2 = append(expected2, testdata2[k].RowRela[i][tp]...)
		}
		for _, i := range expected2 {
			ctx.So(resm[i], convey.ShouldNotBeNil)
		}

		res3, err := s.directions(cntx, testdata2[k].ObjID)
		resm2 := map[int64]*net.Direction{}
		for _, item := range res3 {
			resm2[item.ID] = item
		}
		ctx.So(err, convey.ShouldBeNil)
		for _, i := range testdata2[k].ObjID {
			ctx.So(resm2[i], convey.ShouldNotBeNil)
		}
	})
}
