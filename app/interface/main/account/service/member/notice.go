package member

import (
	"context"
	"net"

	"go-common/app/interface/main/account/model"
	secmodel "go-common/app/service/main/secure/model"

	"go-common/library/log"
)

// NoticeV2 notice v2.
func (s *Service) NoticeV2(c context.Context, mid int64, uuid string, pf string, b int64) (res *model.Notice2, err error) {
	var msg *secmodel.Msg
	if msg, err = s.secureRPC.Status(c, &secmodel.ArgSecure{Mid: mid, UUID: uuid}); err != nil {
		log.Error("service.secureRPC.Status(%d) uuid(%s) error(%v)", mid, uuid, err)
		return
	}
	if msg == nil || !msg.Notify {
		log.Info("s.NoticeV2(%d) msg(%v) continue", mid, msg)
		res = &model.Notice2{Status: model.NoticeStatusNotNotify}
		return
	}
	res = &model.Notice2{
		Status: model.NoticeStatusNotify,
		Type:   model.NoticeTypeSecurity,
		Security: &model.Security{
			Location: msg.Log.Location,
			Time:     msg.Log.Time,
			IP:       inetNtoA(msg.Log.IP),
			Mid:      msg.Log.Mid,
		}}
	return
}

func inetNtoA(sum uint32) string {
	ip := make(net.IP, net.IPv4len)
	ip[0] = byte((sum >> 24) & 0xFF)
	ip[1] = byte((sum >> 16) & 0xFF)
	ip[2] = byte((sum >> 8) & 0xFF)
	ip[3] = byte(sum & 0xFF)
	return ip.String()
}

// CloseNoticeV2 close notice v2.
func (s *Service) CloseNoticeV2(c context.Context, mid int64, uuid string) (err error) {
	arg := &secmodel.ArgSecure{Mid: mid, UUID: uuid}
	if err = s.secureRPC.CloseNotify(c, arg); err != nil {
		log.Error("service.secureRPC.CloseNotify(%v) error(%v)", arg, err)
		return
	}
	return
}
