# Braive

Braive is a naive Kubernetes CNI implementation. It works based on Linux bridge and etcd.

## Deployment

Braive consists of two parts: an CNI plugin and an IPAM plugin.</br>
You have to build the plugins first, this requires an environment with Golang version 1.9 or above.</br>
To deploy the CNI plugin, run the following commands under `bridge` folder

```
go build
mv bridge braive
mkdir -p /opt/cni/bin
mv braive /opt/cni/bin
```

To deploy the IPAM plugin, run the following commands under `ipam` folder

```
go build
mv ipam braive-ipam
mv braive-ipam /opt/cni/bin
```

Copy configuration files under `conf` folder to `/etc/cni/net.d`

```
mkdir -p /etc/cni/net.d
cp /conf/* /etc/cni/net.d
```

## How It Works

The IPAM is based on etcd version 3. When a pod is created, CRI calls CNI and CNI calls our IPAM to allocate IP for the newly created pod. The allocated IP will be store in etcd. Our CNI plugin will create interfaces and write routes and iptables rules for the pod.
