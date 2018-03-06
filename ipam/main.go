package main

import (
	"github.com/containernetworking/cni/pkg/skel"
	"github.com/containernetworking/cni/pkg/version"
	"github.com/huangjiuyuan/braive/ipam/allocator"
)

func main() {
	skel.PluginMain(cmdAdd, cmdDel, version.All)
}

func cmdAdd(args *skel.CmdArgs) error {
	config, err := allocator.LoadIPAllocator()
	if err != nil {
		return err
	}

	return nil
}

func cmdDel(args *skel.CmdArgs) error {
	config, err := allocator.LoadIPAllocator()
	if err != nil {
		return err
	}

	return nil
}
