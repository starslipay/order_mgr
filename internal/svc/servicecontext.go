package svc

import (
	"database/sql"
	"time"

	"github.com/starslipay/order_mgr/internal/config"

	"github.com/starslipay/order_mgr/model/mysql"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config            config.Config
	SqlMasterConn     sqlx.SqlConn
	SqlSlaveConn      sqlx.SqlConn
	TOrderModelMaster mysql.TOrderModel
	TOrderModelSlave  mysql.TOrderModel
}

// newDBConn 自建 *sql.DB 以支持自定义连接池参数
func newDBConn(dataSource string, maxOpen, maxIdle, lifetimeSec int) sqlx.SqlConn {
	db, err := sql.Open("mysql", dataSource)
	if err != nil {
		logx.Must(err)
	}
	db.SetMaxOpenConns(maxOpen)
	db.SetMaxIdleConns(maxIdle)
	db.SetConnMaxLifetime(time.Duration(lifetimeSec) * time.Second)
	if err = db.Ping(); err != nil {
		logx.Must(err)
	}
	// 关闭db熔断保护
	return sqlx.NewSqlConnFromDB(db, sqlx.WithAcceptable(func(err error) bool { return true }))
}

func NewServiceContext(c config.Config) *ServiceContext {
	SqlMasterConn := newDBConn(c.MasterDBConfig.DataSource,
		c.MasterDBConfig.MaxOpenConns, c.MasterDBConfig.MaxIdleConns, c.MasterDBConfig.ConnMaxLifetimeSec)
	SqlSlaveConn := newDBConn(c.SlaveDBConfig.DataSource,
		c.SlaveDBConfig.MaxOpenConns, c.SlaveDBConfig.MaxIdleConns, c.SlaveDBConfig.ConnMaxLifetimeSec)

	return &ServiceContext{
		Config:            c,
		SqlMasterConn:     SqlMasterConn,
		SqlSlaveConn:      SqlSlaveConn,
		TOrderModelMaster: mysql.NewTOrderModel(SqlMasterConn),
		TOrderModelSlave:  mysql.NewTOrderModel(SqlSlaveConn),
	}
}
