
package service

import(
	"context"

	pb "blog/api/blog/v1"
)

type BlogService struct {
	pb.UnimplementedBlogServer
}

func NewBlogService() pb.BlogServer {
	return &BlogService{}
}

func (s *BlogService) CreateBlog(ctx context.Context, req *pb.CreateBlogRequest) (*pb.CreateBlogReply, error) {
	return &pb.CreateBlogReply{}, nil
}
func (s *BlogService) UpdateBlog(ctx context.Context, req *pb.UpdateBlogRequest) (*pb.UpdateBlogReply, error) {
	return &pb.UpdateBlogReply{}, nil
}
func (s *BlogService) DeleteBlog(ctx context.Context, req *pb.DeleteBlogRequest) (*pb.DeleteBlogReply, error) {
	return &pb.DeleteBlogReply{}, nil
}
func (s *BlogService) GetBlog(ctx context.Context, req *pb.GetBlogRequest) (*pb.GetBlogReply, error) {
	return &pb.GetBlogReply{}, nil
}
func (s *BlogService) ListBlog(ctx context.Context, req *pb.ListBlogRequest) (*pb.ListBlogReply, error) {
	return &pb.ListBlogReply{}, nil
}
