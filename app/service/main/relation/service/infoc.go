package service

import (
	"go-common/library/log"
	"go-common/library/log/infoc"
	"strconv"
)

const (
	_prefix         = "/x/internal/relation"
	addFollowingURL = _prefix + "/following/add"
	delFollowingURL = _prefix + "/following/del"
	addWhisperURL   = _prefix + "/whisper/add"
	delWhisperURL   = _prefix + "/whisper/del"
	addBlackURL     = _prefix + "/black/add"
	delBlackURL     = _prefix + "/black/del"
	delFollowerURL  = _prefix + "/follower/del"
	// RelInfocIP map key string
	RelInfocIP = "ip"
	// RelInfocSid RelInfocSid
	RelInfocSid = "sid"
	// RelInfocBuvid RelInfocBuvid
	RelInfocBuvid = "buvid"
	// RelInfocHeaderBuvid RelInfocHeaderBuvid
	RelInfocHeaderBuvid = "Buvid"
	// RelInfocCookieBuvid RelInfocCookieBuvid
	RelInfocCookieBuvid = "buvid3"
	// RelInfocReferer RelInfocReferer
	RelInfocReferer = "Referer"
	// RelInfocUA RelInfocUA
	RelInfocUA = "User-Agent"
)

type relationInfoc struct {
	ip    string //账号进行关注的ip
	mid   string //进行关注行为的账号
	fid   string //被关注的账号id
	ts    string //关注行为发生时的服务器时间
	sid   string //cookie里记录的sid，标识一次登录访问
	buvid string //移动端上报，再请求header里有，标识设备
	url   string //请求的关注接口
	refer string //浏览器上报的请求上级
	ua    string //访问的浏览器或客户端版本
	src   string //业务页面来源编号
}

// RelationInfoc relation related info to BigData for anti-spam
func (s *Service) RelationInfoc(mid, fid, ts int64, ip, sid, buvid, url, refer, ua string, src uint8) {
	s.infoc(relationInfoc{
		ip,
		strconv.FormatInt(mid, 10),
		strconv.FormatInt(fid, 10),
		strconv.FormatInt(ts, 10),
		sid,
		buvid,
		url,
		refer,
		ua,
		strconv.FormatUint(uint64(src), 10),
	})
}

//infoc
func (s *Service) infoc(i interface{}) {
	select {
	case s.inCh <- i:
	default:
		log.Warn("infocproc chan full")
	}
}

// infocproc
func (s *Service) infocproc() {
	var infoc2 = infoc.New(s.c.Infoc)
	for {
		i := <-s.inCh
		switch v := i.(type) {
		case relationInfoc:
			infoc2.Info(v.ip, v.mid, v.fid, v.ts, v.sid, v.buvid, v.url, v.refer, v.ua, v.src)
		default:
			log.Warn("infocproc can't process the type")
		}
	}
}
