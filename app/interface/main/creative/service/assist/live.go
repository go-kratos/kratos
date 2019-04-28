package assist

import (
	"context"
	"go-common/library/log"
)

// LiveStatus get user assist rights about liveRoom
func (s *Service) LiveStatus(c context.Context, mid int64, ip string) (open int8, err error) {
	if open, err = s.assist.LiveStatus(c, mid, ip); err != nil {
		log.Error("s.assist.HasLiveRight mid(%d), ip(%s)", mid, ip)
		return
	}
	return
}

// liveAddAssist add assist to live
// Notice: 这里是新账号系统的Demo
func (s *Service) liveAddAssist(c context.Context, mid, assistMid int64, ak, ck, ip string) (err error) {
	identified, _ := s.acc.IdentifyInfo(c, mid, 1, ip)
	if err = s.acc.CheckIdentify(identified); err != nil {
		log.Error("s.acc.IdentifyInfo mid(%d),ip(%s)", mid, ip)
		return
	}
	if err = s.assist.LiveAddAssist(c, mid, assistMid, ck, ip); err != nil {
		return
	}
	return
}

// liveDelAssist del assist to live
func (s *Service) liveDelAssist(c context.Context, mid, assistMid int64, ck, ip string) (err error) {
	if err = s.assist.LiveDelAssist(c, mid, assistMid, ck, ip); err != nil {
		return
	}
	return
}

// LiveCheckAssist check if is assist in live
func (s *Service) LiveCheckAssist(c context.Context, mid, assistMid int64, ip string) (isAss int8, err error) {
	if isAss, err = s.assist.LiveCheckAssist(c, mid, assistMid, ip); err != nil {
		return
	}
	return
}

// LiveRevocBanned revoke banned in live
func (s *Service) LiveRevocBanned(c context.Context, mid int64, banID, ck, ip string) (err error) {
	if err = s.assist.LiveBannedRevoc(c, mid, banID, ck, ip); err != nil {
		return
	}
	return
}
