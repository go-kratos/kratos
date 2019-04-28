package service

const (
	_range0_30      float64 = 0.0
	_range30_50     float64 = 0.05
	_range50_100    float64 = 0.10
	_range100_150   float64 = 0.15
	_range150_300   float64 = 0.20
	_range300_500   float64 = 0.25
	_range500_1000  float64 = 0.30
	_range1000_2000 float64 = 0.40
	_range2000_3000 float64 = 0.50
	// _range3000_end  float64 = 0.60
)

// TaxRate tax rate
type TaxRate struct {
	start float64
	end   float64
	rate  float64
}

// Tax tax
func Tax(income float64) (tax float64) {
	rs := rates(income)
	for _, r := range rs {
		if income >= r.end {
			tax += Mul(float64(r.end-r.start), r.rate)
		} else {
			tax += Mul(float64(income-r.start), r.rate)
		}
	}
	return
}

func rates(income float64) (rs []*TaxRate) {
	if income > 0 {
		r := &TaxRate{
			start: 0,
			end:   30,
			rate:  _range0_30,
		}
		rs = append(rs, r)
	}
	if income > 30 {
		r := &TaxRate{
			start: 30,
			end:   50,
			rate:  _range30_50,
		}
		rs = append(rs, r)
	}
	if income > 50 {
		r := &TaxRate{
			start: 50,
			end:   100,
			rate:  _range50_100,
		}
		rs = append(rs, r)
	}
	if income > 100 {
		r := &TaxRate{
			start: 100,
			end:   150,
			rate:  _range100_150,
		}
		rs = append(rs, r)
	}
	if income > 150 {
		r := &TaxRate{
			start: 150,
			end:   300,
			rate:  _range150_300,
		}
		rs = append(rs, r)
	}
	if income > 300 {
		r := &TaxRate{
			start: 300,
			end:   500,
			rate:  _range300_500,
		}
		rs = append(rs, r)
	}
	if income > 500 {
		r := &TaxRate{
			start: 500,
			end:   1000,
			rate:  _range500_1000,
		}
		rs = append(rs, r)
	}
	if income > 1000 {
		r := &TaxRate{
			start: 1000,
			end:   2000,
			rate:  _range1000_2000,
		}
		rs = append(rs, r)
	}
	if income > 2000 {
		r := &TaxRate{
			start: 2000,
			end:   3000,
			rate:  _range2000_3000,
		}
		rs = append(rs, r)
	}
	if income > 3000 {
		r := &TaxRate{
			start: 3000,
			end:   1<<31 - 1,
			rate:  _range2000_3000,
		}
		rs = append(rs, r)
	}
	return
}

// FIXME
// func mul(x, y float64) (z float64) {
// 	f := new(big.Float)
// 	xb := big.NewFloat(x)
// 	yb := big.NewFloat(y)
// 	z, _ = f.Mul(xb, yb).Float64()
// 	return
// }
