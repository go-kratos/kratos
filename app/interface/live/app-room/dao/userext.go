package dao

import (
	"context"
	"go-common/app/interface/live/app-room/model"
	"go-common/app/service/live/userext/api/liverpc/v0"
	"go-common/library/ecode"
	"go-common/library/log"
	"strconv"
)

//GetUserConf 获取用户配置
func (d *Dao) GetUserConf(c context.Context, mid int64, targetId int64, keys []int64) (values model.ConfValues, err error) {
	req := &v0.ConfGetReq{
		Uid:      mid,
		TargetId: targetId,
		Type:     keys,
	}
	resp, err := d.UserExtAPI.V0Conf.Get(c, req)
	if err != nil {
		log.Error("call getUserConf sys error %v", err.Error())
		return
	}
	if resp.Code != 0 {
		log.Error("call getUserConf failed %d %v", resp.Code, resp.Msg)
		err = ecode.ParamInvalid
		return
	}
	o := resp.Data
	values = make(map[int64]string, len(o))
	for i, v := range o {
		s, e := strconv.ParseInt(i, 10, 64)
		if e != nil {
			log.Error("call getUserConf string2int64 failed %v", e.Error())
			err = e
			return
		}
		values[s] = v
	}
	return
}

//AsyncSetUserConf 异步设置用户一个配置 "存在"
func (d *Dao) AsyncSetUserConf(c context.Context, mid int64, targetId int64, key int64) (err error) {
	value := model.Set
	vs := make(map[int64]string, 1)
	vs[key] = value
	err = d.giftCache.Save(func() {
		d.SetMoreUserConf(context.Background(), mid, targetId, vs)
	})
	if err != nil {
		err = d.SetMoreUserConf(c, mid, targetId, vs)
	}
	return
}

//SetUserConf 设置用户一个配置 "存在"
func (d *Dao) SetUserConf(c context.Context, mid int64, targetId int64, key int64) (err error) {
	value := model.Set
	vs := make(map[int64]string, 1)
	vs[key] = value
	err = d.SetMoreUserConf(c, mid, targetId, vs)
	return
}

//DelUserConf 删除用户一个配置 "不存在"
func (d *Dao) DelUserConf(c context.Context, mid int64, targetId int64, key int64) (err error) {
	value := model.Empty
	vs := make(map[int64]string, 1)
	vs[key] = value
	err = d.SetMoreUserConf(c, mid, targetId, vs)
	return
}

//SetMoreUserConf 设置用户多个配置
func (d *Dao) SetMoreUserConf(c context.Context, mid int64, targetId int64, values map[int64]string) (err error) {
	varList := make([]*v0.ConfSetReq_Var, len(values))

	ic := 0
	for k, v := range values {
		rv := &v0.ConfSetReq_Var{
			Uid:      mid,
			TargetId: targetId,
			Type:     k,
			Content:  v,
		}
		varList[ic] = rv
		ic++
	}

	req := &v0.ConfSetReq{
		VarList: varList,
	}
	r, err := d.UserExtAPI.V0Conf.Set(c, req)
	if err != nil {
		log.Error("call setUserConf sys error %v", err.Error())
		return
	}
	if r.Code != 0 {
		log.Error("call setUserConf failed %d %v", r.Code, r.Msg)
		err = ecode.ParamInvalid
		return
	}
	return
}
