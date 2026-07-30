package xerr

// 错误码  10000 0000 ~~99999 9999
// 模块id  70000
// 错误码 = 模块id + 业务错误码
var (
	ModuleId        = int64(70000)
	ModuleErrorBase = ModuleId * 10000
)

var (
	// 系统错误 0000-0999
	ErrCodeDB             = ModuleErrorBase + 0
	ErrCodeServerInternal = ModuleErrorBase + 1

	// 业务错误码 1000-1999
	ErrCodeOrderNotFound                        = ModuleErrorBase + 1000
	ErrCodeOrderAlreadySuccess                  = ModuleErrorBase + 1001
	ErrCodeOrderAlreadyClosed                   = ModuleErrorBase + 1002
	ErrCodeOrderTradeStateInvalid               = ModuleErrorBase + 1003
	ErrCodeOrderDuplicate                       = ModuleErrorBase + 1004
	ErrCodeInsertCloseOrderButOrderAlreadyExist = ModuleErrorBase + 1005
	ErrCodeOrderInfoNotMatch                    = ModuleErrorBase + 1006
)
