package medal

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"go-common/app/job/main/usersuit/model"
	"go-common/library/log"
)

// UpInfoData .
func (d *Dao) UpInfoData(c context.Context) (res *model.UpInfo, err error) {
	params := url.Values{}
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	if err = d.client.Get(c, fmt.Sprintf(d.updateInfo, yesterday), "", params, &res); err != nil {
		log.Error("GetWearedfansMedal(%s) error(%v)", d.updateInfo+"?"+params.Encode(), err)
		return
	}
	if res.Code != 0 {
		log.Error("GetWearedfansMedal(%s) res(%d)", d.updateInfo+"?"+params.Encode(), res.Code)
	}
	log.Info("GetWearedfansMedal(%s) res(%+v)", d.updateInfo+"?"+params.Encode(), res)
	return
}
