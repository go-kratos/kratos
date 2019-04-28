package external

import (
	"context"
	"encoding/json"

	"go-common/library/log"
)

// DynamicNew .
func (s *Service) DynamicNew(c context.Context, params string) (res json.RawMessage, err error) {
	if res, err = s.dynamic.DynamicNew(c, params); err != nil {
		log.Error("dynamic.service.DynamicNew.error(%v)", err.Error())
	}
	return
}

// DynamicCount .
func (s *Service) DynamicCount(c context.Context, params string) (res json.RawMessage, err error) {
	if res, err = s.dynamic.DynamicCount(c, params); err != nil {
		log.Error("dynamic.service.DynamicCount.error(%v)", err.Error())
	}
	return
}

// DynamicHistory .
func (s *Service) DynamicHistory(c context.Context, params string) (res json.RawMessage, err error) {
	if res, err = s.dynamic.DynamicHistory(c, params); err != nil {
		log.Error("dynamic.service.DynamicHistory.error(%v)", err.Error())
	}
	return
}
