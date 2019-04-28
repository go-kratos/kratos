package toview

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/binary"
	"fmt"
	"strconv"
	"time"

	"go-common/app/interface/main/history/model"
	"go-common/library/log"

	"github.com/tsuna/gohbase/hrpc"
)

var (
	tableInfo  = "ugc:history_to_view"
	familyAid  = "aid"
	familyAidB = []byte(familyAid)
)

// hashRowKey create rowkey(md5(mid)[:2]+mid) for histroy by mid .
func hashRowKey(mid int64) string {
	var bs = make([]byte, 8)
	binary.LittleEndian.PutUint64(bs, uint64(mid))
	rk := md5.Sum(bs)
	return fmt.Sprintf("%x%d", rk[:2], mid)
}

// Add add one toview.
func (d *Dao) Add(ctx context.Context, mid, aid, now int64) (err error) {
	var (
		timeB  = make([]byte, 8)
		key    = hashRowKey(mid)
		column = strconv.FormatInt(aid, 10)
	)
	ctx, cancel := context.WithTimeout(ctx, time.Duration(d.conf.Info.WriteTimeout))
	defer cancel()
	binary.BigEndian.PutUint64(timeB, uint64(now))
	values := map[string]map[string][]byte{familyAid: map[string][]byte{column: timeB}}
	if _, err = d.info.PutStr(ctx, tableInfo, key, values); err != nil {
		log.Error("toview info.Put error(%v)", err)
	}
	return
}

// Adds add some toview.
func (d *Dao) Adds(ctx context.Context, mid int64, aids []int64, now int64) (err error) {
	var (
		timeB = make([]byte, 8)
		key   = hashRowKey(mid)
	)
	ctx, cancel := context.WithTimeout(ctx, time.Duration(d.conf.Info.WriteTimeout))
	defer cancel()
	binary.BigEndian.PutUint64(timeB, uint64(now))
	aidValues := make(map[string][]byte, len(aids))
	for _, aid := range aids {
		aidValues[strconv.FormatInt(aid, 10)] = timeB
	}
	values := map[string]map[string][]byte{familyAid: aidValues}
	if _, err = d.info.PutStr(ctx, tableInfo, key, values); err != nil {
		log.Error("toview info.Put error(%v)", err)
	}
	return
}

// AddMap add some toview.
func (d *Dao) AddMap(ctx context.Context, mid int64, views map[int64]*model.ToView) (err error) {
	var (
		timeB = make([]byte, 8)
		key   = hashRowKey(mid)
	)
	ctx, cancel := context.WithTimeout(ctx, time.Duration(d.conf.Info.WriteTimeout))
	defer cancel()
	aidValues := make(map[string][]byte, len(views))
	for _, v := range views {
		binary.BigEndian.PutUint64(timeB, uint64(v.Unix))
		aidValues[strconv.FormatInt(v.Aid, 10)] = timeB
	}
	values := map[string]map[string][]byte{familyAid: aidValues}
	if _, err = d.info.PutStr(ctx, tableInfo, key, values); err != nil {
		log.Error("toview info.Put error(%v)", err)
	}
	return
}

// ListInfo get all ToViews from hbase.
func (d *Dao) ListInfo(ctx context.Context, mid int64, aids []int64) (res []*model.ToView, err error) {
	var (
		result *hrpc.Result
		key    = hashRowKey(mid)
	)
	ctx, cancel := context.WithTimeout(ctx, time.Duration(d.conf.Info.ReadTimeout))
	defer cancel()
	var options []func(hrpc.Call) error
	if len(aids) != 0 {
		colunms := make([]string, 0, len(aids))
		for _, aid := range aids {
			colunms = append(colunms, fmt.Sprintf("%d", aid))
		}
		options = append(options, hrpc.Families(map[string][]string{familyAid: colunms}))
	}
	result, err = d.info.GetStr(ctx, tableInfo, key, options...)
	if err != nil && result == nil {
		log.Error("d.info.Get error(%v)", err)
		return
	}
	res = make([]*model.ToView, 0)
	for _, c := range result.Cells {
		if c != nil && len(c.Value) == 8 && bytes.Equal(c.Family, familyAidB) {
			aid, err := strconv.ParseInt(string(c.Qualifier), 10, 64)
			if err != nil {
				log.Error("strconv.ParseInt error(%v)", err)
				continue
			}
			t := &model.ToView{Aid: aid}
			t.Unix = int64(binary.BigEndian.Uint64(c.Value))
			res = append(res, t)
		}
	}
	return
}

// MapInfo get all ToViews from hbase.
func (d *Dao) MapInfo(ctx context.Context, mid int64, aids []int64) (res map[int64]*model.ToView, err error) {
	var (
		result *hrpc.Result
		key    = hashRowKey(mid)
	)
	ctx, cancel := context.WithTimeout(ctx, time.Duration(d.conf.Info.ReadTimeout))
	defer cancel()
	var options []func(hrpc.Call) error
	if len(aids) != 0 {
		colunms := make([]string, 0, len(aids))
		for _, aid := range aids {
			colunms = append(colunms, fmt.Sprintf("%d", aid))
		}
		options = append(options, hrpc.Families(map[string][]string{familyAid: colunms}))
	}
	result, err = d.info.GetStr(ctx, tableInfo, key, options...)
	if err != nil {
		log.Error("d.info.Get error(%v)", err)
		return
	}
	res = make(map[int64]*model.ToView, len(aids))
	if result == nil {
		return
	}
	for _, c := range result.Cells {
		if c != nil && len(c.Value) == 8 && bytes.Equal(c.Family, familyAidB) {
			aid, err := strconv.ParseInt(string(c.Qualifier), 10, 64)
			if err != nil {
				log.Error("strconv.ParseInt error(%v)", err)
				continue
			}
			t := &model.ToView{Aid: aid}
			t.Unix = int64(binary.BigEndian.Uint64(c.Value))
			res[aid] = t
		}
	}
	return
}

// Del delete more toview.
func (d *Dao) Del(ctx context.Context, mid int64, aids []int64) (err error) {
	key := hashRowKey(mid)
	ctx, cancel := context.WithTimeout(ctx, time.Duration(d.conf.Info.WriteTimeout))
	defer cancel()
	columns := make(map[string][]byte, len(aids))
	for _, aid := range aids {
		columns[strconv.FormatInt(aid, 10)] = []byte{}
	}
	values := map[string]map[string][]byte{familyAid: columns}
	if _, err = d.info.Delete(ctx, tableInfo, key, values); err != nil {
		log.Error("toview info.Delete error(%v)", err)
	}
	return
}

// Clear clear ToView.
func (d *Dao) Clear(ctx context.Context, mid int64) (err error) {
	key := hashRowKey(mid)
	ctx, cancel := context.WithTimeout(ctx, time.Duration(d.conf.Info.WriteTimeout))
	defer cancel()
	if _, err = d.info.Delete(ctx, tableInfo, key, nil); err != nil {
		log.Error("toview info.Delete error(%v)", err)
	}
	return
}
