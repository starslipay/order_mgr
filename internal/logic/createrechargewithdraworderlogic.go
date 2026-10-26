package logic

import (
	"context"

	"github.com/starslipay/order_mgr/internal/consts"
	"github.com/starslipay/order_mgr/internal/svc"
	"github.com/starslipay/order_mgr/internal/xerr"
	"github.com/starslipay/order_mgr/model/mysql"
	"github.com/starslipay/order_mgr/order_mgr_pb"
	"github.com/starslipay/paycomm/xerror"
	"google.golang.org/grpc/codes"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateRechargeWithdrawOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateRechargeWithdrawOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateRechargeWithdrawOrderLogic {
	return &CreateRechargeWithdrawOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateRechargeWithdrawOrderLogic) CreateRechargeWithdrawOrder(in *order_mgr_pb.CreateRechargeWithdrawOrderReq) (*order_mgr_pb.CreateRechargeWithdrawOrderRsp, error) {
	order := &mysql.TRechargeWithdrawOrder{
		TransactionId:        in.TransactionId,
		ChannelTransactionId: in.ChannelTransactionId,
		UserId:               in.UserId,
		Uid:                  in.Uid,
		BizType:              int64(in.BizType),
		BankType:             int64(in.BankType),
		Amount:               in.Amount,
		CurType:              int64(in.CurType),
		TradeState:           consts.RechargeWithdrawTradeStateInit,
		Memo:                 in.Memo,
	}

	_, err := l.svcCtx.TRechargeWithdrawOrderModelMaster.Insert(l.ctx, order)
	if err != nil {
		if isDuplicateKeyError(err) {
			return nil, xerror.NewBizError(codes.Internal, xerr.ErrCodeRechargeWithdrawOrderDuplicate, "recharge withdraw order already exists")
		}
		return nil, xerror.NewBizError(codes.Internal, xerr.ErrCodeDB, err.Error())
	}

	return &order_mgr_pb.CreateRechargeWithdrawOrderRsp{
		TransactionId: order.TransactionId,
	}, nil
}
