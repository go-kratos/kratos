package service

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"strings"

	"go-common/app/interface/openplatform/article/model"
	account "go-common/app/service/main/account/model"
)

// AnniversaryInfo get reader and author info in passed year
func (s *Service) AnniversaryInfo(c context.Context, mid int64) (res *model.AnniversaryInfo, err error) {
	if res, err = s.dao.CacheAnniversary(c, mid); err != nil {
		return
	}
	if res == nil {
		res = new(model.AnniversaryInfo)
	}
	res.Mid = mid
	user, _ := s.accountRPC.Info3(c, &account.ArgMid{Mid: mid})
	if user != nil {
		res.Uname = user.Name
		res.Face = user.Face
	}
	if res.AuthorInfo != nil && res.AuthorInfo.ReaderMid != 0 {
		user, _ := s.accountRPC.Info3(c, &account.ArgMid{Mid: res.AuthorInfo.ReaderMid})
		if user != nil {
			res.AuthorInfo.ReaderUname = user.Name
			res.AuthorInfo.ReaderFace = user.Face
		}
	}
	if res.ReaderInfo != nil && res.ReaderInfo.AuthorMid != 0 {
		wordsFloat := float64(res.ReaderInfo.Words) / 1000
		res.ReaderInfo.Words = int64(math.Pow(wordsFloat, 2))
		rankFloat, _ := strconv.ParseFloat(strings.Split(res.ReaderInfo.Rank, "%")[0], 10)
		rankFloat = rankFloat / 100
		rankFloat = math.Pow(rankFloat, 2)
		res.ReaderInfo.Rank = fmt.Sprintf("%.2f", rankFloat*100) + "%"
		user, _ := s.accountRPC.Info3(c, &account.ArgMid{Mid: res.ReaderInfo.AuthorMid})
		if user != nil {
			res.ReaderInfo.AuthorUname = user.Name
		}
	}
	return
}
