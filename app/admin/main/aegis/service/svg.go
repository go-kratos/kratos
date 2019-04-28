package service

import (
	"context"

	"go-common/app/admin/main/aegis/model/net"
	"go-common/app/admin/main/aegis/model/svg"
)

//NetSVG 绘制网的svg图
func (s *Service) NetSVG(c context.Context, netid int64) (nv *svg.NetView, err error) {
	flows, err := s.gorm.FlowsByNet(c, []int64{netid})
	if err != nil {
		return
	}

	dirs, err := s.gorm.DirectionByNet(c, netid)
	if err != nil {
		return
	}

	trans, err := s.gorm.TranByNet(c, netid, true)
	if err != nil {
		return
	}

	var binds []*net.TokenBind
	for _, flow := range flows {
		abs, _, _ := s.fetchOldBindAndLog(c, flow.ID, []int8{1})
		for _, bind := range abs {
			binds = append(binds, bind)
		}
	}
	for _, tran := range trans {
		abs, _, _ := s.fetchOldBindAndLog(c, tran.ID, []int8{2})
		for _, bind := range abs {
			binds = append(binds, bind)
		}
	}

	tbds, err := s.tokenBindDetail(c, binds)
	if err != nil {
		return
	}

	dot := svg.NewDot()
	dot.StartDot().AddTokenBinds(tbds...).
		AddFlow(flows...).
		AddTransitions(trans...).
		AddDirections(dirs...).
		End()
	nv = svg.NewNetView()
	nv.SetDot(dot)
	return
}
