package view

import (
	"bytes"
	"fmt"
	"strconv"
	"time"

	"go-common/app/interface/main/app-view/conf"
	"go-common/app/interface/main/app-view/model/view"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/log/infoc"
)

type viewInfoc struct {
	mid       string
	client    string
	build     string
	buvid     string
	disid     string
	ip        string
	api       string
	now       string
	aid       string
	err       string
	from      string
	trackid   string
	autoplay  string
	fromSpmid string
	spmid     string
}

type relateInfoc struct {
	mid         string
	aid         string
	client      string
	buvid       string
	disid       string
	ip          string
	api         string
	now         string
	isRcmmnd    int8
	rls         []*view.Relate
	trackid     string
	build       string
	returnCode  string
	userFeature string
	from        string
}

//ViewInfoc view infoc
func (s *Service) ViewInfoc(mid int64, plat int, trackid, aid, ip, api, build, buvid, disid, from string, now time.Time, err error, autoplay int, spmid, fromSpmid string) {
	s.infoc(viewInfoc{strconv.FormatInt(mid, 10), strconv.Itoa(plat), build, buvid, disid, ip, api, strconv.FormatInt(now.Unix(), 10), aid, strconv.Itoa(ecode.Cause(err).Code()), from, trackid, strconv.Itoa(autoplay), fromSpmid, spmid})
}

// RelateInfoc Relate Infoc
func (s *Service) RelateInfoc(mid, aid int64, plat int, trackid, build, buvid, disid, ip, api, returnCode, userFeature, from string, rls []*view.Relate, now time.Time, isRec int8) {
	s.infoc(relateInfoc{strconv.FormatInt(mid, 10), strconv.FormatInt(aid, 10), strconv.Itoa(plat), buvid, disid, ip, api, strconv.FormatInt(now.Unix(), 10), isRec, rls, trackid, build, returnCode, userFeature, from})
}

func (s *Service) infoc(i interface{}) {
	select {
	case s.inCh <- i:
	default:
		log.Warn("cacheproc chan full")
	}
}

// WriteViewInfoc Write View Infoc
func (s *Service) infocproc() {
	const (
		noItem = `{"section":{"id":"相关视频","pos":1,"from_item":"%s","items":[]}}`
	)
	var (
		msg1      = []byte(`{"section":{"id":"相关视频","pos":1,"from_item":"`)
		msg2      = []byte(`","items":[`)
		msg3      = []byte(`{"id":`)
		msg4      = []byte(`,"pos":`)
		msg5      = []byte(`,"goto":"`)
		msg6      = []byte(`","from":"`)
		msg7      = []byte(`","source":"`)
		msg8      = []byte(`","av_feature":`)
		msg9      = []byte(`,"type":1,"url":""},`)
		infView   = infoc.New(conf.Conf.InfocView)
		infRelate = infoc.New(conf.Conf.InfocRelate)
		buf       bytes.Buffer
		list      string
	)
	for {
		i := <-s.inCh
		switch v := i.(type) {
		case viewInfoc:
			infView.Info(v.ip, v.now, v.api, v.buvid, v.mid, v.client, v.aid, v.disid, v.err, v.from, v.build, v.trackid, v.autoplay, v.fromSpmid, v.spmid)
		case relateInfoc:
			var trackID string
			if len(v.rls) > 0 {
				buf.Write(msg1)
				buf.WriteString(v.aid)
				buf.Write(msg2)
				for key, value := range v.rls {
					// trackid
					if value.TrackID != "" {
						trackID = value.TrackID
					}
					//list
					id, _ := strconv.ParseInt(value.Param, 10, 64)
					buf.Write(msg3)
					buf.WriteString(strconv.FormatInt(id, 10))
					buf.Write(msg4)
					buf.WriteString(strconv.Itoa(key + 1))
					buf.Write(msg5)
					buf.WriteString(value.Goto)
					buf.Write(msg6)
					buf.WriteString(value.From)
					buf.Write(msg7)
					buf.WriteString(value.Source)
					buf.Write(msg8)
					if value.AvFeature != nil {
						buf.Write(value.AvFeature)
					} else {
						buf.Write([]byte(`""`))
					}
					buf.Write(msg9)
				}
				buf.Truncate(buf.Len() - 1)
				buf.WriteString(`]}}`)
				list = buf.String()
				buf.Reset()
			} else {
				list = fmt.Sprintf(noItem, v.aid)
			}
			infRelate.Info(v.ip, v.now, v.api, v.buvid, v.mid, v.client, "2", list, v.disid, v.isRcmmnd, trackID, v.build, v.returnCode, v.userFeature, v.from)
		default:
			log.Warn("infocproc can't process the type")
		}
	}
}
