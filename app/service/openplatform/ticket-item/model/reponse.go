package model

// Response 标准返回结构
type Response struct {
	Errno int64   `json:"errno"`
	Msg   string  `json:"msg"`
	Data  []int64 `json:"data"`
}

// ListResp pagination
type ListResp struct {
	Total int         `json:"total"`
	List  interface{} `json:"list"`
}

// Page pagination
type Page struct {
	Num   int `json:"num"`
	Size  int `json:"size"`
	Total int `json:"total"`
}

// VenueData 场馆
type VenueData struct {
	ID            int64  `json:"id"`
	Name          string `json:"name"`
	CityID        int    `json:"city"`
	CityName      string `json:"city_name"`
	ProvinceID    int    `json:"province"`
	ProvinceName  string `json:"province_name"`
	DistrictID    int    `json:"district"`
	DistrictName  string `json:"district_name"`
	AddressDetail string `json:"address_detail"`
	PlaceNum      int    `json:"place_num"`
	Status        int    `json:"status"`
	Coordinate    string `json:"coordinate"`
	Traffic       string `json:"traffic"`
	Ctime         string `json:"ctime"`
	Mtime         string `json:"mtime"`
}

// VersionSearchList .
type VersionSearchList struct {
	Result []*Version `json:"result"`
	Page   *Page      `json:"page"`
}

// GuestSearchList .
type GuestSearchList struct {
	Result []*Guest `json:"result"`
	Page   *Page    `json:"page"`
}

// VenueSearchList .
type VenueSearchList struct {
	Result []*VenueData `json:"result"`
	Page   *Page        `json:"page"`
}
