package logic

import (
	"context"
	"strings"

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

type BanPayPreOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBanPayPreOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BanPayPreOrderLogic {
	return &BanPayPreOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *BanPayPreOrderLogic) BanPayPreOrder(in *order_mgr_pb.BanPayPreOrderReq) (*order_mgr_pb.BanPayPreOrderRsp, error) {
	// 先无锁查询判断订单是否存在
	order, err := l.svcCtx.TOrderModelMaster.FindOne(l.ctx, in.TransactionId)
	if err != nil && err != mysql.ErrNotFound {
		return nil, xerror.NewBizError(codes.Internal, xerr.ErrCodeDB, err.Error())
	}

	if order != nil {
		return l.handleExistingOrder(in)
	}

	// 订单不存在，插入一笔 init 状态的新订单
	return l.handleNewOrder(in)
}

func (l *BanPayPreOrderLogic) handleExistingOrder(in *order_mgr_pb.BanPayPreOrderReq) (rsp *order_mgr_pb.BanPayPreOrderRsp, err error) {
	err = l.svcCtx.SqlMasterConn.TransactCtx(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		// 在事务内创建绑定到当前会话的 model，所有操作在同一事务中执行
		model := l.svcCtx.TOrderModelMaster.WithSession(session)

		// 加锁查
		order, err := model.FindOneForUpdate(ctx, in.TransactionId)
		if err != nil {
			if err == mysql.ErrNotFound {
				return xerror.NewBizError(codes.Internal, xerr.ErrCodeOrderNotFound, "for update, but order not found")
			}
			return xerror.NewBizError(codes.Internal, xerr.ErrCodeDB, err.Error())
		}

		switch order.TradeState {
		case consts.OrderTradeStateInit:
			// 初始状态：用输入参数更新订单信息(换支付方式时需要更新)
			order.OutOrderNo = in.OutOrderNo
			order.MerchantId = in.MerchantId
			order.MerchantUid = in.MerchantUid
			order.MerchantName = in.MerchantName
			order.UserId = in.UserId
			order.Uid = in.Uid
			order.Amount = in.Amount
			order.CurType = consts.OrderCurTypeCNY
			order.PayType = consts.OrderPayTypeBanPay

			if err := model.Update(ctx, order); err != nil {
				return xerror.NewBizError(codes.Internal, xerr.ErrCodeDB, err.Error())
			}

			rsp = &order_mgr_pb.BanPayPreOrderRsp{
				TransactionId:    order.TransactionId,
				IsAlreadySuccess: 0,
			}

		case consts.OrderTradeStateSuccess:
			// 生成订单成功凭证
			orderSuccessToken := util.GenOrderSuccessToken(order.TransactionId, order.OutOrderNo,
				order.MerchantUid, order.Uid, order.Amount)

			rsp = &order_mgr_pb.BanPayPreOrderRsp{
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

		case consts.OrderTradeStateClose:
			return xerror.NewBizError(codes.Internal, xerr.ErrCodeOrderAlreadyClosed, "order already closed")

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

func (l *BanPayPreOrderLogic) handleNewOrder(in *order_mgr_pb.BanPayPreOrderReq) (*order_mgr_pb.BanPayPreOrderRsp, error) {
	order := &mysql.TOrder{
		TransactionId: in.TransactionId,
		OutOrderNo:    in.OutOrderNo,
		MerchantId:    in.MerchantId,
		MerchantUid:   in.MerchantUid,
		MerchantName:  in.MerchantName,
		UserId:        in.UserId,
		Uid:           in.Uid,
		Amount:        in.Amount,
		PayType:       consts.OrderPayTypeBanPay,
		TradeState:    consts.OrderTradeStateInit,
	}

	_, err := l.svcCtx.TOrderModelMaster.Insert(l.ctx, order)
	if err != nil {
		if isDuplicateKeyError(err) {
			return nil, xerror.NewBizError(codes.Internal, xerr.ErrCodeOrderDuplicate, "insert order, order already exists")
		}
		return nil, xerror.NewBizError(codes.Internal, xerr.ErrCodeDB, err.Error())
	}

	return &order_mgr_pb.BanPayPreOrderRsp{
		TransactionId:    order.TransactionId,
		IsAlreadySuccess: 0,
	}, nil
}

func isDuplicateKeyError(err error) bool {
	return err != nil && strings.Contains(err.Error(), "Duplicate entry")
}
