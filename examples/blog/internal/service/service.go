package service

import (
	pb "github.com/go-kratos/kratos/examples/blog/api/blog/v1"
	"github.com/go-kratos/kratos/examples/blog/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(NewBlogService)

type BlogService struct {
	pb.UnimplementedBlogServiceServer

	log *log.Helper

	article *biz.ArticleUsecase
}
