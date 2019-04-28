package income

import (
	"bytes"
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	model "go-common/app/job/main/growup/model/income"
	task "go-common/app/job/main/growup/service"

	"go-common/library/log"
)

// BGMs map[avid][]*model.BGM,
func (s *Service) BGMs(c context.Context, limit int64) (bm map[int64][]*model.BGM, err error) {
	var id int64
	bm = make(map[int64][]*model.BGM)
	for {
		var bs []*model.BGM
		bs, id, err = s.dao.GetBGM(c, id, limit)
		if err != nil {
			return
		}
		if len(bs) == 0 {
			break
		}
		for _, b := range bs {
			if _, ok := bm[b.AID]; ok {
				bm[b.AID] = append(bm[b.AID], b)
			} else {
				bm[b.AID] = []*model.BGM{b}
			}
		}
	}
	return
}

func (s *Service) batchInsertBGMs(c context.Context, bs []*model.BGM) (err error) {
	var (
		buff    = make([]*model.BGM, 2000)
		buffEnd = 0
	)
	for _, b := range bs {
		buff[buffEnd] = b
		buffEnd++
		if buffEnd >= 2000 {
			values := bgmValues(buff[:buffEnd])
			buffEnd = 0
			_, err = s.dao.InsertBGM(c, values)
			if err != nil {
				return
			}
		}
	}
	if buffEnd > 0 {
		values := bgmValues(buff[:buffEnd])
		buffEnd = 0
		_, err = s.dao.InsertBGM(c, values)
	}
	return
}

func bgmValues(bs []*model.BGM) (values string) {
	var buf bytes.Buffer
	for _, b := range bs {
		buf.WriteString("(")
		buf.WriteString(strconv.FormatInt(b.MID, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(b.SID, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(b.AID, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(b.CID, 10))
		buf.WriteByte(',')
		buf.WriteString("'" + b.JoinAt + "'")
		buf.WriteByte(',')
		buf.WriteString("\"" + strings.Replace(b.Title, "\"", "\\\"", -1) + "\"")
		buf.WriteString(")")
		buf.WriteByte(',')
	}
	if buf.Len() > 0 {
		buf.Truncate(buf.Len() - 1)
	}
	values = buf.String()
	buf.Reset()
	return
}

const query = "{\"select\": [],\"page\": {\"skip\": %d,\"limit\": %d}}"

func (s *Service) getBGM(c context.Context) (bgms []*model.BGM, err error) {
	from, limit := 0, 1000
	var info []*model.BGM
	for {
		info, err = s.dp.SendBGMRequest(c, fmt.Sprintf(query, from, limit))
		if err != nil {
			log.Error("s.getBGM error(%v)", err)
			return
		}
		if len(info) == 0 {
			break
		}
		bgms = append(bgms, info...)
		from += len(info)
		time.Sleep(15 * time.Second) // 数据平台限流，api 10s内只能请求一次
	}
	log.Info("get bgm data total (%d) rows", from)
	return
}

// SyncBgmInfo sync bgm info from data platform
func (s *Service) SyncBgmInfo(c context.Context) (err error) {
	defer func() {
		task.GetTaskService().SetTaskStatus(c, task.TaskBgmSync, time.Now().AddDate(0, 0, -1).Format(_layout), err)
	}()
	err = s.delBGM(c)
	if err != nil {
		return
	}
	bgms, err := s.getBGM(c)
	if err != nil {
		return
	}
	return s.batchInsertBGMs(c, bgms)
}

func (s *Service) delBGM(c context.Context) (err error) {
	var limit int64 = 2000
	for {
		var rows int64
		rows, err = s.dao.DelBGM(c, limit)
		if err != nil {
			return
		}
		if rows < limit {
			break
		}
	}
	return
}
