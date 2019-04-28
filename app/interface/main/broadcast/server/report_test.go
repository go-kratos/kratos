package server

import "testing"

func TestReportCh(t *testing.T) {
	ch := &Channel{Key: "test_key", Mid: 19158909, IP: "127.0.0.1", Platform: "web", Room: &Room{ID: "video://111/222", AllOnline: 100}}
	reportCh(1, ch)
}
