package consts

const (
	OrderTradeStateInit    = 1  // 初始状态
	OrderTradeStateSuccess = 2  // 成功
	OrderTradeStateRefund  = 3  // 已退款
	OrderTradeStateClose   = 99 // 关闭
)

const (
	OrderPayTypeBanPay = 1 // 余额支付
)

const (
	OrderCurTypeCNY = 1 // 人民币
)

const (
	RechargeWithdrawOrderTypeRecharge = 1 // 充值(Bank2C)
	RechargeWithdrawOrderTypeWithdraw = 2 // 提现(C2Bank)
)

const (
	RechargeWithdrawTradeStateInit    = 1 // 初始状态
	RechargeWithdrawTradeStateSuccess = 2 // 成功
	RechargeWithdrawTradeStateFail    = 3 // 失败
)
