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
		OutOrderNo:    order.OutOrderNo,
		MerchantId:    order.MerchantId,
		MerchantUid:   order.MerchantUid,
		MerchantName:  order.MerchantName,
		UserId:        order.UserId,
		Uid:           order.Uid,
		Amount:        order.Amount,
		PayTime:       util.FormatBeijingTime(order.PayTime),
		TradeState:    int32(order.TradeState),
	}

	// 如果订单状态为成功，生成订单成功凭证
	if consts.OrderTradeStateSuccess == order.TradeState {
		orderInfo.OrderSuccessToken = util.GenOrderSuccessToken(order.TransactionId, order.OutOrderNo,
			order.MerchantUid, order.Uid, order.Amount)
	}

	return &order_mgr_pb.QueryOrderRsp{
		OrderInfo: orderInfo,
	}, nil
}
