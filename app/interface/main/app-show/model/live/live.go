package live

import "encoding/json"

// Feed is live feed
type Feed struct {
	Count int     `json:"count"`
	Lives []*Room `json:"lives"`
}

// Recommend is live recommend
type Recommend struct {
	Count int `json:"count"`
	Lives struct {
		Subject []*Room `json:"subject"`
		Hot     []*Room `json:"hot"`
	} `json:"lives"`
}

type Room struct {
	Owner struct {
		Face string `json:"face"`
		Mid  int    `json:"mid"`
		Name string `json:"name"`
	} `json:"owner"`
	Cover struct {
		Src    string `json:"src"`
		Height int    `json:"height"`
		Width  int    `json:"width"`
	} `json:"cover"`
	Title  string `json:"title"`
	ID     int64  `json:"room_id"`
	Online int    `json:"online"`
	Area   string `json:"area"`
	AreaID int    `json:"area_id"`
}

type TopicHot struct {
	TID      int    `json:"topic_id"`
	TName    string `json:"topic_name"`
	Picture  string `json:"picture"`
	ImageURL string `json:"-"`
}

type TopicImage struct {
	ImageSrc    string `json:"image_src"`
	ImageWidth  int    `json:"image_width"`
	ImageHeight int    `json:"image_height"`
}

func (t *TopicHot) TopicJSONChange() (err error) {
	var tmp TopicImage
	if err = json.Unmarshal([]byte(t.Picture), &tmp); err != nil {
		return
	}
	t.ImageURL = tmp.ImageSrc
	return
}
