export CNI_PATH=/opt/cni/bin/
export NETCONFPATH=/etc/cni/net.d

./docker-run.sh --rm busybox ifconfig
