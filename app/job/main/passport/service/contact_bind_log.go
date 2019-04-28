package service

import (
	"context"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"strings"
	"time"

	"go-common/app/job/main/passport/model"
	"go-common/library/log"
	"go-common/library/queue/databus"
)

func (s *Service) contactBindLogconsumeproc() {
	mergeRoutineNum := int64(s.c.Group.ContactBindLog.Num)
	for {
		msg, ok := <-s.dsContactBindLog.Messages()
		if !ok {
			log.Error("s.telBindlogconsumeproc closed")
			return
		}
		m := &message{data: msg}
		p := &model.BMsg{}
		if err := json.Unmarshal(msg.Value, p); err != nil {
			log.Error("json.Unmarshal(%s) error(%v)", string(msg.Value), err)
			continue
		}

		//m.object = p
		mid := int64(0)
		switch {
		case strings.HasPrefix(p.Table, _telBindTable):
			t := new(model.TelBindLog)
			if err := json.Unmarshal(p.New, t); err != nil {
				log.Error("json.Unmarshal(%s) error(%v)", string(p.New), err)
				continue
			}
			mid = t.Mid
			m.object = p
			log.Info("contactBindLogconsumeproc table:%s key:%s partition:%d offset:%d", p.Table, msg.Key, msg.Partition, msg.Offset)
		case strings.HasPrefix(p.Table, _emailBindTable):
			t := new(model.EmailBindLog)
			if err := json.Unmarshal(p.New, t); err != nil {
				log.Error("json.Unmarshal(%s) error(%v)", string(p.New), err)
				continue
			}
			mid = t.Mid
			m.object = p
			log.Info("contactBindLogconsumeproc table:%s key:%s partition:%d offset:%d", p.Table, msg.Key, msg.Partition, msg.Offset)
		default:
			log.Warn("unrecognized message: %+v", p)
			continue
		}

		if mid == 0 {
			log.Warn("invalid message: %+v", p)
			continue
		}
		s.contactBindLogMu.Lock()
		if s.contactBindLogHead == nil {
			s.contactBindLogHead = m
			s.contactBindLogLast = m
		} else {
			s.contactBindLogLast.next = m
			s.contactBindLogLast = m
		}
		s.contactBindLogMu.Unlock()

		// use specify goroutine to merge messages
		s.contactBindLogMergeChans[mid%mergeRoutineNum] <- m
		log.Info("contactBindLogconsumeproc key:%s partition:%d offset:%d", msg.Key, msg.Partition, msg.Offset)
	}
}

func (s *Service) contactBindLogcommitproc() {
	commits := make(map[int32]*databus.Message, s.c.Group.Log.Size)
	for {
		done := <-s.contactBindLogDoneChan
		// merge partitions to commit offset
		for _, d := range done {
			d.done = true
		}
		s.contactBindLogMu.Lock()
		for ; s.contactBindLogHead != nil && s.contactBindLogHead.done; s.contactBindLogHead = s.contactBindLogHead.next {
			commits[s.contactBindLogHead.data.Partition] = s.contactBindLogHead.data
		}
		s.contactBindLogMu.Unlock()
		for k, m := range commits {
			log.Info("logcommitproc committed, key:%s partition:%d offset:%d", m.Key, m.Partition, m.Offset)
			m.Commit()
			delete(commits, k)
		}
	}
}

func (s *Service) contactBindLogMergeproc(c chan *message) {
	var (
		max    = s.c.Group.ContactBindLog.Size
		merges = make([]*model.BMsg, 0, max)
		marked = make([]*message, 0, max)
		ticker = time.NewTicker(time.Duration(s.c.Group.Log.Ticker))
	)
	for {
		select {
		case msg, ok := <-c:
			if !ok {
				log.Error("s.contactBindLogMergeproc closed")
				return
			}
			p, assertOk := msg.object.(*model.BMsg)
			if !assertOk {
				log.Warn("s.contactBindLogMergeproc cannot convert BMsg")
				continue
			}
			//if p.Action != "insert" {
			//	continue
			//}
			if p.Action == "delete" {
				continue
			}
			log.Info("s.contactBindLogMergeproc: %+v", msg)
			switch {
			case strings.HasPrefix(p.Table, _telBindTable) || strings.HasPrefix(p.Table, _emailBindTable):
				merges = append(merges, p)
			default:
				log.Warn("unrecognized the message: %+v", p)
			}
			marked = append(marked, msg)
			if len(marked) < max && len(merges) < max {
				continue
			}
		case <-ticker.C:
		}
		if len(merges) > 0 {
			s.contactBindLogProcessMerges(merges)
			merges = make([]*model.BMsg, 0, max)
		}
		if len(marked) > 0 {
			s.contactBindLogDoneChan <- marked
			marked = make([]*message, 0, max)
		}
	}
}

func (s *Service) contactBindLogProcessMerges(bmsgs []*model.BMsg) {
	for _, msg := range bmsgs {
		log.Info("contactBindLogProcessMerges: %+v", msg.Table)
		switch {
		case strings.HasPrefix(msg.Table, _telBindTable):
			t := new(model.TelBindLog)
			if err := json.Unmarshal(msg.New, t); err != nil {
				log.Error("json.Unmarshal(%s) error(%v)", string(msg.New), err)
				continue
			}
			s.handleTelBindLog(t)
		case strings.HasPrefix(msg.Table, _emailBindTable):
			t := new(model.EmailBindLog)
			if err := json.Unmarshal(msg.New, t); err != nil {
				log.Error("json.Unmarshal(%s) error(%v)", string(msg.New), err)
				continue
			}
			s.handleEmailBindLog(t)
		}
	}
}

type userLogExtra struct {
	EncryptTel   string `json:"tel"`
	EncryptEmail string `json:"email"`
}

type userLog struct {
	Action    string `json:"action"`
	Mid       int64  `json:"mid"`
	Str0      string `json:"str_0"`
	ExtraData string `json:"extra_data"`
	Business  int    `json:"business"`
	CTime     string `json:"ctime"`
}

func (s *Service) handleTelBindLog(telLog *model.TelBindLog) (err error) {
	var bindLog *model.TelBindLog
	for {
		bindLog, err = s.d.QueryTelBindLog(telLog.ID)
		if err != nil {
			log.Error("QueryTelBindLog (%v) err(%v)", telLog, err)
			time.Sleep(100 * time.Millisecond)
			continue
		}
		break
	}
	if bindLog == nil || bindLog.ID == 0 {
		log.Warn("telephone log (%v) nil", bindLog)
		return
	}
	rt, err := s.encrypt(bindLog.Tel)
	if err != nil {
		log.Error("aesEncrypt(%v) error(%v)", bindLog, err)
		return
	}
	extraData := userLogExtra{
		EncryptTel: rt,
	}
	hash := sha1.New()
	hash.Write([]byte(bindLog.Tel))
	extraDataBytes, err := json.Marshal(extraData)
	if err != nil {
		log.Error("extraData (%v) json marshal err(%v)", extraData, err)
		return
	}
	uLog := userLog{
		Action:    "telBindLog",
		Mid:       bindLog.Mid,
		Str0:      base64.StdEncoding.EncodeToString(hash.Sum(s.hashSalt)),
		ExtraData: string(extraDataBytes),
		Business:  54,
		CTime:     time.Unix(bindLog.Timestamp, 0).Format("2006-01-02 15:04:05"),
	}

	for {
		if err = s.userLogPub.Send(context.Background(), bindLog.Tel, uLog); err != nil {
			log.Error("databus send(%v) error(%v)", uLog, err)
			time.Sleep(100 * time.Millisecond)
			continue
		}
		log.Info("uselog pub uLog: %+v", uLog)
		break
	}
	return
}

func (s *Service) handleEmailBindLog(emailLog *model.EmailBindLog) (err error) {
	var bindLog *model.EmailBindLog
	for {
		bindLog, err = s.d.QueryEmailBindLog(emailLog.ID)
		if err != nil {
			log.Error("QueryEmailBindLog (%v) err(%v)", emailLog, err)
			time.Sleep(100 * time.Millisecond)
			continue
		}
		break
	}
	if bindLog == nil || bindLog.ID == 0 {
		log.Warn("email log (%v) nil", bindLog)
		return
	}
	rt, err := s.encrypt(bindLog.Email)
	if err != nil {
		log.Error("aesEncrypt(%v) error(%v)", bindLog, err)
		return
	}
	extraData := userLogExtra{
		EncryptEmail: rt,
	}
	hash := sha1.New()
	hash.Write([]byte(bindLog.Email))
	extraDataBytes, err := json.Marshal(extraData)
	if err != nil {
		log.Error("extraData (%v) json marshal err(%v)", extraData, err)
		return
	}
	uLog := userLog{
		Action:    "emailBindLog",
		Mid:       bindLog.Mid,
		Str0:      base64.StdEncoding.EncodeToString(hash.Sum(s.hashSalt)),
		ExtraData: string(extraDataBytes),
		Business:  54,
		CTime:     time.Unix(bindLog.Timestamp, 0).Format("2006-01-02 15:04:05"),
	}

	for {
		if err = s.userLogPub.Send(context.Background(), bindLog.Email, uLog); err != nil {
			log.Error("databus send(%v) error(%v)", uLog, err)
			time.Sleep(100 * time.Millisecond)
			continue
		}
		log.Info("uselog pub uLog: %+v", uLog)
		break
	}
	return
}
