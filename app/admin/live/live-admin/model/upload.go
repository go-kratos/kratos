package model

import "go-common/library/database/bfs"

const (
	_defaultPaddingX = 10
	_defaultPaddingY = 10
	_defaultScale    = 0.035
)

// TweakWatermark makes some attributes of watermark default if they are not legal.
func TweakWatermark(req *bfs.Request) {
	if req.WMKey != "" || req.WMText != "" {
		if req.WMPaddingX == 0 {
			req.WMPaddingX = _defaultPaddingX
		}

		if req.WMPaddingY == 0 {
			req.WMPaddingY = _defaultPaddingY
		}

		if req.WMScale <= 0 || req.WMScale >= 1 {
			req.WMScale = _defaultScale
		}
	}
}
