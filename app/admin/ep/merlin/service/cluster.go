package service

import (
	"context"
	"strings"

	"go-common/app/admin/ep/merlin/model"
)

// QueryCluster query cluster.
func (s *Service) QueryCluster(c context.Context) (clusters []*model.Cluster, err error) {
	var tmpClusters []*model.Cluster
	if tmpClusters, err = s.dao.QueryClusters(c); err != nil {
		return
	}

	for _, tmpCluster := range tmpClusters {
		if !strings.Contains(tmpCluster.Name, "uat") {
			clusters = append(clusters, tmpCluster)
		}
	}

	for _, cluster := range clusters {
		for _, supportName := range s.c.BiliHub.SupportNetWork {
			if strings.Contains(cluster.Name, supportName) {
				cluster.IsSupportSnapShot = true
				continue
			}
		}
	}
	return
}
