package model

import (
	"time"

	"go-common/library/log"
)

// DBOldElecPayTradeInfo .
type DBOldElecPayTradeInfo struct {
	ID      int64
	OrderID string
	AVID    string
}

// DBOldElecPayOrder .
type DBOldElecPayOrder struct {
	ID       int64  `json:"id"`
	APPID    int    `json:"app_id"`
	UPMID    int64  `json:"mid"`
	PayMID   int64  `json:"pay_mid"`
	OrderID  string `json:"order_no"`
	ElecNum  int64  `json:"elec_num"`
	Status   int    `json:"status"` // 订单状态，1.消费中 2.消费成功 3.消费失败
	CTimeStr string `json:"ctime"`
	MTimeStr string `json:"mtime"`
	CTime    time.Time
	MTime    time.Time
}

// IsPaid .
func (d *DBOldElecPayOrder) IsPaid() bool {
	return d.Status == 2
}

// IsHiddnRank .
func (d *DBOldElecPayOrder) IsHiddnRank() bool {
	return d.APPID == 19 // 动态互推
}

// ParseCTime .
func (d *DBOldElecPayOrder) ParseCTime() (t time.Time) {
	if !d.CTime.IsZero() {
		return d.CTime
	}
	var err error
	if t, err = time.ParseInLocation("2006-01-02 15:04:05", d.CTimeStr, time.Local); err != nil {
		log.Error("DBOldElecPayOrder ctime parse failed: %s, err: %+v", d.CTimeStr, err)
		t = time.Now()
	}
	return
}

// ParseMTime .
func (d *DBOldElecPayOrder) ParseMTime() (t time.Time) {
	if !d.MTime.IsZero() {
		return d.MTime
	}
	var err error
	if t, err = time.ParseInLocation("2006-01-02 15:04:05", d.MTimeStr, time.Local); err != nil {
		log.Error("DBOldElecPayOrder mtime parse failed: %s, err: %+v", d.MTimeStr, err)
		t = time.Now()
	}
	return
}

// DBOldElecMessage .
type DBOldElecMessage struct {
	ID       int64  `json:"id"`
	MID      int64  `json:"mid"`
	RefMID   int64  `json:"ref_mid"`
	RefID    int64  `json:"ref_id"`
	Message  string `json:"message"`
	AVID     string `json:"av_no"`
	DateVer  string `json:"date_version"` // yyyy-MM格式，年-月
	Type     int    `json:"type"`         // 留言类型, 1.用户对up主留言 2.up回复用户留言
	State    int    `json:"state"`        // 留言状态 0.未回复 1.已回复 2 已屏蔽
	CTimeStr string `json:"ctime"`
	MTimeStr string `json:"mtime"`
	CTime    time.Time
	MTime    time.Time
}

// ParseCTime .
func (d *DBOldElecMessage) ParseCTime() (t time.Time) {
	if !d.CTime.IsZero() {
		return d.CTime
	}
	var err error
	if t, err = time.ParseInLocation("2006-01-02 15:04:05", d.CTimeStr, time.Local); err != nil {
		log.Error("DBOldElecMessage ctime parse failed: %s, err: %+v", d.CTimeStr, err)
		t = time.Now()
	}
	return
}

// ParseMTime .
func (d *DBOldElecMessage) ParseMTime() (t time.Time) {
	if !d.MTime.IsZero() {
		return d.MTime
	}
	var err error
	if t, err = time.ParseInLocation("2006-01-02 15:04:05", d.MTimeStr, time.Local); err != nil {
		log.Error("DBOldElecMessage mtime parse failed: %s, err: %+v", d.MTimeStr, err)
		t = time.Now()
	}
	return
}

// DBOldElecUserSetting .
type DBOldElecUserSetting struct {
	ID        int64 `json:"id"`
	MID       int64 `json:"mid"`
	SettingID int   `json:"setting_id"`
	Status    int   `json:"status"`
}

// BitValue 返回该配置位==1的数值
func (d *DBOldElecUserSetting) BitValue() int32 {
	switch d.SettingID {
	case 1:
		return 0x1
	case 2:
		return 0x2
	default:
		log.Error("DBOldElecUserSetting unknown SettingID:%d, %+v", d.SettingID, d)
	}
	return 0
}

// DBElecMessage .
type DBElecMessage struct {
	ID      int64
	Ver     int64
	AVID    int64
	UPMID   int64
	PayMID  int64
	Message string
	Replied bool
	Hidden  bool
	CTime   time.Time
	MTime   time.Time
}

// DBElecReply .
type DBElecReply struct {
	ID     int64
	MSGID  int64
	Reply  string
	Hidden bool
	CTime  time.Time
	MTime  time.Time
}
