package form

import (
	"github.com/go-kratos/kratos/v2/encoding"
	"github.com/stretchr/testify/require"
	"testing"
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
