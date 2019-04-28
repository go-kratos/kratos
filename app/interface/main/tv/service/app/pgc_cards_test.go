package service

import (
	"testing"

	"fmt"

	. "github.com/smartystreets/goconvey/convey"
)

func TestService_PgcSeasonCards(t *testing.T) {
	Convey("PgcCards", t, WithService(func(s *Service) {
		ids := []int64{}
		for i := int64(1); i < 7000; i++ {
			ids = append(ids, i)
		}
		res, err := s.PgcCards(ids)
		So(err, ShouldBeNil)
		fmt.Println(len(res))
		for k, v := range res {
			fmt.Println("Key: ", k)
			fmt.Println("Value: ", v)
		}
	}))
}
