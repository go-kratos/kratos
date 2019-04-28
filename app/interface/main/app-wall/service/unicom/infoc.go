package unicom

import (
	"strconv"
	"time"

	"go-common/app/interface/main/app-wall/conf"
	log "go-common/library/log"
	binfoc "go-common/library/log/infoc"
)

type orderInfoc struct {
	usermob   string
	orderType string
	ip        string
	mobiApp   string
	build     string
	now       string
}

type ipInfoc struct {
	usermob  string
	isValide string
	ip       string
	mobiApp  string
	build    string
	now      string
}

type packInfoc struct {
	usermob      string
	phone        string
	mid          string
	requestNo    string
	packName     string
	packIntegral string
	packType     string
	now          string
}

// Infoc write data for Hadoop do analytics
func (s *Service) unicomInfoc(mobiApp, usermob, ip string, build, orderType int, now time.Time) {
	select {
	case s.logCh <- orderInfoc{usermob, strconv.Itoa(orderType), ip, mobiApp, strconv.Itoa(build), strconv.FormatInt(now.Unix(), 10)}:
	default:
		log.Warn("unicomInfoc log buffer is full")
	}
}

// Infoc write data for Hadoop do analytics
func (s *Service) ipInfoc(mobiApp, usermob, ip string, build int, isValide bool, now time.Time) {
	select {
	case s.logCh <- ipInfoc{usermob, strconv.FormatBool(isValide), ip, mobiApp, strconv.Itoa(build), strconv.FormatInt(now.Unix(), 10)}:
	default:
		log.Warn("ipInfoc log buffer is full")
	}
}

func (s *Service) unicomInfocproc() {
	var (
		unicominf2 = binfoc.New(conf.Conf.UnicomUserInfoc2)
		ipinf2     = binfoc.New(conf.Conf.UnicomIpInfoc2)
	)
	for {
		i, ok := <-s.logCh
		if !ok {
			log.Warn("infoc proc exit")
			return
		}
		switch v := i.(type) {
		case orderInfoc:
			unicominf2.Info(v.now, "0", v.usermob, v.orderType, v.ip, v.mobiApp, v.build)
		case ipInfoc:
			ipinf2.Info(v.now, "0", v.isValide, v.ip, v.usermob, v.mobiApp, v.build)
		}
	}
}

// unicomPackInfoc unicom pack infoc
func (s *Service) unicomPackInfoc(usermob, packName, orderNumber string, phone, packIntegral, packType int, mid int64, now time.Time) {
	select {
	case s.packCh <- packInfoc{usermob, strconv.Itoa(phone), strconv.FormatInt(mid, 10),
		orderNumber, packName, strconv.Itoa(packIntegral), strconv.Itoa(packType), strconv.FormatInt(now.Unix(), 10)}:
	default:
		log.Warn("unicomPackInfoc log buffer is full")
	}
}

func (s *Service) unicomPackInfocproc() {
	var (
		packinf = binfoc.New(conf.Conf.UnicomPackInfoc)
	)
	for {
		i, ok := <-s.packCh
		if !ok {
			log.Warn("infoc proc exit")
			return
		}
		switch v := i.(type) {
		case packInfoc:
			packinf.Info(v.now, "0", v.usermob, v.phone, v.mid, v.requestNo, v.packName, v.packIntegral, v.packType)
		}
	}
}
