package api

const (
	_normalVip = 1 //月度大会员
	_annualVip = 2 //年度会员

	_statusAvailable = 1 //未过期
	_statusFrozen    = 2 //冻结
)

// IsValid decide the user is valid vip or not.
func (v *VipInfo) IsValid() bool {
	return v.Status == _statusAvailable && (v.Type == _normalVip || v.Type == _annualVip)
}

// IsFrozen decide the user is frozen vip or not.
func (v *VipInfo) IsFrozen() bool {
	return v.Status == _statusFrozen
}
