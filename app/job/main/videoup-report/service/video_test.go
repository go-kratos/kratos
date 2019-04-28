package service

import (
	"go-common/app/job/main/videoup-report/conf"
	"testing"
)

// TestHdlTraffic Test calculate video check spend time function.
func TestHdlTraffic(t *testing.T) {
	err := conf.Init()
	if err != nil {
		return
	}
	s := New(conf.Conf)
	s.hdlTraffic()
}
