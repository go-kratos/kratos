package service

import (
	"context"

	"go-common/app/service/main/filter/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

// rubbishFilter 反垃圾
func (s *Service) rubbishFilter(c context.Context, area *model.Area, content string, oid int64, contentID int64, senderID int64) (limitType string, hits []string, err error) {
	var (
		hit         string
		rubbishArea string
	)
	rubbishArea = area.RubbishName()
	if rubbishArea == "message" {
		senderID = oid
	}
	hits = make([]string, 0)
	if hit, limitType, err = s.dao.AntispamFilter(c, rubbishArea, content, oid, contentID, senderID); err != nil {
		if ecode.ServiceUnavailable.Equal(err) {
			log.Errorv(c, log.KV("log", "antispam degrade occure"), log.KV("area", rubbishArea), log.KV("content", content), log.KV("oid", oid), log.KV("contentID", contentID), log.KV("senderID", senderID))
			err = nil
			limitType = model.LimitTypeOK
		}
		return
	}
	hits = append(hits, hit)
	return
}
