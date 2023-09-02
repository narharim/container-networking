package network

import (
	"fmt"

	"github.com/vishvananda/netlink"
)

type Bridge struct {
	Name string
}

func (b Bridge) Create() error {
	bridge := &netlink.Bridge{
		LinkAttrs: netlink.LinkAttrs{Name: b.Name},
	}
	if err := netlink.LinkAdd(bridge); err != nil {
		return fmt.Errorf("failed to create bridge: %w", err)
	}
	return nil
}

func (b Bridge) ConfigureBridge() error {

	if err := b.Create(); err != nil {
		return fmt.Errorf("failed to create bridge: %w", err)
	}

	if err := SetLinkUp(b.Name); err != nil {
		return fmt.Errorf("failed to set link up for %s: %w", b.Name, err)
	}

	return nil
}
