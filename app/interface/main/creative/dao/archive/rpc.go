package archive

import (
	"context"

	"go-common/app/service/main/archive/api"
	"go-common/app/service/main/archive/model/archive"
	"go-common/library/ecode"
	"go-common/library/log"
)

// Archive get archive.
func (d *Dao) Archive(c context.Context, aid int64, ip string) (a *api.Arc, err error) {
	var arg = &archive.ArgAid2{Aid: aid, RealIP: ip}
	if a, err = d.arc.Archive3(c, arg); err != nil {
		log.Error("rpc archive (%d) error(%v)", aid, err)
		err = ecode.CreativeArcServiceErr
	}
	return
}

// Archives get archive list.
func (d *Dao) Archives(c context.Context, aids []int64, ip string) (a map[int64]*api.Arc, err error) {
	var arg = &archive.ArgAids2{Aids: aids, RealIP: ip}
	if a, err = d.arc.Archives3(c, arg); err != nil {
		log.Error("rpc archive (%v) error(%v)", aids, err)
		err = ecode.CreativeArcServiceErr
	}
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

// UpCount get archives count.
func (d *Dao) UpCount(c context.Context, mid int64) (count int, err error) {
	var arg = &archive.ArgUpCount2{Mid: mid}
	if count, err = d.arc.UpCount2(c, arg); err != nil {
		log.Error("rpc UpCount2 (%v) error(%v)", mid, err)
		err = ecode.CreativeArcServiceErr
	}
	return
}

// Video get video.
func (d *Dao) Video(c context.Context, aid, cid int64, ip string) (v *api.Page, err error) {
	var arg = &archive.ArgVideo2{Aid: aid, Cid: cid, RealIP: ip}
	if v, err = d.arc.Video3(c, arg); err != nil {
		if ecode.Cause(err) == ecode.NothingFound {
			err = nil
			return
		}
		log.Error("rpc video3 (%d) error(%v)", aid, err)
		err = ecode.CreativeArcServiceErr
	}
	return
}
