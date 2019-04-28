package service

import (
	"encoding/json"
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_Relate(t *testing.T) {
	Convey("get app relate", t, WithService(func(s *Service) {
		tmp := s.relateCache
		byte, _ := json.Marshal(tmp)
		fmt.Println(string(byte))
		tmp2 := s.relatePgcMapCache
		byte2, _ := json.Marshal(tmp2)
		fmt.Println(string(byte2))
	}))
}
