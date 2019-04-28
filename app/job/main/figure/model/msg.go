package model

import (
	"encoding/json"
	"strconv"

	"go-common/library/log"

	"github.com/pkg/errors"
)

// MsgCanal canal message struct
type MsgCanal struct {
	Action string          `json:"action"`
	Table  string          `json:"table"`
	New    json.RawMessage `json:"new"`
	Old    json.RawMessage `json:"old"`
}

// MsgVipInfo message for user vip staus
type MsgVipInfo struct {
	Mid       int64 `json:"mid"`
	Type      int8  `json:"type"`
	Timestamp int64 `json:"ts"`
}

type MsgAccountLog struct {
	Mid     int64             `json:"mid"`
	IP      string            `json:"ip"`
	TS      int64             `json:"ts"`
	Content map[string]string `json:"content"`
}

func (m *MsgAccountLog) ExpFrom() (exp int) {
	var (
		fromExp = m.Content["from_exp"]
		err     error
	)
	if exp, err = strconv.Atoi(fromExp); err != nil {
		err = errors.Wrapf(err, "fromExp (%s)", fromExp)
		log.Error("%+v", err)
		exp = 0
	}
	return
}

func (m *MsgAccountLog) ExpTo() (exp int) {
	var (
		toExp = m.Content["to_exp"]
		err   error
	)
	if exp, err = strconv.Atoi(toExp); err != nil {
		err = errors.Wrapf(err, "toExp (%s)", toExp)
		log.Error("%+v", err)
		exp = 0
	}
	return
}

func (m *MsgAccountLog) IsViewExp() bool {
	var (
		operater = m.Content["operater"]
	)
	return operater == "watch"
}
