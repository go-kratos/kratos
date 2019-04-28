package archive

import (
	"context"
	"strconv"

	"go-common/app/interface/openplatform/article/model"
	"go-common/app/service/main/archive/api"
	"go-common/app/service/main/archive/model/archive"

	"go-common/library/ecode"
	"go-common/library/log"
)

// UpCount get archives count.
func (d *Dao) UpCount(c context.Context, mid int64) (count int, err error) {
	var arg = &archive.ArgUpCount2{Mid: mid}
	if count, err = d.arc.UpCount2(c, arg); err != nil {
		log.Error("rpc UpCount2 (%v) error(%v)", mid, err)
		err = ecode.CreativeArcServiceErr
	}
	return
}

// Archives get archive list.
func (d *Dao) Archives(c context.Context, aids []int64) (a map[int64]*api.Arc, err error) {
	var arg = &archive.ArgAids2{Aids: aids, RealIP: ""}
	if a, err = d.arc.Archives3(c, arg); err != nil {
		log.Error("d.arc.Archive3 aids(%v+)|error(%v)", aids, err)
		err = ecode.CreativeArcServiceErr
	}
	return
}

// Archive get archive info.
func (d *Dao) Archive(c context.Context, aid int64) (a *api.Arc, err error) {
	var arg = &archive.ArgAid2{Aid: aid, RealIP: ""}
	if a, err = d.arc.Archive3(c, arg); err != nil {
		log.Error("d.arc.Archive3 aid(%d)|error(%v)", aid, err)
		err = ecode.CreativeArcServiceErr
	}
	return
}

// Stats get archives stat.
func (d *Dao) Stats(c context.Context, aids []int64) (a map[int64]*api.Stat, err error) {
	var arg = &archive.ArgAids2{Aids: aids, RealIP: ""}
	if a, err = d.arc.Stats3(c, arg); err != nil {
		log.Error("rpc Stats (%v) error(%v)", aids, err)
		err = ecode.CreativeArcServiceErr
	}
	return
}

// ArticleMetas batch get articles by aids.
func (d *Dao) ArticleMetas(c context.Context, aids []int64) (res map[int64]*model.Meta, err error) {
	arg := &model.ArgAids{Aids: aids, RealIP: ""}
	if res, err = d.art.ArticleMetas(c, arg); err != nil {
		log.Error("d.art.ArticleMetas aids(%+v)|error(%v)", arg, err)
		if _, er := strconv.ParseInt(err.Error(), 10, 64); er != nil {
			err = ecode.CreativeArticleRPCErr
		}
	}
	return
}
