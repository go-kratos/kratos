package service

import (
	"context"
	"go-common/app/service/main/tv/internal/model"
	"go-common/library/log"
)

func (s *Service) TokenInfos(c context.Context, tokens []string) (tis []*model.TokenInfo, err error) {
	payParams, err := s.dao.CachePayParamsByTokens(c, tokens)
	if err != nil {
		log.Error("s.dao.CachePayParamsByTokens(%v) err(%v)", tokens, err)
		return
	}
	tis = make([]*model.TokenInfo, 0, len(tokens))
	for tn, pp := range payParams {
		token := &model.TokenInfo{Token: tn}
		token.CopyFromPayParam(pp)
		tis = append(tis, token)
	}
	return
}
