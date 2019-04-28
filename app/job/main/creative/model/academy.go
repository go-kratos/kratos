package model

const (
	//BusinessForArchvie 稿件
	BusinessForArchvie = 1
	//BusinessForArticle 专栏
	BusinessForArticle = 2
)

//OArchive for academy.
type OArchive struct {
	ID       int64
	OID      int64
	Business int
}
