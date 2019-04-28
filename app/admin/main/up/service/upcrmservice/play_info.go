package upcrmservice

import (
	"context"
	"errors"
	"go-common/app/admin/main/up/model/upcrmmodel"
	"go-common/library/ecode"
	"go-common/library/log"
)

//PlayQueryInfo handle PlayQueryInfo
func (s *Service) PlayQueryInfo(context context.Context, arg *upcrmmodel.PlayQueryArgs) (result upcrmmodel.PlayQueryResult, err error) {

	var types []int
	switch arg.BusinessType {
	case 0:
		types = append(types, upcrmmodel.BusinessTypeVideo, upcrmmodel.BusinessTypeAudio, upcrmmodel.BusinessTypeArticle)
	case upcrmmodel.BusinessTypeVideo, upcrmmodel.BusinessTypeAudio, upcrmmodel.BusinessTypeArticle:
		types = append(types, arg.BusinessType)
	default:
		err = errors.New("business type not support")
		return
	}

	playData, e := s.crmdb.QueryPlayInfo(arg.Mid, types)
	err = e
	if err != nil && err != ecode.NothingFound {
		log.Error("fail to get from db, err=%+v", err)
		return
	}
	if err == ecode.NothingFound {
		err = nil
		log.Warn("up not found in play info db, mid=%d", arg.Mid)
	}

	for _, v := range playData {
		result.BusinessData = append(result.BusinessData, upcrmmodel.CastUpPlayInfoToPlayInfo(v))
	}
	// 去查up_base_info以获取总稿件量
	baseInfo, e := s.crmdb.QueryUpBaseInfo(arg.Mid, "article_count_accumulate, article_count_30day")
	err = e
	if err == nil {
		result.ArticleCount30Day = baseInfo.ArticleCount30day
		result.ArticleCountAccumulate = baseInfo.ArticleCountAccumulate
	} else if err == ecode.NothingFound {
		err = nil
		log.Warn("up base info not found, mid=%d", arg.Mid)
	}
	log.Info("get play info ok, mid=%d", arg.Mid)
	return
}
