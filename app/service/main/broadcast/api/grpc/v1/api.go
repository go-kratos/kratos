package v1

import (
	"fmt"
)

func (m *ConnectReq) String() string {
	return fmt.Sprintf("server:%s serverKey:%s token:%s", m.Server, m.ServerKey, m.Token)
}

func (m *OnlineReq) String() string {
	return fmt.Sprintf("server:%s", m.Server)
}

func (m *OnlineReply) String() string {
	return fmt.Sprintf("rooms:%d", len(m.RoomCount))
}
