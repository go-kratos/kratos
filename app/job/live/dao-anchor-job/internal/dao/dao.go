package dao

import (
	"context"

	"go-common/app/job/live/dao-anchor-job/internal/conf"
	av_api "go-common/app/service/live/av/api/liverpc"
	daoAnchor_api_v0 "go-common/app/service/live/dao-anchor/api/grpc/v0"
	daoAnchor_api "go-common/app/service/live/dao-anchor/api/grpc/v1"
	video_api "go-common/app/service/video/stream-mng/api/v1"
	"go-common/library/database/bfs"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/rpc/liverpc"
)

// Dao dao
type Dao struct {
	c              *conf.Config
	AvApi          *av_api.Client
	daoAnchorApi   *daoAnchor_api.Client
	VideoApi       video_api.StreamClient
	BfsClient      *bfs.BFS
	HttpClient     *bm.Client
	daoAnchorApiV0 *daoAnchor_api_v0.Client
}

// New init mysql db
func New(c *conf.Config) (dao *Dao) {
	dao = &Dao{
		c:          c,
		AvApi:      av_api.New(getConf("av")),
		BfsClient:  bfs.New(c.BfsCli),
		HttpClient: bm.NewClient(c.HttpCli),
	}
	daoAnchorApi, err := daoAnchor_api.NewClient(c.GrpcCli)
	if err != nil {
		panic(err)
	}
	dao.daoAnchorApi = daoAnchorApi
	videoApi, err := video_api.NewClient(c.GrpcCli)
	if err != nil {
		panic(err)
	}
	dao.VideoApi = videoApi
	daoAnchorApiV0, err := daoAnchor_api_v0.NewClient(c.GrpcCli)
	if err != nil {
		panic(err)
	}
	dao.daoAnchorApiV0 = daoAnchorApiV0
	return
}

// Close close the resource.
func (d *Dao) Close() {

}

// Ping dao ping
func (d *Dao) Ping(ctx context.Context) error {
	// TODO: add mc,redis... if you use
	return nil
}

func getConf(appName string) *liverpc.ClientConfig {
	c := conf.Conf.LiveRpc
	if c != nil {
		return c[appName]
	}
	return nil
}
