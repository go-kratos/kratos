package region

import (
	"bytes"
	"strconv"
	"time"

	"go-common/app/interface/main/app-show/conf"
	"go-common/library/log"
	binfoc "go-common/library/log/infoc"
)

type infoc struct {
	mid     string
	rid     string
	tid     string
	pn      string
	hotavid []int64
	newavid []int64
	now     string
}

// Infoc write data for Hadoop do analytics
func (s *Service) infoc(mid int64, hotavid, newavid []int64, rid, tid int, pull bool, now time.Time) {
	var pn string
	if pull {
		pn = "1"
	} else {
		pn = "2"
	}
	select {
	case s.logCh <- infoc{strconv.FormatInt(mid, 10), strconv.Itoa(rid), strconv.Itoa(tid), pn, hotavid, newavid, strconv.FormatInt(now.Unix(), 10)}:
	default:
		log.Warn("infoc log buffer is full")
	}
}

// writeInfoc
func (s *Service) infocproc() {
	const (
		noItem1 = `{"section":{"rid":`
		noItem2 = `,"tagid":`
		noItem3 = `,"mid":`
		noItem4 = `,"pn":`
		noItem5 = `,"hot_avids":[],"item_avids":[]}}`
	)
	var (
		msg1 = []byte(`{"section":{"rid":`)
		msg2 = []byte(`,"tagid":`)
		msg3 = []byte(`,"mid":`)
		msg4 = []byte(`,"pn":`)
		msg5 = []byte(`,"hot_avids":[`)
		msg6 = []byte(`,`)
		msg7 = []byte(`],"item_avids":[`)
		msg8 = []byte(`,`)
		msg9 = []byte(`]}}`)
		inf2 = binfoc.New(conf.Conf.FeedInfoc2)
		buf  bytes.Buffer
		list string
	)
	for {
		i, ok := <-s.logCh
		if !ok {
			log.Warn("infoc proc exit")
			return
		}
		switch v := i.(type) {
		case infoc:
			if len(v.newavid) > 0 {
				buf.Write(msg1)
				buf.WriteString(v.rid)
				buf.Write(msg2)
				buf.WriteString(v.tid)
				buf.Write(msg3)
				buf.WriteString(v.mid)
				buf.Write(msg4)
				buf.WriteString(v.pn)
				buf.Write(msg5)
				for _, v := range v.hotavid {
					buf.WriteString(strconv.FormatInt(v, 10))
					buf.Write(msg6)
				}
				if len(v.hotavid) > 0 {
					buf.Truncate(buf.Len() - 1)
				}
				buf.Write(msg7)
				for _, v := range v.newavid {
					buf.WriteString(strconv.FormatInt(v, 10))
					buf.Write(msg8)
				}
				buf.Truncate(buf.Len() - 1)
				buf.Write(msg9)
				list = buf.String()
				buf.Reset()
			} else {
				list = noItem1 + v.rid + noItem2 + v.tid + noItem3 + v.mid + noItem4 + v.pn + noItem5
			}
			inf2.Info(v.now, list)
		}
	}
}
