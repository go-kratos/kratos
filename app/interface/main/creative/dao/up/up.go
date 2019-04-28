package up

import (
	"context"
	"sync"

	"go-common/app/interface/main/creative/model/up"
	upapi "go-common/app/service/main/up/api/v1"
	upmdl "go-common/app/service/main/up/model"
	"go-common/library/log"
	"go-common/library/sync/errgroup"
	"net/url"
	"strconv"
)

// UpInfo rpc
func (d *Dao) UpInfo(c context.Context, mid int64, from int, ip string) (res *upmdl.UpInfo, err error) {
	var arg = &upmdl.ArgInfo{
		Mid:  mid,
		From: from,
	}
	if res, err = d.up.Info(c, arg); err != nil {
		log.Error("d.up.Info error(%v)", err)
	}
	return
}

// UpSwitch get switch
func (d *Dao) UpSwitch(c context.Context, mid int64, from int, ip string) (res *upmdl.PBUpSwitch, err error) {
	var arg = &upmdl.ArgUpSwitch{
		Mid:  mid,
		From: from,
		IP:   ip,
	}
	if res, err = d.up.UpSwitch(c, arg); err != nil {
		log.Error("d.up.UpSwitch error(%v)", err)
	}
	return
}

// SetUpSwitch set switch
func (d *Dao) SetUpSwitch(c context.Context, mid int64, state, from int, ip string) (res *upmdl.PBSetUpSwitchRes, err error) {
	var arg = &upmdl.ArgUpSwitch{
		Mid:   mid,
		From:  from,
		State: state,
		IP:    ip,
	}
	if res, err = d.up.SetUpSwitch(c, arg); err != nil {
		log.Error("d.up.SetUpSwitch error(%v)", err)
	}
	return
}

// UpSpecialGroups 获取UP主的特殊用户组
func (d *Dao) UpSpecialGroups(c context.Context, mid int64) (groups map[int64]*up.SpecialGroup, err error) {
	groups = make(map[int64]*up.SpecialGroup)
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	var res struct {
		Code int                `json:"code"`
		Msg  string             `json:"message"`
		Data []*up.SpecialGroup `json:"data"`
	}
	if err = d.httpClient.Get(c, d.c.Host.API+_upSpecialGroupURI, "", params, &res); err != nil {
		log.Error("d.UpSpecialGroups() error(%v)", err)
		return
	}
	if res.Data == nil {
		log.Warn("UpSpecialGroups(%d) error when get up groups", mid)
		return
	}
	for _, v := range res.Data {
		groups[v.GroupID] = v
	}
	return
}

// UpSpecial 获取UP主的特殊用户组
func (d *Dao) UpSpecial(c context.Context, gpid int64) (ups map[int64]int64, err error) {
	var (
		res  *upapi.UpGroupMidsReply
		page int
		g    errgroup.Group
		l    sync.RWMutex
	)
	if res, err = d.UpClient.UpGroupMids(c, &upapi.UpGroupMidsReq{
		GroupID: gpid,
		Pn:      1,
		Ps:      1,
	}); err != nil {
		log.Error("UpSpecial d.UpSpecial gpid(%d)|error(%v)", gpid, err)
		return
	}
	log.Warn("UpSpecial get total: gpid(%d)|total(%d)", gpid, res.Total)
	if res.Total <= 0 {
		return
	}
	ups = make(map[int64]int64, res.Total)
	ps := int(10000)
	pageNum := res.Total / ps
	if res.Total%ps != 0 {
		pageNum++
	}
	for page = 1; page <= pageNum; page++ {
		tmpPage := page
		g.Go(func() (err error) {
			resgg, err := d.UpClient.UpGroupMids(c, &upapi.UpGroupMidsReq{
				GroupID: gpid,
				Pn:      tmpPage,
				Ps:      ps,
			})
			if err != nil {
				log.Error("d.UpGroupMids gg (%d,%d,%d) error(%v) ", gpid, tmpPage, ps, err)
				err = nil
				return
			}
			for _, mid := range resgg.Mids {
				l.Lock()
				ups[mid] = mid
				l.Unlock()
			}
			return
		})
	}
	g.Wait()
	log.Warn("UpSpecial get result: gpid,total,midslen,upslens (%d)|(%d)|(%d)", gpid, res.Total, len(ups))
	return
}
