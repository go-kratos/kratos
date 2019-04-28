package service

import (
	"context"
	"image"
	"mime/multipart"

	httpV1 "go-common/app/service/bbq/video-image/api/http/v1"
	"go-common/app/service/bbq/video-image/conf"
	"go-common/app/service/bbq/video-image/dao"
)

// Service struct
type Service struct {
	c   *conf.Config
	dao *dao.Dao
}

// New init
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:   c,
		dao: dao.New(c),
	}
	return s
}

// Ping Service
func (s *Service) Ping(c context.Context) (err error) {
	return s.dao.Ping(c)
}

// Close Service
func (s *Service) Close() {
	s.dao.Close()
}

// VideoCoverScore .
func (s *Service) VideoCoverScore(c context.Context, name string, f multipart.File) (resp *httpV1.ScoreResponse, err error) {
	img, _, err := image.Decode(f)
	if err != nil {
		return
	}

	imgScorer, err := NewImageScorer(img, 100)
	if err != nil {
		return
	}
	imgScorer.CompareWithPure()

	resp = &httpV1.ScoreResponse{
		Name:  name,
		Score: imgScorer.Score,
	}
	return
}
