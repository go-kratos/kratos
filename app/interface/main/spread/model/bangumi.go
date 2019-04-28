package model

import "encoding/json"

// Bangumi struct
type Bangumi struct {
	Actors          string `json:"actors,omitempty"`
	Akira           string `json:"akira,omitempty"`
	Alias           string `json:"alias,omitempty"`
	Country         string `json:"country,omitempty"`
	CoverImage      string `json:"cover_image,omitempty"`
	DisplayAddress  string `json:"display_address,omitempty"`
	DownloadAddress string `json:"download_address,omitempty"`
	Duration        int64  `json:"duration"`
	Episodes        []struct {
		Index   int64  `json:"index"`
		PlayURL string `json:"play_url"`
	} `json:"episodes"`
	Intro        string `json:"intro,omitempty"`
	Name         string `json:"name,omitempty"`
	PlayCount    int64  `json:"play_count"`
	Premieredate string `json:"premieredate,omitempty"`
	Season       struct {
		ID            int64  `json:"id"`
		Index         int64  `json:"index"`
		Paymentstatus int64  `json:"paymentstatus"`
		Title         string `json:"title"`
	} `json:"season"`
	Staff struct {
		AnimationProduction    string `json:"animation_production,omitempty"`
		AnimationScript        string `json:"animation_script,omitempty"`
		ArtSupervisor          string `json:"art_supervisor,omitempty"`
		CharacterDisign        string `json:"character_disign,omitempty"`
		ChiefContributor       string `json:"chief_contributor,omitempty"`
		ChiefDirector          string `json:"chief_director,omitempty"`
		ChiefExecutiveProducer string `json:"chief_executive_producer,omitempty"`
		ChiefProducer          string `json:"chief_producer,omitempty"`
		ChiefProductionManager string `json:"chief_production_manager,omitempty"`
		ChiefScenarist         string `json:"chief_scenarist,omitempty"`
		ColorDesign            string `json:"color_design,omitempty"`
		Director               string `json:"director,omitempty"`
		DocumentaryProduction  string `json:"documentary_production,omitempty"`
		ExecutiveProducer      string `json:"executive_producer,omitempty"`
		JointProduction        string `json:"joint_production,omitempty"`
		Music                  string `json:"music,omitempty"`
		OriginalAutor          string `json:"original_autor,omitempty"`
		OriginalCharacter      string `json:"original_character,omitempty"`
		PaintingSupervisor     string `json:"painting_supervisor,omitempty"`
		Performer              string `json:"performer,omitempty"`
		Produce                string `json:"produce,omitempty"`
		Producer               string `json:"producer,omitempty"`
		Production             string `json:"production,omitempty"`
		ProductionManager      string `json:"production_manager,omitempty"`
		Publisher              string `json:"publisher,omitempty"`
		Screenwriter           string `json:"screenwriter,omitempty"`
		SeriesComposition      string `json:"series_composition,omitempty"`
		Star                   string `json:"star,omitempty"`
		Storyboard             string `json:"storyboard,omitempty"`
		Supervisor             string `json:"supervisor,omitempty"`
	} `json:"staff"`
	Type int64 `json:"type"`
}

// BangumiResp .
type BangumiResp struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Result  json.RawMessage `json:"result"`
	Total   int64           `json:"total"`
}

// BangumiOffResp .
type BangumiOffResp struct {
	Code    int           `json:"code"`
	Message string        `json:"message"`
	Total   int64         `json:"total"`
	Result  []*BangumiOff `json:"result"`
}

// BangumiOff .
type BangumiOff struct {
	Name     string `json:"name"`
	Seasonid int64  `json:"seasonid"`
	Type     int64  `json:"type"`
}
