package academy

import (
	"context"
	"fmt"

	"go-common/app/interface/main/creative/model/academy"
	"go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/xstr"
)

const (
	// select
	_getTagListSQL  = "SELECT id, parent_id, type, business, name, `desc`, state, ctime,mtime FROM academy_tag WHERE state=0 ORDER BY rank asc"
	_getTagByIDSQL  = "SELECT id, parent_id, type, business, name, `desc`, state, ctime,mtime FROM academy_tag WHERE state=0 AND id=?"
	_getTagInIDsSQL = "SELECT id, parent_id, type, business, name, `desc`, state, ctime,mtime FROM academy_tag WHERE state=0 AND id in (%s)"
	_getTagLinkSQL  = "SELECT id, tid, link_id FROM academy_tag_link WHERE tid in (%s)"
)

// TagList get all tag from academy_tag.
func (d *Dao) TagList(c context.Context) (res map[string][]*academy.Tag, tagMap, parentChildMap map[int64]*academy.Tag, err error) {
	rows, err := d.db.Query(c, _getTagListSQL)
	if err != nil {
		log.Error("d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	res = make(map[string][]*academy.Tag)        //最终返回标签列表
	all := make(map[int8]map[int64]*academy.Tag) //含有分类的一二级标签
	top := make(map[int64]*academy.Tag)
	tagMap = make(map[int64]*academy.Tag)         //不包含二级标签
	parentChildMap = make(map[int64]*academy.Tag) //包含一二级标签
	allTag := make([]*academy.Tag, 0)
	for rows.Next() {
		t := &academy.Tag{}
		if err = rows.Scan(&t.ID, &t.ParentID, &t.Type, &t.Business, &t.Name, &t.Desc, &t.State, &t.CTime, &t.MTime); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		nt := &academy.Tag{}
		*nt = *t
		tagMap[t.ID] = nt        //nt对象 不会被修改
		parentChildMap[t.ID] = t //t对象 会被修改
		allTag = append(allTag, t)
	}
	for _, t := range allTag {
		if t.ParentID == 0 {
			top[t.ID] = t
			all[t.Type] = top
			tyName := academy.TagClassMap(int(t.Type))
			res[tyName] = append(res[tyName], t)
		}
	}
	for _, t := range allTag {
		a, ok := all[t.Type][t.ParentID]
		if ok && a != nil && a.Type == t.Type {
			a.Children = append(a.Children, t)
		}
	}
	return
}

//Tag get one tag.
func (d *Dao) Tag(c context.Context, id int64) (t *academy.Tag, err error) {
	row := d.db.QueryRow(c, _getTagByIDSQL, id)
	t = &academy.Tag{}
	if err = row.Scan(&t.ID, &t.ParentID, &t.Type, &t.Business, &t.Name, &t.Desc, &t.State, &t.CTime, &t.MTime); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			return
		}
		log.Error("row.Scan error(%v)", err)
	}
	return
}

//Tags get some tags by ids.
func (d *Dao) Tags(c context.Context, ids []int64) (res []*academy.Tag, err error) {
	rows, err := d.db.Query(c, fmt.Sprintf(_getTagInIDsSQL, xstr.JoinInts(ids)))
	if err != nil {
		log.Error("d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	res = make([]*academy.Tag, 0)
	for rows.Next() {
		t := &academy.Tag{}
		if err = rows.Scan(&t.ID, &t.ParentID, &t.Type, &t.Business, &t.Name, &t.Desc, &t.State, &t.CTime, &t.MTime); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		res = append(res, t)
	}
	return
}

//LinkTags get link tags by h5 tids.
func (d *Dao) LinkTags(c context.Context, ids []int64) (res []*academy.LinkTag, err error) {
	rows, err := d.db.Query(c, fmt.Sprintf(_getTagLinkSQL, xstr.JoinInts(ids)))
	if err != nil {
		log.Error("d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	res = make([]*academy.LinkTag, 0)
	for rows.Next() {
		t := &academy.LinkTag{}
		if err = rows.Scan(&t.ID, &t.TID, &t.LinkID); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		res = append(res, t)
	}
	return
}
