package mysql

import (
	"errors"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var ErrNotFound = sqlx.ErrNotFound

// ErrNoRowsAffected 表示修改类 SQL（Insert/Update/Delete）执行成功但实际影响行数不符合预期，
// 通常意味着目标记录不存在、并发条件未命中或数据未发生变化。
var ErrNoRowsAffected = errors.New("no rows affected")
