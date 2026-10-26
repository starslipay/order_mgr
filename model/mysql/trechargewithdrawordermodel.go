package mysql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ TRechargeWithdrawOrderModel = (*customTRechargeWithdrawOrderModel)(nil)

type (
	// TRechargeWithdrawOrderModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTRechargeWithdrawOrderModel.
	TRechargeWithdrawOrderModel interface {
		tRechargeWithdrawOrderModel
		WithSession(session sqlx.Session) TRechargeWithdrawOrderModel
		UpdateTradeState(ctx context.Context, transactionId string, tradeState int64) error
	}

	customTRechargeWithdrawOrderModel struct {
		*defaultTRechargeWithdrawOrderModel
	}
)

// NewTRechargeWithdrawOrderModel returns a model for the database table.
func NewTRechargeWithdrawOrderModel(conn sqlx.SqlConn) TRechargeWithdrawOrderModel {
	return &customTRechargeWithdrawOrderModel{
		defaultTRechargeWithdrawOrderModel: newTRechargeWithdrawOrderModel(conn),
	}
}

func (m *customTRechargeWithdrawOrderModel) WithSession(session sqlx.Session) TRechargeWithdrawOrderModel {
	return NewTRechargeWithdrawOrderModel(sqlx.NewSqlConnFromSession(session))
}

// Insert 重写生成的 Insert，校验影响行数，确保确实插入了一行。
func (m *customTRechargeWithdrawOrderModel) Insert(ctx context.Context, data *TRechargeWithdrawOrder) (sql.Result, error) {
	ret, err := m.defaultTRechargeWithdrawOrderModel.Insert(ctx, data)
	if err != nil {
		return nil, err
	}
	if err := checkRowsAffected(ret); err != nil {
		return nil, err
	}
	return ret, nil
}

// UpdateTradeState 仅更新交易状态，校验影响行数，避免更新未命中任何记录时被静默忽略。
func (m *customTRechargeWithdrawOrderModel) UpdateTradeState(ctx context.Context, transactionId string, tradeState int64) error {
	query := fmt.Sprintf("update %s set `trade_state` = ? where `transaction_id` = ?", m.table)
	ret, err := m.conn.ExecCtx(ctx, query, tradeState, transactionId)
	if err != nil {
		return err
	}
	return checkRowsAffected(ret)
}
