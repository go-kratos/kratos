package databus

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"go-common/app/admin/main/aegis/model"
	"go-common/app/admin/main/aegis/model/resource"
	"go-common/library/conf/env"
	"go-common/library/log"
	"go-common/library/queue/databus"
)

//var
var (
	ErrInfo     = errors.New("error info")
	ErrNilAgent = errors.New("nil agent")

	_formtTime = "2006-01-02 15:04:05"
	agent      *databus.Databus
)

//RscMsg .
type RscMsg struct {
	Action string          `json:"action"`
	BizID  int64           `json:"business_id"`
	Raw    json.RawMessage `json:"raw"`
}

//AddInfo info for add
type AddInfo struct {
	BusinessID int64  `json:"business_id"`
	NetID      int64  `json:"net_id"`
	OID        string `json:"oid"`
	MID        int64  `json:"mid"`
	Content    string `json:"content"`
	Extra1     int64  `json:"extra1"`
	Extra2     int64  `json:"extra2"`
	Extra3     int64  `json:"extra3"`
	Extra4     int64  `json:"extra4"`
	Extra5     int64  `json:"extra5"`
	Extra6     int64  `json:"extra6"`
	Extra1s    string `json:"extra1s"`
	Extra2s    string `json:"extra2s"`
	Extra3s    string `json:"extra3s"`
	Extra4s    string `json:"extra4s"`
	MetaData   string `json:"metadata"`
	ExtraTime1 time.Time
	OCtime     time.Time
	Ptime      time.Time
}

//UpdateInfo info for update
type UpdateInfo = model.UpdateOption

//CancelInfo info for cancel
type CancelInfo = model.CancelOption

//InitAegis .
func InitAegis(c *databus.Config) {
	if c == nil {
		c = _defaultConfig
		if d, ok := _defaultAddrConfig[env.DeployEnv]; ok {
			c.Key = d.Key
			c.Secret = d.Secret
			c.Addr = d.Addr
		}
	}
	if agent != nil {
		agent.Close()
	}
	agent = databus.New(c)
}

//CloseAegis .
func CloseAegis() {
	if agent != nil {
		agent.Close()
	}
}

//Add .
func Add(m *AddInfo) (err error) {
	if agent == nil {
		return ErrNilAgent
	}
	if m == nil || m.BusinessID <= 0 || m.NetID <= 0 || len(m.OID) == 0 {
		return ErrInfo
	}
	opt := &model.AddOption{
		Resource: resource.Resource{
			BusinessID: m.BusinessID,
			OID:        m.OID,
			MID:        m.MID,
			Content:    m.Content,
			Extra1:     m.Extra1,
			Extra2:     m.Extra2,
			Extra3:     m.Extra3,
			Extra4:     m.Extra4,
			Extra5:     m.Extra5,
			Extra6:     m.Extra6,
			Extra1s:    m.Extra1s,
			Extra2s:    m.Extra2s,
			Extra3s:    m.Extra3s,
			Extra4s:    m.Extra4s,
			ExtraTime1: m.ExtraTime1.Format(_formtTime),
			OCtime:     m.OCtime.Format(_formtTime),
			Ptime:      m.Ptime.Format(_formtTime),
		},

		NetID: m.NetID,
	}

	msg, err := formatMsg(opt, m.BusinessID, "add")
	if err != nil {
		return ErrInfo
	}
	return aegisSend(context.Background(), m.OID, msg)
}

//Update .
func Update(m *UpdateInfo) (err error) {
	if agent == nil {
		return ErrNilAgent
	}
	if m == nil || m.BusinessID <= 0 || m.NetID <= 0 || len(m.OID) == 0 || len(m.Update) == 0 {
		return ErrInfo
	}

	msg, err := formatMsg(m, m.BusinessID, "update")
	if err != nil {
		return ErrInfo
	}
	return aegisSend(context.Background(), m.OID, msg)
}

//Cancel .
func Cancel(m *CancelInfo) (err error) {
	if agent == nil {
		return ErrNilAgent
	}
	if m == nil || m.BusinessID <= 0 || len(m.Oids) == 0 {
		return ErrInfo
	}

	msg, err := formatMsg(m, m.BusinessID, "cancel")
	if err != nil {
		return ErrInfo
	}
	return aegisSend(context.Background(), m.Oids[0], msg)
}

//aegisSend .
func aegisSend(c context.Context, key string, msg interface{}) (err error) {
	log.Info("start to send key(%s) msg(%+v)", key, msg)

	for retry := 0; retry < 3; retry++ {
		if err = agent.Send(c, key, msg); err == nil {
			break
		}
	}
	if err != nil {
		log.Error("s.aegisPub.Send(%s) error(%v) msg(%+v) ", key, err, msg)
	}

	return
}

func formatMsg(m interface{}, bizid int64, action string) (*RscMsg, error) {
	raw, err := json.Marshal(m)
	if err != nil {
		return nil, ErrInfo
	}
	msg := &RscMsg{
		Action: action,
		BizID:  bizid,
		Raw:    raw,
	}
	return msg, nil
}
