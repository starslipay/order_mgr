package mysql

import (
	"context"
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
