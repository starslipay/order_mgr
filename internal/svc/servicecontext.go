package svc

import (
	"github.com/starslipay/order_mgr/internal/config"

	"github.com/starslipay/order_mgr/model/mysql"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config            config.Config
	SqlMasterConn     sqlx.SqlConn
	SqlSlaveConn      sqlx.SqlConn
	TOrderModelMaster mysql.TOrderModel
	TOrderModelSlave  mysql.TOrderModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 关闭db熔断保护
	SqlMasterConn := sqlx.NewMysql(c.MasterDBConfig.DataSource, sqlx.WithAcceptable(func(err error) bool { return true }))
	SqlSlaveConn := sqlx.NewMysql(c.SlaveDBConfig.DataSource, sqlx.WithAcceptable(func(err error) bool { return true }))

	return &ServiceContext{
		Config:            c,
		SqlMasterConn:     SqlMasterConn,
		SqlSlaveConn:      SqlSlaveConn,
		TOrderModelMaster: mysql.NewTOrderModel(SqlMasterConn),
		TOrderModelSlave:  mysql.NewTOrderModel(SqlSlaveConn),
	}
}
