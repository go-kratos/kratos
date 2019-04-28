package service

import (
	"context"
	"flag"
	"fmt"
	"net"
	"path/filepath"
	"testing"
	"time"

	"go-common/app/interface/openplatform/article/conf"
	"go-common/library/cache/redis"

	"github.com/golang/mock/gomock"
	. "github.com/smartystreets/goconvey/convey"
)

var (
	dataID   = int64(175)
	noDataID = int64(1000000000)
	dataMID  = int64(27515309)
	s        *Service
	c        = context.TODO()
)

func CleanCache() {
	pool := redis.NewPool(conf.Conf.Redis)
	pool.Get(c).Do("FLUSHDB")
	conn, _ := net.Dial("tcp", conf.Conf.Memcache.Addr)
	fmt.Fprintf(conn, "flush_all\n")
	conn.Close()
}

func init() {
	dir, _ := filepath.Abs("../cmd/convey-test.toml")
	flag.Set("conf", dir)
	conf.Init()
	s = New(conf.Conf)
	time.Sleep(time.Second)
}

func WithService(f func(s *Service)) func() {
	return func() {
		Reset(func() { CleanCache() })
		f(s)
	}
}

func WithMock(t *testing.T, f func(mock *gomock.Controller)) func() {
	return func() {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		f(mockCtrl)
	}
}

// func httpMock(method, url string) *gock.Request {
// 	r := gock.New(url)
// 	r.Method = strings.ToUpper(method)
// 	return r
// }

func WithCleanCache(f func()) func() {
	return func() {
		Reset(func() { CleanCache() })
		f()
	}
}

/* mysql
INSERT INTO `bilibili_article`.`article_likes_aid_01`(`article_id`, `mid`, `type`) VALUES (2, 1, 1)
INSERT INTO `bilibili_article`.`filtered_articles`(`id`, `article_id`, `mtime`, `ctime`, `category_id`, `title`, `summary`, `banner_url`, `template_id`, `state`, `mid`, `reprint`, `image_urls`, `publish_time`, `attributes`, `words`, `dynamic_intro`) VALUES (473, 175, '2017-06-29 12:42:05', '2017-06-19 19:04:51', 39, '夏至未至', '在感情的围城里，他是个彻底的失败者。对于鲍小姐，他失败与无力抵抗;对于苏小姐，他失败于优柔寡断;对于唐小姐，他失败于无所作为;对于孙小姐，他失败于不能决断。虽不知为何他能受到如此多女子的青睐，但他终究', '/bfs/test/e7b129f2ba8fa59337cbcea2f651b0dd4919fae3.jpg', 4, 0, 175, 1, '/bfs/test/e7b129f2ba8fa59337cbcea2f651b0dd4919fae3.jpg', 1497595087, 0, 0, '');
INSERT INTO `bilibili_article`.`filtered_articles`(`id`, `article_id`, `mtime`, `ctime`, `category_id`, `title`, `summary`, `banner_url`, `template_id`, `state`, `mid`, `reprint`, `image_urls`, `publish_time`, `attributes`, `words`, `dynamic_intro`) VALUES (474, 176, '2017-06-27 16:47:25', '2017-06-19 19:12:12', 39, '鲍小姐12', '在感情的围城里，他是个彻底的失败者。对于鲍小姐，他失败与无力抵抗;对于苏小姐，他失败于优柔寡断;对于唐小姐，他失败于无所作为;对于孙小姐，他失败于不能决断。虽不知为何他能受到如此多女子的青睐，但他终', '/bfs/archive/48f3e2f7b955c190d218ed1e42868469aebaa5a0.png@0-284-1999-1185a_75q.webp', 4, 0, 1, 0, '/bfs/test/1ce440a94a3e2b4f0d81e44d093715cc1e4eae26.jpg', 1498553268, 0, 0, '');
INSERT INTO `bilibili_article`.`article_likes_mid_01`(`id`, `mtime`, `ctime`, `deleted_time`, `mid`, `article_id`, `type`) VALUES (1, '2017-11-29 15:07:08', '2017-11-29 15:07:08', 0, 1, 2, 1);

INSERT INTO `bilibili_article`.`filtered_articles`(`id`, `article_id`, `mtime`, `ctime`, `category_id`, `title`, `summary`, `banner_url`, `template_id`, `state`, `mid`, `reprint`, `image_urls`, `publish_time`, `attributes`, `words`, `dynamic_intro`) VALUES (165, 1, '2017-06-29 12:42:05', '2017-06-19 19:04:51', 39, '夏至未至', '在感情的围城里，他是个彻底的失败者。对于鲍小姐，他失败与无力抵抗;对于苏小姐，他失败于优柔寡断;对于唐小姐，他失败于无所作为;对于孙小姐，他失败于不能决断。虽不知为何他能受到如此多女子的青睐，但他终究', '/bfs/test/e7b129f2ba8fa59337cbcea2f651b0dd4919fae3.jpg', 4, 0, 27515309, 1, '/bfs/test/e7b129f2ba8fa59337cbcea2f651b0dd4919fae3.jpg', 1497595087, 0, 0, '');
INSERT INTO `bilibili_article`.`filtered_articles`(`id`, `article_id`, `mtime`, `ctime`, `category_id`, `title`, `summary`, `banner_url`, `template_id`, `state`, `mid`, `reprint`, `image_urls`, `publish_time`, `attributes`, `words`, `dynamic_intro`) VALUES (476, 2, '2017-06-22 00:11:19', '2017-06-19 19:13:48', 29, '《弹幕音乐绘》：从没玩过STG的人，也能体验到弹幕的乐趣', 'AlphaGO战胜柯洁、Google翻译准确度大幅提升、自动驾驶技术日趋成熟；这些无一不在告诉我们人工智能技术正飞速发展。那么如果让人工智能来玩STG会怎么样呢？\nSTG也称为弹幕游戏。往往由于满屏幕', '/bfs/test/4d1516ddad93cdb1fc31210b0e2f3526f65d0d6b.jpg', 4, 0, 482, 0, '/bfs/test/4d1516ddad93cdb1fc31210b0e2f3526f65d0d6b.jpg', 1497871027, 0, 0, '');
INSERT INTO `bilibili_article`.`filtered_articles`(`id`, `article_id`, `mtime`, `ctime`, `category_id`, `title`, `summary`, `banner_url`, `template_id`, `state`, `mid`, `reprint`, `image_urls`, `publish_time`, `attributes`, `words`, `dynamic_intro`) VALUES (477, 3, '2017-10-19 10:34:07', '2017-06-19 19:21:28', 30, '从没玩过STG的人，也能体验到弹幕的乐趣1', 'AlphaGO战胜柯洁、Google翻译准确度大幅提升、自动驾驶技术日趋成熟；这些无一不在告诉我们人工智能技术正飞速发展。那么如果让人工智能来玩STG会怎么样呢？\n\n视频链接\nSTG也称为弹幕游戏。', '/bfs/test/7797d025e1ded27521f10be42d247e5bc0c85ac1.jpg', 4, 0, 482, 1, '/bfs/test/7797d025e1ded27521f10be42d247e5bc0c85ac1.jpg', 1497871310, 0, 522, '');
INSERT INTO `bilibili_article`.`filtered_articles`(`id`, `article_id`, `mtime`, `ctime`, `category_id`, `title`, `summary`, `banner_url`, `template_id`, `state`, `mid`, `reprint`, `image_urls`, `publish_time`, `attributes`, `words`, `dynamic_intro`) VALUES (478, 4, '2017-06-22 00:11:19', '2017-06-19 19:25:31', 38, 'Product managers for the digital world', 'The role of the product manager is expanding due to the growing importance of data in decision makin', '/bfs/test/fafe56c7488e7def062ae9d77c0fb507b119ea1d.jpg', 4, 0, 2089809, 0, '/bfs/test/fafe56c7488e7def062ae9d77c0fb507b119ea1d.jpg', 1497871547, 0, 0, '');
INSERT INTO `bilibili_article`.`article_notices`(`id`, `title`, `url`, `stime`, `etime`, `state`, `ctime`, `mtime`) VALUES (2, 'wuhao test2 edit', 'http://www.bilibilii.com', '2017-12-12 00:00:00', '2017-12-29 00:00:00', 1, '2017-12-12 18:26:12', '2017-12-15 15:28:47');

INSERT INTO `bilibili_article`.`articles`(`id`, `mtime`, `ctime`, `deleted_time`, `category_id`, `title`, `summary`, `banner_url`, `template_id`, `state`, `mid`, `reprint`, `image_urls`, `publish_time`, `reason`, `attributes`, `words`, `dynamic_intro`, `origin_image_urls`) VALUES (599, '2017-12-13 16:23:39', '2017-09-12 17:55:56', 0, 38, '规划局个乖宝宝？啊啊啊', '巴巴大热天团团圆圆和和美美的团团圆圆一样大热天团团圆圆和和美美的团团圆圆一样大热天团团圆圆和和美美的团团圆圆一样大热天团团圆圆和和美美的团团圆圆一样大热天团团圆圆和和美美的团团圆圆一样大热天团团圆圆和', '/bfs/article/09fad83873aeec34fe7885ffe71c896c4565a8ca.jpg', 4, -2, 88888929, 0, '/bfs/article/09fad83873aeec34fe7885ffe71c896c4565a8ca.jpg', 0, '', 0, 218, '', '/bfs/article/09fad83873aeec34fe7885ffe71c896c4565a8ca.jpg');
INSERT INTO `bilibili_article`.`articles`(`id`, `mtime`, `ctime`, `deleted_time`, `category_id`, `title`, `summary`, `banner_url`, `template_id`, `state`, `mid`, `reprint`, `image_urls`, `publish_time`, `reason`, `attributes`, `words`, `dynamic_intro`, `origin_image_urls`) VALUES (600, '2017-12-13 16:23:39', '2017-09-12 18:11:45', 0, 25, '啦啦啦', '大热天团团圆圆和和美美的团团圆圆一样大热天团团圆圆和和美美的团团圆圆一样大热天团团圆圆和和美美的团团圆圆一样大热天团团圆圆和和美美的团团圆圆一样大热天团团圆圆和和美美的团团圆圆一样大热天团团圆圆和和美', '/bfs/article/61ba2a38fce91de39608cf4f57309b53839534a8.jpg', 4, -2, 88888929, 0, '/bfs/article/61ba2a38fce91de39608cf4f57309b53839534a8.jpg', 0, '', 0, 216, '', '/bfs/article/61ba2a38fce91de39608cf4f57309b53839534a8.jpg');
INSERT INTO `bilibili_article`.`lists`(`id`, `mtime`, `ctime`, `deleted_time`, `image_url`, `mid`, `name`, `update_time`) VALUES (1, '2018-01-26 15:44:53', '2018-01-26 15:44:53', 0, '', 100, 'name', '0000-00-00 00:00:00');
INSERT INTO `bilibili_article`.`article_lists`(`id`, `mtime`, `ctime`, `deleted_time`, `article_id`, `list_id`, `position`) VALUES (1, '2018-01-27 13:39:03', '2018-01-27 13:39:03', 0, 165, 1, 0);
INSERT INTO `bilibili_article`.`article_lists`(`id`, `mtime`, `ctime`, `deleted_time`, `article_id`, `list_id`, `position`) VALUES (2, '2018-01-27 13:39:09', '2018-01-27 13:39:09', 0, 476, 1, 1);
INSERT INTO `bilibili_article`.`lists`(`id`, `mtime`, `ctime`, `deleted_time`, `image_url`, `mid`, `name`, `update_time`) VALUES (8, '2018-01-27 18:32:05', '2018-01-26 16:06:57', 0, '', 88888929, 'name', '2018-01-26 16:06:57');
INSERT INTO `bilibili_article`.`article_lists`(`id`, `mtime`, `ctime`, `deleted_time`, `article_id`, `list_id`, `position`) VALUES (9, '2018-01-28 14:57:53', '2018-01-28 14:57:53', 0, 165, 5, 0);
INSERT INTO `bilibili_article`.`article_lists`(`id`, `mtime`, `ctime`, `deleted_time`, `article_id`, `list_id`, `position`) VALUES (10, '2018-01-28 14:58:36', '2018-01-28 14:58:36', 0, 476, 5, 0);
*/
