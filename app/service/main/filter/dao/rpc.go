package dao

import (
	"context"

	"go-common/app/service/main/antispam/model"
	"go-common/library/ecode"

	"github.com/pkg/errors"
)

// AntispamFilter .
func (d *Dao) AntispamFilter(c context.Context, area string, content string, oID int64, contentID int64, senderID int64) (hit string, limitType string, err error) {
	var (
		arg = &model.Suspicious{
			Area:     area,
			Content:  content,
			OId:      oID,
			Id:       contentID,
			SenderId: senderID,
		}
		res *model.SuspiciousResp
	)
	if res, err = d.antispamRPC.Filter(c, arg); err != nil {
		err = errors.WithStack(err)
		return
	}
	if res == nil {
		err = ecode.RequestErr
		return
	}
	hit = res.Content
	limitType = res.LimitType
	return
}
