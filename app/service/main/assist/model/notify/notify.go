package notify

// AddTitle
const (
	AddTitle                = "骑士任命"
	DelTitle                = "解除骑士"
	DelFollowerTitle        = "拒绝加入您的骑士团"
	AddContent              = `你已经被up主#{%s>>}{"http://space.bilibili.com/%d/#!/"} 任命为骑士, 若想取消任命, 可前往关注列表处退出up主#{%s>>}{"http://space.bilibili.com/%d/#!/fans/follow"}的骑士团`
	DelContent              = `你已经被up主#{%s>>}{"http://space.bilibili.com/%d/#!/"} 取消了骑士权限`
	DelFollowerContent      = `#{%s>>}{"http://space.bilibili.com/%d/#!/"} 拒绝成为你的骑士团。退出成功!`
	Mc                      = "1_8_1"
	AddAssNotifyAct         = 1
	DelAssNotifyAct         = 2
	DelAssNotifyFollowerAct = 3
)
