package service

import (
	"fmt"
	"go-common/app/interface/main/videoup/model/archive"
	"go-common/library/log"
	"go-common/library/queue/databus/report"
	"net"
	"strconv"
	"time"
)

// SendArchiveLog send to log service
func (s *Service) SendArchiveLog(aid int64, build, buvid, action, platform string, ap *archive.ArcParam, videoerr error) {
	var (
		LogType = archive.LogTypeSuccess
		index   = []interface{}{ap.TypeID, ap.UpFrom, ap.OrderID, ap.Title, ap.BizFrom, strconv.Itoa(ap.MissionID), fmt.Sprint(videoerr)}
	)
	if videoerr != nil {
		LogType = archive.LogTypeFail
	}
	buildNum, _ := strconv.Atoi(build)
	ip := net.IP(ap.IPv6[:]).String()
	uInfo := &report.UserInfo{
		Mid:      ap.Mid,
		Platform: platform,
		Build:    int64(buildNum),
		Buvid:    buvid,
		Business: archive.ArchiveAddLogID,
		Type:     LogType,
		Oid:      aid,
		Action:   action,
		Ctime:    time.Now(),
		Index:    index,
		IP:       ip,
	}
	ap.Aid = aid

	// es 最多存32k,我们日志就最多记100p吧
	if len(ap.Videos) > 100 {
		ap.Videos = ap.Videos[:100]
	}

	uInfo.Content = map[string]interface{}{
		"content": ap,
		"errmsg":  videoerr,
	}
	report.User(uInfo)
	log.Info("sendLog data(%+v)", uInfo)
}
