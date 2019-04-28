package collector

import (
	"testing"

	"go-common/app/service/main/dapper/model"
)

func TestPeerServiceDetect(t *testing.T) {
	detect := NewPeerServiceDetectProcesser(map[string]map[string]struct{}{
		"s_a": {
			"o_1":  struct{}{},
			"o_11": struct{}{},
		},
		"s_b": {"o_2": struct{}{}},
		"s_c": {"o_21": struct{}{}},
	})
	sp1, _ := model.FromProtoSpan(&model.ProtoSpan{ServiceName: "xxx", OperationName: "o_1"}, false)
	detect.Process(sp1)
	if sp1.StringTag("peer.service") != "s_a" && sp1.BoolTag("_auto.peer.service") {
		t.Errorf("expect get s_a get %s", sp1.StringTag("peer.service"))
	}
	if err := detect.Process(&model.Span{ServiceName: "hh", OperationName: ""}); err != nil {
		t.Errorf("expect get noting")
	}
}
