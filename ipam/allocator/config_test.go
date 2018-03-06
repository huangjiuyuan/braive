package allocator

import (
	"net"
	"reflect"
	"testing"
)

func TestLoadIPAllocator(t *testing.T) {
	subnet, err := parseCIDR("10.22.0.0/16")
	if err != nil {
		t.Errorf("expected no error, got %s", err)
	}

	expected := &IPAMConfig{
		Range: &Range{
			RangeStart: net.ParseIP("10.22.0.2"),
			RangeEnd:   net.ParseIP("10.22.255.255"),
			Subnet:     *subnet,
			Gateway:    net.ParseIP("10.22.0.1"),
		},
	}

	config, err := LoadIPAllocator()
	if err != nil {
		t.Errorf("expected no error, got %s", err)
	}

	if !reflect.DeepEqual(expected, config) {
		t.Errorf("expected %#v, got %#v", expected, config)
	}
}
