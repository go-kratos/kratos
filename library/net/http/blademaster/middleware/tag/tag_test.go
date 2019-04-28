package tag_test

import (
	"context"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"go-common/library/log"
	"go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/tag"

	"github.com/stretchr/testify/assert"
)

func init() {
	log.Init(nil)
}

func makeBlueGreenTag() *tag.Tag {
	var pf tag.PolicyFunc
	pf = func(ctx *blademaster.Context) string {
		if ctx.Request.Form.Get("color") == "blue" {
			return "blue"
		}
		return "green"
	}

	t := tag.New("BlueGreen", pf)
	return t
}

func TestBlueGreen(t *testing.T) {
	tg := makeBlueGreenTag()
	engine := blademaster.Default()
	engine.Use(tg)
	engine.GET("/bgget", func(ctx *blademaster.Context) {
		color, ok := tag.Value(ctx, "BlueGreen")
		if !ok {
			ctx.Abort()
			return
		}
		ctx.String(200, "color is: "+color)
	})

	go func() {
		engine.Run(":18080")
	}()
	defer func() {
		engine.Server().Shutdown(context.TODO())
	}()

	time.Sleep(1 * time.Second)
	client := new(http.Client)
	resp, err := client.Get("http://127.0.0.1:18080/bgget?color=blue")
	assert.Nil(t, err)
	defer resp.Body.Close()
	assert.Equal(t, resp.StatusCode, 200)

	body, err := ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)
	assert.Equal(t, string(body), "color is: blue")

	resp, err = client.Get("http://127.0.0.1:18080/bgget?color=green")
	assert.Nil(t, err)
	assert.Equal(t, resp.StatusCode, 200)
	body, err = ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)
	assert.Equal(t, string(body), "color is: green")
}
