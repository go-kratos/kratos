package service

import (
	"context"
	"fmt"
	"strings"

	"go-common/app/admin/main/search/dao"
	"go-common/app/admin/main/search/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

// Upsert upsert docs.
func (s *Service) Upsert(c context.Context, up *model.UpsertParams, dataBody map[string][]model.MapData) (err error) {
	app, ok := s.queryConf[up.Business]
	if app2, ok2 := model.QueryConf[up.Business]; ok2 {
		app = app2
		ok = true
	}
	if !ok {
		err = fmt.Errorf("up.Business(%s) not exists in queryConf", up.Business)
		return
	}
	if app.ESCluster == "" {
		err = fmt.Errorf("app(%+v) escluster is empty", app)
		return
	}
	// dataBody to upsertBody
	up.UpsertBody = []model.UpsertBody{}
	for indexName, docs := range dataBody {
		if !strings.Contains(indexName, app.IndexPrefix) {
			log.Error("invalid indexName (%s)", indexName)
			continue
		}
		for _, doc := range docs {
			indexID := doc.StrID(app.IndexID)
			//TODO 提前告知base不对
			if indexID == "" {
				continue
			}
			upsert := model.UpsertBody{IndexName: indexName, IndexType: app.IndexType, IndexID: indexID, Doc: doc}
			up.UpsertBody = append(up.UpsertBody, upsert)
		}
	}
	if err = s.dao.UpsertBulk(c, app.ESCluster, up); err != nil {
		dao.PromError(fmt.Sprintf("es:%s 更新失败", app.ESCluster), "s.dao.UpsertBulk error(%v) ", err)
		err = ecode.SearchUpdateIndexFailed
	}
	return
}
