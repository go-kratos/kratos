package service

import (
	"context"
	"encoding/json"

	"go-common/app/admin/main/esports/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/xstr"

	"github.com/jinzhu/gorm"
)

var _emptyTreeList = make([][]*model.TreeList, 0)

// AddTree .
func (s *Service) AddTree(c context.Context, param *model.Tree) (rs model.Tree, err error) {
	var rootID int64
	db := s.dao.DB.Model(&model.Tree{}).Create(param)
	if err = db.Error; err != nil {
		log.Error("AddTree s.dao.DB.Model Create(%+v) error(%v)", param, err)
		return
	}
	rootID = (db.Value.(*model.Tree)).RootID
	rs.ID = (db.Value.(*model.Tree)).ID
	rs.MaID = (db.Value.(*model.Tree)).MaID
	rs.MadID = (db.Value.(*model.Tree)).MadID
	rs.Pid = (db.Value.(*model.Tree)).Pid
	rs.Mid = (db.Value.(*model.Tree)).Mid
	// update detail online.
	s.UpOnline(c, rs.MadID, _downLine)
	// update first root_id.
	if param.RootID == 0 || rs.Pid == 0 {
		upDB := s.dao.DB.Model(&model.Tree{})
		if err = upDB.Where("id=?", rs.ID).Updates(map[string]interface{}{"root_id": rs.ID}).Error; err != nil {
			log.Error("AddTree db.Model  Updates(%+v) error(%v)", rs.ID, err)
			return
		}
		rootID = (upDB.Value.(*model.Tree)).RootID
	}
	rs.RootID = rootID
	return
}

// EditTree .
func (s *Service) EditTree(c context.Context, param *model.TreeEditParam) (err error) {
	var (
		paramIDs         []int64
		trees, paramTree []*model.Tree
		delIDs           []int64
		nodes            map[int64]int64
	)
	if err = json.Unmarshal([]byte(param.Nodes), &nodes); err != nil {
		err = ecode.RequestErr
		return
	}
	if len(nodes) == 0 {
		err = ecode.EsportsTreeEmptyErr
		return
	}
	for id := range nodes {
		if id > 0 {
			paramIDs = append(paramIDs, id)
		} else {
			err = ecode.RequestErr
			return
		}
	}
	// check tree maid same .
	if err = s.dao.DB.Model(&model.Tree{}).Model(&model.Tree{}).Where("id IN (?)", paramIDs).Find(&paramTree).Error; err != nil {
		log.Error("EditTree s.dao.DB.Model  Find(%+v) error(%v)", paramIDs, err)
	}
	for _, tmpTree := range paramTree {
		if tmpTree.MadID != param.MadID {
			err = ecode.EsportsTreeNodeErr
			return
		}
	}
	ParamIDCount := len(paramIDs)
	treeDB := s.dao.DB.Model(&model.Tree{})
	if err = treeDB.Where("mad_id=?", param.MadID).Where("is_deleted=0").Find(&trees).Error; err != nil {
		log.Error("EditTree Error (%v)", err)
		return
	}
	treeCount := len(trees)
	if ParamIDCount != treeCount {
		for _, tree := range trees {
			if _, ok := nodes[tree.ID]; !ok {
				delIDs = append(delIDs, tree.ID)
			}
		}
	}
	// Prevent multiplayer editing
	if ParamIDCount != (treeCount - len(delIDs)) {
		err = ecode.EsportsMultiEdit
		return
	}
	tx := s.dao.DB.Begin()
	if err = tx.Error; err != nil {
		log.Error("s.dao.DB.Begin error(%v)", err)
		return
	}
	if err = upTree(tx, delIDs, nodes, param.MadID); err != nil {
		err = tx.Rollback().Error
		return
	}
	if err = tx.Commit().Error; err != nil {
		log.Error("tx.Commit error(%v)", err)
		return
	}
	return
}

func upTree(tx *gorm.DB, delIDs []int64, nodes map[int64]int64, madID int64) (err error) {
	if err = tx.Model(&model.Tree{}).Where("id IN (?)", delIDs).Updates(map[string]interface{}{"is_deleted": _deleted}).Error; err != nil {
		log.Error("upTree tx.Model  Updates(%+v) error(%v)", delIDs, err)
		return
	}
	if err = tx.Model(&model.Tree{}).Exec(model.BatchEditTreeSQL(nodes)).Error; err != nil {
		log.Error("upTree tx.Model Exec(%+v) error(%v)", nodes, err)
	}
	if err = tx.Model(&model.MatchDetail{}).Where("id=?", madID).Updates(map[string]interface{}{"online": _downLine}).Error; err != nil {
		log.Error("upTree MatchDetail tx.Model Updates(%d) error(%v)", madID, err)
	}
	return
}

// DelTree .
func (s *Service) DelTree(c context.Context, param *model.TreeDelParam) (err error) {
	var ids []int64
	if ids, err = xstr.SplitInts(param.IDs); err != nil {
		err = ecode.RequestErr
		return
	}
	if err = s.dao.DB.Model(&model.Tree{}).Model(&model.Tree{}).Where("id IN (?)", ids).Updates(map[string]interface{}{"is_deleted": _deleted}).Error; err != nil {
		log.Error("DelTree treeDB.Model  Updates(%+v) error(%v)", ids, err)
	}
	// update detail online.
	s.UpOnline(c, param.MadID, _downLine)
	return
}

// TreeList .
func (s *Service) TreeList(c context.Context, param *model.TreeListParam) (res *model.TreeDetailList, err error) {
	var (
		trees       []*model.Tree
		mids        []int64
		contests    []*model.Contest
		mapContests map[int64]*model.ContestInfo
		sTree       []*model.TreeList
		cInfos      []*model.ContestInfo
		detail      model.MatchDetail
		rs          [][]*model.TreeList
	)
	db := s.dao.DB.Model(&model.MatchDetail{}).Where("id=?", param.MadID).First(&detail)
	if err = db.Error; err != nil {
		log.Error("TreeList MatchDetail db.Model Find error(%v)", err)
	}
	if detail.ID == 0 {
		err = ecode.EsportsTreeDetailErr
		return
	}
	treeDB := s.dao.DB.Model(&model.Tree{})
	treeDB = treeDB.Where("mad_id=?", param.MadID).Where("is_deleted=0")
	if err = treeDB.Limit(s.c.Rule.MaxTreeContests).Order("root_id ASC,pid ASC,game_rank ASC").Find(&trees).Error; err != nil {
		log.Error("TreeList Error (%v)", err)
		return
	}
	for _, tree := range trees {
		mids = append(mids, tree.Mid)
	}
	contestsDB := s.dao.DB.Model(&model.Contest{})
	contestsDB = contestsDB.Where("id IN (?)", mids)
	count := len(mids)
	if err = contestsDB.Offset(0).Limit(count).Find(&contests).Error; err != nil {
		log.Error("ContestList Error (%v)", err)
		return
	}
	if cInfos, err = s.contestInfos(c, contests, false); err != nil {
		log.Error("s.contestInfos Error (%v)", err)
		return
	}
	mapContests = make(map[int64]*model.ContestInfo, count)
	for _, info := range cInfos {
		mapContests[info.ID] = info
	}
	for _, tree := range trees {
		if tree.Pid == 0 {
			if len(sTree) > 0 {
				rs = append(rs, sTree)
			}
			sTree = nil
			if cInfo, ok := mapContests[tree.Mid]; ok {
				sTree = append(sTree, &model.TreeList{Tree: tree, ContestInfo: cInfo})
			} else {
				sTree = append(sTree, &model.TreeList{Tree: tree})
			}
		} else {
			if cInfo, ok := mapContests[tree.Mid]; ok {
				sTree = append(sTree, &model.TreeList{Tree: tree, ContestInfo: cInfo})
			} else {
				sTree = append(sTree, &model.TreeList{Tree: tree})
			}
		}
	}
	if len(sTree) > 0 {
		rs = append(rs, sTree)
	}
	res = new(model.TreeDetailList)
	if len(trees) == 0 {
		res.Tree = _emptyTreeList
	} else {
		res.Tree = rs
	}
	res.Detail = &detail
	return
}
