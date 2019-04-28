package v1

import "strconv"

func (r *RoomsReply) String() string {
	return strconv.FormatInt(int64(len(r.Rooms)), 10)
}
