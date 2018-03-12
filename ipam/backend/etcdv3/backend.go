package etcdv3

import (
	"context"
	"net"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/golang/glog"
)

var (
	etcdDialTimeout    = 2 * time.Second
	etcdRequestTimeout = 5 * time.Second
	lastIPPrefix       = "lastReservedIP."
)

type Store struct {
	EtcdClient *clientv3.Client
	Ctx        context.Context
	Cancel     context.CancelFunc
}

func New(network string, endpoints []string) (*Store, error) {
	// Set default endpoints if not specified.
	if len(endpoints) == 0 {
		endpoints = append(endpoints, "10.9.96.4:2379")
	}

	// Initialize etcd client.
	ctx, cancel := context.WithTimeout(context.Background(), etcdRequestTimeout)
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: etcdDialTimeout,
	})
	if err != nil {
		glog.Errorf("Failed to start etcd client: %s", err)
		return nil, err
	}

	return &Store{cli, ctx, cancel}, nil
}

func (s *Store) Reserve(id string, ip net.IP, rangeID string) (bool, error) {
	// Check if an IP has been reserved before.
	gr, err := s.EtcdClient.KV.Get(s.Ctx, ip.String())
	if err != nil {
		glog.Errorf("Failed to get IP %s", ip.String())
		return false, err
	}
	if len(gr.Kvs) != 0 {
		glog.Warningf("IP %s has been reserved before", string(gr.Kvs[0].Value))
		return false, nil
	}

	// Reserve an IP.
	_, err = s.EtcdClient.KV.Put(s.Ctx, ip.String(), id)
	if err != nil {
		glog.Errorf("Failed to reserve IP %s", ip.String())
		return false, err
	}
	glog.Info("Succeed to reserve IP %s with ID %s", ip.String(), id)

	// Update last reserved IP.
	_, err = s.EtcdClient.KV.Put(s.Ctx, lastIPPrefix+rangeID, ip.String())
	if err != nil {
		glog.Errorf("Failed to update last reserved IP %s", ip.String())
		return false, err
	}
	glog.Info("Last reserved IP %s has been updated with ID %s", ip.String(), id)

	return true, nil
}

func (s *Store) LastReservedIP(rangeID string) (net.IP, error) {
	// Get last reserved IP.
	gr, err := s.EtcdClient.KV.Get(s.Ctx, lastIPPrefix+rangeID)
	if err != nil {
		glog.Errorf("Failed to get last reserved IP")
		return nil, err
	}

	if len(gr.Kvs) != 0 {
		return gr.Kvs[0].Value, nil
	}
	return nil, nil
}

func (s *Store) Release(ip net.IP) error {
	// Release an IP.
	dr, err := s.EtcdClient.Delete(s.Ctx, ip.String(), clientv3.WithPrevKV())
	if err != nil {
		glog.Errorf("Failed to release IP %s", ip.String())
		return err
	}
	glog.Infof("Succeed to release IP %s with ID %s", ip.String(), dr.PrevKvs[0].Value)

	return nil
}

func (s *Store) ReleaseByID(id string) error {
	// Get all reserved IPs.
	gr, err := s.EtcdClient.Get(s.Ctx, "", clientv3.WithPrefix())
	if err != nil {
		glog.Errorf("Failed to get all reserved IPs")
		return err
	}

	exist := false
	for _, kv := range gr.Kvs {
		if string(kv.Value) == id {
			// Release an IP by ID.
			dr, err := s.EtcdClient.Delete(s.Ctx, string(kv.Key), clientv3.WithPrevKV())
			if err != nil {
				glog.Errorf("Failed to release IP %s by ID %s", string(kv.Key), dr.PrevKvs[0].Value)
				return err
			}
			exist = true
			glog.Infof("Succeed to release IP %s by ID %s", string(kv.Key), dr.PrevKvs[0].Value)
		}
	}

	if !exist {
		glog.Warningf("IP with ID %s not found", id)
	}

	return nil
}

func (s *Store) Close() error {
	err := s.EtcdClient.Close()
	if err != nil {
		return err
	}

	return nil
}
