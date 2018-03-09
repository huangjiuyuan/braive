package etcdv3

import (
	"net"
	"reflect"
	"testing"
)

func TestReserve(t *testing.T) {
	tests := []struct {
		ip       string
		id       string
		rangeID  string
		expected bool
	}{
		{
			"192.168.0.100",
			"100",
			"0",
			true,
		},
	}

	s, err := New("cni", nil)
	if err != nil {
		t.Errorf("error happened when initializing etcd store: %s", err)
	}

	for _, test := range tests {
		ok, err := s.Reserve(test.id, net.ParseIP(test.ip), test.rangeID)
		if err != nil {
			t.Errorf("error happened when reserving: %s", err)
		}

		if !reflect.DeepEqual(test.expected, ok) {
			t.Errorf("expected %#v, got %#v", test.expected, ok)
		}

		s.EtcdClient.Delete(s.Ctx, test.ip)
		s.EtcdClient.Delete(s.Ctx, lastIPPrefix+test.rangeID)
	}
}

func TestLastReservedIP(t *testing.T) {
	tests := []struct {
		ip       string
		rangeID  string
		expected string
	}{
		{
			"192.168.0.100",
			"0",
			"192.168.0.100",
		},
	}

	s, err := New("cni", nil)
	if err != nil {
		t.Errorf("error happened when initializing etcd store: %s", err)
	}

	for _, test := range tests {
		s.EtcdClient.KV.Put(s.Ctx, lastIPPrefix+test.rangeID, test.ip)

		ip, err := s.LastReservedIP(test.rangeID)
		if err != nil {
			t.Errorf("error happened when reserving: %s", err)
		}

		if !reflect.DeepEqual(test.expected, string(ip)) {
			t.Errorf("expected %s, got %s", test.expected, string(ip))
		}

		s.EtcdClient.Delete(s.Ctx, lastIPPrefix+test.rangeID)
	}
}

func TestRelease(t *testing.T) {
	tests := []struct {
		ip string
		id string
	}{
		{
			"192.168.0.100",
			"100",
		},
	}

	s, err := New("cni", nil)
	if err != nil {
		t.Errorf("error happened when initializing etcd store: %s", err)
	}

	for _, test := range tests {
		s.EtcdClient.KV.Put(s.Ctx, test.ip, test.id)

		err := s.Release(net.ParseIP(test.ip))
		if err != nil {
			t.Errorf("error happened when reserving: %s", err)
		}

		gr, err := s.EtcdClient.KV.Get(s.Ctx, test.ip)
		if err != nil {
			t.Errorf("Failed to get IP: %s", err)
		}

		if len(gr.Kvs) != 0 {
			t.Errorf("expected 0, got %d", len(gr.Kvs))
		}
	}
}

func TestReleaseByID(t *testing.T) {
	tests := []struct {
		ip string
		id string
	}{
		{
			"192.168.0.100",
			"100",
		},
	}

	s, err := New("cni", nil)
	if err != nil {
		t.Errorf("error happened when initializing etcd store: %s", err)
	}

	for _, test := range tests {
		s.EtcdClient.KV.Put(s.Ctx, test.ip, test.id)

		err := s.ReleaseByID(test.id)
		if err != nil {
			t.Errorf("error happened when reserving: %s", err)
		}

		gr, err := s.EtcdClient.KV.Get(s.Ctx, test.ip)
		if err != nil {
			t.Errorf("Failed to get IP: %s", err)
		}

		if len(gr.Kvs) != 0 {
			t.Errorf("expected 0, got %d", len(gr.Kvs))
		}
	}
}
