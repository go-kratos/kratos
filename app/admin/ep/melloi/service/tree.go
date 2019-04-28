package service

import (
	"context"

	"go-common/app/admin/ep/melloi/model"
	"go-common/library/ecode"
)

// QueryServiceTreeToken get tree token by user sessionID
func (s *Service) QueryServiceTreeToken(c context.Context, sessionID string) (token string, err error) {
	return s.dao.QueryServiceTreeToken(c, sessionID)
}

// QueryUserRoleNode QueryUserRoleNode
func (s *Service) QueryUserRoleNode(c context.Context, sessionID string) (apps []string, err error) {
	var (
		roleApp []*model.RoleApp
		token   string
	)
	if token, err = s.QueryServiceTreeToken(c, sessionID); err != nil {
		return
	}

	if roleApp, err = s.dao.QueryUserRoleApp(c, token); err != nil {
		err = ecode.MerlinGetUserTreeFailed
		return
	}

	for _, app := range roleApp {
		// Role 1:管理员 2:研发 3:测试 4:运维 5:访客
		// Type 1.公司 2.部门 3.项目 4. 应用 5.环境 6.挂载点
		if app.Role <= 3 && app.Type == 4 {
			apps = append(apps, app.Name)
		}
	}
	return
}

// QueryUserTree get user tree
func (s *Service) QueryUserTree(c context.Context, sessionID string) (firstRetMap []map[string]interface{}, err error) {
	var (
		treeMap *model.UserTree
		token   string
	)
	if token, err = s.QueryServiceTreeToken(c, sessionID); err != nil {
		err = ecode.MerlinGetUserTreeFailed
		return
	}

	if treeMap, err = s.dao.QueryUserTree(c, token); err != nil {
		err = ecode.MerlinGetUserTreeFailed
		return
	}
	if treeMap.Bilibili == nil {
		return
	}
	firstLevelMapTmp := treeMap.Bilibili
	firstLevelMap := firstLevelMapTmp["children"].(map[string]interface{})
	if firstLevelMap == nil {
		return
	}

	for firstLevelKey, firstLevelChildren := range firstLevelMap {
		var secondRetMap []map[string]interface{}

		secondLevelMapTmp := firstLevelChildren.(map[string]interface{})
		if secondLevelMapTmp["children"] == nil {
			continue
		}
		secondLevelMap := secondLevelMapTmp["children"].(map[string]interface{})

		for secondLevelKey, secondLevelChildren := range secondLevelMap {
			var thirdRetMap []map[string]interface{}
			thirdLevelMapTmp := secondLevelChildren.(map[string]interface{})
			if thirdLevelMapTmp["children"] == nil {
				continue
			}
			thirdLevelMap := thirdLevelMapTmp["children"].(map[string]interface{})
			for thirdLevelKey, thirdLevelChildren := range thirdLevelMap {
				if thirdLevelKey != "" && thirdLevelChildren != nil {
					thirdTmp := make(map[string]interface{})
					thirdTmp["value"] = thirdLevelKey
					thirdTmp["label"] = thirdLevelKey
					thirdRetMap = append(thirdRetMap, thirdTmp)
				}
			}

			if secondLevelKey != "" {
				secondTmp := make(map[string]interface{})
				secondTmp["value"] = secondLevelKey
				secondTmp["label"] = secondLevelKey
				secondTmp["children"] = thirdRetMap
				secondRetMap = append(secondRetMap, secondTmp)
			}
		}

		if firstLevelKey != "" {
			firstTmp := make(map[string]interface{})
			firstTmp["value"] = firstLevelKey
			firstTmp["label"] = firstLevelKey
			firstTmp["children"] = secondRetMap
			firstRetMap = append(firstRetMap, firstTmp)
		}
	}
	return
}

// QueryTreeAdmin get tree admin
func (s *Service) QueryTreeAdmin(c context.Context, path string, sessionID string) (roles []*model.TreeRole, err error) {
	var (
		tar   *model.TreeAdminResponse
		token string
	)
	if token, err = s.QueryServiceTreeToken(c, sessionID); err != nil {
		return
	}

	if tar, err = s.dao.QueryTreeAdmin(c, path, token); err != nil {
		err = ecode.MerlinGetUserTreeFailed
		return
	}
	//todo: 如果没有服务树权限需要申请
	for _, v := range tar.Data {
		if v.Role == 1 {
			roles = append(roles, v)
		}
	}
	return
}
