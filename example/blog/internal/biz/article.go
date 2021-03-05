package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"time"
)

type Article struct {
	Id        int64
	Title     string
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
	Like      int64
}

type ArticleRepo interface {
	// db
	ListArticle(ctx context.Context) ([]*Article, error)
	GetArticle(ctx context.Context, id int64) (*Article, error)
	CreateArticle(ctx context.Context, article *Article) error
	UpdateArticle(ctx context.Context, id int64, article *Article) error
	DeleteArticle(ctx context.Context, id int64) error

	// redis
	GetArticleLike(ctx context.Context, id int64) (rv int64, err error)
	IncArticleLike(ctx context.Context, id int64) error
}

type ArticleUsecase struct {
	repo ArticleRepo
}

func NewArticleUsecase(repo ArticleRepo, logger log.Logger) *ArticleUsecase {
	return &ArticleUsecase{repo: repo}
}

func (uc *ArticleUsecase) List(ctx context.Context) (ps []*Article, err error) {
	ps, err = uc.repo.ListArticle(ctx)
	if err != nil {
		return
	}
	return
}

func (uc *ArticleUsecase) Get(ctx context.Context, id int64) (p *Article, err error) {
	p, err = uc.repo.GetArticle(ctx, id)
	if err != nil {
		return
	}
	err = uc.repo.IncArticleLike(ctx, id)
	if err != nil {
		return
	}
	p.Like, err = uc.repo.GetArticleLike(ctx, id)
	if err != nil {
		return
	}
	return
}

func (uc *ArticleUsecase) Create(ctx context.Context, article *Article) error {
	return uc.repo.CreateArticle(ctx, article)
}

func (uc *ArticleUsecase) Update(ctx context.Context, id int64, article *Article) error {
	return uc.repo.UpdateArticle(ctx, id, article)
}

func (uc *ArticleUsecase) Delete(ctx context.Context, id int64) error {
	return uc.repo.DeleteArticle(ctx, id)
}
