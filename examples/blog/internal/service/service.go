package service

import (
	pb "github.com/SeeMusic/kratos/examples/blog/api/blog/v1"
	"github.com/SeeMusic/kratos/examples/blog/internal/biz"

	"github.com/SeeMusic/kratos/v2/log"
	"github.com/google/wire"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(NewBlogService)

type BlogService struct {
	pb.UnimplementedBlogServiceServer

	log *log.Helper

	article *biz.ArticleUsecase
}
