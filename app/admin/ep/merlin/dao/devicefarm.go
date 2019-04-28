package dao

import (
	"context"
	"fmt"
	"net/http"

	"go-common/app/admin/ep/merlin/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

const (
	_deviceCode        = 0
	_devicesListURL    = "/apis/devices/list"
	_deviceURL         = "/apis/devices/get"
	_deviceBootURL     = "/apis/devices/boot"
	_deviceShutDownURL = "/apis/devices/shutdown"
)

// MobileDeviceList Get Device Farm List .
func (d *Dao) MobileDeviceList(c context.Context) (resTotal map[string][]*model.Device, err error) {
	var (
		req      *http.Request
		hostList = d.c.DeviceFarm.HostList
	)

	resTotal = make(map[string][]*model.Device)

	for _, host := range hostList {

		var res *model.DeviceListResponse

		url := fmt.Sprintf("http://%s", host+_devicesListURL)
		if req, err = d.newRequest(http.MethodGet, url, nil); err != nil {
			log.Error("d.MobileDeviceList newRequest url(%s) err(%+v)", _devicesListURL, err)
			continue
		}

		if err = d.httpClient.Do(c, req, &res); err != nil {
			log.Error("d.MobileDeviceList httpClient url(%s) err(%+v)", _devicesListURL, err)
			continue
		}
		if res.Status.Code != _deviceCode {
			err = ecode.MerlinDeviceFarmErr
			log.Error("Status url(%s) res(%s),err(%+v)", _devicesListURL, res, err)
			continue
		}

		resTotal[host] = res.Data.Devices
	}
	return
}

// MobileDeviceDetail Get Mobile Device Detail.
func (d *Dao) MobileDeviceDetail(c context.Context, host, serial string) (device *model.Device, err error) {
	var (
		req *http.Request
		res *model.DeviceListDetailResponse
	)

	url := fmt.Sprintf("http://%s?serial=%s", host+_deviceURL, serial)
	if req, err = d.newRequest(http.MethodGet, url, nil); err != nil {
		log.Error("d.MobileDeviceDetail newRequest url(%s) err(%+v)", _deviceURL, err)
		return
	}

	if err = d.httpClient.Do(c, req, &res); err != nil {
		log.Error("d.MobileDeviceDetail url(%s) err(%+v)", _deviceURL, err)
		return
	}
	if res.Status.Code != _deviceCode {
		err = ecode.MerlinDeviceFarmErr
		log.Error("Status url(%s) res(%s) err(%+v)", _deviceURL, res, err)
		return
	}
	device = res.Data.Devices

	return
}

// BootMobileDevice Boot Mobile Device.
func (d *Dao) BootMobileDevice(c context.Context, host, serial string) (deviceBootData *model.DeviceBootData, err error) {
	var (
		req *http.Request
		res *model.DeviceBootResponse
	)

	reqModel := &model.DeviceBootRequest{
		Serial: serial,
	}

	url := fmt.Sprintf("http://%s", host+_deviceBootURL)
	if req, err = d.newRequest(http.MethodPost, url, reqModel); err != nil {
		log.Error("d.BootMobileDevice newRequest url(%s) err(%+v)", _deviceBootURL, err)
		return
	}

	req.Header.Set("content-type", "application/json")

	if err = d.httpClient.Do(c, req, &res); err != nil {
		log.Error("d.BootMobileDevice url(%s) err(%+v)", _deviceBootURL, err)
		return
	}
	if res.Status.Code != _deviceCode {
		err = ecode.MerlinDeviceFarmErr
		log.Error("Status url(%s) res(%s) err(%+v)", _deviceBootURL, res, err)
		return
	}
	deviceBootData = res.Data
	return
}

// ShutdownMobileDevice Shutdown Mobile Device.
func (d *Dao) ShutdownMobileDevice(c context.Context, host, serial string) (err error) {
	var (
		req *http.Request
		res *model.DeviceShutDownResponse
	)

	reqModel := &model.DeviceBootRequest{
		Serial: serial,
	}

	url := fmt.Sprintf("http://%s", host+_deviceShutDownURL)

	if req, err = d.newRequest(http.MethodPost, url, reqModel); err != nil {
		log.Error("d.ShutdownMobileDevice newRequest url(%s) err(%+v)", _deviceShutDownURL, err)
		return
	}
	req.Header.Set("content-type", "application/json")

	if err = d.httpClient.Do(c, req, &res); err != nil {
		log.Error("d.ShutdownMobileDevice url(%s) err(%+v)", _deviceShutDownURL, err)
		return
	}
	if res.Status.Code != _deviceCode {
		err = ecode.MerlinDeviceFarmErr
		log.Error("Status url(%s) res(%s) err(%+v)", _deviceShutDownURL, res, err)
		return
	}
	return
}
