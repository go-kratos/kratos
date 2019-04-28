package v1

import (
	"github.com/gogo/protobuf/proto"
)

// String GetCookieInfoReq string
func (m *GetCookieInfoReq) String() string {
	return ""
}

// String GetCookieInfoReply string
func (m *GetCookieInfoReply) String() string {
	return proto.MarshalTextString(m)
}

// String GetTokenInfoReq string
func (m *GetTokenInfoReq) String() string {
	return proto.MarshalTextString(m)
}

// String GetTokenInfoReply string
func (m *GetTokenInfoReply) String() string {
	return proto.MarshalTextString(m)
}
