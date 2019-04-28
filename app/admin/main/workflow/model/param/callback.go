package param

// AddCallbackParam is the model to create callback
type AddCallbackParam struct {
	Business    int8   `json:"business"`
	URL         string `json:"url"`
	IsSobot     bool   `json:"is_sobot"`
	State       int8   `json:"state"`
	ExternalAPI string `json:"external_api"`
	SourceAPI   string `json:"source_api"`
}
