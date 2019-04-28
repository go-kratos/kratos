package v1

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

// RelationInfoc relation related info to BigData for anti-spam
func (s *GiftService) bagLogInfoc(uid, bagID, giftID, num, afterNum int64, source string) {
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

// MakeID 生成上报id
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
func (s *GiftService) infoc(i interface{}) {
	select {
	case inCh <- i:
	default:
		log.Warn("infocproc chan full")
	}
}

// infocproc
func (s *GiftService) infocproc() {
	var infoc2 = infoc.New(s.conf.Infoc["bagLog"])
	for {
		i := <-inCh
		switch v := i.(type) {
		case bagLogInfoc:
			err := infoc2.Info(v.id, v.uid, v.bagID, v.giftID, v.num, v.afterNum, v.source, v.infoType, v.ctime)
			log.Info("infocproc info %v,ret:%v", v, err)
		default:
			log.Warn("infocproc can't process the type")
		}
	}
}
