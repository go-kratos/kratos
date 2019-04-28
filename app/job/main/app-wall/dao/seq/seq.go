package seq

import (
	"context"

	"go-common/app/job/main/app-wall/conf"
	seq "go-common/app/service/main/seq-server/model"
	seqrpc "go-common/app/service/main/seq-server/rpc/client"
	"go-common/library/log"
)

type Dao struct {
	c          *conf.Config
	seqRPC     *seqrpc.Service2
	businessID int64
	token      string
}

func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:          c,
		seqRPC:     seqrpc.New2(c.SeqRPC),
		businessID: c.Seq.BusinessID,
		token:      c.Seq.Token,
	}
	return
}

// SeqID
func (d *Dao) SeqID(ctx context.Context) (requestNo int64, err error) {
	arg := &seq.ArgBusiness{
		BusinessID: d.businessID,
		Token:      d.token,
	}
	if requestNo, err = d.seqRPC.ID(ctx, arg); err != nil {
		log.Error("d.seqRPC.ID error (%v)", err)
		return
	}
	return
}
