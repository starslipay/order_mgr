$MODULE_NAME = "order_mgr"
$VERSION = "v1.0.0"
$IMAGE_NAME = "${MODULE_NAME}:${VERSION}"

docker rm -f $MODULE_NAME
docker rmi -f $IMAGE_NAME
docker build -t $IMAGE_NAME .
docker run -d --name $MODULE_NAME --network dev_pay_net -p 30883:8080 -p 9095:9090 $IMAGE_NAME
docker ps
docker logs $MODULE_NAME
