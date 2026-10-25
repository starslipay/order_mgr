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
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type RefundOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRefundOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefundOrderLogic {
	return &RefundOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RefundOrderLogic) RefundOrder(in *order_mgr_pb.RefundOrderReq) (*order_mgr_pb.RefundOrderRsp, error) {
	order, err := l.svcCtx.TOrderModelMaster.FindOne(l.ctx, in.TransactionId)
	if err != nil {
		if err == mysql.ErrNotFound {
			return nil, xerror.NewBizError(codes.Internal, xerr.ErrCodeOrderNotFound, "order not found")
		}
		return nil, xerror.NewBizError(codes.Internal, xerr.ErrCodeDB, err.Error())
	}

	switch order.TradeState {
	case consts.OrderTradeStateSuccess, consts.OrderTradeStateRefund:
		// 成功状态或已退款状态：允许退款/重入
	default:
		return nil, xerror.NewBizError(codes.Internal, xerr.ErrCodeOrderTradeStateInvalid, "order trade state invalid for refund")
	}

	var rsp *order_mgr_pb.RefundOrderRsp
	err = l.svcCtx.SqlMasterConn.TransactCtx(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		model := l.svcCtx.TOrderModelMaster.WithSession(session)

		order, err := model.FindOneForUpdate(ctx, in.TransactionId)
		if err != nil {
			return xerror.NewBizError(codes.Internal, xerr.ErrCodeDB, err.Error())
		}

		switch order.TradeState {
		case consts.OrderTradeStateSuccess:
			order.TradeState = consts.OrderTradeStateRefund
			if err := model.Update(ctx, order); err != nil {
				return xerror.NewBizError(codes.Internal, xerr.ErrCodeDB, err.Error())
			}
			rsp = &order_mgr_pb.RefundOrderRsp{
				TransactionId: order.TransactionId,
			}

		case consts.OrderTradeStateRefund:
			// 已退款：幂等重入
			rsp = &order_mgr_pb.RefundOrderRsp{
				TransactionId: order.TransactionId,
			}

		default:
			return xerror.NewBizError(codes.Internal, xerr.ErrCodeOrderTradeStateInvalid, "order trade state invalid")
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return rsp, nil
}
