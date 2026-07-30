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

func (l *BanPaySuccessOrderLogic) BanPaySuccessOrder(in *order_mgr_pb.BanPaySuccessOrderReq) (rsp *order_mgr_pb.BanPaySuccessOrderRsp, err error) {
	err = l.svcCtx.SqlMasterConn.TransactCtx(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		// 在事务内创建绑定到当前会话的 model，所有操作在同一事务中执行
		model := l.svcCtx.TOrderModelMaster.WithSession(session)

		// 使用 FOR UPDATE 锁定订单行，事务持有行锁直到提交
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
			if order.Amount != in.Amount {
				return xerror.NewBizError(codes.Internal, xerr.ErrCodeOrderAlreadySuccessButInfoNotMatch, "order already success, but amount not match")
			}
			if order.UserId != in.UserId {
				return xerror.NewBizError(codes.Internal, xerr.ErrCodeOrderAlreadySuccessButInfoNotMatch, "order already success, but user id not match")
			}
			if order.Uid != in.Uid {
				return xerror.NewBizError(codes.Internal, xerr.ErrCodeOrderAlreadySuccessButInfoNotMatch, "order already success, but uid not match")
			}
			if order.Spid != in.Spid {
				return xerror.NewBizError(codes.Internal, xerr.ErrCodeOrderAlreadySuccessButInfoNotMatch, "order already success, but spid not match")
			}

			order.TradeState = consts.OrderTradeStateSuccess
			payTime, err := time.Parse("2006-01-02 15:04:05", in.PayTime)
			if err != nil {
				return xerror.NewBizError(codes.Internal, xerr.ErrCodeOrderTradeStateInvalid, "pay time format error")
			}
			order.PayTime = payTime
			err = model.Update(ctx, order)
			if err != nil {
				return xerror.NewBizError(codes.Internal, xerr.ErrCodeDB, err.Error())
			}

		case consts.OrderTradeStateSuccess:
			// 已支付成功：校验关键信息一致性
			if order.Amount != in.Amount {
				return xerror.NewBizError(codes.Internal, xerr.ErrCodeOrderTradeStateInvalid, "order already success, but amount not match")
			}
			if order.UserId != in.UserId {
				return xerror.NewBizError(codes.Internal, xerr.ErrCodeOrderTradeStateInvalid, "order already success, but user id not match")
			}
			if order.Uid != in.Uid {
				return xerror.NewBizError(codes.Internal, xerr.ErrCodeOrderTradeStateInvalid, "order already success, but uid not match")
			}
			if order.Spid != in.Spid {
				return xerror.NewBizError(codes.Internal, xerr.ErrCodeOrderTradeStateInvalid, "order already success, but spid not match")
			}
		case consts.OrderTradeStateClose:
			return xerror.NewBizError(codes.Internal, xerr.ErrCodeOrderAlreadyClosed, "order already closed")

		default:
			return xerror.NewBizError(codes.Internal, xerr.ErrCodeOrderTradeStateInvalid, "order trade state invalid")
		}

		// 生成订单成功凭证
		orderSuccessToken := util.GenOrderSuccessToken(order.TransactionId, order.Spid, order.UserId, order.Uid, order.Amount)
		rsp = &order_mgr_pb.BanPaySuccessOrderRsp{
			OrderInfo: &order_mgr_pb.OrderInfo{
				TransactionId:     order.TransactionId,
				Spid:              order.Spid,
				UserId:            order.UserId,
				Uid:               order.Uid,
				Amount:            order.Amount,
				PayTime:           order.PayTime.Format("2006-01-02 15:04:05"),
				TradeState:        consts.OrderTradeStateSuccess,
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
