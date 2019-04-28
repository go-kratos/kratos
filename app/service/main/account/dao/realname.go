package dao

import (
	"context"

	member "go-common/app/service/main/member/model"
)

// RealnameDetail is.
func (d *Dao) RealnameDetail(c context.Context, mid int64) (detail *member.RealnameDetail, err error) {
	req := &member.ArgMemberMid{Mid: mid}
	if detail, err = d.mRPC.RealnameDetail(c, req); err != nil {
		return
	}
	return
}
