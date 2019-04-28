package pendant

import (
	"context"
	"strconv"

	"go-common/app/service/main/usersuit/model"
	"go-common/library/ecode"

	"github.com/pkg/errors"
)

// VipInfo get identify info by calling api.
func (d *Dao) VipInfo(c context.Context, mid int64, ip string) (idt *model.VipInfo, err error) {
	var res struct {
		Code int
		Data *model.VipInfo
	}
	if err = d.client.Get(c, d.vipInfoURL+strconv.FormatInt(mid, 10), ip, nil, &res); err != nil {
		return
	}
	if res.Code != 0 {
		err = errors.WithStack(ecode.Int(res.Code))
		return
	}
	idt = res.Data
	return
}
