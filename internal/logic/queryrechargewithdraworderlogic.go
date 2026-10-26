package logic

import (
	"context"

	"github.com/starslipay/order_mgr/internal/svc"
	"github.com/starslipay/order_mgr/internal/util"
	"github.com/starslipay/order_mgr/internal/xerr"
	"github.com/starslipay/order_mgr/model/mysql"
	"github.com/starslipay/order_mgr/order_mgr_pb"
	"github.com/starslipay/paycomm/xerror"
	"google.golang.org/grpc/codes"

	"github.com/zeromicro/go-zero/core/logx"
)

type QueryRechargeWithdrawOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewQueryRechargeWithdrawOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryRechargeWithdrawOrderLogic {
	return &QueryRechargeWithdrawOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *QueryRechargeWithdrawOrderLogic) QueryRechargeWithdrawOrder(in *order_mgr_pb.QueryRechargeWithdrawOrderReq) (*order_mgr_pb.QueryRechargeWithdrawOrderRsp, error) {
	order, err := l.svcCtx.TRechargeWithdrawOrderModelSlave.FindOne(l.ctx, in.TransactionId)
	if err != nil {
		if err == mysql.ErrNotFound {
			return nil, xerror.NewBizError(codes.Internal, xerr.ErrCodeRechargeWithdrawOrderNotFound, "recharge withdraw order not found")
		}
		return nil, xerror.NewBizError(codes.Internal, xerr.ErrCodeDB, err.Error())
	}

	orderInfo := &order_mgr_pb.RechargeWithdrawOrderInfo{
		TransactionId:        order.TransactionId,
		ChannelTransactionId: order.ChannelTransactionId,
		UserId:               order.UserId,
		Uid:                  order.Uid,
		BizType:              int32(order.BizType),
		BankType:             int32(order.BankType),
		Amount:               order.Amount,
		CurType:              int32(order.CurType),
		TradeState:           int32(order.TradeState),
		Memo:                 order.Memo,
		CreateTime:           util.FormatBeijingTime(order.CreateTime),
		UpdateTime:           util.FormatBeijingTime(order.UpdateTime),
	}

	return &order_mgr_pb.QueryRechargeWithdrawOrderRsp{
		OrderInfo: orderInfo,
	}, nil
}
