package service

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"time"

	"go-common/library/log"
	"go-common/library/queue/databus"
	"go-common/library/queue/databus/report"
)

// aiReq is a producer send file data to eroticism inspect serve
func (j *Job) aiReq(ctx context.Context, meta *Meta) {
	url := fmt.Sprintf(_downloadFmt, j.downloadHost, meta.Bucket, meta.Filename)
	msg := &AIReqMessage{URL: url, FileName: meta.Filename, Bucket: meta.Bucket}
	if err := j.AIYellowingProducer.Send(ctx, meta.Filename, msg); err != nil {
		log.Error("d.databus.Send(%s,%v) error(%v)", meta.Bucket, *msg, err)
		return
	}
	log.Info("send msg to ai success: key:(%s),value:(%+v)", meta.Filename, msg)
}

// aiResp is a consumer receive message from eroticism inspect serve
func (j *Job) aiResp(ctx context.Context) {
	var (
		msgs  <-chan *databus.Message
		ok    bool
		msg   *databus.Message
		aiMsg *AIRespMessage
		err   error
	)
	msgs = j.AIYellowingConsumer.Messages() //ai response message
	for {
		if msg, ok = <-msgs; !ok {
			log.Info("databus Consumer exit")
			break
		}
		if err = msg.Commit(); err != nil {
			log.Error("msg.Commit() error(%v)", err)
			continue

		}
		aiMsg = new(AIRespMessage)
		// todo: parse kafka message from ai and judge if add to db
		if err = json.Unmarshal(msg.Value, aiMsg); err != nil {
			log.Error("Job.run.Unmarshal(key:%s, value:%s),err(%v)", msg.Key, string(msg.Value), err)
			continue
		}
		log.Info("consume msg from ai success: key:(%s),value:(%+v)", msg.Key, aiMsg)
		if aiMsg.ErrorCode != 0 {
			log.Error("ai audit failed code: %d, message: %s", aiMsg.ErrorCode, aiMsg.ErrorMsg)
			continue
		}

		// deal ai result
		// sex > threshold will save in db
		sex := int(math.Round(aiMsg.Sex * 10000))
		politics := int(math.Round(aiMsg.Politics * 10000))
		blood := int(math.Round(aiMsg.Blood * 10000))
		violent := int(math.Round(aiMsg.Violent * 10000))
		if sex > j.Threshold.Sex || politics > j.Threshold.Politics || blood > j.Threshold.Blood || violent > j.Threshold.Violent {
			log.Warn("url:%s need save db points(sex:%d,politics:%d,blood:%d,violent:%d)", aiMsg.URL, sex, politics, blood, violent)
			if err = j.RetryAddRecord(aiMsg); err != nil {
				log.Error("j.RetryAddRecord(%+v) failed(%v)", aiMsg, err)
				continue
			} else {
				log.Info("call upload-admin success: %+v", aiMsg)
			}
		}
		//todo: send to audit log platform
		m := &report.ManagerInfo{
			Uname:    "ai",
			UID:      0,
			Business: 161, //bfs 上传审核
			Type:     1,   //ai审核
			Action:   "ai",
			Index:    []interface{}{aiMsg.URL, aiMsg.Bucket, aiMsg.FileName, int(math.Round(aiMsg.Sex * 10000)), int(math.Round(aiMsg.Politics * 10000))},
			Ctime:    time.Now(),
			Content: map[string]interface{}{ //分数
				"sex":       aiMsg.Sex,
				"violent":   aiMsg.Violent,
				"politics":  aiMsg.Politics,
				"blood":     aiMsg.Blood,
				"is_yellow": aiMsg.IsYellow,
			},
		}
		if err = report.Manager(m); err != nil {
			log.Error("report.Manager() m(%+v) error:%v", m, err)
			continue
		}
		log.Info("send log report success: (%+v) \n", m)
	}
}
