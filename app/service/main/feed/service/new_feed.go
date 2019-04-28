package service

// ///////

// struct Resource {
//  Key []string
//  FoldFunc func
// //  列表接口
//  list func(mid int64) ids
//  详情接口(ids []int64 ) // 根据id获取详情
//  expire int // 过期时间
//  len int // 缓存长度
// }
// Resources []Resource

/*
struct Item {
 Key string
 FoldFunc func([]*feed.Feed) ([]*feed.Feed)
 List func(mid int64) (ids []int64) //  列表接口
 Items(ids []int64 ) // 根据id获取详情
 expire int // 过期时间
 len int // 缓存长度
}

struct Resource {
	[]*Item
}
func (s *Service) Feed(c context.Context, tid int64, mid int64, pn, ps int, ip string) (res []*feed.Feed, err error) {
	switch tid {
	case 1:
		feed(c, &Resource{"article"}, mid, pn, ps)
	}
}

func (s *Service) feed(c context.Context, *Resource, mid int64, pn, ps int) (res []*feed.Feed, err error) {
}
// app (archive + bangumi + article)
// app (archive + bangumi)
// web (archive + bangumi)
// article
// archive
// bangumi
// live?
*/
