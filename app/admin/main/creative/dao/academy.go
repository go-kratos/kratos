package dao

import (
	"context"
	"strconv"

	"go-common/app/interface/openplatform/article/model"
	"go-common/app/service/main/archive/api"
	"go-common/app/service/main/archive/model/archive"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/metadata"
)

// Archive get archive.
func (d *Dao) Archive(c context.Context, aid int64) (a *api.Arc, err error) {
	var ip = metadata.String(c, metadata.RemoteIP)
	var arg = &archive.ArgAid2{Aid: aid, RealIP: ip}
	if a, err = d.arc.Archive3(c, arg); err != nil {
		log.Error("d.arc.Archive3 aid(%d)|error(%v)", aid, err)
		err = ecode.CreativeArcServiceErr
	}
	return
}

// Archives get archive list.
func (d *Dao) Archives(c context.Context, aids []int64) (a map[int64]*api.Arc, err error) {
	var ip = metadata.String(c, metadata.RemoteIP)
	var arg = &archive.ArgAids2{Aids: aids, RealIP: ip}
	if a, err = d.arc.Archives3(c, arg); err != nil {
		log.Error("d.arc.Archive3 aids(%v+)|error(%v)", aids, err)
		err = ecode.CreativeArcServiceErr
	}
	return
}

// ArticleMetas batch get articles by aids.
func (d *Dao) ArticleMetas(c context.Context, aids []int64) (res map[int64]*model.Meta, err error) {
	var ip = metadata.String(c, metadata.RemoteIP)
	arg := &model.ArgAids{Aids: aids, RealIP: ip}
	if res, err = d.art.ArticleMetas(c, arg); err != nil {
		log.Error("d.art.ArticleMetas aids(%+v)|error(%v)", arg, err)
		if _, er := strconv.ParseInt(err.Error(), 10, 64); er != nil {
			err = ecode.CreativeArticleRPCErr
		}
	}
	log.Info("d.art.ArticleMetas aids(%v)|res(%+v)", aids, res)
	return
}

// Stats get archives stat.
func (d *Dao) Stats(c context.Context, aids []int64, ip string) (a map[int64]*api.Stat, err error) {
	var arg = &archive.ArgAids2{Aids: aids, RealIP: ip}
	if a, err = d.arc.Stats3(c, arg); err != nil {
		log.Error("rpc Stats (%v) error(%v)", aids, err)
		err = ecode.CreativeArcServiceErr
	}
	return
}
