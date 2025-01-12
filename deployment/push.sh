DOMAIN=$1
VERSION=$2

docker tag coolcar/$DOMAIN ccr.ccs.tencentyun.com/coolcarc/$DOMAIN:$VERSION
docker push ccr.ccs.tencentyun.com/coolcarc/$DOMAIN:$VERSION