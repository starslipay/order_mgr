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

type CloseOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCloseOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CloseOrderLogic {
	return &CloseOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CloseOrderLogic) CloseOrder(in *order_mgr_pb.CloseOrderReq) (*order_mgr_pb.CloseOrderRsp, error) {
	// 先无锁查询判断订单是否存在
	order, err := l.svcCtx.TOrderModelMaster.FindOne(l.ctx, in.TransactionId)
	if err != nil && err != mysql.ErrNotFound {
		return nil, xerror.NewBizError(codes.Internal, xerr.ErrCodeDB, err.Error())
	}

	if order != nil {
		return l.closeExistingOrder(in)
	}

	// 订单不存在，插入一笔关闭状态的订单
	return l.closeNonExistOrder(in)
}

func (l *CloseOrderLogic) closeExistingOrder(in *order_mgr_pb.CloseOrderReq) (*order_mgr_pb.CloseOrderRsp, error) {
	var rsp *order_mgr_pb.CloseOrderRsp

	err := l.svcCtx.SqlMasterConn.TransactCtx(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		// 在事务内创建绑定到当前会话的 model
		model := l.svcCtx.TOrderModelMaster.WithSession(session)

		order, err := model.FindOneForUpdate(ctx, in.TransactionId)
		if err != nil {
			if err == mysql.ErrNotFound {
				// 理论上不应该发生，因为订单存在时才会调用此方法
				return xerror.NewBizError(codes.Internal, xerr.ErrCodeOrderNotFound, "order not found for update")
			}
			return xerror.NewBizError(codes.Internal, xerr.ErrCodeDB, err.Error())
		}

		switch order.TradeState {
		case consts.OrderTradeStateInit:
			// 初始状态：修改为关单
			order.TradeState = consts.OrderTradeStateClose

			if err := model.Update(ctx, order); err != nil {
				return xerror.NewBizError(codes.Internal, xerr.ErrCodeDB, err.Error())
			}

			rsp = &order_mgr_pb.CloseOrderRsp{
				TransactionId: order.TransactionId,
			}

		case consts.OrderTradeStateClose:
			// 已关单：直接返回成功（幂等）
			rsp = &order_mgr_pb.CloseOrderRsp{
				TransactionId: order.TransactionId,
			}

		case consts.OrderTradeStateSuccess:
			// 并发错误：订单已成功状态, 让业务重试
			return xerror.NewBizError(codes.Internal, xerr.ErrCodeOrderAlreadySuccess, "order already success")

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

func (l *CloseOrderLogic) closeNonExistOrder(in *order_mgr_pb.CloseOrderReq) (*order_mgr_pb.CloseOrderRsp, error) {
	order := &mysql.TOrder{
		TransactionId: in.TransactionId,
		TradeState:    consts.OrderTradeStateClose,
	}

	_, err := l.svcCtx.TOrderModelMaster.Insert(l.ctx, order)
	if err != nil {
		// 并发场景：其他线程已插入同一订单（唯一键冲突）
		if isDuplicateKeyError(err) {
			// 并发冲突，无单关单订单已存在，让业务调关补接口重试
			return nil, xerror.NewBizError(codes.Internal, xerr.ErrCodeOrderInsertOrderDuplicate, "insert close order, order already exist")
		}
		return nil, xerror.NewBizError(codes.Internal, xerr.ErrCodeDB, err.Error())
	}

	return &order_mgr_pb.CloseOrderRsp{
		TransactionId: order.TransactionId,
	}, nil
}
