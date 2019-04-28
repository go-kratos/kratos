package service

import (
	"context"

	"github.com/jinzhu/gorm"
	"go-common/app/admin/main/aegis/model/net"
	"go-common/library/log"
)

func (s *Service) addFlowResources(tx *gorm.DB, netID int64, rids []int64, flowID int64, state int8) (err error) {
	fr := &net.FlowResource{
		FlowID: flowID,
		State:  state,
		NetID:  netID,
	}
	for _, rid := range rids {
		fr.ID = 0
		fr.RID = rid
		if err = s.gorm.AddItem(context.TODO(), tx, fr); err != nil {
			log.Error("addFlowResources s.gorm.AddItem error(%v) fr(%+v)", err, fr)
			return
		}
	}

	return
}

/**
 * updateFlowResources 正常流转到新节点
 * 指定网 & 指定资源 & 新节点
 * 不需要指定现状flowid，不同资源可能不在同一现状节点上
 * 已被取消运行的资源现状，不能被更新
 */
func (s *Service) updateFlowResources(c context.Context, tx *gorm.DB, netID int64, rid int64, newFlowID int64) (err error) {
	var (
		frs []*net.FlowResource
	)

	//资源的运行现状
	if frs, err = s.gorm.FRByNetRID(c, []int64{netID}, []int64{rid}, false); err != nil {
		log.Error("updateFlowResources s.gorm.FRByNetRID error(%v)", err)
		return
	}
	if len(frs) == 0 {
		if err = s.addFlowResources(tx, netID, []int64{rid}, newFlowID, net.FRStateRunning); err != nil {
			log.Error("updateFlowResources s.addFlowResources error(%v)", err)
		}
		return
	}

	//确定state & 记录数目：注意单线->并发和并发->单线的转折点
	//todo--只有新节点，如何确定要不要并发拆分或并发合并？--并发要怎么存储（在配置时候确认并发拆分点、并发分支、并发结合点）？
	//todo--需要上游transitionid确认上游dir.order吗？能出现transition->flow的多对多的情况吗（不能）？
	//todo--单独通过flow去确定是在单线上还是并发线上？
	for _, item := range frs {
		if err = s.gorm.UpdateFields(context.TODO(), tx, net.TableFlowResource, item.ID,
			map[string]interface{}{"flow_id": newFlowID, "state": net.FRStateRunning}); err != nil {
			log.Error("updateFlowResources s.gorm.UpdateFields error(%v)", err)
			return
		}
	}
	return
}
