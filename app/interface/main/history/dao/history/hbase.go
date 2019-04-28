package history

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"go-common/app/interface/main/history/model"
	"go-common/library/log"

	"github.com/tsuna/gohbase/hrpc"
)

var (
	tableInfo = "ugc:history"
	family    = "info"
	familyB   = []byte(family)
	columnSW  = "sw"
	columnSWB = []byte(columnSW)
)

// hashRowKey create rowkey(md5(mid)[:2]+mid) for histroy by mid .
func hashRowKey(mid int64) string {
	var bs = make([]byte, 8)
	binary.LittleEndian.PutUint64(bs, uint64(mid))
	rk := md5.Sum(bs)
	return fmt.Sprintf("%x%d", rk[:2], mid)
}

func (d *Dao) column(aid int64, typ int8) string {
	if typ < model.TypeArticle {
		return strconv.FormatInt(aid, 10)
	}
	return fmt.Sprintf("%d_%d", aid, typ)
}

// Add add history list.
func (d *Dao) Add(ctx context.Context, mid int64, h *model.History) (err error) {
	var (
		valueByte []byte
		column    string
		key       = hashRowKey(mid)
		fValues   = make(map[string][]byte)
	)
	if h.Aid == 0 {
		return
	}
	column = d.column(h.Aid, h.TP)
	if valueByte, err = json.Marshal(h); err != nil {
		log.Error("json.Marshal(%v) error(%v)", h, err)
		return
	}
	fValues[column] = valueByte
	values := map[string]map[string][]byte{family: fValues}
	ctx, cancel := context.WithTimeout(ctx, time.Duration(d.conf.Info.WriteTimeout))
	defer cancel()
	if _, err = d.info.PutStr(ctx, tableInfo, key, values); err != nil {
		log.Error("info.PutStr error(%v)", err)
	}
	return
}

// AddMap add history list.
func (d *Dao) AddMap(ctx context.Context, mid int64, hs map[int64]*model.History) (err error) {
	var (
		timeB   []byte
		column  string
		key     = hashRowKey(mid)
		fValues = make(map[string][]byte)
	)
	for _, h := range hs {
		if h.Aid == 0 {
			continue
		}
		// TODO typ and h.type is or not consistent .
		column = d.column(h.Aid, h.TP)
		if timeB, err = json.Marshal(h); err != nil {
			log.Error("json.Marshal(%v) error(%v)", h, err)
			continue
		}
		fValues[column] = timeB
	}
	values := map[string]map[string][]byte{family: fValues}
	ctx, cancel := context.WithTimeout(ctx, time.Duration(d.conf.Info.WriteTimeout))
	defer cancel()
	if _, err = d.info.PutStr(ctx, tableInfo, key, values); err != nil {
		log.Error("info.PutStr error(%v)", err)
	}
	return
}

// AidsMap get all Historys from hbase by aids.
func (d *Dao) AidsMap(ctx context.Context, mid int64, aids []int64) (his map[int64]*model.History, err error) {
	key := hashRowKey(mid)
	ctx, cancel := context.WithTimeout(ctx, time.Duration(d.conf.Info.ReadTimeout))
	defer cancel()
	var options []func(hrpc.Call) error
	if len(aids) != 0 {
		columns := make([]string, 0, len(aids))
		for _, aid := range aids {
			columns = append(columns, fmt.Sprintf("%d", aid))
		}
		options = append(options, hrpc.Families(map[string][]string{family: columns}))
	}
	if _, his, err = d.get(ctx, tableInfo, key, options...); err != nil {
		log.Error("hbase get() error:%+v", err)
	}
	return
}

// Map get all Historys from hbase.
func (d *Dao) Map(ctx context.Context, mid int64) (his map[string]*model.History, err error) {
	key := hashRowKey(mid)
	ctx, cancel := context.WithTimeout(ctx, time.Duration(d.conf.Info.ReadTimeout))
	defer cancel()
	if his, _, err = d.get(ctx, tableInfo, key); err != nil {
		log.Error("d.get() err:%+v", err)
	}
	return
}

// get hbase get op.
func (d *Dao) get(ctx context.Context, table, key string, options ...func(hrpc.Call) error) (his map[string]*model.History, av map[int64]*model.History, err error) {
	var result *hrpc.Result
	if result, err = d.info.GetStr(ctx, table, key, options...); err != nil {
		log.Error("d.info.Get error(%v)", err)
		return
	}
	if result == nil {
		return
	}
	expire := time.Now().Unix() - 60*60*24*90 // NOTO hbase 90days validity.
	delColumn := make(map[string][]byte)
	his = make(map[string]*model.History, len(result.Cells))
	av = make(map[int64]*model.History)
	for _, c := range result.Cells {
		if c != nil && bytes.Equal(c.Family, familyB) {
			if bytes.Equal(c.Qualifier, columnSWB) {
				continue
			}
			columns := strings.Split(string(c.Qualifier), "_")
			if len(columns) == 0 {
				continue
			}
			h := &model.History{}
			h.Aid, err = strconv.ParseInt(columns[0], 10, 64)
			if err != nil {
				log.Error("strconv.ParseInt err(%v)", err)
				continue
			}
			if err = json.Unmarshal(c.Value, h); err != nil {
				log.Error("json.Unmarshal(%s) error(%v)", string(c.Value), err)
				continue
			}
			if h.TP == model.TypeOffline {
				h.TP = model.TypeUGC
			}
			if h.Unix == 0 { // live 默认时间 7月19号删除
				h.Unix = int64(*(c.Timestamp)) / 1000
			}
			if h.Unix < expire {
				delColumn[string(c.Qualifier)] = []byte{}
				continue
			}
			his[d.column(h.Aid, h.TP)] = h
			h.FillBusiness()
			if h.TP < model.TypeArticle {
				av[h.Aid] = h
			}
		}
	}
	if len(delColumn) > 0 {
		d.delChan.Save(func() {
			log.Warn("delete hbase key:%s", key)
			d.delete(context.Background(), table, key, delColumn)
		})
	}
	return
}

// Delete delete more history.
func (d *Dao) delete(ctx context.Context, table, key string, delColumn map[string][]byte) (err error) {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(d.conf.Info.WriteTimeout))
	defer cancel()
	values := map[string]map[string][]byte{family: delColumn}
	if _, err = d.info.Delete(ctx, table, key, values); err != nil {
		log.Error("info.Delete() delColumn:%+v,error(%v)", delColumn, err)
	}
	return
}

// DelAids delete more history.
func (d *Dao) DelAids(ctx context.Context, mid int64, aids []int64) (err error) {
	columns := make(map[string][]byte, len(aids))
	for _, aid := range aids {
		columns[strconv.FormatInt(aid, 10)] = []byte{}
	}
	columns[strconv.FormatInt(mid, 10)] = []byte{}
	columns[d.column(mid, model.TypeLive)] = []byte{}
	ctx, cancel := context.WithTimeout(ctx, time.Duration(d.conf.Info.WriteTimeout))
	defer cancel()
	key := hashRowKey(mid)
	values := map[string]map[string][]byte{family: columns}
	if _, err = d.info.Delete(ctx, tableInfo, key, values); err != nil {
		log.Error("info.Delete()len:%d,error(%v)", len(aids), err)
	}
	return
}

// Delete delete more history.
func (d *Dao) Delete(ctx context.Context, mid int64, his []*model.History) (err error) {
	columns := make(map[string][]byte, len(his))
	for _, h := range his {
		columns[d.column(h.Aid, h.TP)] = []byte{}
	}
	ctx, cancel := context.WithTimeout(ctx, time.Duration(d.conf.Info.WriteTimeout))
	defer cancel()
	values := map[string]map[string][]byte{family: columns}
	if _, err = d.info.Delete(ctx, tableInfo, hashRowKey(mid), values); err != nil {
		log.Error("info.Delete()len:%d,error(%v)", len(his), err)
	}
	return
}

// Clear clear history.
func (d *Dao) Clear(ctx context.Context, mid int64) (err error) {
	var sw int64
	ctx, cancel := context.WithTimeout(ctx, time.Duration(d.conf.Info.ReadTimeout))
	defer cancel()
	if sw, err = d.InfoShadow(ctx, mid); err != nil {
		return
	}
	key := hashRowKey(mid)
	if _, err = d.info.Delete(ctx, tableInfo, key, nil); err != nil {
		log.Error("info.Delete error(%v)", err)
	}
	if sw == model.ShadowOn {
		d.SetInfoShadow(ctx, mid, sw)
	}
	return
}

// SetInfoShadow set the user switch to hbase.
func (d *Dao) SetInfoShadow(ctx context.Context, mid, value int64) (err error) {
	var key = hashRowKey(mid)
	ctx, cancel := context.WithTimeout(ctx, time.Duration(d.conf.Info.ReadTimeout))
	defer cancel()
	valueB := []byte(strconv.FormatInt(value, 10))
	values := map[string]map[string][]byte{family: {columnSW: valueB}}
	if _, err = d.info.PutStr(ctx, tableInfo, key, values); err != nil {
		log.Error("info.Put error(%v)", err)
	}
	return
}

// InfoShadow return the user switch from hbase.
func (d *Dao) InfoShadow(ctx context.Context, mid int64) (sw int64, err error) {
	var key = hashRowKey(mid)
	ctx, cancel := context.WithTimeout(ctx, time.Duration(d.conf.Info.ReadTimeout))
	defer cancel()
	result, err := d.info.GetStr(ctx, tableInfo, key, hrpc.Families(map[string][]string{family: {columnSW}}))
	if err != nil {
		log.Error("d.info.Get error(%v)", err)
		return
	}
	if result == nil {
		err = errors.New("info sw data null")
		return
	}
	for _, c := range result.Cells {
		if c != nil && bytes.Equal(c.Family, familyB) {
			if columnSW == string(c.Qualifier) {
				sw, _ = strconv.ParseInt(string(c.Value), 10, 0)
				return
			}
		}
	}
	return
}
