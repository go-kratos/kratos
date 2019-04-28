package service

import (
	"context"
	"encoding/json"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRecoverIdx(t *testing.T) {
	var c = context.Background()
	Convey("recover idx", t, WithService(func(s *Service) {
		var d struct {
			Oid   int64 `json:"oid"`
			Tp    int8  `json:"tp"`
			Count int   `json:"count"`
			Floor int   `json:"floor"`
		}
		d.Oid = 78
		d.Tp = 2
		d.Count = 20
		content, _ := json.Marshal(&d)
		s.acionRecoverFloorIdx(c, &consumerMsg{Data: content})
	}))
}
