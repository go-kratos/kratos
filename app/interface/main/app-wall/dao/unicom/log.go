package unicom

import (
	"context"
	"fmt"
	"time"

	"go-common/app/interface/main/app-wall/model/unicom"
	"go-common/library/database/elastic"
	"go-common/library/log"
)

const (
	_logUserKey = "log_user_action_91_%s"
)

func keyLogUser(now time.Time) string {
	return fmt.Sprintf(_logUserKey, now.Format("2006_01"))
}

func timeToString(now time.Time) string {
	return now.Format("2006-01-02 15:04:05")
}

// SearchUserBindLog unicom user bind log
func (d *Dao) SearchUserBindLog(ctx context.Context, mid int64, now time.Time) (ulogs []*unicom.UserLog, err error) {
	var data struct {
		Result []struct {
			Ctime     string `json:"ctime"`
			ExtraData string `json:"extra_data"`
		} `json:"result"`
	}
	var (
		endtime     = now.AddDate(0, -1, 0)
		endWeekTime = now.AddDate(0, 0, -7)
	)
	req := d.es.NewRequest("log_user_action")
	req.Index(keyLogUser(now), keyLogUser(endtime)).WhereEq("mid", mid).WhereIn("action", []string{"unicom_userpack_add", "unicom_userpack_deduct"}).
		WhereRange("ctime", timeToString(endWeekTime), timeToString(now), elastic.RangeScopeLcRc).Order("ctime", "desc").
		Pn(1).Ps(2000)
	if err = req.Scan(ctx, &data); err != nil {
		log.Error("search user bind params(%s) error(%v)", req.Params(), err)
		return
	}
	for _, d := range data.Result {
		ulog := &unicom.UserLog{}
		if err = ulog.UserLogJSONChange(d.ExtraData); err != nil {
			log.Error("ulog json change error(%v)", err)
			continue
		}
		ulog.Ctime = d.Ctime
		ulogs = append(ulogs, ulog)
	}
	return
}
