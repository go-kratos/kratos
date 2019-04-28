package ugc

import (
	"context"

	"go-common/library/log"
)

// FinishArc updates the submit status from 1 to 0
func (d *Dao) FinishArc(c context.Context, aid int64) (err error) {
	if _, err = d.DB.Exec(c, _finishArc, aid); err != nil {
		log.Error("FinishVideos Error: %v", aid, err)
	}
	return
}
