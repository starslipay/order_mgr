package xerr

// 错误码  10000 0000 ~~99999 9999
// 模块id  70000
// 错误码 = 模块id + 业务错误码
var (
	ModuleId        = int64(455906)
	ModuleErrorBase = ModuleId * 100
)

var (
	// 系统错误 0000-0999
	ErrCodeDB             = ModuleErrorBase + 0
	ErrCodeServerInternal = ModuleErrorBase + 1

	// 业务错误码 1000-1999
	ErrCodeOrderNotFound             = ModuleErrorBase + 100
	ErrCodeOrderAlreadySuccess       = ModuleErrorBase + 101
	ErrCodeOrderAlreadyClosed        = ModuleErrorBase + 102
	ErrCodeOrderTradeStateInvalid    = ModuleErrorBase + 103
	ErrCodeOrderInsertOrderDuplicate = ModuleErrorBase + 104
	ErrCodeOrderInfoNotMatch         = ModuleErrorBase + 105
	ErrCodeCheckDeductToken          = ModuleErrorBase + 106
)
