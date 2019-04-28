package view

import (
	"context"
	"time"

	"go-common/app/interface/main/tv/model"
	upMdl "go-common/app/interface/main/tv/model/upper"
	"go-common/app/interface/main/tv/model/view"
	arcwar "go-common/app/service/main/archive/api"
	"go-common/app/service/main/archive/model/archive"
	"go-common/library/ecode"
	"go-common/library/log"
)

// View  all view data.
func (s *Service) View(c context.Context, mid, aid int64, ak, ip string, now time.Time) (v *view.View, isok bool, errMsg string, err error) {
	var (
		arcMeta *model.ArcCMS
		faved   bool
		reginfo *arcwar.Tp
	)
	if arcMeta, isok, errMsg, err = s.ArcMsg(aid); err != nil { // arc auth msg, if not ok, we return the err msg
		log.Info("View ArcMsg Aid:%d, Err:%v", aid, err)
		return
	} else if !isok {
		return
	}
	if v, err = s.ViewPage(c, mid, aid, ak, ip, now); err != nil { // View page
		if err == ecode.AccessDenied || err == ecode.NothingFound {
			log.Warn("s.ViewPage() mid(%d) aid(%d) ak(%s) ip(%s) error(%v)",
				mid, aid, ak, ip, err)
			return
		}
		if err == ecode.TvAllDataAuditing { // the err is used to transport the message that all the data is being audited
			isok, errMsg = s.cmsDao.AuditingMsg()
			err = nil
			return
		}
		log.Error("s.ViewPage() mid(%d) aid(%d) ak(%s) ip(%s) error(%v)",
			mid, aid, ak, ip, err)
		return
	}
	if reginfo, err = s.arcDao.TypeInfo(int32(arcMeta.TypeID)); err != nil {
		log.Warn("s.arcDao.TypeInfo tid(%d) error(%v)", arcMeta.TypeID, err)
	}
	if reginfo != nil {
		v.PID = reginfo.Pid
	}
	s.arcMetaInterv(v, arcMeta)  // cms interv
	s.initRelates(c, v, ip, now) // get relates
	v.ReqUser = &view.ReqUser{}
	if faved, err = s.favDao.InDefault(c, mid, aid); err != nil {
		log.Warn("s.favDao InDefault Mid %d, Aid %d, Err %v", mid, aid, err)
	}
	if faved {
		v.ReqUser.Favorite = 1
	}
	return
}

// arcMetaInterv replaces view's archive info by arc cms meta info
func (s *Service) arcMetaInterv(v *view.View, arcMeta *model.ArcCMS) {
	if arcMeta == nil {
		return
	}
	if arcMeta.Title != "" {
		v.Static.Title = arcMeta.Title
	}
	if arcMeta.Cover != "" {
		v.Static.Pic = arcMeta.Cover
	}
	if arcMeta.Content != "" {
		v.Static.Desc = arcMeta.Content
	}
	// todo 分p的干预
}

// ViewPage view page data.
func (s *Service) ViewPage(c context.Context, mid, aid int64, ak, ip string, now time.Time) (v *view.View, err error) {
	var (
		vs    *view.Static
		vp    *arcwar.ViewReply
		upper *upMdl.Upper
	)
	if vp, err = s.arcDao.GetView(c, aid); err != nil {
		log.Error("ViewPage getView Aid:%d, Err:%v", aid, err)
		return
	}
	if upper, err = s.upDao.LoadUpMeta(c, vp.Arc.Author.Mid); err != nil || !upper.CanShow() { // if upper can't be found or upper is not valid, hide it
		vp.Arc.Author.Face = ""
		vp.Arc.Author.Name = ""
	} else { // if upper is valid, use cms info to show
		vp.Arc.Author.Face = upper.CMSFace
		vp.Arc.Author.Name = upper.CMSName
	}
	vs = &view.Static{Arc: vp.Arc}
	if err = s.initPages(c, vs, vp.Pages); err != nil {
		log.Error("ViewPage initPages Aid %d, Err %v", aid, err)
		return
	}
	v = &view.View{Static: vs}
	if v.AttrVal(archive.AttrBitIsPGC) != archive.AttrYes {
		// check access
		if err = s.checkAceess(c, mid, v.Aid, int(v.State), int(v.Access), ak, ip); err != nil {
			// archive is ForbitFixed and Transcoding and StateForbitDistributing need analysis history body .
			if v.State != archive.StateForbidFixed {
				return
			}
			err = nil
		}
		if v.Access > 0 {
			v.Stat.View = 0
		}
	}
	if mid != 0 {
		v.History, _ = s.arcDao.Progress(ctx, v.Aid, mid)
	}
	return
}
