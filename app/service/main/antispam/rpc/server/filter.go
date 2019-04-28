package server

import (
	"fmt"

	"go-common/app/service/main/antispam/conf"
	"go-common/app/service/main/antispam/model"
	"go-common/app/service/main/antispam/service"

	"go-common/library/log"
	"go-common/library/net/rpc/context"
)

// Filter .
type Filter struct {
	svr service.Service
}

func valid(susp *model.Suspicious) error {
	if susp == nil {
		err := fmt.Errorf("nil request params susp(%v)", susp)
		log.Error("%v", err)
		return err
	}
	if susp.OId <= 0 {
		err := fmt.Errorf("OId(%d) must be greater than 0", susp.SenderId)
		log.Error("%v", err)
		return err
	}
	if susp.Area == "message" {
		// for backward compatibility(history reason)
		susp.Area = model.AreaIMessage
		if susp.SenderId <= 0 {
			err := fmt.Errorf("senderId(%d) must be greater than 0", susp.SenderId)
			log.Error("%v", err)
			return err
		}
		return nil
	}
	if _, ok := conf.Areas[susp.Area]; !ok {
		err := fmt.Errorf("invalid area(%s)", susp.Area)
		log.Error("%v", err)
		return err
	}
	return nil
}

// Check .
func (f *Filter) Check(ctx context.Context, susp *model.Suspicious, resp *model.SuspiciousResp) error {
	if err := valid(susp); err != nil {
		return err
	}
	ret, err := f.svr.Filter(ctx, susp)
	if err != nil {
		return err
	}
	*resp = *ret
	return nil
}

// Ping .
func (f *Filter) Ping(ctx context.Context, arg *struct{}, res *struct{}) error {
	return nil
}
