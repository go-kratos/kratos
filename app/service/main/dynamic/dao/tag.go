package dao

import (
	"context"
	"net/url"

	"go-common/library/ecode"
)

// Hot Second type hot tag ids.
func (d *Dao) Hot(c context.Context) (res map[int32][]int64, err error) {
	params := url.Values{}
	var rs struct {
		Code int               `json:"code"`
		Data map[int32][]int64 `json:"data"`
	}
	if err = d.httpR.Get(c, d.hotURI, "", params, &rs); err != nil {
		PromError("二级分区热门Tag接口", "d.httpR.Get(%s) error(%v)", d.hotURI, err)
		return
	}
	if rs.Code != ecode.OK.Code() {
		PromError("二级分区热门Tag接口", "tag hotmap url(%s) res code(%d) or res.data(%v)", d.hotURI, rs.Code, rs.Data)
		err = ecode.Int(rs.Code)
		return
	}
	res = rs.Data
	return
}

// Rids first type ids.
func (d *Dao) Rids(c context.Context) (res []int32, err error) {
	params := url.Values{}
	var rs struct {
		Code int     `json:"code"`
		Data []int32 `json:"data"`
	}
	if err = d.httpR.Get(c, d.pridURI, "", params, &rs); err != nil {
		PromError("一级分区ID接口", "d.httpR.Get(%s) error(%v)", d.pridURI, err)
		return
	}
	if rs.Code != ecode.OK.Code() {
		PromError("一级分区ID接口", "tag prids url(%s) res code(%d) or res.data(%v)", d.pridURI, rs.Code, rs.Data)
		err = ecode.Int(rs.Code)
		return
	}
	res = rs.Data
	return
}
