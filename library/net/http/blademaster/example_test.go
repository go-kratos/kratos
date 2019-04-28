package blademaster_test

import (
	"io/ioutil"
	"log"
	"time"

	"go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/binding"
	"go-common/library/net/http/blademaster/middleware/auth"
	"go-common/library/net/http/blademaster/middleware/verify"
	"go-common/library/net/http/blademaster/tests"
)

// This example start a http server and listen at port 8080,
// it will handle '/ping' and return response in html text
func Example() {
	engine := blademaster.Default()
	engine.GET("/ping", func(c *blademaster.Context) {
		c.String(200, "%s", "pong")
	})
	engine.Run(":8080")
}

// This example use `RouterGroup` to separate different requests,
// it will handle ('/group1/ping', '/group2/ping') and return response in json
func ExampleRouterGroup() {
	engine := blademaster.Default()

	group := engine.Group("/group1", blademaster.CORS())
	group.GET("/ping", func(c *blademaster.Context) {
		c.JSON(map[string]string{"message": "hello"}, nil)
	})

	group2 := engine.Group("/group2", blademaster.CORS())
	group2.GET("/ping", func(c *blademaster.Context) {
		c.JSON(map[string]string{"message": "welcome"}, nil)
	})

	engine.Run(":8080")
}

// This example add two middlewares in the root router by `Use` method,
// it will add CORS headers in response and log total consumed time
func ExampleEngine_Use() {
	timeLogger := func() blademaster.HandlerFunc {
		return func(c *blademaster.Context) {
			start := time.Now()
			c.Next()
			log.Printf("total consume: %v", time.Since(start))
		}
	}

	engine := blademaster.Default()

	engine.Use(blademaster.CORS())
	engine.Use(timeLogger())

	engine.GET("/ping", func(c *blademaster.Context) {
		c.String(200, "%s", "pong")
	})
	engine.Run(":8080")
}

// This example add two middlewares in the root router by `UseFunc` method,
// it will log total consumed time
func ExampleEngine_UseFunc() {
	engine := blademaster.Default()

	engine.UseFunc(func(c *blademaster.Context) {
		start := time.Now()
		c.Next()
		log.Printf("total consume: %v", time.Since(start))
	})

	engine.GET("/ping", func(c *blademaster.Context) {
		c.String(200, "%s", "pong")
	})
	engine.Run(":8080")
}

// This example start a http server through the specified unix socket,
// it will handle '/ping' and return reponse in html text
func ExampleEngine_RunUnix() {
	engine := blademaster.Default()
	engine.GET("/ping", func(c *blademaster.Context) {
		c.String(200, "%s", "pong")
	})

	unixs, err := ioutil.TempFile("", "engine.sock")
	if err != nil {
		log.Fatalf("Failed to create temp file: %s", err)
	}

	if err := engine.RunUnix(unixs.Name()); err != nil {
		log.Fatalf("Failed to serve with unix socket: %s", err)
	}
}

// This example show how to render response in json format,
// it will render structures as json like: `{"code":0,"message":"0","data":{"Time":"2017-11-14T23:03:22.0523199+08:00"}}`
func ExampleContext_JSON() {
	type Data struct {
		Time time.Time
	}

	engine := blademaster.Default()

	engine.GET("/ping", func(c *blademaster.Context) {
		var d Data
		d.Time = time.Now()
		c.JSON(d, nil)
	})

	engine.Run(":8080")
}

// This example show how to render response in protobuf format
// it will marshal whole response content to protobuf
func ExampleContext_Protobuf() {
	engine := blademaster.Default()
	engine.GET("/ping.pb", func(c *blademaster.Context) {
		t := &tests.Time{
			Now: time.Now().Unix(),
		}
		c.Protobuf(t, nil)
	})

	engine.Run(":8080")
}

// This example show how to render response in XML format,
// it will render structure as XML like: `<Data><Time>2017-11-14T23:03:49.2231458+08:00</Time></Data>`
func ExampleContext_XML() {
	type Data struct {
		Time time.Time
	}

	engine := blademaster.Default()

	engine.GET("/ping", func(c *blademaster.Context) {
		var d Data
		d.Time = time.Now()
		c.XML(d, nil)
	})

	engine.Run(":8080")
}

// This example show how to protect your handlers by HTTP basic auth,
// it will validate the baisc auth and abort with status 403 if authentication is invalid
func ExampleContext_Abort() {
	engine := blademaster.Default()
	engine.UseFunc(func(c *blademaster.Context) {
		user, pass, isok := c.Request.BasicAuth()
		if !isok || user != "root" || pass != "root" {
			c.AbortWithStatus(403)
			return
		}
	})

	engine.GET("/auth", func(c *blademaster.Context) {
		c.String(200, "%s", "Welcome")
	})

	engine.Run(":8080")
}

// This example show how to using the default parameter binding to parse the url param from get request,
// it will validate the request and abort with status 400 if params is invalid
func ExampleContext_Bind() {
	engine := blademaster.Default()
	engine.GET("/bind", func(c *blademaster.Context) {
		v := new(struct {
			// This mark field `mids` should exist and every element should greater than 1
			Mids    []int64 `form:"mids" validate:"dive,gt=1,required"`
			Title   string  `form:"title" validate:"required"`
			Content string  `form:"content"`
			// This mark field `cid` should between 1 and 10
			Cid int `form:"cid" validate:"min=1,max=10"`
		})

		err := c.Bind(v)
		if err != nil {
			// Do not call any write response method in this state,
			// the response body is already written in `c.BindWith` method
			return
		}
		c.String(200, "parse params by bind %+v", v)
	})

	engine.Run(":8080")
}

// This example show how to using the json binding to parse the json param from post request body,
// it will validate the request and abort with status 400 if params is invalid
func ExampleContext_BindWith() {
	engine := blademaster.Default()
	engine.POST("/bindwith", func(c *blademaster.Context) {
		v := new(struct {
			// This mark field `mids` should exist and every element should greater than 1
			Mids    []int64 `json:"mids" validate:"dive,gt=1,required"`
			Title   string  `json:"title" validate:"required"`
			Content string  `json:"content"`
			// This mark field `cid` should between 1 and 10
			Cid int `json:"cid" validate:"min=1,max=10"`
		})

		err := c.BindWith(v, binding.JSON)
		if err != nil {
			// Do not call any write response method in this state,
			// the response body is already written in `c.BindWith` method
			return
		}
		c.String(200, "parse params by bindwith %+v", v)
	})

	engine.Run(":8080")
}

func ExampleEngine_Inject() {
	v := verify.New(nil)
	auth := auth.New(nil)
	engine := blademaster.Default()
	engine.Inject("^/index", v.Verify, auth.User)
	engine.POST("/index/hello", func(c *blademaster.Context) {
		c.JSON("hello, world", nil)
	})
	engine.Run(":8080")
}
