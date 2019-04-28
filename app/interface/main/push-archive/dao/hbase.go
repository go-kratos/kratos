package dao

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"strconv"
	"sync"

	"go-common/app/interface/main/push-archive/model"
	"go-common/library/log"
	"go-common/library/sync/errgroup"

	"github.com/tsuna/gohbase/hrpc"
)

const _hbaseShard = 200

var (
	hbaseTable   = "ugc:PushArchive"
	hbaseFamily  = "relation"
	hbaseFamilyB = []byte(hbaseFamily)
)

func _rowKey(upper, fans int64) string {
	k := fmt.Sprintf("%d_%d", upper, fans%_hbaseShard)
	key := fmt.Sprintf("%x", md5.Sum([]byte(k)))
	return key
}

// Fans gets the upper's fans.
func (d *Dao) Fans(c context.Context, upper int64, isPGC bool) (res map[int64]int, err error) {
	var mutex sync.Mutex
	res = make(map[int64]int)
	group := errgroup.Group{}
	for i := 0; i < _hbaseShard; i++ {
		shard := int64(i)
		group.Go(func() (e error) {
			key := _rowKey(upper, shard)
			relations, e := d.fansByKey(context.TODO(), key)
			if e != nil {
				return
			}
			mutex.Lock()
			for fans, tp := range relations {
				// pgc稿件，屏蔽非特殊关注粉丝
				if isPGC && tp != model.RelationSpecial {
					continue
				}
				res[fans] = tp
			}
			mutex.Unlock()
			return
		})
	}
	group.Wait()
	return
}

// AddFans add upper's fans.
func (d *Dao) AddFans(c context.Context, upper, fans int64, tp int) (err error) {
	key := _rowKey(upper, fans)
	relations, err := d.fansByKey(c, key)
	if err != nil {
		return
	}
	relations[fans] = tp
	err = d.saveRelation(c, key, upper, relations)
	return
}

// DelFans del fans.
func (d *Dao) DelFans(c context.Context, upper, fans int64) (err error) {
	key := _rowKey(upper, fans)
	relations, err := d.fansByKey(c, key)
	if err != nil {
		return
	}
	delete(relations, fans)
	err = d.saveRelation(c, key, upper, relations)
	return
}

// DelSpecialAttention del special attention.
func (d *Dao) DelSpecialAttention(c context.Context, upper, fans int64) (err error) {
	key := _rowKey(upper, fans)
	relations, err := d.fansByKey(c, key)
	if err != nil {
		return
	}
	if relations[fans] != model.RelationSpecial {
		return
	}
	relations[fans] = model.RelationAttention
	err = d.saveRelation(c, key, upper, relations)
	return
}

func (d *Dao) fansByKey(c context.Context, key string) (relations map[int64]int, err error) {
	var (
		result      *hrpc.Result
		ctx, cancel = context.WithTimeout(c, d.relationHBaseReadTimeout)
	)
	defer cancel()
	relations = make(map[int64]int)

	if result, err = d.relationHBase.Get(ctx, []byte(hbaseTable), []byte(key)); err != nil {
		log.Error("d.relationHBase.Get error(%v) querytable(%v)", err, hbaseTable)
		PromError("hbase:Get")
		return
	} else if result == nil {
		return
	}
	for _, c := range result.Cells {
		if c != nil && bytes.Equal(c.Family, hbaseFamilyB) {
			if err = json.Unmarshal(c.Value, &relations); err != nil {
				log.Error("json.Unmarshal() error(%v)", err)
				return
			}
			break
		}
	}
	return
}

func (d *Dao) saveRelation(c context.Context, key string, upper int64, relations map[int64]int) (err error) {
	var (
		column      = strconv.FormatInt(upper, 10)
		ctx, cancel = context.WithTimeout(c, d.relationHBaseWriteTimeout)
	)
	defer cancel()
	value, err := json.Marshal(relations)
	if err != nil {
		return
	}
	values := map[string]map[string][]byte{hbaseFamily: {column: value}}
	if _, err = d.relationHBase.PutStr(ctx, hbaseTable, key, values); err != nil {
		log.Error("d.relationHBase.PutStr error(%v), table(%s), values(%+v)", err, hbaseTable, values)
		PromError("hbase:Put")
	}
	return
}

// filterFanByUpper 根据fans在hbase存储的up主列表，筛选出upper主在up主列表中的粉丝
func (d *Dao) filterFanByUpper(c context.Context, fan int64, up interface{}, table string, family []string) (included bool, err error) {
	var (
		res         *hrpc.Result
		key         string
		ctx, cancel = context.WithTimeout(c, d.fanHBaseReadTimeout)
	)
	defer cancel()
	upper := up.(int64)
	rowKeyMD := md5.Sum([]byte(strconv.FormatInt(fan, 10)))
	key = fmt.Sprintf("%x", rowKeyMD)
	if res, err = d.fanHBase.Get(ctx, []byte(table), []byte(key)); err != nil {
		log.Error("d.fanHBase.Get error(%v) querytable(%v) key(%s), fan(%d), upper(%d)", err, table, key, fan, upper)
		PromError("hbase:Get")
		return
	} else if res == nil {
		return
	}
	for _, c := range res.Cells {
		if c == nil || !existFamily(c.Family, family) {
			continue
		}
		upID := int64(binary.BigEndian.Uint32(c.Value))
		if upID != upper || upID <= 0 {
			continue
		}
		included = true
		log.Info("filter fan: included by hbase, fan(%d) upper(%d) table(%s)", fan, upper, table)
		return
	}
	if !included {
		log.Info("filter fan: excluded by hbase, fan(%d) upper(%d) table(%s)", fan, upper, table)
	}
	return
}

// FilterFans 批量筛选
func (d *Dao) FilterFans(fans *[]int64, params map[string]interface{}) (err error) {
	base := params["base"]
	table := params["table"].(string)
	family := params["family"].([]string)
	result := params["result"].(*[]int64)
	excluded := params["excluded"].(*[]int64)
	handler := params["handler"].(func(context.Context, int64, interface{}, string, []string) (bool, error))
	mutex := sync.Mutex{}
	group := errgroup.Group{}
	l := len(*fans)
	for i := 0; i < l; i++ {
		shared := (*fans)[i]
		group.Go(func() (e error) {
			included, e := handler(context.TODO(), shared, base, table, family)
			if e != nil {
				log.Error("FilterFans error(%v) fan(%d) base(%d) table(%s) family(%v)", e, shared, base, table, family)
			}
			mutex.Lock()
			if included {
				*result = append(*result, shared)
			} else {
				*excluded = append(*excluded, shared)
			}
			mutex.Unlock()
			return
		})
	}
	group.Wait()
	return
}

// existFamily 某个hbase列族是否存在于指定列族中
func existFamily(actual []byte, family []string) bool {
	for _, f := range family {
		if bytes.Equal(actual, []byte(f)) {
			return true
		}
	}
	return false
}

// filterFanByActive 根据用户的活跃时间段，过滤不在活跃期内更新的粉丝; 若无活跃列表，从默认活跃时间内过滤
func (d *Dao) filterFanByActive(ctx context.Context, fan int64, oneHour interface{}, table string, family []string) (included bool, err error) {
	var (
		b          []byte
		result     *hrpc.Result
		c, cancel  = context.WithTimeout(ctx, d.fanHBaseReadTimeout)
		activeHour int
	)
	defer cancel()
	hour := oneHour.(int)
	if _, included = d.ActiveDefaultTime[hour]; included {
		return
	}
	rowKey := md5.Sum(strconv.AppendInt(b, fan, 10))
	key := fmt.Sprintf("%x", rowKey)
	if result, err = d.fanHBase.Get(c, []byte(table), []byte(key)); err != nil {
		log.Error("filterFanByActive d.fanHBase.Get error(%v) table(%s) key(%s) fan(%d)", err, table, key, fan)
		PromError("hbase:Get")
		return
	} else if result == nil {
		return
	}
	included = false
	for _, cell := range result.Cells {
		if cell != nil && existFamily(cell.Family, family) {
			activeHour, err = strconv.Atoi(string(cell.Value))
			if err != nil {
				log.Error("filterFanByActive strconv.Atoi error(%v) fan(%d) value(%s)", err, fan, string(cell.Value))
				break
			}
			if activeHour == hour {
				included = true
				break
			}
		}
	}
	if !included {
		log.Info("filter fan：excluded by active time from table, fan(%d)", fan)
	}
	return
}

// ExistsInBlacklist 按黑名单过滤用户
func (d *Dao) ExistsInBlacklist(ctx context.Context, upper int64, mids []int64) (exists, notExists []int64) {
	var (
		mutex sync.Mutex
		group = errgroup.Group{}
	)
	for _, mid := range mids {
		mid := mid
		group.Go(func() error {
			include, _ := d.filterFanByUpper(context.Background(), mid, upper, d.c.Abtest.HbaseBlacklistTable, d.c.Abtest.HbaseBlacklistFamily)
			mutex.Lock()
			if include {
				exists = append(exists, mid)
			} else {
				notExists = append(notExists, mid)
			}
			mutex.Unlock()
			return nil
		})
	}
	group.Wait()
	return
}

// ExistsInWhitelist 按白名单过滤用户
func (d *Dao) ExistsInWhitelist(ctx context.Context, upper int64, mids []int64) (exists, notExists []int64) {
	var (
		mutex sync.Mutex
		group = errgroup.Group{}
	)
	for _, mid := range mids {
		mid := mid
		group.Go(func() error {
			include, _ := d.filterFanByUpper(context.Background(), mid, upper, d.c.Abtest.HbaseeWhitelistTable, d.c.Abtest.HbaseWhitelistFamily)
			mutex.Lock()
			if include {
				exists = append(exists, mid)
			} else {
				notExists = append(notExists, mid)
			}
			mutex.Unlock()
			return nil
		})
	}
	group.Wait()
	return
}
