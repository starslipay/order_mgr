package logic

import (
	"context"

	"github.com/starslipay/order_mgr/internal/consts"
	"github.com/starslipay/order_mgr/internal/svc"
	"github.com/starslipay/order_mgr/internal/util"
	"github.com/starslipay/order_mgr/internal/xerr"
	"github.com/starslipay/order_mgr/model/mysql"
	"github.com/starslipay/order_mgr/order_mgr_pb"
	"github.com/starslipay/paycomm/xerror"
	"google.golang.org/grpc/codes"

	"github.com/zeromicro/go-zero/core/logx"
)

type QueryOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewQueryOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryOrderLogic {
	return &QueryOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *QueryOrderLogic) QueryOrder(in *order_mgr_pb.QueryOrderReq) (*order_mgr_pb.QueryOrderRsp, error) {
	order, err := l.svcCtx.TOrderModelMaster.FindOne(l.ctx, in.TransactionId)
	if err != nil {
		if err == mysql.ErrNotFound {
			return nil, xerror.NewBizError(codes.Internal, xerr.ErrCodeOrderNotFound, "order not found")
		} else {
			return nil, xerror.NewBizError(codes.Internal, xerr.ErrCodeDB, err.Error())
		}
	}

	orderInfo := &order_mgr_pb.OrderInfo{
		TransactionId: order.TransactionId,
		Spid:          order.Spid,
		UserId:        order.UserId,
		Uid:           order.Uid,
		Amount:        order.Amount,
		PayTime:       order.PayTime.Format("2006-01-02 15:04:05"),
		TradeState:    int32(order.TradeState),
	}

	if consts.OrderTradeStateSuccess == order.TradeState {
		orderInfo.OrderSuccessToken = util.GenOrderSuccessToken(order.TransactionId, order.Spid, order.UserId, order.Uid, order.Amount)
	}

	return &order_mgr_pb.QueryOrderRsp{
		OrderInfo: orderInfo,
	}, nil
}
