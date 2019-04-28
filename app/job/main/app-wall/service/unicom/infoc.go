package unicom

import (
	"strconv"
	"time"

	"go-common/app/job/main/app-wall/conf"
	"go-common/library/log"
	binfoc "go-common/library/log/infoc"
)

type infoc struct {
	usermob    string
	phone      string
	mid        string
	leve       string
	integral   string
	flow       string
	unicomtype string
	now        string
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

// unicomInfoc unicom user infoc
func (s *Service) unicomInfoc(usermob string, phone, leve, integral, flow int, utype string, mid int64, now time.Time) {
	select {
	case s.logCh[mid%s.c.ChanDBNum] <- infoc{usermob, strconv.Itoa(phone), strconv.FormatInt(mid, 10),
		strconv.Itoa(leve), strconv.Itoa(integral), strconv.Itoa(flow), utype, strconv.FormatInt(now.Unix(), 10)}:
	default:
		log.Warn("unicomInfoc log buffer is full")
	}
}

func (s *Service) unicomInfocproc(i int64) {
	var (
		cliChan = s.logCh[i]
		packinf = binfoc.New(conf.Conf.UnicomUserInfoc2)
	)
	for {
		i, ok := <-cliChan
		if !ok {
			log.Warn("infoc proc exit")
			return
		}
		switch v := i.(type) {
		case infoc:
			packinf.Info(v.now, "0", v.usermob, v.phone, v.mid, v.leve, v.integral, v.flow, v.unicomtype)
			log.Info("unicomInfocproc log mid(%v) phone(%v)", v.mid, v.phone)
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
