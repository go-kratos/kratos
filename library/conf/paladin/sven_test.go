package paladin

import (
	"context"
	"testing"
	"time"

	"go-common/library/conf/env"

	"github.com/naoina/toml"
	"github.com/stretchr/testify/assert"
)

type testObj struct {
	Bool   bool
	Int    int64
	Float  float64
	String string
}

func (t *testObj) Set(text string) error {
	return toml.Unmarshal([]byte(text), t)
}

type testConf struct {
	Bool   bool
	Int    int64
	Float  float64
	String string
	Object *testObj
}

func (t *testConf) Set(text string) error {
	return toml.Unmarshal([]byte(text), t)
}

func TestSven(t *testing.T) {
	svenHost = "config.bilibili.co"
	svenVersion = "server-1"
	svenPath = "/tmp"
	svenToken = "1afe5efaf45e11e7b3f8c6cd4f230d8c"
	svenAppoint = ""
	svenTreeid = "2888"
	env.Region = "sh"
	env.Zone = "sh001"
	env.Hostname = "test"
	env.DeployEnv = "dev"
	env.AppID = "main.common-arch.msm-service"

	sven, err := NewSven()
	assert.Nil(t, err)
	testSvenMap(t, sven)
	testSvenValue(t, sven)
	testWatch(t, sven)
}

func testSvenMap(t *testing.T, cli Client) {
	m := Map{}
	text, err := cli.Get("test.toml").String()
	assert.Nil(t, err)
	assert.Nil(t, m.Set(text), text)
	b, err := m.Get("bool").Bool()
	assert.Nil(t, err)
	assert.Equal(t, b, true, "bool")
	// int64
	i, err := m.Get("int").Int64()
	assert.Nil(t, err)
	assert.Equal(t, i, int64(100), "int64")
	// float64
	f, err := m.Get("float").Float64()
	assert.Nil(t, err)
	assert.Equal(t, f, 100.1, "float64")
	// string
	s, err := m.Get("string").String()
	assert.Nil(t, err)
	assert.Equal(t, s, "text", "string")
	// error
	n, err := m.Get("not_exsit").String()
	assert.NotNil(t, err)
	assert.Equal(t, n, "", "not_exsit")

	obj := new(testObj)
	text, err = m.Get("object").Raw()
	assert.Nil(t, err)
	assert.Nil(t, obj.Set(text))
	assert.Equal(t, obj.Bool, true, "bool")
	assert.Equal(t, obj.Int, int64(100), "int64")
	assert.Equal(t, obj.Float, 100.1, "float64")
	assert.Equal(t, obj.String, "text", "string")
}

func testSvenValue(t *testing.T, cli Client) {
	v := new(testConf)
	text, err := cli.Get("test.toml").Raw()
	assert.Nil(t, err)
	assert.Nil(t, v.Set(text))
	assert.Equal(t, v.Bool, true, "bool")
	assert.Equal(t, v.Int, int64(100), "int64")
	assert.Equal(t, v.Float, 100.1, "float64")
	assert.Equal(t, v.String, "text", "string")
	assert.Equal(t, v.Object.Bool, true, "bool")
	assert.Equal(t, v.Object.Int, int64(100), "int64")
	assert.Equal(t, v.Object.Float, 100.1, "float64")
	assert.Equal(t, v.Object.String, "text", "string")
}

func testWatch(t *testing.T, cli Client) {
	ch := cli.WatchEvent(context.Background())
	select {
	case <-time.After(time.Second):
		t.Log("watch timeout")
	case e := <-ch:
		s, err := cli.Get("static").String()
		assert.Nil(t, err)
		assert.Equal(t, s, e.Value, "watch value")

		t.Logf("watch event:%+v", e)
	}
}
