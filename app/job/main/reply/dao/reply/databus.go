package reply

import (
	"context"
	"fmt"

	model "go-common/app/job/main/reply/model/reply"
	"go-common/library/log"
)

type event struct {
	Action  string         `json:"action"`
	Mid     int64          `json:"mid"`
	Subject *model.Subject `json:"subject"`
	Reply   *model.Reply   `json:"reply"`
	Report  *model.Report  `json:"report,omitempty"`
}

// PubEvent pub reply event.
func (d *Dao) PubEvent(c context.Context, action string, mid int64, sub *model.Subject, rp *model.Reply, report *model.Report) error {
	e := &event{
		Action:  action,
		Mid:     mid,
		Subject: sub,
		Reply:   rp,
		Report:  report,
	}
	if sub == nil {
		log.Error("PubEvent failed,sub is nil!value: %v %v %v %v", action, mid, rp, report)
		return nil
	}
	return d.eventBus.Send(c, fmt.Sprint(sub.Oid), &e)
}
