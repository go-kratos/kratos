package v1

import (
	"context"
	"time"

	v1pb "go-common/app/admin/live/live-admin/api/http/v1"
	"go-common/app/admin/live/live-admin/conf"
	"go-common/app/admin/live/live-admin/dao"
	"go-common/library/ecode"
	"go-common/library/log"
)

// TokenService struct
type TokenService struct {
	conf *conf.Config
	// optionally add other properties here, such as dao
	dao *dao.Dao
}

//NewTokenService init
func NewTokenService(c *conf.Config, d *dao.Dao) (s *TokenService) {
	s = &TokenService{
		conf: c,
		dao:  d,
	}
	return s
}

// New implementation
// Request for a token for upload.
// `method:"POST" internal:"true"`
func (s *TokenService) New(ctx context.Context, req *v1pb.NewTokenReq) (resp *v1pb.NewTokenResp, err error) {
	// Must be live's bucket.
	_, ok := s.conf.Bucket[req.Bucket]
	if !ok {
		err = ecode.UploadBucketErr
		return
	}

	var token string
	if token, err = s.dao.RequestUploadToken(ctx, req.Bucket, req.Operator, time.Now().Unix()); err != nil {
		log.Error("New a upload token failure: %v", err)
		err = ecode.UploadTokenGenErr
		return
	}

	resp = &v1pb.NewTokenResp{
		Token: token,
	}

	return
}
