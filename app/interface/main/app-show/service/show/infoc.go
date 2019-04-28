package show

import (
	"bytes"
	"strconv"
	"time"

	"go-common/app/interface/main/app-show/conf"
	"go-common/app/interface/main/app-show/model"
	"go-common/app/interface/main/app-show/model/feed"
	"go-common/app/interface/main/app-show/model/show"
	"go-common/library/log"
	binfoc "go-common/library/log/infoc"
)

type infoc struct {
	mid      string
	client   string
	buvid    string
	disid    string
	ip       string
	api      string
	now      string
	isRcmmnd string
	items    []*show.Item
}

type feedInfoc struct {
	mobiApp    string
	device     string
	build      string
	now        string
	loginEvent string
	mid        string
	buvid      string
	page       string
	feed       []*feed.Item
}

// Infoc write data for Hadoop do analytics
func (s *Service) Infoc(mid int64, plat int8, buvid, disid, ip, api string, items []*show.Item, now time.Time) {
	select {
	case s.logCh <- infoc{strconv.FormatInt(mid, 10), strconv.Itoa(int(plat)), buvid, disid, ip, api, strconv.FormatInt(now.Unix(), 10), "1", items}:
	default:
		log.Warn("infoc log buffer is full")
	}
}

func (s *Service) infocfeed(i interface{}) {
	select {
	case s.logFeedCh <- i:
	default:
		log.Warn("infocfeed chan full")
	}
}

// writeInfoc
func (s *Service) infocproc() {
	const (
		// infoc format {"section":{"id":"热门推荐","pos":1,"items":[{"id":%s,"pos":%d,"type":1,"url":""}]}}
		noItem = `{"section":{"id":"热门推荐","pos":1,"items":[""]}}`
	)
	var (
		msg1 = []byte(`{"section":{"id":"热门推荐","pos":1,"items":[`)
		msg2 = []byte(`{"id":`)
		msg3 = []byte(`,"pos":`)
		msg4 = []byte(`,"type":1,"url":""},`)

		inf2 = binfoc.New(conf.Conf.Infoc2)
		buf  bytes.Buffer
		list string
	)
	for {
		i := <-s.logCh
		if len(i.items) > 0 {
			buf.Write(msg1)
			for i, v := range i.items {
				if v.Goto != model.GotoAv {
					continue
				}
				buf.Write(msg2)
				buf.WriteString(v.Param)
				buf.Write(msg3)
				buf.WriteString(strconv.Itoa(i + 1))
				buf.Write(msg4)
			}
			buf.Truncate(buf.Len() - 1)
			buf.WriteString(`]}}`)
			list = buf.String()
			buf.Reset()
		} else {
			list = noItem
		}
		inf2.Info(i.ip, i.now, i.api, i.buvid, i.mid, i.client, "1", list, i.disid, i.isRcmmnd)
	}
}

func (s *Service) infocfeedproc() {
	const (
		noItem = `[]`
	)
	var (
		msg1    = []byte(`[`)
		msg2    = []byte(`{"goto":"`)
		msg3    = []byte(`","param":"`)
		msg4    = []byte(`","uri":"`)
		msg5    = []byte(`","r_pos":`)
		msg6    = []byte(`,"from_type":"`)
		msg9    = []byte(`","corner_mark":`)
		msg10   = []byte(`,"rcmd_content":"`)
		msg11   = []byte(`","card_style":`)
		msg12   = []byte(`,"items":[`)
		msg13   = []byte(`{"goto":"`)
		msg14   = []byte(`","param":"`)
		msg17   = []byte(`","pos":`)
		msg15   = []byte(`},`)
		msg16   = []byte(`]`)
		msg7    = []byte(`},`)
		msg8    = []byte(`]`)
		buf     bytes.Buffer
		list    string
		feedInf = binfoc.New(s.c.FeedTabInfoc)
	)
	for {
		i, ok := <-s.logFeedCh
		if !ok {
			log.Warn("infoc proc exit")
			return
		}
		switch l := i.(type) {
		case *feedInfoc:
			if f := l.feed; len(f) == 0 {
				list = noItem
			} else {
				buf.Write(msg1)
				for _, item := range f {
					buf.Write(msg2)
					buf.WriteString(item.Goto)
					buf.Write(msg3)
					buf.WriteString(item.Param)
					buf.Write(msg4)
					buf.WriteString(item.URI)
					buf.Write(msg5)
					buf.WriteString(strconv.FormatInt(item.Idx, 10))
					buf.Write(msg6)
					buf.WriteString(item.FromType)
					buf.Write(msg9)
					buf.WriteString(strconv.Itoa(int(item.CornerMark)))
					buf.Write(msg10)
					buf.WriteString(item.RcmdContent)
					buf.Write(msg11)
					buf.WriteString(strconv.Itoa(int(item.CardStyle)))
					if len(item.Item) > 0 {
						buf.Write(msg12)
						for pos, it := range item.Item {
							buf.Write(msg13)
							buf.WriteString(it.Goto)
							buf.Write(msg14)
							buf.WriteString(it.Param)
							buf.Write(msg17)
							buf.WriteString(strconv.Itoa(pos + 1))
							buf.Write(msg15)
						}
						buf.Truncate(buf.Len() - 1)
						buf.Write(msg16)
					}
					buf.Write(msg7)
				}
				buf.Truncate(buf.Len() - 1)
				buf.Write(msg8)
				list = buf.String()
				buf.Reset()
			}
			log.Info("showtab_infoc_index(%s,%s,%s,%s,%s,%s,%s,%s)_list(%s)", l.mobiApp, l.device, l.build, l.now, l.loginEvent, l.mid, l.buvid, l.page, list)
			feedInf.Info(l.mobiApp, l.device, l.build, l.now, l.loginEvent, l.mid, l.buvid, list, l.page)
		}
	}
}
