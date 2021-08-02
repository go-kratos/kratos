package kratos

import (
	"testing"
	"time"

	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/stretchr/testify/assert"
)

func TestApp(t *testing.T) {
	hs := http.NewServer()
	gs := grpc.NewServer()
	app := New(
		Name("kratos"),
		Version("v1.0.0"),
		Server(hs, gs),
	)
	time.AfterFunc(time.Second, func() {
		app.Stop()
	})
	if err := app.Run(); err != nil {
		t.Fatal(err)
	}
}

func TestApp_ID(t *testing.T) {
	v := "123"
	o := New(ID(v))
	assert.Equal(t, v, o.ID())
}

func TestApp_Name(t *testing.T) {
	v := "123"
	o := New(Name(v))
	assert.Equal(t, v, o.Name())
}

func TestApp_Version(t *testing.T) {
	v := "123"
	o := New(Version(v))
	assert.Equal(t, v, o.Version())
}

func TestApp_Metadata(t *testing.T) {
	v := map[string]string{
		"a": "1",
		"b": "2",
	}
	o := New(Metadata(v))
	assert.Equal(t, v, o.Metadata())
}
