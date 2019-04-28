package server

import (
	"go-common/app/service/main/resource/conf"
	"go-common/app/service/main/resource/model"
	"go-common/app/service/main/resource/service"
	"go-common/library/ecode"
	"go-common/library/net/rpc"
	"go-common/library/net/rpc/context"
)

// RPC struct
type RPC struct {
	s *service.Service
}

// New init rpc server.
func New(c *conf.Config, s *service.Service) (svr *rpc.Server) {
	r := &RPC{s: s}
	svr = rpc.NewServer(c.RPCServer)
	if err := svr.Register(r); err != nil {
		panic(err)
	}
	return
}

// Ping check connection success.
func (r *RPC) Ping(c context.Context, arg *struct{}, res *struct{}) (err error) {
	return
}

// ResourceAll get all resource.
func (r *RPC) ResourceAll(c context.Context, a *struct{}, res *[]*model.Resource) (err error) {
	*res = r.s.ResourceAll(c)
	return
}

// AssignmentAll get all assignment.
func (r *RPC) AssignmentAll(c context.Context, a *struct{}, res *[]*model.Assignment) (err error) {
	*res = r.s.AssignmentAll(c)
	return
}

// DefBanner get default banner.
func (r *RPC) DefBanner(c context.Context, a *struct{}, as *model.Assignment) (err error) {
	res := r.s.DefBanner(c)
	if res == nil {
		err = ecode.NothingFound
		return
	}
	*as = *res
	return
}

// Resource get resource.
func (r *RPC) Resource(c context.Context, arg *model.ArgRes, res *model.Resource) (err error) {
	rs := r.s.Resource(c, arg.ResID)
	if rs == nil {
		err = ecode.NothingFound
		return
	}
	*res = *rs
	return
}

// Resources get resource.
func (r *RPC) Resources(c context.Context, as *model.ArgRess, res *map[int]*model.Resource) (err error) {
	*res = r.s.Resources(c, as.ResIDs)
	return
}

// Banners get banners.
func (r *RPC) Banners(c context.Context, ab *model.ArgBanner, res *model.Banners) (err error) {
	// func Banners already new rs, rs can not be nil.
	rs := r.s.Banners(c, ab.Plat, ab.Build, ab.AID, ab.MID, ab.ResIDs, ab.Channel, ab.IP, ab.Buvid, ab.Network, ab.MobiApp, ab.Device, ab.OpenEvent, ab.AdExtra, ab.Version, ab.IsAd)
	*res = *rs
	return
}

// PasterAPP get paster for APP.
func (r *RPC) PasterAPP(c context.Context, arg *model.ArgPaster, res *model.Paster) (err error) {
	var rs *model.Paster
	if rs, err = r.s.PasterAPP(c, arg.Platform, arg.AdType, arg.Aid, arg.TypeId, arg.Buvid); err == nil {
		*res = *rs
	}
	return
}

// IndexIcon get index icon.
func (r *RPC) IndexIcon(c context.Context, a *struct{}, res *map[string][]*model.IndexIcon) (err error) {
	*res = r.s.IndexIcon(c)
	return
}

// PlayerIcon get player icon config.
func (r *RPC) PlayerIcon(c context.Context, arg *struct{}, res *model.PlayerIcon) (err error) {
	var rs *model.PlayerIcon
	rs, err = r.s.PlayerIcon(c)
	if err == nil {
		*res = *rs
	}
	return
}

// Cmtbox get live box.
func (r *RPC) Cmtbox(c context.Context, cb *model.ArgCmtbox, res *model.Cmtbox) (err error) {
	rs, err := r.s.Cmtbox(c, cb.ID)
	if err == nil {
		*res = *rs
	}
	return
}

// SideBars get sode bar.
func (r *RPC) SideBars(c context.Context, a *struct{}, res *model.SideBars) (err error) {
	sbs := r.s.SideBars(c)
	*res = *sbs
	return
}

// AbTest get abtest.
func (r *RPC) AbTest(c context.Context, ab *model.ArgAbTest, res *map[string]*model.AbTest) (err error) {
	*res = r.s.AbTest(c, ab.Groups, ab.IP)
	return
}

// PasterCID get all Paster's cid.
func (r *RPC) PasterCID(c context.Context, a *struct{}, res *map[int64]int64) (err error) {
	*res, err = r.s.PasterCID(c)
	return
}
