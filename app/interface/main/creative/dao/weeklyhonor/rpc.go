package weeklyhonor

import (
	"context"
	up "go-common/app/service/main/up/api/v1"
	"go-common/library/log"
)

const fromWeeklyHonor = 1

// ChangeUpSwitch change up switch on/off
func (d *Dao) ChangeUpSwitch(c context.Context, mid int64, state uint8) (err error) {
	req := up.UpSwitchReq{
		Mid:   mid,
		From:  fromWeeklyHonor,
		State: state,
	}
	if _, err = d.upClient.SetUpSwitch(c, &req); err != nil {
		log.Error("d.upClient.SetUpSwitch req(%+v),err(%v)", req, err)
	}
	return
}

// GetUpSwitch get up switch state
func (d *Dao) GetUpSwitch(c context.Context, mid int64) (state uint8, err error) {
	req := up.UpSwitchReq{
		Mid:  mid,
		From: fromWeeklyHonor,
	}
	reply, err := d.upClient.UpSwitch(c, &req)
	if err != nil {
		log.Error("d.upClient.UpSwitch req(%+v),err(%v)", req, err)
		return
	}
	state = reply.GetState()
	return
}
