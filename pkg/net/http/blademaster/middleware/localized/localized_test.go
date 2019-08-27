package localized

import (
	"context"
	"net/http"
	"testing"

	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"

	"github.com/stretchr/testify/assert"
)

func TestLocalizedMiddleware(t *testing.T) {
	md := make(metadata.MD)
	ctx := metadata.NewContext(context.Background(), md)
	c1 := &bm.Context{
		Request: &http.Request{Header: http.Header{"Accept-Language": {"zh-CN,zh;q=0.9,en;q=0.8"}}},
		Context: ctx,
	}
	Localized(c1)
	nmd, ok := metadata.FromContext(c1)
	if !ok {
		t.Fatal("get metadata error")
	}
	assert.Contains(t, nmd, "locale")
	assert.Equal(t, nmd["locale"], []string{"zh-CN", "zh", "en"})

	md2 := make(metadata.MD)
	c2 := &bm.Context{
		Request: &http.Request{Header: http.Header{}},
		Context: metadata.NewContext(context.Background(), md2),
	}
	Localized(c2)
	nmd2, ok := metadata.FromContext(c2)
	if !ok {
		t.Fatal("get metadata error")
	}
	assert.Contains(t, nmd2, "locale")
	assert.NotContains(t, nmd2["locale"], "zh-CN")

	md3 := make(metadata.MD)
	c3 := &bm.Context{
		Request: &http.Request{Header: http.Header{"Accept-Language": {"en-US,en;q=0.9,zh-CN;q=0.8,zh-TW;q=0.7,zh-HK;q=0.6,zh;q=0.5"}}},
		Context: metadata.NewContext(context.Background(), md3),
	}
	Localized(c3)
	nmd3, ok := metadata.FromContext(c3)
	if !ok {
		t.Fatal("get metadata error")
	}
	assert.NotContains(t, nmd3, "en-US")
}
