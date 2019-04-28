package http

import (
	"go-common/app/interface/bbq/app-bbq/api/http/v1"
	topic "go-common/app/service/bbq/topic/api"
	bm "go-common/library/net/http/blademaster"

	"github.com/pkg/errors"
)

func topicDetail(c *bm.Context) {
	arg := new(topic.TopicVideosReq)
	if err := c.Bind(arg); err != nil {
		errors.Wrap(err, "参数验证失败")
		return
	}
	mid := int64(0)
	midValue, exists := c.Get("mid")
	if exists {
		mid = midValue.(int64)
	}

	c.JSON(srv.TopicDetail(c, mid, arg))
}

func discoveryList(c *bm.Context) {
	arg := new(v1.DiscoveryReq)
	if err := c.Bind(arg); err != nil {
		errors.Wrap(err, "参数验证失败")
		return
	}

	mid := int64(0)
	midValue, exists := c.Get("mid")
	if exists {
		mid = midValue.(int64)
	}

	c.JSON(srv.Discovery(c, mid, arg))
}

func topicSearch(c *bm.Context) {
	arg := new(v1.TopicSearchReq)
	if err := c.Bind(arg); err != nil {
		errors.Wrap(err, "参数验证失败")
		return
	}

	c.JSON(srv.TopicSearch(c, arg))
}
