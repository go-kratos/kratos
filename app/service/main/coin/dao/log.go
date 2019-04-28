package dao

import (
	"context"
	"encoding/json"
	"time"

	pb "go-common/app/service/main/coin/api"
	"go-common/app/service/main/coin/model"
	"go-common/library/database/elastic"
	"go-common/library/log"
	"go-common/library/queue/databus/report"
)

// AddLog .
func (d *Dao) AddLog(mid, ts int64, from, to float64, reason, ip, operator string, oid int64, typ int) {
	report.User(&report.UserInfo{
		Mid:      mid,
		Business: model.ReportType,
		IP:       ip,
		Type:     typ,
		Ctime:    time.Unix(ts, 0),
		Index:    []interface{}{operator},
		Oid:      oid,
		Content: map[string]interface{}{
			"from":   int64(from * _multi),
			"to":     int64(to * _multi),
			"reason": reason,
		},
	})
}

// CoinLog .
func (d *Dao) CoinLog(c context.Context, mid int64) (ls []*pb.ModelLog, err error) {
	t := time.Now()
	from := t.Add(-time.Hour * 24 * 7)
	var res struct {
		Page struct {
			Num   int `json:"num"`
			Size  int `json:"size"`
			Total int `json:"total"`
		} `json:"page"`
		Result []*report.UserActionLog `json:"result"`
	}
	err = d.es.NewRequest("log_user_action").
		IndexByTime("log_user_action_21", elastic.IndexTypeWeek, from, time.Now()).
		WhereEq("mid", mid).
		WhereRange("ctime", from.Format("2006-01-02 15:04:05"), "", elastic.RangeScopeLcRc).
		Pn(1).Ps(1000).Order("ctime", elastic.OrderDesc).
		Scan(c, &res)
	if err != nil {
		log.Error("coineslog mid: %v, err: %v", mid, err)
		PromError("log:es")
		return
	}
	for _, r := range res.Result {
		ts, _ := time.ParseInLocation("2006-01-02 15:04:05", r.Ctime, time.Local)
		var ex struct {
			From   int64
			To     int64
			Reason string
		}
		json.Unmarshal([]byte(r.Extra), &ex)
		ls = append(ls, &pb.ModelLog{
			From:      float64(ex.From) / _multi,
			To:        float64(ex.To) / _multi,
			IP:        r.IP,
			Desc:      ex.Reason,
			TimeStamp: ts.Unix(),
		})
	}
	return
}
