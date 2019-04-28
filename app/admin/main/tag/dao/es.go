package dao

import (
	"context"
	"time"

	"go-common/app/admin/main/tag/model"
	"go-common/library/database/elastic"
	"go-common/library/log"
)

const (
	_esBusinessID      = "tag_list"
	_esArchive         = "archive"
	_esIndex           = "tag"
	_esArchiveInterval = 100
)

// ESearchTag search tag from elastic.
func (d *Dao) ESearchTag(c context.Context, esTag *model.ESTag) (tags *model.MngSearchTagList, err error) {
	req := d.es.NewRequest(_esBusinessID).Index(_esIndex).Fields("id", "name", "content", "state", "tag_type", "verify_state", "atten_count", "use_count", "ctime", "mtime").Order(esTag.Order, esTag.Sort).Pn(int(esTag.Pn)).Ps(int(esTag.Ps))
	if esTag.Keyword != "" {
		req.OrderScoreFirst(false).WhereLike([]string{"name"}, []string{esTag.Keyword}, true, elastic.LikeLevelMiddle)
	}
	if len(esTag.IDs) != 0 {
		req.WhereIn("id", esTag.IDs)
	}
	if esTag.TagType != model.TypeUnknow {
		req.WhereEq("tag_type", esTag.TagType)
	}
	if esTag.State != model.StateUnknown {
		req.WhereEq("state", esTag.State)
	}
	if esTag.Vstate != model.VerifyUnknown {
		req.WhereEq("verify_state", esTag.Vstate)
	}
	tags = new(model.MngSearchTagList)
	if err = req.Scan(c, &tags); err != nil {
		log.Error("ESearchTag req.Scan(%v) error(%v)", req.Params(), err)
	}
	return
}

// UpdateESearchTag update es search tag.
func (d *Dao) UpdateESearchTag(c context.Context, tag *model.Tag) (err error) {
	mst := &model.UpdateESearchTag{
		ID: tag.ID,
	}
	if tag.State > model.StateUnknown {
		mst.State = &tag.State
	}
	if tag.Verify > model.VerifyUnknown {
		mst.Verify = &tag.Verify
	}
	if tag.Type > model.TypeUnknow {
		mst.Type = &tag.Type
	}
	req := d.es.NewUpdate(_esBusinessID).AddData(_esIndex, mst)
	if err = req.Do(c); err != nil {
		log.Error("d.dao.UpdateESearchTag(%v) error(%v)", req.Params(), err)
	}
	return
}

// ESearchArchives get archive infos through es search service.
func (d *Dao) ESearchArchives(c context.Context, aids []int64) (arcMap map[int64]*model.SearchRes, err error) {
	var n = _esArchiveInterval
	arcMap = make(map[int64]*model.SearchRes, len(aids))
	for len(aids) > 0 {
		if n > len(aids) {
			n = len(aids)
		}
		req := d.es.NewRequest(_esArchive).Index(_esArchive).Fields("id", "typeid", "title", "mission_id", "mid", "pubtime", "ctime", "copyright", "state").WhereIn("id", aids[:n]).Pn(1).Ps(n)
		res := &struct {
			Archives []*model.SearchRes `json:"result"`
		}{}
		if err = req.Scan(c, &res); err != nil {
			log.Error("ESearchArchives req.Scan(%v) error(%v)", req.Params(), err)
			return
		}
		for _, v := range res.Archives {
			arcMap[v.ID] = v
		}
		aids = aids[n:]
		if len(aids) == 0 {
			return
		}
		time.Sleep(time.Second)
	}
	return
}
