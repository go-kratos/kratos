package context

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestContext(t *testing.T) {
	ctx1 := context.WithValue(context.Background(), "go-kratos", "https://github.com/go-kratos/")
	ctx2 := context.WithValue(context.Background(), "kratos", "https://go-kratos.dev/")

	ctx, cancel := Merge(ctx1, ctx2)
	defer cancel()

	got := ctx.Value("go-kratos")
	value1, ok := got.(string)
	assert.Equal(t, ok, true)
	assert.Equal(t, value1,"https://github.com/go-kratos/")
	//
	got2 := ctx.Value("kratos")
	value2, ok := got2.(string)
	assert.Equal(t, ok, true)
	assert.Equal(t, value2,"https://go-kratos.dev/")

	t.Log(value1)
	t.Log(value2)
}
