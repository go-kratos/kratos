package http

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"testing"
	"time"

	"go-common/app/interface/main/favorite/conf"
	"go-common/app/interface/main/favorite/service"
	xhttp "go-common/library/net/http/blademaster"
)

const (
	_mid       = 88888894
	_vmid      = 12345
	_aid       = 5463438
	_aids      = "5463438,5463439"
	_fid       = 1852
	_delFid    = 123
	_oldFid    = 1791
	_newFid    = 1792
	_fidsSort  = "1107,1852,1792,1791"
	_name      = "folder-name-test"
	_rename    = "folder-rename-test"
	_public    = 1
	_tpid      = 2659 // TopicID
	_type      = 1    // Article
	_v3Fid     = 0
	_searchFid = 0
	_oid       = 123
	_pn        = 1
	_ps        = 30
	// video folder
	_videoFolders       = "http://127.0.0.1:6010/x/internal/v2/fav/folder"
	_addVideoFolder     = "http://127.0.0.1:6010/x/internal/v2/fav/folder/add"
	_delVideoFolder     = "http://127.0.0.1:6010/x/internal/v2/fav/folder/del"
	_renameVideoFolder  = "http://127.0.0.1:6010/x/internal/v2/fav/folder/rename"
	_upStateVideoFolder = "http://127.0.0.1:6010/x/internal/v2/fav/folder/public"
	_sortVideoFolders   = "http://127.0.0.1:6010/x/internal/v2/fav/folder/sort"
	// video
	_favVideo       = "http://127.0.0.1:6010/x/internal/v2/fav/video"
	_favVideoNewest = "http://127.0.0.1:6010/x/internal/v2/fav/video/newest"
	_addFavVideo    = "http://127.0.0.1:6010/x/internal/v2/fav/video/add"
	_delFavVideo    = "http://127.0.0.1:6010/x/internal/v2/fav/video/del"
	_delFavVideos   = "http://127.0.0.1:6010/x/internal/v2/fav/video/mdel"
	_moveFavVideos  = "http://127.0.0.1:6010/x/internal/v2/fav/video/move"
	_copyFavVideos  = "http://127.0.0.1:6010/x/internal/v2/fav/video/copy"
	_isFavoureds    = "http://127.0.0.1:6010/x/internal/v2/fav/video/favoureds"
	_isFavoured     = "http://127.0.0.1:6010/x/internal/v2/fav/video/favoured"
	_inDefaultFav   = "http://127.0.0.1:6010/x/internal/v2/fav/video/default"
	// topic
	_favTopics       = "http://127.0.0.1:6010/x/internal/v2/fav/topic"
	_addFavTopic     = "http://127.0.0.1:6010/x/internal/v2/fav/topic/add"
	_delFavTopic     = "http://127.0.0.1:6010/x/internal/v2/fav/topic/del"
	_isTopicFavoured = "http://127.0.0.1:6010/x/internal/v2/fav/topic/favoured"
	// fav v3
	_favorites = "http://127.0.0.1:6010/x/internal/v3/fav"
	_addFav    = "http://127.0.0.1:6010/x/internal/v3/fav/add"
	_isFavored = "http://127.0.0.1:6010/x/internal/v3/fav/favored"
	_delFav    = "http://127.0.0.1:6010/x/internal/v3/fav/del"
)

func TestHttp(t *testing.T) {
	if err := conf.Init(); err != nil {
		t.Fatalf("conf.Init() error(%v)", err)
	}
	svr := service.New(conf.Conf)
	client := xhttp.NewClient(conf.Conf.HTTPClient)
	Init(conf.Conf, svr)

	// video foler
	testVideoFolders(client, t, _mid, _vmid, _aid)
	testRenameVideoFolder(client, t, _mid, _fid, _rename)
	testAddVideoFolder(client, t, _mid, _public, _name)
	testSortVideoFolder(client, t, _mid, _fidsSort)
	testUpStateVideoFolder(client, t, _mid, _fid, _public)
	testDelVideoFolder(client, t, _mid, _delFid)
	// video
	testVideos(client, t, _mid, _vmid, _fid)
	testFavVideoNewest(client, t, _mid, _vmid, _searchFid, _pn, _ps)
	testAddFavVideo(client, t, _mid, _fid, _aid)
	testMoveFavVideos(client, t, _mid, _oldFid, _newFid, _aids)
	testCopyFavVideos(client, t, _mid, _oldFid, _newFid, _aids)
	testIsFavoured(client, t, _mid, _aid)
	testIsFavoureds(client, t, _mid, _aids)
	testInDefaultFav(client, t, _mid, _aid)
	testDelVideo(client, t, _mid, _fid, _aid)
	testDelVideos(client, t, _mid, _fid, _aids)
	// topic
	testFavTopics(client, t, _mid, _pn, _ps)
	testAddFavTopic(client, t, _mid, _tpid)
	testIsTopicFavoured(client, t, _mid, _tpid)
	testDelFavTopic(client, t, _mid, _tpid)
	// fav v3
	testFavorites(client, t, _type, _mid, _vmid, _v3Fid)
	testAddFav(client, t, _type, _mid, _v3Fid, _oid)
	testIsFavored(client, t, _type, _mid, _v3Fid, _oid)
	testDelFav(client, t, _type, _mid, _v3Fid, _oid)
}

func testVideoFolders(client *xhttp.Client, t *testing.T, mid, vmid, aid int64) {
	params := &url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("appkey", conf.Conf.App.Key)
	params.Set("ts", strconv.FormatInt(time.Now().Unix(), 10))
	params.Set("sign", createSign(params.Encode()))
	// send
	req, err := http.NewRequest("GET", _videoFolders+"?"+params.Encode(), nil)
	t.Log(req.URL.String())
	if err != nil {
		t.Errorf("NewRequest() error(%v)", err)
	}
	res := map[string]interface{}{}
	if err = client.Do(context.TODO(), req, &res); err != nil {
		t.Errorf("client.Do() error(%v)", err)
	} else {
		result("VideoFolders", t, res)
	}
}

func testAddVideoFolder(client *xhttp.Client, t *testing.T, mid, public int64, name string) {
	params := &url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("public", strconv.FormatInt(public, 10))
	params.Set("name", name)
	params.Set("appkey", conf.Conf.App.Key)
	params.Set("ts", strconv.FormatInt(time.Now().Unix(), 10))
	params.Set("sign", createSign(params.Encode()))
	// send
	req, err := http.NewRequest("POST", _addVideoFolder+"?"+params.Encode(), nil)
	t.Log(req.URL.String())
	if err != nil {
		t.Errorf("NewRequest() error(%v)", err)
	}
	res := map[string]interface{}{}
	if err = client.Do(context.TODO(), req, &res); err != nil {
		t.Errorf("client.Do() error(%v)", err)
	} else {
		result("DelVideoFolder", t, res)
	}
}

func testRenameVideoFolder(client *xhttp.Client, t *testing.T, mid, fid int64, name string) {
	params := &url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("fid", strconv.FormatInt(fid, 10))
	params.Set("name", name)
	params.Set("appkey", conf.Conf.App.Key)
	params.Set("ts", strconv.FormatInt(time.Now().Unix(), 10))
	params.Set("sign", createSign(params.Encode()))
	// send
	req, err := http.NewRequest("POST", _renameVideoFolder+"?"+params.Encode(), nil)
	t.Log(req.URL.String())
	if err != nil {
		t.Errorf("NewRequest() error(%v)", err)
	}
	res := map[string]interface{}{}
	if err = client.Do(context.TODO(), req, &res); err != nil {
		t.Errorf("client.Do() error(%v)", err)
	} else {
		result("RenameVideoFolder", t, res)
	}
}

func testUpStateVideoFolder(client *xhttp.Client, t *testing.T, mid, fid, public int64) {
	params := &url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("fid", strconv.FormatInt(fid, 10))
	params.Set("public", strconv.FormatInt(public, 10))
	params.Set("appkey", conf.Conf.App.Key)
	params.Set("ts", strconv.FormatInt(time.Now().Unix(), 10))
	params.Set("sign", createSign(params.Encode()))
	// send
	req, err := http.NewRequest("POST", _upStateVideoFolder+"?"+params.Encode(), nil)
	t.Log(req.URL.String())
	if err != nil {
		t.Errorf("NewRequest() error(%v)", err)
	}
	res := map[string]interface{}{}
	if err = client.Do(context.TODO(), req, &res); err != nil {
		t.Errorf("client.Do() error(%v)", err)
	} else {
		result("RenameVideoFolder", t, res)
	}
}

func testSortVideoFolder(client *xhttp.Client, t *testing.T, mid int64, fids string) {
	params := &url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("fids", fids)
	params.Set("appkey", conf.Conf.App.Key)
	params.Set("ts", strconv.FormatInt(time.Now().Unix(), 10))
	params.Set("sign", createSign(params.Encode()))
	// send
	req, err := http.NewRequest("POST", _sortVideoFolders+"?"+params.Encode(), nil)
	t.Log(req.URL.String())
	if err != nil {
		t.Errorf("NewRequest() error(%v)", err)
	}
	res := map[string]interface{}{}
	if err = client.Do(context.TODO(), req, &res); err != nil {
		t.Errorf("client.Do() error(%v)", err)
	} else {
		result("RenameVideoFolder", t, res)
	}
}

func testDelVideoFolder(client *xhttp.Client, t *testing.T, mid, fid int64) {
	params := &url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("fid", strconv.FormatInt(fid, 10))
	params.Set("appkey", conf.Conf.App.Key)
	params.Set("ts", strconv.FormatInt(time.Now().Unix(), 10))
	params.Set("sign", createSign(params.Encode()))
	// send
	req, err := http.NewRequest("POST", _delVideoFolder+"?"+params.Encode(), nil)
	t.Log(req.URL.String())
	if err != nil {
		t.Errorf("NewRequest() error(%v)", err)
	}
	res := map[string]interface{}{}
	if err = client.Do(context.TODO(), req, &res); err != nil {
		t.Errorf("client.Do() error(%v)", err)
	} else {
		result("DelVideoFolder", t, res)
	}
}

func testVideos(client *xhttp.Client, t *testing.T, mid, vmid, fid int64) {
	params := &url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("fid", strconv.FormatInt(fid, 10))
	params.Set("appkey", conf.Conf.App.Key)
	params.Set("ts", strconv.FormatInt(time.Now().Unix(), 10))
	params.Set("sign", createSign(params.Encode()))
	// send
	req, err := http.NewRequest("GET", _favVideo+"?"+params.Encode(), nil)
	t.Log(req.URL.String())
	if err != nil {
		t.Errorf("NewRequest() error(%v)", err)
	}
	res := map[string]interface{}{}
	if err = client.Do(context.TODO(), req, &res); err != nil {
		t.Errorf("client.Do() error(%v)", err)
	} else {
		result("Videos", t, res)
	}
}

func testFavVideoNewest(client *xhttp.Client, t *testing.T, mid, vmid, fid, pn, ps int64) {
	params := &url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("fid", strconv.FormatInt(fid, 10))
	params.Set("pn", strconv.FormatInt(pn, 10))
	params.Set("ps", strconv.FormatInt(ps, 10))
	params.Set("appkey", conf.Conf.App.Key)
	params.Set("ts", strconv.FormatInt(time.Now().Unix(), 10))
	params.Set("sign", createSign(params.Encode()))
	// send
	req, err := http.NewRequest("GET", _favVideoNewest+"?"+params.Encode(), nil)
	t.Log(req.URL.String())
	if err != nil {
		t.Errorf("NewRequest() error(%v)", err)
	}
	res := map[string]interface{}{}
	if err = client.Do(context.TODO(), req, &res); err != nil {
		t.Errorf("client.Do() error(%v)", err)
	} else {
		result("FavVideoNewest", t, res)
	}
}

func testAddFavVideo(client *xhttp.Client, t *testing.T, mid, fid, aid int64) {
	params := &url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("fid", strconv.FormatInt(fid, 10))
	params.Set("aid", strconv.FormatInt(aid, 10))
	params.Set("appkey", conf.Conf.App.Key)
	params.Set("ts", strconv.FormatInt(time.Now().Unix(), 10))
	params.Set("sign", createSign(params.Encode()))
	// send
	req, err := http.NewRequest("POST", _addFavVideo+"?"+params.Encode(), nil)
	t.Log(req.URL.String())
	if err != nil {
		t.Errorf("NewRequest() error(%v)", err)
	}
	res := map[string]interface{}{}
	if err = client.Do(context.TODO(), req, &res); err != nil {
		t.Errorf("client.Do() error(%v)", err)
	} else {
		result("AddFavVideo", t, res)
	}
}

func testMoveFavVideos(client *xhttp.Client, t *testing.T, mid, oldFid, newFid int64, aids string) {
	params := &url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("old_fid", strconv.FormatInt(oldFid, 10))
	params.Set("new_fid", strconv.FormatInt(newFid, 10))
	params.Set("aids", aids)
	params.Set("appkey", conf.Conf.App.Key)
	params.Set("ts", strconv.FormatInt(time.Now().Unix(), 10))
	params.Set("sign", createSign(params.Encode()))
	// send
	req, err := http.NewRequest("POST", _moveFavVideos+"?"+params.Encode(), nil)
	t.Log(req.URL.String())
	if err != nil {
		t.Errorf("NewRequest() error(%v)", err)
	}
	res := map[string]interface{}{}
	if err = client.Do(context.TODO(), req, &res); err != nil {
		t.Errorf("client.Do() error(%v)", err)
	} else {
		result("MoveFavVideos", t, res)
	}
}

func testCopyFavVideos(client *xhttp.Client, t *testing.T, mid, oldFid, newFid int64, aids string) {
	params := &url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("old_fid", strconv.FormatInt(oldFid, 10))
	params.Set("new_fid", strconv.FormatInt(newFid, 10))
	params.Set("aids", aids)
	params.Set("appkey", conf.Conf.App.Key)
	params.Set("ts", strconv.FormatInt(time.Now().Unix(), 10))
	params.Set("sign", createSign(params.Encode()))
	// send
	req, err := http.NewRequest("POST", _copyFavVideos+"?"+params.Encode(), nil)
	t.Log(req.URL.String())
	if err != nil {
		t.Errorf("NewRequest() error(%v)", err)
	}
	res := map[string]interface{}{}
	if err = client.Do(context.TODO(), req, &res); err != nil {
		t.Errorf("client.Do() error(%v)", err)
	} else {
		result("CopyFavVideos", t, res)
	}
}

func testIsFavoured(client *xhttp.Client, t *testing.T, mid, aid int64) {
	params := &url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("aid", strconv.FormatInt(aid, 10))
	params.Set("appkey", conf.Conf.App.Key)
	params.Set("ts", strconv.FormatInt(time.Now().Unix(), 10))
	params.Set("sign", createSign(params.Encode()))
	// send
	req, err := http.NewRequest("GET", _isFavoured+"?"+params.Encode(), nil)
	t.Log(req.URL.String())
	if err != nil {
		t.Errorf("NewRequest() error(%v)", err)
	}
	res := map[string]interface{}{}
	if err = client.Do(context.TODO(), req, &res); err != nil {
		t.Errorf("client.Do() error(%v)", err)
	} else {
		result("IsFavoured", t, res)
	}
}

func testIsFavoureds(client *xhttp.Client, t *testing.T, mid int64, aids string) {
	params := &url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("aids", aids)
	params.Set("appkey", conf.Conf.App.Key)
	params.Set("ts", strconv.FormatInt(time.Now().Unix(), 10))
	params.Set("sign", createSign(params.Encode()))
	// send
	req, err := http.NewRequest("GET", _isFavoureds+"?"+params.Encode(), nil)
	t.Log(req.URL.String())
	if err != nil {
		t.Errorf("NewRequest() error(%v)", err)
	}
	res := map[string]interface{}{}
	if err = client.Do(context.TODO(), req, &res); err != nil {
		t.Errorf("client.Do() error(%v)", err)
	} else {
		result("IsFavoureds", t, res)
	}
}

func testInDefaultFav(client *xhttp.Client, t *testing.T, mid, aid int64) {
	params := &url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("aid", strconv.FormatInt(aid, 10))
	params.Set("appkey", conf.Conf.App.Key)
	params.Set("ts", strconv.FormatInt(time.Now().Unix(), 10))
	params.Set("sign", createSign(params.Encode()))
	// send
	req, err := http.NewRequest("GET", _inDefaultFav+"?"+params.Encode(), nil)
	t.Log(req.URL.String())
	if err != nil {
		t.Errorf("NewRequest() error(%v)", err)
	}
	res := map[string]interface{}{}
	if err = client.Do(context.TODO(), req, &res); err != nil {
		t.Errorf("client.Do() error(%v)", err)
	} else {
		result("IsFavoureds", t, res)
	}
}

func testDelVideo(client *xhttp.Client, t *testing.T, mid, fid, aid int64) {
	params := &url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("aid", strconv.FormatInt(aid, 10))
	params.Set("appkey", conf.Conf.App.Key)
	params.Set("ts", strconv.FormatInt(time.Now().Unix(), 10))
	params.Set("sign", createSign(params.Encode()))
	// send
	req, err := http.NewRequest("POST", _delFavVideo+"?"+params.Encode(), nil)
	t.Log(req.URL.String())
	if err != nil {
		t.Errorf("NewRequest() error(%v)", err)
	}
	res := map[string]interface{}{}
	if err = client.Do(context.TODO(), req, &res); err != nil {
		t.Errorf("client.Do() error(%v)", err)
	} else {
		result("DelVideo", t, res)
	}
}

func testDelVideos(client *xhttp.Client, t *testing.T, mid, fid int64, aids string) {
	params := &url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("fid", strconv.FormatInt(fid, 10))
	params.Set("aids", aids)
	params.Set("appkey", conf.Conf.App.Key)
	params.Set("ts", strconv.FormatInt(time.Now().Unix(), 10))
	params.Set("sign", createSign(params.Encode()))
	// send
	req, err := http.NewRequest("POST", _delFavVideos+"?"+params.Encode(), nil)
	t.Log(req.URL.String())
	if err != nil {
		t.Errorf("NewRequest() error(%v)", err)
	}
	res := map[string]interface{}{}
	if err = client.Do(context.TODO(), req, &res); err != nil {
		t.Errorf("client.Do() error(%v)", err)
	} else {
		result("DelVideos", t, res)
	}
}

func testFavTopics(client *xhttp.Client, t *testing.T, mid, pn, ps int64) {
	params := &url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("pn", strconv.FormatInt(pn, 10))
	params.Set("ps", strconv.FormatInt(ps, 10))
	params.Set("appkey", conf.Conf.App.Key)
	params.Set("ts", strconv.FormatInt(time.Now().Unix(), 10))
	params.Set("sign", createSign(params.Encode()))
	// send
	req, err := http.NewRequest("GET", _favTopics+"?"+params.Encode(), nil)
	t.Log(req.URL.String())
	if err != nil {
		t.Errorf("NewRequest() error(%v)", err)
	}
	res := map[string]interface{}{}
	if err = client.Do(context.TODO(), req, &res); err != nil {
		t.Errorf("client.Do() error(%v)", err)
	} else {
		result("FavTopics", t, res)
	}
}

func testAddFavTopic(client *xhttp.Client, t *testing.T, mid, tpid int64) {
	params := &url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("tpid", strconv.FormatInt(tpid, 10))
	params.Set("appkey", conf.Conf.App.Key)
	params.Set("ts", strconv.FormatInt(time.Now().Unix(), 10))
	params.Set("sign", createSign(params.Encode()))
	// send
	req, err := http.NewRequest("POST", _addFavTopic+"?"+params.Encode(), nil)
	t.Log(req.URL.String())
	if err != nil {
		t.Errorf("NewRequest() error(%v)", err)
	}
	res := map[string]interface{}{}
	if err = client.Do(context.TODO(), req, &res); err != nil {
		t.Errorf("client.Do() error(%v)", err)
	} else {
		result("AddFavTopic", t, res)
	}
}

func testIsTopicFavoured(client *xhttp.Client, t *testing.T, mid, tpid int64) {
	params := &url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("tpid", strconv.FormatInt(tpid, 10))
	params.Set("appkey", conf.Conf.App.Key)
	params.Set("ts", strconv.FormatInt(time.Now().Unix(), 10))
	params.Set("sign", createSign(params.Encode()))
	// send
	req, err := http.NewRequest("GET", _isTopicFavoured+"?"+params.Encode(), nil)
	t.Log(req.URL.String())
	if err != nil {
		t.Errorf("NewRequest() error(%v)", err)
	}
	res := map[string]interface{}{}
	if err = client.Do(context.TODO(), req, &res); err != nil {
		t.Errorf("client.Do() error(%v)", err)
	} else {
		result("IsTopicFavoured", t, res)
	}
}

func testDelFavTopic(client *xhttp.Client, t *testing.T, mid, tpid int64) {
	params := &url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("tpid", strconv.FormatInt(tpid, 10))
	params.Set("appkey", conf.Conf.App.Key)
	params.Set("ts", strconv.FormatInt(time.Now().Unix(), 10))
	params.Set("sign", createSign(params.Encode()))
	// send
	req, err := http.NewRequest("POST", _delFavTopic+"?"+params.Encode(), nil)
	t.Log(req.URL.String())
	if err != nil {
		t.Errorf("NewRequest() error(%v)", err)
	}
	res := map[string]interface{}{}
	if err = client.Do(context.TODO(), req, &res); err != nil {
		t.Errorf("client.Do() error(%v)", err)
	} else {
		result("DelFavTopic", t, res)
	}
}

func testFavorites(client *xhttp.Client, t *testing.T, tp, mid, vmid, fid int64) {
	params := &url.Values{}
	params.Set("type", strconv.FormatInt(tp, 10))
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("fid", strconv.FormatInt(fid, 10))
	params.Set("appkey", conf.Conf.App.Key)
	params.Set("ts", strconv.FormatInt(time.Now().Unix(), 10))
	params.Set("sign", createSign(params.Encode()))
	// send
	req, err := http.NewRequest("GET", _favorites+"?"+params.Encode(), nil)
	t.Log(req.URL.String())
	if err != nil {
		t.Errorf("NewRequest() error(%v)", err)
	}
	res := map[string]interface{}{}
	if err = client.Do(context.TODO(), req, &res); err != nil {
		t.Errorf("client.Do() error(%v)", err)
	} else {
		result("Favorites", t, res)
	}
}

func testAddFav(client *xhttp.Client, t *testing.T, tp, mid, fid, oid int64) {
	params := &url.Values{}
	params.Set("type", strconv.FormatInt(tp, 10))
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("fid", strconv.FormatInt(fid, 10))
	params.Set("oid", strconv.FormatInt(oid, 10))
	params.Set("appkey", conf.Conf.App.Key)
	params.Set("ts", strconv.FormatInt(time.Now().Unix(), 10))
	params.Set("sign", createSign(params.Encode()))
	// send
	req, err := http.NewRequest("POST", _addFav+"?"+params.Encode(), nil)
	t.Log(req.URL.String())
	if err != nil {
		t.Errorf("NewRequest() error(%v)", err)
	}
	res := map[string]interface{}{}
	if err = client.Do(context.TODO(), req, &res); err != nil {
		t.Errorf("client.Do() error(%v)", err)
	} else {
		result("AddFav", t, res)
	}
}

func testIsFavored(client *xhttp.Client, t *testing.T, tp, mid, fid, oid int64) {
	params := &url.Values{}
	params.Set("type", strconv.FormatInt(tp, 10))
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("fid", strconv.FormatInt(fid, 10))
	params.Set("oid", strconv.FormatInt(oid, 10))
	params.Set("appkey", conf.Conf.App.Key)
	params.Set("ts", strconv.FormatInt(time.Now().Unix(), 10))
	params.Set("sign", createSign(params.Encode()))
	// send
	req, err := http.NewRequest("GET", _isFavored+"?"+params.Encode(), nil)
	t.Log(req.URL.String())
	if err != nil {
		t.Errorf("NewRequest() error(%v)", err)
	}
	res := map[string]interface{}{}
	if err = client.Do(context.TODO(), req, &res); err != nil {
		t.Errorf("client.Do() error(%v)", err)
	} else {
		result("IsFavored", t, res)
	}
}

func testDelFav(client *xhttp.Client, t *testing.T, tp, mid, fid, oid int64) {
	params := &url.Values{}
	params.Set("type", strconv.FormatInt(tp, 10))
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("fid", strconv.FormatInt(fid, 10))
	params.Set("oid", strconv.FormatInt(oid, 10))
	params.Set("appkey", conf.Conf.App.Key)
	params.Set("ts", strconv.FormatInt(time.Now().Unix(), 10))
	params.Set("sign", createSign(params.Encode()))
	// send
	req, err := http.NewRequest("POST", _delFav+"?"+params.Encode(), nil)
	t.Log(req.URL.String())
	if err != nil {
		t.Errorf("NewRequest() error(%v)", err)
	}
	res := map[string]interface{}{}
	if err = client.Do(context.TODO(), req, &res); err != nil {
		t.Errorf("client.Do() error(%v)", err)
	} else {
		result("DelFav", t, res)
	}
}

func createSign(params string) string {
	mh := md5.Sum([]byte(params + conf.Conf.App.Secret))
	return hex.EncodeToString(mh[:])
}

func result(name string, t *testing.T, res map[string]interface{}) {
	t.Log("[==========" + name + " Testing Result==========]")
	if rs, ok := res["code"]; ok {
		t.Log(fmt.Sprintf("code:%v, message:%s, data:%v", rs, res["message"], res["data"]))
	}
	t.Log("[↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑]\r\n")
}
