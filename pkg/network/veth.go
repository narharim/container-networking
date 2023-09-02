package network

import (
	"fmt"

	"github.com/vishvananda/netlink"
	"github.com/vishvananda/netns"
)

type VethPair struct {
	Name1 string
	Name2 string
}

func (v VethPair) Create() error {
	ethPair := &netlink.Veth{
		LinkAttrs: netlink.LinkAttrs{Name: v.Name1},
		PeerName:  v.Name2,
	}
	if err := netlink.LinkAdd(ethPair); err != nil {
		return fmt.Errorf("failed to create veth pair: %w", err)
	}
	return nil
}

func (v VethPair) MoveOneToNamespace(ns netns.NsHandle) error {
	link, err := GetLinkByName(v.Name2)
	if err != nil {
		return err
	}
	if err := netlink.LinkSetNsFd(link, int(ns)); err != nil {
		return fmt.Errorf("failed to move %s to namespace: %w", v.Name2, err)
	}
	return nil
}

func ConfigureIP(name, ipAddress string) error {
	link, err := GetLinkByName(name)
	if err != nil {
		return err
	}
	addr, err := netlink.ParseAddr(ipAddress)
	if err != nil {
		return fmt.Errorf("failed to parse address: %v\n", err)
	}

	if err := netlink.AddrAdd(link, addr); err != nil {
		return fmt.Errorf("failed to assign address to %s: %v\n", link, err)

	}
	return nil
}

func (v VethPair) ConfigureVethPair(ns netns.NsHandle, confIp bool, ip1, ip2 string) error {

	if err := v.Create(); err != nil {
		return fmt.Errorf("failed to create veth pair: %w", err)
	}

	if err := v.MoveOneToNamespace(ns); err != nil {
		return fmt.Errorf("failed to move veth pair to namespace: %w", err)
	}

	if err := SetLinkUp(v.Name1); err != nil {
		return fmt.Errorf("failed to set link up for %s: %w", v.Name1, err)
	}

	if confIp {
		if err := ConfigureIP(v.Name1, ip1); err != nil {
			return fmt.Errorf("failed to configure IP for %s: %w", v.Name1, err)
		}
	}

	if err := netns.Set(ns); err != nil {
		return fmt.Errorf("failed to set new namespace:%w", err)
	}

	if err := SetLinkUp(v.Name2); err != nil {
		return fmt.Errorf("failed to set link up for %s: %w", v.Name2, err)
	}

	if err := ConfigureIP(v.Name2, ip2); err != nil {
		return fmt.Errorf("failed to configure IP for %s: %w", v.Name2, err)
	}

	return nil
}

func AttachToBridge(v string, b Bridge) error {

	vlink, err := GetLinkByName(v)
	if err != nil {
		return err
	}

	blink, err := GetLinkByName(b.Name)
	if err != nil {
		return err
	}

	if err := netlink.LinkSetMaster(vlink, blink); err != nil {
		return fmt.Errorf("failed to set %s to master %s : %w", vlink, blink, err)
	}
	return nil
}
