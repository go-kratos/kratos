package game

type Download struct {
	ID          int64
	Title       string
	Desc        string
	Icon        string
	Cover       string
	URLType     int
	URLValue    string
	BtnTxt      int
	ReType      int
	ReValue     string
	ButtonText  string
	DoubleCover string
	Number      int
}

func (d *Download) CardChange() {
	switch d.BtnTxt {
	case 0:
		d.ButtonText = "下载"
	case 1:
		d.ButtonText = "预约"
	case 2:
		d.ButtonText = "查看详情"
	}
}
