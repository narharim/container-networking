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

func GetVethLinkByName(name string) (netlink.Link, error) {
	link, err := netlink.LinkByName(name)
	if err != nil {
		return nil, fmt.Errorf("failed to get veth link by name %s: %w", name, err)
	}
	return link, nil
}

func (v VethPair) MoveOneToNamespace(ns netns.NsHandle) error {
	link, err := GetVethLinkByName(v.Name2)
	if err != nil {
		return err
	}
	if err := netlink.LinkSetNsFd(link, int(ns)); err != nil {
		return fmt.Errorf("failed to move %s to namespace: %w", v.Name2, err)
	}
	return nil
}

func SetLinkUp(name string) error {
	link, err := GetVethLinkByName(name)
	if err != nil {
		return err
	}
	if err := netlink.LinkSetUp(link); err != nil {
		return fmt.Errorf("failed to set interface up: %v\n", err)
	}
	return nil
}

func ConfigureIP(name, ipAddress string) error {
	link, err := GetVethLinkByName(name)
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
