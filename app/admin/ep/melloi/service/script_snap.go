package service

import (
	"context"
	"encoding/json"

	"go-common/app/admin/ep/melloi/model"
	"go-common/library/log"
)

//AddSciptSnap add scriptSnap
func (s *Service) AddSciptSnap(ptestParam model.DoPtestParam, excuteId string) (scriptSnapIds []int, err error) {
	var (
		scripts      []*model.Script
		scriptSnapId int
	)

	//场景脚本快照保存逻辑
	if ptestParam.Type == model.SCENE_SCRIPT_TYPE {
		for _, script := range ptestParam.Scripts {
			script.SceneID = ptestParam.SceneID
			if scriptSnapId, err = s.AddScriptSnapInfo(script, excuteId); err != nil {
				log.Error("s.AddScriptSnapInfo err :(%v)", err)
				return
			}
			scriptSnapIds = append(scriptSnapIds, scriptSnapId)
		}
		return
	}

	//单场景http 脚本的快照保存
	if scripts, err = s.QueryScripts(&model.Script{ID: ptestParam.ScriptID}, 1, 5); err != nil {
		log.Error(" s.QueryScripts err :(%v)", err)
		return
	}
	if scriptSnapId, err = s.AddScriptSnapInfo(scripts[0], excuteId); err != nil {
		log.Error("s.AddScriptSnapInfo err :(%v)", err)
		return
	}
	scriptSnapIds = append(scriptSnapIds, scriptSnapId)
	return
}

//AddScriptSnapInfo Add scriptSnap Info
func (s *Service) AddScriptSnapInfo(script *model.Script, excuteId string) (scriptSnapId int, err error) {
	var JSON []byte
	scriptSnap := model.ScriptSnap{ScriptID: script.ID}
	script.ID = 0

	if JSON, err = json.Marshal(script); err != nil {
		log.Error(" json.Marshal(script) error :(%v)", err)
		return
	}
	if err = json.Unmarshal([]byte(string(JSON)), &scriptSnap); err != nil {
		return
	}

	scriptSnap.ExecuteID = excuteId
	if scriptSnapId, err = s.dao.AddScriptSnap(&scriptSnap); err != nil {
		log.Error("s.dao.AddScriptSnap error :(%v)", err)
		return
	}
	return
}

//AddSnap add snap
func (s *Service) AddSnap(c context.Context, ptestParam model.DoPtestParam, executeID, jobName, jobNamed string) (scriptSnapIDs []int, err error) {
	if ptestParam.Type == model.PROTOCOL_GRPC {
		if scriptSnapIDs, err = s.AddGRPCSnap(ptestParam.ScriptID, executeID); err != nil {
			log.Error("save grpc snap failed,(%v)", err)
			s.DeleteJob(context.TODO(), jobName)
			return
		}
	} else {
		//用jobNamed 表示执行id，即是 ExcuteId
		if scriptSnapIDs, err = s.AddSciptSnap(ptestParam, jobNamed); err != nil {
			log.Error("s.AddSciptSnap err :(%v)", err)
			s.DeleteJob(context.TODO(), jobName)
			return
		}
	}
	return
}

// AddGRPCSnap Add GRPC Snap
func (s *Service) AddGRPCSnap(grpcID int, executeID string) (snapIDs []int, err error) {
	var (
		grpc     *model.GRPC
		grpcSnap *model.GRPCSnap
		j        []byte
	)
	if grpc, err = s.dao.QueryGRPCByID(grpcID); err != nil {
		return
	}
	grpcSnap = &model.GRPCSnap{GRPCID: grpc.ID, ExecuteID: executeID}
	grpc.ID = 0
	if j, err = json.Marshal(grpc); err != nil {
		return
	}
	if err = json.Unmarshal([]byte(string(j)), grpcSnap); err != nil {
		return
	}
	if err = s.dao.CreateGRPCSnap(grpcSnap); err != nil {
		return
	}
	return []int{grpcSnap.ID}, nil
}
