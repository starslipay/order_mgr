package xerr

// 错误码  10000 0000 ~~99999 9999
// 模块id  70000
// 错误码 = 模块id + 业务错误码
var (
	ModuleId = int64(455906)
)

var (
	// 系统错误 0000-0999
	ErrCodeDB             = int64(455906000)
	ErrCodeServerInternal = int64(455906001)

	// 业务错误码 1000-1999
	ErrCodeOrderNotFound             = int64(455906100)
	ErrCodeOrderAlreadySuccess       = int64(455906101)
	ErrCodeOrderAlreadyClosed        = int64(455906102)
	ErrCodeOrderTradeStateInvalid    = int64(455906103)
	ErrCodeOrderInsertOrderDuplicate = int64(455906104)
	ErrCodeOrderInfoNotMatch         = int64(455906105)
	ErrCodeCheckDeductToken          = int64(455906106)
)
