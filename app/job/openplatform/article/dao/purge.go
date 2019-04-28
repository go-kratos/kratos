package dao

import (
	"context"
	"net/url"
	"time"

	"go-common/library/log"
)

const _purgeURL = "http://cp.bilibili.co/api_purge.php"

// PurgeCDN purges cdn.
func (d *Dao) PurgeCDN(c context.Context, file string) (err error) {
	defer func() {
		if err == nil {
			return
		}
		time.Sleep(time.Second)
		if e := d.PushCDN(c, file); e != nil {
			log.Error("d.PushCDN(%s) error(%+v)", file, e)
		}
	}()
	params := url.Values{}
	params.Set("sid", "2d0586d2c63fb82a69b20c8992811055")
	params.Set("file", file)
	if err = d.httpClient.Get(c, _purgeURL, "", params, nil); err != nil {
		log.Error("d.httpClient.Get() error(%+v)", err)
		PromError("purge:刷新CDN")
	}
	return
}
