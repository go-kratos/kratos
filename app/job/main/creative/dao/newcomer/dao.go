package newcomer

import (
	"context"

	grpc "go-common/app/interface/main/creative/api"
	"go-common/app/job/main/creative/conf"
	"go-common/library/database/sql"
	httpx "go-common/library/net/http/blademaster"
)

// Dao is search dao.
type Dao struct {
	c *conf.Config
	//db
	db *sql.DB
	//grpc
	grpcClient grpc.CreativeClient
	// http
	httpClient *httpx.Client
	msgURI     string
}

// New new a search dao.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:          c,
		db:         sql.NewMySQL(c.DB.Creative),
		httpClient: httpx.NewClient(c.HTTPClient.Slow),
		msgURI:     c.Host.Message + "/api/notify/send.user.notify.do", //发送用户通知消息接口
	}

	var err error
	if d.grpcClient, err = grpc.NewClient(c.CreativeGRPClient); err != nil {
		panic(err)
	}
	return d
}

// CheckTaskState call grpc client.
func (d *Dao) CheckTaskState(c context.Context, MID, TaskID int64) (*grpc.TaskReply, error) {
	return d.grpcClient.CheckTaskState(c, &grpc.TaskRequest{Mid: MID, TaskId: TaskID})
}

// Ping ping grpc.
func (d *Dao) Ping(c context.Context) (err error) {
	d.grpcClient.Ping(c, nil)
	return
}

// Close close grpc.
func (d *Dao) Close(c context.Context) (err error) {
	d.grpcClient.Close(c, nil)
	return
}
