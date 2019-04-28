package flowmonitor

import (
	"time"
	"errors"

	xtime "go-common/library/time"
)

type Config struct {
	Interval xtime.Duration
	Addr     string
}

func (fm *FlowMonitor) checkConfig() (err error) {
	if fm.conf == nil {
		return errors.New("config for flowmonitor is nil")
	}
	if fm.conf.Interval == 0 {
		fm.conf.Interval = xtime.Duration(time.Second * 5)
	}
	if fm.conf.Addr == "" {
		return errors.New("addr of flowmonitor is nil")
	}
	return
}
