package wall

type Wall struct {
	Id       int    `json:"-"`
	Name     string `json:"name"`
	Package  string `json:"package"`
	Size     string `json:"size"`
	Logo     string `json:"logo"`
	Download string `json:"download"`
	Remark   string `json:"remark"`
}
