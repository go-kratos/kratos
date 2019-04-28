package dao

import (
	"context"
	"fmt"

	"go-common/app/admin/main/tv/model"
	"go-common/library/database/elastic"
	"go-common/library/log"
)

// ArcES treats the ugc index request and call the ES to get the result
func (d *Dao) ArcES(c context.Context, req *model.ReqArcES) (data *model.EsUgcResult, err error) {
	var (
		cfg = d.c.Cfg.EsIdx.UgcIdx
		r   = d.esClient.NewRequest(cfg.Business).Index(cfg.Index).WhereEq("deleted", 0)
	)
	if req.Valid != "" {
		r = r.WhereEq("valid", req.Valid)
	}
	if req.AID != "" {
		r = r.WhereEq("aid", req.AID)
	}
	if req.Result != "" {
		r = r.WhereEq("result", req.Result)
	}
	if len(req.Typeids) != 0 {
		r = r.WhereIn("typeid", req.Typeids)
	}
	if req.Title != "" {
		r = r.WhereLike([]string{"title"}, []string{req.Title}, true, elastic.LikeLevelMiddle)
	}
	if len(req.Mids) != 0 {
		r = r.WhereIn("mid", req.Mids)
	}
	r.Ps(req.Ps).Pn(int(req.Pn))
	if req.MtimeOrder != "" {
		r = r.Order("mtime", req.MtimeSort())
	}
	if req.PubtimeOrder != "" {
		r = r.Order("pubtime", req.PubtimeSort())
	}
	if err = r.Scan(c, &data); err != nil {
		log.Error("ArcES:Scan params(%s) error(%v)", r.Params(), err)
		return
	}
	if data == nil || data.Page == nil {
		err = fmt.Errorf("data or data.Page nil")
		log.Error("ArcES params(%s) error(%v)", r.Params(), err)
		return
	}
	return
}
