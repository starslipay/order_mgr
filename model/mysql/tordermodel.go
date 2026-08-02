package mysql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ TOrderModel = (*customTOrderModel)(nil)

type (
	// TOrderModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTOrderModel.
	TOrderModel interface {
		tOrderModel
		WithSession(session sqlx.Session) TOrderModel
		FindOneForUpdate(ctx context.Context, transactionId string) (*TOrder, error)
	}

	customTOrderModel struct {
		*defaultTOrderModel
	}
)

// NewTOrderModel returns a model for the database table.
func NewTOrderModel(conn sqlx.SqlConn) TOrderModel {
	return &customTOrderModel{
		defaultTOrderModel: newTOrderModel(conn),
	}
}

func (m *customTOrderModel) WithSession(session sqlx.Session) TOrderModel {
	return NewTOrderModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customTOrderModel) FindOneForUpdate(ctx context.Context, transactionId string) (*TOrder, error) {
	query := fmt.Sprintf("select %s from %s where `transaction_id` = ? limit 1 for update", tOrderRows, m.table)
	var resp TOrder
	err := m.conn.QueryRowCtx(ctx, &resp, query, transactionId)
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

// Insert 重写生成的 Insert，校验影响行数，确保确实插入了一行。
func (m *customTOrderModel) Insert(ctx context.Context, data *TOrder) (sql.Result, error) {
	ret, err := m.defaultTOrderModel.Insert(ctx, data)
	if err != nil {
		return nil, err
	}
	if err := checkRowsAffected(ret); err != nil {
		return nil, err
	}
	return ret, nil
}

// Update 重写生成的 Update，校验影响行数，避免更新未命中任何记录时被静默忽略。
func (m *customTOrderModel) Update(ctx context.Context, data *TOrder) error {
	query := fmt.Sprintf("update %s set %s where `transaction_id` = ?", m.table, tOrderRowsWithPlaceHolder)
	ret, err := m.conn.ExecCtx(ctx, query, data.OutOrderNo, data.MerchantId, data.MerchantUid, data.MerchantName, data.UserId, data.Uid, data.TradeState, data.Amount, data.CurType, data.PayType, data.PayTime, data.TransactionId)
	if err != nil {
		return err
	}
	return checkRowsAffected(ret)
}

// Delete 重写生成的 Delete，校验影响行数，避免删除未命中任何记录时被静默忽略。
func (m *customTOrderModel) Delete(ctx context.Context, transactionId string) error {
	query := fmt.Sprintf("delete from %s where `transaction_id` = ?", m.table)
	ret, err := m.conn.ExecCtx(ctx, query, transactionId)
	if err != nil {
		return err
	}
	return checkRowsAffected(ret)
}

// checkRowsAffected 校验一次修改类 SQL 的影响行数，行数为 0 时返回 ErrNoRowsAffected。
func checkRowsAffected(ret sql.Result) error {
	affected, err := ret.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrNoRowsAffected
	}
	return nil
}
