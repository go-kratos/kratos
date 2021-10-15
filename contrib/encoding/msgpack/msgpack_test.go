package msgpack

import (
	"testing"

	testData "github.com/go-kratos/kratos/v2/internal/testdata/encoding"
	"github.com/stretchr/testify/assert"
)

type loginRequest struct {
	UserName string
	Password string
}

type testModel struct {
	ID   int32
	Name string
}

func TestName(t *testing.T) {
	c := new(codec)
	assert.Equal(t, c.Name(), "msgpack")
}

func TestCodec(t *testing.T) {
	c := new(codec)

	model := testData.TestModel{
		Id:    1,
		Name:  "kratos",
		Hobby: []string{"study", "eat", "play"},
	}

	m, err := c.Marshal(&model)
	assert.Nil(t, err)

	var res testData.TestModel

	err = c.Unmarshal(m, &res)
	assert.Nil(t, err)

	assert.Equal(t, res.Id, model.Id)
	assert.Equal(t, res.Name, model.Name)
	assert.Equal(t, res.Hobby, model.Hobby)

	t2 := testModel{ID: 1, Name: "name"}
	m, err = c.Marshal(&t2)
	assert.Nil(t, err)
	var t3 testModel
	err = c.Unmarshal(m, &t3)
	assert.Nil(t, err)
	assert.Equal(t, t3.ID, t2.ID)
	assert.Equal(t, t3.Name, t2.Name)

	request := loginRequest{
		UserName: "username",
		Password: "password",
	}
	m, err = c.Marshal(&request)
	assert.Nil(t, err)
	var req loginRequest
	err = c.Unmarshal(m, &req)
	assert.Nil(t, err)
	assert.Equal(t, req.Password, request.Password)
	assert.Equal(t, req.UserName, request.UserName)

}
