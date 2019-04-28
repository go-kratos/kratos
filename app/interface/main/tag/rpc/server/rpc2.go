package rpc

import (
	"go-common/app/interface/main/tag/model"
	"go-common/library/log"
	"go-common/library/net/rpc/context"
)

// UpBind .
func (r *RPC) UpBind(c context.Context, arg *model.ArgBind, res *struct{}) (err error) {
	var checked []string
	for _, name := range arg.Names {
		if name, err = r.srv.CheckName(name); err == nil && name != "" {
			checked = append(checked, name)
		}
	}
	if len(checked) > r.c.Tag.ArcTagMaxNum {
		log.Error("len checked(%d)", len(checked))
		checked = checked[0:r.c.Tag.ArcTagMaxNum]
	}
	err = r.srv.UpResBind(c, arg.Oid, arg.Mid, checked, arg.Type, c.Now())
	return
}

// AdminBind .
func (r *RPC) AdminBind(c context.Context, arg *model.ArgBind, res *struct{}) (err error) {
	var checked []string
	for _, name := range arg.Names {
		if name, err = r.srv.CheckName(name); err == nil && name != "" {
			checked = append(checked, name)
		}
	}
	if len(checked) > r.c.Tag.ArcTagMaxNum {
		log.Error("len checked(%d)", len(checked))
		checked = checked[0:r.c.Tag.ArcTagMaxNum]
	}
	err = r.srv.ResAdminBind(c, arg.Oid, arg.Mid, checked, arg.Type, c.Now())
	return
}

// ResTags returns resource tags by resource id .
func (r *RPC) ResTags(c context.Context, arg *model.ArgResTags, res *map[int64][]*model.Tag) (err error) {
	*res, err = r.srv.ResTags(c, arg.Oids, arg.Mid, arg.Type)
	return
}
