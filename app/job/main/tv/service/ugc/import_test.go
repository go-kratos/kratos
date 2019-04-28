package ugc

import (
	"encoding/json"
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestService_InitUpper(t *testing.T) {
	Convey("TestService_InitUpper", t, WithService(func(s *Service) {
		err := s.InitUpper(10920044)
		So(err, ShouldBeNil)
	}))
}

func TestService_UpArchives(t *testing.T) {
	Convey("TestService_UpArchives", t, WithService(func(s *Service) {
		pMatch, pAids, err := s.UpArchives(10920044, 1, 100)
		So(err, ShouldBeNil)
		data, _ := json.Marshal(pMatch)
		data2, _ := json.Marshal(pAids)
		fmt.Println(string(data))
		fmt.Println(string(data2))
	}))
}
