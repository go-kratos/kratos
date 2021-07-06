package form

import (
	"testing"

	"github.com/go-kratos/kratos/v2/encoding"
	"github.com/go-kratos/kratos/v2/internal/testproto/complex"
	"github.com/stretchr/testify/require"
)

type LoginRequest struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

const contentType = "x-www-form-urlencoded"

func TestFormCodecMarshal(t *testing.T) {
	req := &LoginRequest{
		Username: "kratos",
		Password: "kratos_pwd",
	}
	content, err := encoding.GetCodec(contentType).Marshal(req)
	require.NoError(t, err)
	require.Equal(t, []byte("password=kratos_pwd&username=kratos"), content)

	req = &LoginRequest{
		Username: "kratos",
		Password: "",
	}
	content, err = encoding.GetCodec(contentType).Marshal(req)
	require.NoError(t, err)
	require.Equal(t, []byte("username=kratos"), content)
}

func TestFormCodecUnmarshal(t *testing.T) {
	req := &LoginRequest{
		Username: "kratos",
		Password: "kratos_pwd",
	}
	content, err := encoding.GetCodec(contentType).Marshal(req)
	require.NoError(t, err)

	var bindReq = new(LoginRequest)
	err = encoding.GetCodec(contentType).Unmarshal(content, bindReq)
	require.NoError(t, err)
	require.Equal(t, "kratos", bindReq.Username)
	require.Equal(t, "kratos_pwd", bindReq.Password)
}

func TestProtoEncodeDecode(t *testing.T) {
	in := &complex.Complex{
		Id:      2233,
		NoOne:   "2233",
		Simple:  &complex.Simple{Component: "5566"},
		Simples: []string{"3344", "5566"},
	}
	content, err := encoding.GetCodec(contentType).Marshal(in)
	require.NoError(t, err)
	require.Equal(t, "id=2233&numberOne=2233&simples=3344&simples=5566&very_simple.component=5566", string(content))
	var in2 = &complex.Complex{}
	err = encoding.GetCodec(contentType).Unmarshal(content, in2)
	require.NoError(t, err)
	require.Equal(t, int64(2233), in2.Id)
	require.Equal(t, "2233", in2.NoOne)
	require.NotEmpty(t, in2.Simple)
	require.Equal(t, "5566", in2.Simple.Component)
	require.NotEmpty(t, in2.Simples)
	require.Len(t, in2.Simples, 2)
	require.Equal(t, "3344", in2.Simples[0])
	require.Equal(t, "5566", in2.Simples[1])
}
