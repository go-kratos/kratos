package model

// challenge audit state mapping to platform state map[business][platform_state][flow]chall_state
// business > 100 means feedback flow of business%100
// platform state 0.pending 1.handling 2.done 3.close
const (
	_AuditFlow    = 0
	_FeedbcakFlow = 1
	_PlatPending  = 0
	_PlatHandling = 1
	_PlatDone     = 2
	_PlatClosed   = 3
)

var pMetas map[int8]map[int8]map[int64][]int64

func init() {
	pMetas = map[int8]map[int8]map[int64][]int64{

		// 稿件申诉
		2: {
			_PlatPending: {
				_AuditFlow:    {},
				_FeedbcakFlow: {0, 1},
			},

			_PlatHandling: {
				_AuditFlow:    {},
				_FeedbcakFlow: {2, 6},
			},

			_PlatDone: {
				_AuditFlow:    {},
				_FeedbcakFlow: {3, 4, 5, 7, 8},
			},

			_PlatClosed: { // 关闭状态合并到已办
				_AuditFlow:    {},
				_FeedbcakFlow: {5, 7, 8},
			},
		},

		//稿件审核
		6: {
			_PlatPending: {
				_AuditFlow:    {0},
				_FeedbcakFlow: {},
			},
			_PlatDone: {
				_AuditFlow:    {1, 2, 3},
				_FeedbcakFlow: {},
			},
			_PlatClosed: {
				_AuditFlow:    {3},
				_FeedbcakFlow: {},
			},
		},
	}
}

// PlatformMetas return platform meta data
func PlatformMetas() map[int8]map[int8]map[int64][]int64 {
	return pMetas
}
