package network

import (
	"github.com/vishvananda/netns"
)

func DeleteNameSpace(name string) error {

	err := netns.DeleteNamed(name)
	if err != nil {
		return err
	}

	return nil
}
