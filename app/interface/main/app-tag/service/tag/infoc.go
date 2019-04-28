package tag

import (
	"bytes"
	"strconv"

	"go-common/app/interface/main/app-tag/model/region"
	"go-common/library/log"
	binfoc "go-common/library/log/infoc"
)

type feedInfoc struct {
	mobiApp    string
	device     string
	build      string
	now        string
	pull       string
	loginEvent string
	tagID      string
	tagName    string
	mid        string
	buvid      string
	displayID  string
	feed       *region.Show
	isRec      string
	topChannel string
}

// infoc write data for Hadoop do analytics
func (s *Service) infoc(i interface{}) {
	select {
	case s.logCh <- i:
	default:
		log.Warn("infocproc chan full")
	}
}

func (s *Service) infocfeedproc() {
	const (
		noItem = `[]`
	)
	var (
		msg1     = []byte(`[`)
		msg2     = []byte(`{"goto":"`)
		msg3     = []byte(`","param":"`)
		msg4     = []byte(`","uri":"`)
		msg5     = []byte(`","r_pos":`)
		msg6     = []byte(`,"from_type":"`)
		msg7     = []byte(`"},`)
		msg8     = []byte(`]`)
		buf      bytes.Buffer
		list     string
		feedInf2 = binfoc.New(s.c.FeedInfoc2)
	)
	for {
		i, ok := <-s.logCh
		if !ok {
			log.Warn("infoc proc exit")
			return
		}
		var pos int
		switch l := i.(type) {
		case *feedInfoc:
			if f := l.feed; len(f.New) == 0 {
				list = noItem
			} else {
				buf.Write(msg1)
				for _, item := range f.New {
					buf.Write(msg2)
					buf.WriteString(item.Goto)
					buf.Write(msg3)
					buf.WriteString(item.Param)
					buf.Write(msg4)
					buf.WriteString(item.URI)
					buf.Write(msg5)
					pos++
					buf.WriteString(strconv.Itoa(pos))
					buf.Write(msg6)
					buf.WriteString("recommend")
					buf.Write(msg7)
				}
				buf.Truncate(buf.Len() - 1)
				buf.Write(msg8)
				list = buf.String()
				buf.Reset()
			}
			log.Info("tag_infoc_index(%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s)_list(%s)", l.mobiApp, l.device, l.build, l.now, l.pull,
				l.loginEvent, l.tagID, l.tagName, l.mid, l.buvid, l.displayID, list)
			feedInf2.Info(l.mobiApp, l.device, l.build, l.now, l.pull, l.loginEvent, l.tagID, l.tagName, l.mid, l.buvid, l.displayID,
				list, "54", l.isRec, l.topChannel)
		}
	}
}
