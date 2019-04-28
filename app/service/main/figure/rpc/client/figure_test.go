package figure

import (
	"context"
	"fmt"
	"testing"
	"time"

	"go-common/app/service/main/figure/model"
)

var (
	s *Service
)

func init() {
	s = New(nil)
	time.Sleep(2 * time.Second)
}

func TestUserFigure(t *testing.T) {
	res, _ := s.UserFigure(context.TODO(), &model.ArgUserFigure{Mid: 27515628})
	fmt.Println(res)
}
