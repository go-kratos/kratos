package proto

import (
	"github.com/stretchr/testify/assert"
	"testing"

	testproto "github.com/go-kratos/kratos/v2/internal/testproto/encoding"
)

func TestName(t *testing.T) {
	c := new(codec)
	assert.Equal(t, c.Name(), "proto")
}

func TestCodec(t *testing.T) {
	c := new(codec)

	model := testproto.TestModel{
		Id:    1,
		Name:  "kratos",
		Hobby: []string{"study", "eat", "play"},
	}

	m, err := c.Marshal(&model)
	assert.Nil(t, err)

	var res testproto.TestModel

	err = c.Unmarshal(m, &res)
	assert.Nil(t, err)

	assert.Equal(t, res.Id, model.Id)
	assert.Equal(t, res.Name, model.Name)
	assert.Equal(t, res.Hobby, model.Hobby)
}
