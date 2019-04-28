package region

import (
	"bytes"
	"strconv"
	"time"

	"go-common/app/interface/main/app-feed/model/tag"
	"go-common/library/log"
	binfoc "go-common/library/log/infoc"
)

type tagsInfoc struct {
	typ    string
	mid    string
	client string
	build  string
	buvid  string
	disid  string
	ip     string
	api    string
	now    string
	tags   []*tag.Tag
}

type tagInfoc struct {
	typ    string
	mid    string
	client string
	build  string
	buvid  string
	disid  string
	ip     string
	api    string
	now    string
	rid    string
	tid    string
}

func (s *Service) TagsInfoc(mid int64, plat int8, build int, buvid, disid, ip, api string, tags []*tag.Tag, now time.Time) {
	select {
	case s.logCh <- tagsInfoc{"推荐页入口", strconv.FormatInt(mid, 10), strconv.Itoa(int(plat)), strconv.Itoa(build), buvid, disid, ip, api, strconv.FormatInt(now.Unix(), 10), tags}:
	default:
		log.Warn("infoc log buffer is full")
	}
}

func (s *Service) ChangeTagsInfoc(mid int64, plat int8, build int, buvid, disid, ip, api string, tags []*tag.Tag, now time.Time) {
	select {
	case s.logCh <- tagsInfoc{"换一换", strconv.FormatInt(mid, 10), strconv.Itoa(int(plat)), strconv.Itoa(build), buvid, disid, ip, api, strconv.FormatInt(now.Unix(), 10), tags}:
	default:
		log.Warn("infoc log buffer is full")
	}
}

func (s *Service) AddTagInfoc(mid int64, plat int8, build int, buvid, disid, ip, api string, rid int, tid int64, now time.Time) {
	select {
	case s.logCh <- tagInfoc{"订阅标签", strconv.FormatInt(mid, 10), strconv.Itoa(int(plat)), strconv.Itoa(build), buvid, disid, ip, api, strconv.FormatInt(now.Unix(), 10), strconv.Itoa(rid), strconv.FormatInt(tid, 10)}:
	default:
		log.Warn("infoc log buffer is full")
	}
}

func (s *Service) CancelTagInfoc(mid int64, plat int8, build int, buvid, disid, ip, api string, rid int, tid int64, now time.Time) {
	select {
	case s.logCh <- tagInfoc{"取消标签", strconv.FormatInt(mid, 10), strconv.Itoa(int(plat)), strconv.Itoa(build), buvid, disid, ip, api, strconv.FormatInt(now.Unix(), 10), strconv.Itoa(rid), strconv.FormatInt(tid, 10)}:
	default:
		log.Warn("infoc log buffer is full")
	}
}

// writeInfoc
func (s *Service) infocproc() {
	const (
		// infoc format {"section":{"id":"%s","pos":1,"items":[{"id":%s,"pos":%d,"url":""}]}}
		noItem1 = `{"section":{"id":"`
		noItem2 = `","pos":1,"items":[]}}`
	)
	var (
		msg1 = []byte(`{"section":{"id":"`)
		msg2 = []byte(`","pos":1,"items":[`)
		msg3 = []byte(`{"id":`)
		msg4 = []byte(`,"name":"`)
		msg5 = []byte(`","url":""},`)
		msg6 = []byte(`]}}`)
		inf2 = binfoc.New(s.c.TagInfoc2)
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
		case tagsInfoc:
			if len(v.tags) > 0 {
				buf.Write(msg1)
				buf.WriteString(v.typ)
				buf.Write(msg2)
				for _, v := range v.tags {
					buf.Write(msg3)
					buf.WriteString(strconv.FormatInt(v.ID, 10))
					buf.Write(msg4)
					buf.WriteString(v.Name)
					buf.Write(msg5)
				}
				buf.Truncate(buf.Len() - 1)
				buf.Write(msg6)
				list = buf.String()
				buf.Reset()
			} else {
				list = noItem1 + v.typ + noItem2
			}
			inf2.Info(v.ip, v.now, v.api, v.buvid, v.mid, v.client, "1", list, v.disid, "1", v.build)
		case tagInfoc:
			inf2.Info(v.ip, v.now, v.api, v.buvid, v.mid, v.client, "1", list, v.disid, "1", v.build)
		}
	}
}
