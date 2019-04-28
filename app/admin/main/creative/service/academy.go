package service

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"time"

	"go-common/app/admin/main/creative/model/academy"
	"go-common/app/interface/openplatform/article/model"
	"go-common/app/service/main/archive/api"
	"go-common/library/log"

	"github.com/jinzhu/gorm"
)

//Archive get one archive by rpc
func (s *Service) Archive(c context.Context, aid int64) (res *api.Arc, err error) {
	if res, err = s.dao.Archive(c, aid); err != nil {
		log.Error("s.dao.Archive aid(%d)|err(%+v)", aid, err)
	}
	return
}

//Archives get archives by rpc
func (s *Service) Archives(c context.Context, aids []int64) (res map[int64]*api.Arc, err error) {
	if res, err = s.dao.Archives(c, aids); err != nil {
		log.Error("s.dao.Archives aids(%d)|err(%+v)", aids, err)
	}
	return
}

//Articles get articles by rpc
func (s *Service) Articles(c context.Context, aids []int64) (res map[int64]*model.Meta, err error) {
	if res, err = s.dao.ArticleMetas(c, aids); err != nil {
		log.Error("s.dao.ArticleMetas aids(%+v)|err(%+v)", aids, err)
	}
	return
}

//ArchivesWithES for es search
func (s *Service) ArchivesWithES(c context.Context, aca *academy.EsParam) (res *academy.SearchResult, err error) {
	if res, err = s.dao.ArchivesWithES(c, aca); err != nil {
		log.Error("s.dao.ArchivesWithES aca(%+v)|err(%+v)", aca, err)
	}
	return
}

// Stats get archives stat.
func (s *Service) Stats(c context.Context, aids []int64, ip string) (res map[int64]*api.Stat, err error) {
	if res, err = s.dao.Stats(c, aids, ip); err != nil {
		log.Error("s.dao.Archives aids(%d)|ip(%s)|err(%+v)", aids, ip, err)
	}
	return
}

//SearchKeywords for list search keywords.
func (s *Service) SearchKeywords() (res []interface{}, err error) {
	var sks []*academy.SearchKeywords
	if err = s.DB.Where("state=0").Order("rank ASC").Find(&sks).Error; err != nil {
		log.Error("SearchKeywords error(%v)", err)
		return
	}
	if len(sks) == 0 {
		return
	}

	res = s.trees(sks, "ID", "ParentID", "Children")
	return
}

//trees for generate tree data set
// data - orm result set
// idFieldStr - primary key in table map to struct
// pidFieldStr - top parent id in table map to struct
// chFieldStr - struct child nodes

func (s *Service) trees(data interface{}, idFieldStr, pidFieldStr, chFieldStr string) (res []interface{}) {
	if reflect.TypeOf(data).Kind() != reflect.Slice {
		return
	}

	sli := reflect.ValueOf(data)
	top := make(map[int64]interface{})
	res = make([]interface{}, 0, sli.Len())
	for i := 0; i < sli.Len(); i++ {
		v := sli.Index(i).Interface()
		if reflect.TypeOf(v).Kind() != reflect.Ptr {
			continue
		}

		if reflect.ValueOf(v).IsNil() {
			continue
		}

		getValue := reflect.ValueOf(v).Elem()
		getType := reflect.TypeOf(v).Elem()
		pid := getValue.FieldByName(pidFieldStr).Interface().(int64)
		if _, ok := getType.FieldByName(pidFieldStr); ok && pid == 0 {
			id := getValue.FieldByName(idFieldStr).Interface().(int64)
			top[id] = v
			res = append(res, v)
		}
	}

	for i := 0; i < sli.Len(); i++ {
		v := sli.Index(i).Interface()
		if reflect.TypeOf(v).Kind() != reflect.Ptr {
			continue
		}

		if reflect.ValueOf(v).IsNil() {
			continue
		}

		pid := reflect.ValueOf(v).Elem().FieldByName(pidFieldStr).Interface().(int64)
		if pid == 0 {
			continue
		}

		if p, ok := top[pid]; ok {
			ch := reflect.ValueOf(p).Elem().FieldByName(chFieldStr)
			ch.Set(reflect.Append(ch, reflect.ValueOf(v)))
		}
	}
	return
}

//SubSearchKeywords for add search keywords.
func (s *Service) SubSearchKeywords(vs []*academy.SearchKeywords) (err error) {
	if len(vs) == 0 {
		return
	}

	origins := []*academy.SearchKeywords{}
	if err = s.DB.Model(&academy.SearchKeywords{}).Find(&origins).Error; err != nil {
		log.Error("SubSearchKeywords Find error(%v)", err)
		return
	}

	originMap := make(map[int64]*academy.SearchKeywords)
	for _, v := range origins {
		originMap[v.ID] = v
	}

	newParents := make([]*academy.SearchKeywords, 0)
	oldParents := make([]*academy.SearchKeywords, 0)
	children := make([]*academy.SearchKeywords, 0)
	newChildren := make([]*academy.SearchKeywords, 0)
	oldChildren := make([]*academy.SearchKeywords, 0)
	now := time.Now().Format("2006-01-02 15:04:05")
	for _, v := range vs {
		if v == nil {
			continue
		}
		v.Name = strings.TrimSpace(v.Name)              //删除字符串前后空格
		if vv, ok := originMap[v.ID]; ok && vv != nil { //父节点为老的
			v.CTime = vv.CTime
			v.MTime = now
			oldParents = append(oldParents, v)
			if len(v.Children) > 0 {
				for _, vvv := range v.Children {
					vvv.ParentID = vv.ID //追加父节点ID
					vvv.CTime = now
					vvv.MTime = now
					children = append(children, vvv) //新老子节点同时存在
				}
			}
		} else {
			v.CTime = now
			v.MTime = now
			newParents = append(newParents, v)
		}
	}
	oldParents = append(oldParents, newParents...)

	tx := s.DB
	if len(oldParents) > 0 {
		if err = s.insertKeyWords(tx, oldParents); err != nil {
			return
		}
	}

	newParentsNames := make([]string, 0)
	for _, v := range newParents {
		newParentsNames = append(newParentsNames, v.Name)
	}

	var pidMap map[string]*academy.SearchKeywords
	if len(newParentsNames) > 0 {
		if pidMap, err = s.upRanks(tx, newParentsNames); err != nil {
			return
		}
	}

	newChildrenNames := make([]string, 0)
	for _, v := range newParents { //父节点为新的
		if v == nil {
			continue
		}
		if vv, ok := pidMap[v.Name]; ok && vv != nil {
			if len(v.Children) > 0 {
				for _, vvv := range v.Children {
					vvv.ParentID = vv.ID //追加父节点ID
					vvv.CTime = now
					vvv.MTime = now
					newChildren = append(newChildren, vvv)
					newChildrenNames = append(newChildrenNames, vvv.Name)
				}
			}
		}
	}

	for _, v := range children {
		if v == nil {
			continue
		}
		if vv, ok := originMap[v.ID]; ok && vv != nil {
			v.CTime = vv.CTime
			v.MTime = now
			oldChildren = append(oldChildren, v)
		} else {
			v.CTime = now
			v.MTime = now
			newChildren = append(newChildren, v)
			newChildrenNames = append(newChildrenNames, v.Name)
		}
	}
	oldChildren = append(oldChildren, newChildren...)

	if len(oldChildren) > 0 {
		if err = s.insertKeyWords(tx, oldChildren); err != nil {
			return
		}
	}

	if len(newChildrenNames) > 0 {
		if _, err = s.upRanks(tx, newChildrenNames); err != nil {
			return
		}
	}
	return
}

func (s *Service) insertKeyWords(tx *gorm.DB, vs []*academy.SearchKeywords) (err error) {
	if len(vs) == 0 {
		return
	}
	valSearks := make([]string, 0, len(vs))
	valSearksArgs := make([]interface{}, 0)
	for _, v := range vs {
		valSearks = append(valSearks, "(?, ?, ?, ?, ?, ?, ?, ?)")
		valSearksArgs = append(valSearksArgs, v.ID, v.Rank, v.ParentID, v.State, v.Name, v.Comment, v.CTime, v.MTime)
	}
	sqlStr := fmt.Sprintf("INSERT INTO academy_search_keywords (id, rank, parent_id, state, name, comment, ctime, mtime) VALUES %s "+
		"ON DUPLICATE KEY UPDATE id=VALUES(id), rank=VALUES(rank), parent_id=VALUES(parent_id), state=VALUES(state), name=VALUES(name), comment=VALUES(comment), ctime=VALUES(ctime), mtime=VALUES(mtime)", strings.Join(valSearks, ","))
	if err = tx.Exec(sqlStr, valSearksArgs...).Error; err != nil {
		log.Error("insertKeyWords error(%v)", err)
	}
	return
}

func (s *Service) upRanks(tx *gorm.DB, names []string) (pidMap map[string]*academy.SearchKeywords, err error) {
	if len(names) == 0 {
		return
	}

	setRanks := []*academy.SearchKeywords{}
	if err = s.DB.Model(&academy.SearchKeywords{}).Where("name IN(?)", names).Find(&setRanks).Error; err != nil {
		log.Error("upRanks Find error(%v)", err)
		return
	}

	upRankSQL := "UPDATE academy_search_keywords SET rank = CASE id "
	ids := make([]int64, 0)
	pidMap = make(map[string]*academy.SearchKeywords)
	for _, v := range setRanks {
		upRankSQL += fmt.Sprintf("WHEN %d THEN %d ", v.ID, v.ID)
		ids = append(ids, v.ID)
		pidMap[v.Name] = v
	}

	if len(ids) == 0 {
		return
	}
	upRankSQL += "END WHERE id IN (?)"
	if err = tx.Exec(upRankSQL, ids).Error; err != nil {
		log.Error("upRanks update rank ids(%+v) error(%v)", ids, err)
	}
	return
}
