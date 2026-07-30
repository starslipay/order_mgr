# 生成项目代码框架
goctl rpc protoc order_mgr.proto --go_out=. --go-grpc_out=. --zrpc_out=.

pushd
cd model/mysql
# 生成mysql代码
goctl model mysql ddl -src order.sql -dir .
popd