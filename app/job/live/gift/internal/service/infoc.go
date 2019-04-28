package service

import (
	"go-common/library/log"
	"go-common/library/log/infoc"
	"strconv"
	"time"
)

var inCh = make(chan interface{}, 10240)

const maxInt = int(^uint(0) >> 1)

type bagLogInfoc struct {
	id       string
	uid      string
	bagID    string
	giftID   string
	num      string
	afterNum string
	source   string
	infoType string
	ctime    string
}

type giftActionInfoc struct {
	uid       int64
	roomid    int64
	item      int64
	value     int64
	change    int64
	describe  string
	extra     string
	ts        int64
	platform  string
	clientver string
	buvid     string
	ua        string
	referer   string
}

// bagLogInfoc 包裹日志打点
func (s *Service) bagLogInfoc(uid, bagID, giftID, num, afterNum int64, source string) {
	s.infoc(bagLogInfoc{
		id:       MakeID(uid),
		uid:      strconv.FormatInt(uid, 10),
		bagID:    strconv.FormatInt(bagID, 10),
		giftID:   strconv.FormatInt(giftID, 10),
		num:      strconv.FormatInt(num, 10),
		afterNum: strconv.FormatInt(afterNum, 10),
		source:   source,
		infoType: "1",
		ctime:    time.Now().Format("2006-01-02 15:04:05"),
	})
}

//giftActionInfoc 道具打点
func (s *Service) giftActionInfoc(uid, roomid, item, value, change int64, describe, platform string) {
	s.infoc(giftActionInfoc{
		uid:       uid,
		roomid:    roomid,
		item:      item,
		value:     value,
		change:    change,
		describe:  describe,
		extra:     "",
		ts:        time.Now().Unix(),
		platform:  platform,
		clientver: "",
		buvid:     "",
		ua:        "",
		referer:   "",
	})
}

// MakeID MakeID
func MakeID(uid int64) string {
	prefix := strconv.FormatInt(uid%10, 10)
	postfix := strconv.Itoa(maxInt - int(time.Now().Unix()*10000))
	uidStr := strconv.FormatInt(uid, 10)
	l := len(uidStr)
	var middle string
	if l >= 10 {
		middle = uidStr
	} else {
		var s string
		for i := 0; i < (10 - l); i++ {
			s += "0"
		}
		middle = s + uidStr
	}
	return prefix + middle + postfix
}

//infoc
func (s *Service) infoc(i interface{}) {
	select {
	case inCh <- i:
	default:
		log.Warn("infocproc chan full")
	}
}

// infocproc
func (s *Service) infocproc() {
	var bl = infoc.New(s.c.Infoc["bagLog"])
	var ga = infoc.New(s.c.Infoc["giftAction"])
	for {
		i := <-inCh
		switch v := i.(type) {
		case bagLogInfoc:
			err := bl.Info(v.id, v.uid, v.bagID, v.giftID, v.num, v.afterNum, v.source, v.infoType, v.ctime)
			log.Info("bagLogInfoc info %v,ret:%v", v, err)
		case giftActionInfoc:
			err := ga.Info(v.uid, v.roomid, v.item, v.value, v.change, v.describe, v.extra, v.ts, v.platform, v.clientver, v.buvid, v.ua, v.referer)
			log.Info("giftActionInfoc info %v,ret:%v", v, err)
		default:
			log.Warn("infocproc can't process the type")
		}
	}
}
