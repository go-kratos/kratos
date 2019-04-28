package datadao

import (
	"context"
	"encoding/json"
	"fmt"
	"go-common/app/interface/main/mcn/conf"
	"go-common/app/interface/main/mcn/model/mcnmodel"
	"go-common/app/interface/main/mcn/tool/cache"
	"time"
)

var (
	dateFmt = "20060102"
)

//CacheBaseLoader base loader
type CacheBaseLoader struct {
	SignID int64
	Date   time.Time
	Val    interface{}
	desc   string
}

func descHelper(d cache.DataLoader) (key string) {
	return d.Desc()
}

func newCacheBaseLoader(signID int64, date time.Time, val interface{}, desc string) CacheBaseLoader {
	return CacheBaseLoader{SignID: signID, Date: date, Val: val, desc: desc}
}

//Key cache's key
func (s *CacheBaseLoader) Key() (key string) {
	return fmt.Sprintf("%s_%d_%s", descHelper(s), s.SignID, s.Date.Format(dateFmt))
}

//Value cache's value
func (s *CacheBaseLoader) Value() (value interface{}) {
	return s.Val
}

//LoadValue need load value
func (s *CacheBaseLoader) LoadValue(c context.Context) (value interface{}, err error) {
	panic("implement me")
}

//Expire expiration
func (s *CacheBaseLoader) Expire() time.Duration {
	return time.Duration(conf.Conf.Memcache.McnDataCacheExpire)
}

//Desc key desc
func (s *CacheBaseLoader) Desc() string {
	return s.desc
}

// -----------------------------------------

//LoadFuncWithTp sign with type
type LoadFuncWithTp func(c context.Context, signID int64, date time.Time, tp string) (res interface{}, err error)

//LoadFuncOnlySign only sign
type LoadFuncOnlySign func(c context.Context, signID int64, date time.Time) (res interface{}, err error)

//NewCacheMcnDataWithTp reply
func NewCacheMcnDataWithTp(signID int64, date time.Time, tp string, val interface{}, desc string, loadFunc LoadFuncWithTp) *cacheMcnDataWithTp {
	return &cacheMcnDataWithTp{
		CacheBaseLoader: newCacheBaseLoader(signID, date, val, desc),
		Tp:              tp,
		LoadFunc:        loadFunc,
	}
}

//cacheMcnDataWithTp cache
type cacheMcnDataWithTp struct {
	CacheBaseLoader
	Tp       string
	LoadFunc LoadFuncWithTp
}

//Key key
func (s *cacheMcnDataWithTp) Key() (key string) {
	return fmt.Sprintf("%s_%d_%s_%s", s.Desc(), s.SignID, s.Date.Format(dateFmt), s.Tp)
}

//LoadValue load
func (s *cacheMcnDataWithTp) LoadValue(c context.Context) (value interface{}, err error) {
	value, err = s.LoadFunc(c, s.SignID, s.Date, s.Tp)
	if err != nil {
		s.Val = nil
		return
	}

	if sorter, ok := value.(mcnmodel.Sorter); ok {
		sorter.Sort()
	}

	// 如果s.Val存在，则将结果populate到s.Val上，因为外部会直接使用原始传入的s.Val值
	if s.Val != nil {
		var b, _ = json.Marshal(value)
		json.Unmarshal(b, s.Val)
	}
	return
}

// NewCacheMcnDataSignID 请求只有sign id的情况
func NewCacheMcnDataSignID(signID int64, date time.Time, val interface{}, desc string, loadFunc LoadFuncOnlySign) *cacheMcnDataSignID {
	return &cacheMcnDataSignID{
		CacheBaseLoader: newCacheBaseLoader(signID, date, val, desc),
		LoadFunc:        loadFunc,
	}
}

type cacheMcnDataSignID struct {
	CacheBaseLoader
	LoadFunc LoadFuncOnlySign
}

//Key key
func (s *cacheMcnDataSignID) Key() (key string) {
	return fmt.Sprintf("%s_%d_%s", s.Desc(), s.SignID, s.Date.Format(dateFmt))
}

//LoadValue load
func (s *cacheMcnDataSignID) LoadValue(c context.Context) (value interface{}, err error) {
	value, err = s.LoadFunc(c, s.SignID, s.Date)
	if err != nil {
		s.Val = nil
		return
	}

	if sorter, ok := value.(mcnmodel.Sorter); ok {
		sorter.Sort()
	}

	// 如果s.Val存在，则将结果populate到s.Val上，因为外部会直接使用原始传入的s.Val值
	if s.Val != nil {
		var b, _ = json.Marshal(value)
		json.Unmarshal(b, s.Val)
	}
	return
}
