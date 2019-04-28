package model

// Stat struct of Stat.
// type Stat struct {
// 	Mid       int64     `json:"-"`
// 	Following int64     `json:"following"`
// 	Whisper   int64     `json:"whisper"`
// 	Black     int64     `json:"black"`
// 	Follower  int64     `json:"follower"`
// 	CTime     time.Time `json:"-"`
// 	MTime     time.Time `json:"-"`
// }

// Count get count of following, including attr following, whisper.
func (st *Stat) Count() int {
	return int(st.Following + st.Whisper)
}

// BlackCount get count of black, including attr black.
func (st *Stat) BlackCount() int {
	return int(st.Black)
}

// Empty get if the stat is empty.
func (st *Stat) Empty() bool {
	return st.Following == 0 && st.Whisper == 0 && st.Black == 0 && st.Follower == 0
}

// Fill fill by the incoming stat with its non-negative fields.
func (st *Stat) Fill(ost *Stat) {
	if ost.Following >= 0 {
		st.Following = ost.Following
	}
	if ost.Whisper >= 0 {
		st.Whisper = ost.Whisper
	}
	if ost.Black >= 0 {
		st.Black = ost.Black
	}
	if ost.Follower >= 0 {
		st.Follower = ost.Follower
	}
}
