package etcdv3

import (
	"net"
	"reflect"
	"testing"
)

func TestReserve(t *testing.T) {
	tests := []struct {
		id       string
		ip       string
		expected bool
	}{
		{
			"100",
			"192.168.0.100",
			true,
		},
	}

	s, err := New("cni", nil)
	if err != nil {
		t.Errorf("error happened when initializing etcd store: %s", err)
	}

	for _, test := range tests {
		ok, err := s.Reserve(test.id, net.ParseIP(test.ip), "0")
		if err != nil {
			t.Errorf("error happened when reserving: %s", err)
		}

		if !reflect.DeepEqual(test.expected, ok) {
			t.Errorf("expected %#v, got %#v", test.expected, ok)
		}
	}
}

func TestLastReservedIP(t *testing.T) {
	tests := []struct {
		rangeIP  string
		expected string
	}{
		{
			"0",
			"192.168.0.100",
		},
	}

	s, err := New("cni", nil)
	if err != nil {
		t.Errorf("error happened when initializing etcd store: %s", err)
	}

	for _, test := range tests {
		ip, err := s.LastReservedIP(test.rangeIP)
		if err != nil {
			t.Errorf("error happened when reserving: %s", err)
		}

		if !reflect.DeepEqual(test.expected, string(ip)) {
			t.Errorf("expected %s, got %s", test.expected, string(ip))
		}
	}
}

func TestRelease(t *testing.T) {}

func TestReleaseByID(t *testing.T) {}
