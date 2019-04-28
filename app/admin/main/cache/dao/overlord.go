package dao

import (
	"context"
	"net"
	"strconv"

	"go-common/app/admin/main/cache/model"
	"go-common/library/log"
)

var (
	apiserverURI = "http://cache-mng.bilibili.co/api/v1/appids/%s"
)

// OverlordClusters get all overlord clusters.
func (d *Dao) OverlordClusters(c context.Context, zone, appid string) (ocs []*model.OverlordCluster, err error) {
	var res struct {
		Data []*model.OverlordApiserver `json:"grouped_clusters"`
	}
	if err = d.client.RESTfulGet(c, apiserverURI, "", nil, &res, appid); err != nil {
		log.Error("overlord cluster url(%s) appid(%s) error(%v)", apiserverURI, appid, err)
		return
	}
GETALL:
	for _, oa := range res.Data {
		if zone == "" || oa.Group == zone {
			for _, oc := range oa.Clusters {
				cluster := &model.OverlordCluster{
					Name:             oc.Name,
					Type:             oc.Type,
					Zone:             zone,
					HashMethod:       "fnv1a_64",
					HashDistribution: "ketama",
					HashTag:          "{}",
					ListenProto:      "tcp",
					ListenAddr:       net.JoinHostPort("0.0.0.0", strconv.Itoa(oc.FrontEndPort)),
					DailTimeout:      1000,
					ReadTimeout:      1000,
					WriteTimeout:     1000,
					NodeConn:         2,
					PingFailLimit:    3,
					PingAutoEject:    true,
				}
				for _, oci := range oc.Instances {
					if oc.Type == "redis_cluster" && oci.Role != "master" {
						continue
					}
					on := &model.OverlordNode{
						Alias:  oci.Alias,
						Addr:   net.JoinHostPort(oci.IP, strconv.Itoa(oci.Port)),
						Weight: oci.Weight,
					}
					cluster.Nodes = append(cluster.Nodes, on)
				}
				ocs = append(ocs, cluster)
			}
		}
	}
	if len(ocs) == 0 && zone != "" {
		zone = ""
		goto GETALL
	}
	return
}
