package service

import (
	"context"

	"go-common/app/admin/ep/merlin/model"
	"go-common/library/ecode"
)

const (
	_startsWith          = "startswith"
	_merlinHostnameRegex = ".*-[0-9]-0"
)

// UserTreeAsOption get user tree as option.
func (s *Service) UserTreeAsOption(c context.Context, sessionID string) (firstRetMap []map[string]interface{}, err error) {
	var treeMap *model.UserTree

	if treeMap, err = s.dao.UserTree(c, sessionID); err != nil {
		err = ecode.MerlinGetUserTreeFailed
		return
	}
	if treeMap.Bilibili == nil {
		return
	}
	firstLevelMapTmp := treeMap.Bilibili
	if firstLevelMapTmp["children"] == nil {
		return
	}
	firstLevelMap := firstLevelMapTmp["children"].(map[string]interface{})
	for firstLevelKey, firstLevelChildren := range firstLevelMap {
		secondLevelMapTmp := firstLevelChildren.(map[string]interface{})
		if secondLevelMapTmp["children"] == nil {
			continue
		}
		secondLevelMap := secondLevelMapTmp["children"].(map[string]interface{})
		var secondRetMap []map[string]interface{}

		for secondLevelKey, secondLevelChildren := range secondLevelMap {
			thirdLevelMapTmp := secondLevelChildren.(map[string]interface{})
			if thirdLevelMapTmp["children"] == nil {
				continue
			}
			thirdLevelMap := thirdLevelMapTmp["children"].(map[string]interface{})
			var thirdRetMap []map[string]interface{}
			for thirdLevelKey, thirdLevelChildren := range thirdLevelMap {
				if thirdLevelKey != "" && thirdLevelChildren != nil {
					childID := int64(thirdLevelChildren.(map[string]interface{})["id"].(float64))

					thirdTmp := make(map[string]interface{})
					thirdTmp["label"] = thirdLevelKey
					v := new(struct {
						ID   int64  `json:"id"`
						Name string `json:"name"`
					})
					v.ID = childID
					v.Name = thirdLevelKey
					thirdTmp["value"] = v

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

// VerifyTreeContainerNode verify tree containers.
func (s *Service) VerifyTreeContainerNode(c context.Context, sessionID string, tnr *model.TreeNode) (err error) {
	var (
		sonMap map[string]interface{}
		ok     bool
	)
	if sonMap, err = s.dao.TreeSon(c, sessionID, tnr.TreePath()); err != nil {
		err = ecode.MerlinGetUserTreeFailed
		return
	}
	if sonMap, ok = sonMap["dev"].(map[string]interface{}); !ok {
		err = ecode.MerlinLoseTreeContainerNodeErr
		return
	}
	if sonMap, ok = sonMap["children"].(map[string]interface{}); !ok {
		err = ecode.MerlinLoseTreeContainerNodeErr
		return
	}
	if sonMap["containers"] == nil {
		err = ecode.MerlinLoseTreeContainerNodeErr
	}
	return

}

// TreeRoleAsAuditor tree role as auditor.
func (s *Service) TreeRoleAsAuditor(c context.Context, sessionID, firstNode string) (auditors []string, err error) {
	var treeRoles []*model.TreeRole
	if treeRoles, err = s.dao.TreeRoles(c, sessionID, firstNode); err != nil {
		return
	}
	for _, treeRole := range treeRoles {
		if treeRole.Role == model.TreeRoleAdmin {
			auditors = append(auditors, treeRole.UserName)
		}
	}
	return
}

// TreeMachineIsExist tree machine is exist.
func (s *Service) TreeMachineIsExist(c context.Context, podName string, tn *model.TreeNode) (b bool, err error) {
	var (
		pathAndPodNameMap = make(map[string][]string)
		result            map[string]bool
	)
	pathAndPodNameMap[tn.TreePath()] = []string{podName}
	if result, err = s.TreeMachinesIsExist(c, pathAndPodNameMap); err != nil {
		return
	}
	b = result[podName]
	return
}

// TreeMachinesIsExist tree machines is exist.
func (s *Service) TreeMachinesIsExist(c context.Context, pathAndMachinesMap map[string][]string) (machinesIsExist map[string]bool, err error) {
	var (
		treeAppInstances map[string][]*model.TreeAppInstance
		paths            []string
	)

	machinesIsExist = make(map[string]bool)

	for path, machineNames := range pathAndMachinesMap {
		paths = append(paths, path)
		for _, machineName := range machineNames {
			machinesIsExist[machineName] = false
		}
	}

	if treeAppInstances, err = s.dao.TreeAppInstance(c, paths); err != nil {
		return
	}

	for _, machinesInTree := range treeAppInstances {
		for _, machineInTree := range machinesInTree {
			if _, ok := machinesIsExist[machineInTree.HostName]; ok {
				machinesIsExist[machineInTree.HostName] = true
			}
		}
	}
	return
}

// QueryTreeInstanceForMerlin query TreeInstance for merlin.
func (s *Service) QueryTreeInstanceForMerlin(c context.Context, sessionID string, tn *model.TreeNode) (res map[string]*model.TreeInstance, err error) {
	tir := &model.TreeInstanceRequest{
		HostnameRegex: _merlinHostnameRegex,
	}
	if path := tn.TreePathWithoutEmptyField(); path != "" {
		tir.Path = path
		tir.PathFuzzy = _startsWith
	}
	res, err = s.dao.QueryTreeInstances(c, sessionID, tir)
	return
}
