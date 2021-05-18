package service

import (
	"context"
	pb "traces/api/blog/v1"
)

type BlogService struct {
	pb.UnimplementedBlogServer
}

func (s *BlogService) GetBlog(ctx context.Context, request *pb.GetBlogRequest) (*pb.GetBlogReply, error) {
	return &pb.GetBlogReply{Blog: "Teletubbies say hello."}, nil
}

func (s *BlogService) MustEmbedUnimplementedBlogServer() {
	panic("implement me")
}

func NewVehicleService() *BlogService {
	return &BlogService{}
}