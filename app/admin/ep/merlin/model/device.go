package model

// DeviceListResponse Device List Response.
type DeviceListResponse struct {
	Status *DeviceResponseStatus `json:"status"`
	Data   *DeviceListData       `json:"data"`
}

// DeviceListDetailResponse Device List Detail Response.
type DeviceListDetailResponse struct {
	Status *DeviceResponseStatus `json:"status"`
	Data   *DeviceListDetailData `json:"data"`
}

// DeviceBootResponse Device Boot Response.
type DeviceBootResponse struct {
	Status *DeviceResponseStatus `json:"status"`
	Data   *DeviceBootData       `json:"data"`
}

// DeviceShutDownResponse Device Shut Down Response.
type DeviceShutDownResponse struct {
	Status *DeviceResponseStatus `json:"status"`
}

// DeviceBootData Device Boot Data.
type DeviceBootData struct {
	WSRUL     string `json:"wsurl"`
	UploadURL string `json:"uploadURL"`
}

// DeviceListDetailData Device List Detail Data.
type DeviceListDetailData struct {
	Devices *Device `json:"device"`
}

// DeviceListData Device List Data.
type DeviceListData struct {
	Devices []*Device `json:"devices"`
}

// DeviceResponseStatus Device Response Status.
type DeviceResponseStatus struct {
	Code   int    `json:"code"`
	Status string `json:"status"`
}

// Device Device.
type Device struct {
	Serial      string `json:"serial"`
	Name        string `json:"name"`
	State       string `json:"state"`
	Mode        string `json:"mode"`
	Version     string `json:"version"`
	CPU         string `json:"cpu"`
	IsSimulator bool   `json:"isSimulator"`
}

// DeviceBootRequest Device Boot Request.
type DeviceBootRequest struct {
	Serial string `json:"serial"`
}

// MobileDeviceCategoryResponse MobileDeviceCategory Response.
type MobileDeviceCategoryResponse struct {
	Name   string        `json:"name"`
	Label  string        `json:"label"`
	Values []interface{} `json:"values"`
}

// MobileMachineResponse MobileMachine Response.
type MobileMachineResponse struct {
	*MobileMachine
	ImageSrc string `json:"image_src"`
}
