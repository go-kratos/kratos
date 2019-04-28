package service

import (
	"fmt"
	rpc "go-common/app/service/bbq/recsys/api/grpc/v1"

	"github.com/go-ego/murmur"
)

const (
	//ABTestA ...
	ABTestA = "bbq-rec-A"

	//ABTestB ...
	ABTestB = "bbq-rec-B"
)

//DoABTest ...
func (s *Service) DoABTest(request *rpc.RecsysRequest) {
	bucket := -1
	if request.MID > 0 {
		bucket = int(request.MID % 100)
	} else if len(request.BUVID) > 0 {
		hash := murmur.Sum32(request.BUVID)
		level0 := hash % 100
		level1 := hash / 100 % 100
		level2 := hash / 10000 % 100
		bucket = int(level0)
		request.Abtest = fmt.Sprintf("Rank:%d;Recall:%d;Rule:%d", level0, level1, level2)
	}
	if bucket != -1 {
		if bucket < 50 {
			request.Abtest = ABTestA
		} else {
			request.Abtest = ABTestB
		}
	}
	//white list
	if request.MID == 5829468 {
		request.Abtest = ABTestA
	}
	if request.MID == 208259 {
		request.Abtest = ABTestB
	}
}
