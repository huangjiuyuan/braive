export CNI_PATH=$GOPATH/src/github.com/huangjiuyuan/braive/bin
export NETCONFPATH=$GOPATH/src/github.com/huangjiuyuan/braive/conf

./docker-run.sh --rm busybox ifconfig
