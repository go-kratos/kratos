package danmu

import (
	"context"
	"go-common/app/interface/main/creative/model/danmu"

	"go-common/app/interface/main/dm2/model"
	"go-common/library/log"
)

// SubView fn
func (s *Service) SubView(c context.Context, aid int64, ip string) (ret *danmu.SubtitleSubjectReply, err error) {
	var sub *model.SubtitleSubjectReply
	if sub, err = s.sub.View(c, aid); err != nil {
		log.Error("s.sub.View err(%v) | aid(%d), ip(%s)", err, aid, ip)
		return
	}
	if sub != nil {
		ret = &danmu.SubtitleSubjectReply{
			AllowSubmit: sub.AllowSubmit,
			Lan:         sub.Lan,
			LanDoc:      sub.LanDoc,
		}
	}
	return
}
