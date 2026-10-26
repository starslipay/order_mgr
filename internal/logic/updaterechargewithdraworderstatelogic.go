package logic

import (
	"context"

	"github.com/starslipay/order_mgr/internal/svc"
	"github.com/starslipay/order_mgr/internal/xerr"
	"github.com/starslipay/order_mgr/model/mysql"
	"github.com/starslipay/order_mgr/order_mgr_pb"
	"github.com/starslipay/paycomm/xerror"
	"google.golang.org/grpc/codes"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateRechargeWithdrawOrderStateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateRechargeWithdrawOrderStateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateRechargeWithdrawOrderStateLogic {
	return &UpdateRechargeWithdrawOrderStateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateRechargeWithdrawOrderStateLogic) UpdateRechargeWithdrawOrderState(in *order_mgr_pb.UpdateRechargeWithdrawOrderStateReq) (*order_mgr_pb.UpdateRechargeWithdrawOrderStateRsp, error) {
	err := l.svcCtx.TRechargeWithdrawOrderModelMaster.UpdateTradeState(l.ctx, in.TransactionId, int64(in.TradeState))
	if err != nil {
		if err == mysql.ErrNoRowsAffected {
			return nil, xerror.NewBizError(codes.Internal, xerr.ErrCodeRechargeWithdrawOrderNotFound, "recharge withdraw order not found")
		}
		return nil, xerror.NewBizError(codes.Internal, xerr.ErrCodeDB, err.Error())
	}

	return &order_mgr_pb.UpdateRechargeWithdrawOrderStateRsp{
		TransactionId: in.TransactionId,
	}, nil
}
