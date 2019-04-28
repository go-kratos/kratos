package subtitle

import (
	"context"
	"go-common/app/interface/main/dm2/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

// View fn
func (d *Dao) View(c context.Context, aid int64) (res *model.SubtitleSubjectReply, err error) {
	var arg = &model.ArgArchiveID{
		Aid: aid,
	}
	if res, err = d.sub.SubtitleSubjectSubmitGet(c, arg); err != nil {
		log.Error("d.sub.SubtitleSubjectSubmitGet (%+v) error(%v)", arg, err)
		err = ecode.CreativeSubtitleAPIErr
	}
	return
}
