package service

import (
	"context"
	"go-common/app/interface/bbq/app-bbq/api/http/v1"
	"go-common/app/interface/bbq/app-bbq/model"
	"go-common/library/ecode"
	"go-common/library/log"

	"github.com/golang/protobuf/ptypes/empty"
)

//Report ..
func (s *Service) Report(c context.Context, arg *v1.ReportRequest, mid int64, ak string) (res *empty.Empty, err error) {
	res = &empty.Empty{}
	switch arg.Type {
	case model.TypeVideo:
		if err = s.dao.ReportVideo(c, arg.SVID, mid, model.MapReasons[arg.Reason]); err != nil {
			log.Warnv(c, log.KV("event", "report video err"), log.KV("err", err), log.KV("req", arg))
			return
		}
	case model.TypeComment:
		reason := model.BiliReasonsMap[arg.Reason]
		content := ""
		if reason == 0 {
			content = "BBQ其他理由"
		}
		req := &v1.CommentReportReq{
			SvID:      arg.SVID,
			RpID:      arg.RpID,
			Reason:    reason,
			Content:   content,
			Type:      model.DefaultCmType,
			AccessKey: ak,
		}
		if err = s.CommentReport(c, req); err != nil {
			log.Warnv(c, log.KV("event", "report comment err"), log.KV("err", err), log.KV("req", arg))
			return
		}
	case model.TypeDanmu:
		var (
			dm *model.Danmu
		)
		if dm, err = s.dao.RawBullet(c, arg.Danmu); err != nil {
			log.Errorw(c, "event", "query mid", "err", err)
			err = ecode.ReportDanmuError
			return
		}
		if err = s.dao.ReportDanmu(c, arg.Danmu, dm.MID, mid, model.MapReasons[arg.Reason]); err != nil {
			log.Warnv(c, log.KV("event", "report danmu err"), log.KV("err", err), log.KV("req", arg))
			return
		}
	case model.TypeUser:
		var rType int
		if arg.Reason == 5 {
			rType = 0
		} else {
			rType = 1
		}
		if err = s.dao.ReportUser(c, rType, arg.UpMID, mid, model.MapReasons[arg.Reason]); err != nil {
			log.Warnv(c, log.KV("event", "report user err"), log.KV("err", err), log.KV("req", arg))
			return
		}
	}
	return
}
