package network

import (
	"fmt"

	"github.com/vishvananda/netlink"
)

func GetLinkByName(name string) (netlink.Link, error) {
	link, err := netlink.LinkByName(name)
	if err != nil {
		return nil, fmt.Errorf("failed to get link by name %s: %w", name, err)
	}
	return link, nil
}

func SetLinkUp(name string) error {
	link, err := GetLinkByName(name)
	if err != nil {
		return err
	}
	if err := netlink.LinkSetUp(link); err != nil {
		return fmt.Errorf("failed to set interface up: %v\n", err)
	}
	return nil
}
