package util

const (
	// DefaultPerPage .
	DefaultPerPage = 20
)

// SimplePage calculate "from", "to" without total_counts
// "from" index start from 1
func (p *Pagination) SimplePage() (from int64, to int64) {
	if p.CurPage == 0 || p.PerPage == 0 {
		p.CurPage, p.PerPage = 1, DefaultPerPage
	}
	from = (p.CurPage-1)*p.PerPage + 1
	to = from + p.PerPage - 1
	return
}

// Page calculate "from", "to" with total_counts
// index start from 1
func (p *Pagination) Page(total int64) (from int64, to int64) {
	if p.CurPage == 0 {
		p.CurPage = 1
	}
	if p.PerPage == 0 {
		p.PerPage = DefaultPerPage
	}

	if total == 0 || total < p.PerPage*(p.CurPage-1) {
		return
	}
	if total <= p.PerPage {
		return 1, total
	}
	from = (p.CurPage-1)*p.PerPage + 1
	if (total - from + 1) < p.PerPage {
		return from, total
	}
	return from, from + p.PerPage - 1
}

// VagueOffsetLimit calculate "offset", "limit" without total_counts
func (p *Pagination) VagueOffsetLimit() (offset int64, limit int64) {
	from, to := p.SimplePage()
	if to == 0 || from == 0 {
		return 0, 0
	}
	return from - 1, to - from + 1
}

// OffsetLimit calculate "offset" and "start" with total_counts
func (p *Pagination) OffsetLimit(total int64) (offset int64, limit int64) {
	from, to := p.Page(total)
	if to == 0 || from == 0 {
		return 0, 0
	}
	return from - 1, to - from + 1
}

// Pagination perform page algorithm
type Pagination struct {
	CurPage int64
	PerPage int64
}
