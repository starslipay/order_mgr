package logic

import (
	"context"
	"time"

	"github.com/starslipay/order_mgr/internal/consts"
	"github.com/starslipay/order_mgr/internal/svc"
	"github.com/starslipay/order_mgr/internal/util"
	"github.com/starslipay/order_mgr/internal/xerr"
	"github.com/starslipay/order_mgr/model/mysql"
	"github.com/starslipay/order_mgr/order_mgr_pb"
	"github.com/starslipay/paycomm/xerror"
	"google.golang.org/grpc/codes"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type BanPaySuccessOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBanPaySuccessOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BanPaySuccessOrderLogic {
	return &BanPaySuccessOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *BanPaySuccessOrderLogic) checkOrderInfoConsistency(order *mysql.TOrder, in *order_mgr_pb.BanPaySuccessOrderReq) error {
	if order.OutOrderNo != in.OutOrderNo {
		return xerror.NewBizError(codes.Internal, xerr.ErrCodeOrderInfoNotMatch, "order update success, but out_order_no not match")
	}

	if order.UserId != in.UserId {
		return xerror.NewBizError(codes.Internal, xerr.ErrCodeOrderInfoNotMatch, "order update success, but user id not match")
	}
	if order.Uid != in.Uid {
		return xerror.NewBizError(codes.Internal, xerr.ErrCodeOrderInfoNotMatch, "order update success, but uid not match")
	}
	if order.MerchantId != in.MerchantId {
		return xerror.NewBizError(codes.Internal, xerr.ErrCodeOrderInfoNotMatch, "order update success, but merchant id not match")
	}

	if order.MerchantUid != in.MerchantUid {
		return xerror.NewBizError(codes.Internal, xerr.ErrCodeOrderInfoNotMatch, "order update success, but merchant uid not match")
	}

	if order.Amount != in.Amount {
		return xerror.NewBizError(codes.Internal, xerr.ErrCodeOrderInfoNotMatch, "order update success, but amount not match")
	}
	return nil
}

func (l *BanPaySuccessOrderLogic) BanPaySuccessOrder(in *order_mgr_pb.BanPaySuccessOrderReq) (rsp *order_mgr_pb.BanPaySuccessOrderRsp, err error) {
	err = l.svcCtx.SqlMasterConn.TransactCtx(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		// 在事务内创建绑定到当前会话的 model，所有操作在同一事务中执行
		model := l.svcCtx.TOrderModelMaster.WithSession(session)

		order, err := model.FindOneForUpdate(ctx, in.TransactionId)
		if err != nil {
			if err == mysql.ErrNotFound {
				return xerror.NewBizError(codes.Internal, xerr.ErrCodeOrderNotFound, "success order, but order not found")
			}
			return xerror.NewBizError(codes.Internal, xerr.ErrCodeDB, err.Error())
		}

		switch order.TradeState {
		case consts.OrderTradeStateInit:
			// 校验关键信息一致性
			err = l.checkOrderInfoConsistency(order, in)
			if err != nil {
				return err
			}

			payTime, err := time.Parse("2006-01-02 15:04:05", in.PayTime)
			if err != nil {
				return xerror.NewBizError(codes.Internal, xerr.ErrCodeOrderTradeStateInvalid, "pay time format error")
			}
			order.PayTime = payTime
			order.TradeState = consts.OrderTradeStateSuccess

			err = model.Update(ctx, order)
			if err != nil {
				return xerror.NewBizError(codes.Internal, xerr.ErrCodeDB, err.Error())
			}

		case consts.OrderTradeStateSuccess:
			// 校验关键信息一致性
			err = l.checkOrderInfoConsistency(order, in)
			if err != nil {
				return err
			}
		case consts.OrderTradeStateClose:
			// 并发冲突，订单已关闭状态，让业务调关补接口重试
			return xerror.NewBizError(codes.Internal, xerr.ErrCodeOrderAlreadyClosed, "order already closed")

		default:
			return xerror.NewBizError(codes.Internal, xerr.ErrCodeOrderTradeStateInvalid, "order trade state invalid")
		}

		// 生成订单成功凭证
		orderSuccessToken := util.GenOrderSuccessToken(order.TransactionId, order.OutOrderNo, order.MerchantUid, order.Uid, order.Amount)
		rsp = &order_mgr_pb.BanPaySuccessOrderRsp{
			OrderInfo: &order_mgr_pb.OrderInfo{
				TransactionId:     order.TransactionId,
				OutOrderNo:        order.OutOrderNo,
				MerchantId:        order.MerchantId,
				MerchantUid:       order.MerchantUid,
				UserId:            order.UserId,
				Uid:               order.Uid,
				Amount:            order.Amount,
				PayTime:           order.PayTime.Format("2006-01-02 15:04:05"),
				TradeState:        int32(order.TradeState),
				OrderSuccessToken: orderSuccessToken,
			},
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return rsp, nil
}
