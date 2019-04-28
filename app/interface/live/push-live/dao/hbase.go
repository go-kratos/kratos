package dao

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/tsuna/gohbase/hrpc"
	"go-common/app/interface/live/push-live/model"
	"go-common/library/log"
	"go-common/library/sync/errgroup"
	"sync"
)

const _hbaseShard = 200

var (
	hbaseTable   = "ugc:PushArchive"
	hbaseFamily  = "relation"
	hbaseFamilyB = []byte(hbaseFamily)
)

// Fans gets the upper's fans.
func (d *Dao) Fans(c context.Context, upper int64, types int) (fans map[int64]bool, fansSP map[int64]bool, err error) {
	var mutex sync.Mutex
	fans = make(map[int64]bool)
	fansSP = make(map[int64]bool)
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
			for fansID, fansType := range relations {
				switch types {
				// 返回普通关注
				case model.RelationAttention:
					if fansType == types {
						fans[fansID] = true
					}
				// 返回特别关注
				case model.RelationSpecial:
					if fansType == types {
						fansSP[fansID] = true
					}
				// 同时返回普通关注与特别关注
				case model.RelationAll:
					if fansType == model.RelationSpecial {
						fansSP[fansID] = true
					} else if fansType == model.RelationAttention {
						fans[fansID] = true
					}
				default:
					return
				}
			}
			mutex.Unlock()
			return
		})
	}
	group.Wait()
	return
}

// SeparateFans Separate the upper's fans by 1 or 2.
func (d *Dao) SeparateFans(c context.Context, upper int64, fansIn map[int64]bool) (fans map[int64]bool, fansSP map[int64]bool, err error) {
	var mutex sync.Mutex
	special := make(map[int64]bool)
	fans = make(map[int64]bool)
	fansSP = make(map[int64]bool)
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
			// 获取所有特别关注
			for fansID, fansType := range relations {
				if fansType == model.RelationSpecial {
					special[fansID] = true
				}
			}
			mutex.Unlock()
			return
		})
	}
	group.Wait()
	for id := range fansIn {
		if _, ok := special[id]; ok {
			// 特别关注
			fansSP[id] = true
		} else {
			// 不是特别关注就是普通关注
			fans[id] = true
		}
	}
	return
}

func _rowKey(upper, fans int64) string {
	k := fmt.Sprintf("%d_%d", upper, fans%_hbaseShard)
	key := fmt.Sprintf("%x", md5.Sum([]byte(k)))
	return key
}

func (d *Dao) fansByKey(c context.Context, key string) (relations map[int64]int, err error) {
	relations = make(map[int64]int)
	var (
		query       *hrpc.Get
		result      *hrpc.Result
		ctx, cancel = context.WithTimeout(c, d.relationHBaseReadTimeout)
	)
	defer cancel()
	if result, err = d.relationHBase.GetStr(ctx, hbaseTable, key); err != nil {
		log.Error("d.relationHBase.Get error(%v) querytable(%v)", err, string(query.Table()))
		// PromError("hbase:Get")
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
