package main

import (
	"fmt"
	"net"
	"os"

	"github.com/vishvananda/netlink"
)

func main() {

	//Checking if the user is root
	if os.Getuid() != 0 {
		fmt.Println("Error: user must be root!")
		return
	}

	veth := &netlink.Veth{
		LinkAttrs: netlink.LinkAttrs{Name: "veth0"},
		PeerName:  "ceth0",
	}

	if err := netlink.LinkAdd(veth); err != nil {
		fmt.Println("Failed to create veth pair:", err)
		return
	}

	originalNsIfaces, _ := net.Interfaces()
	fmt.Printf("Root Namespace Interfaces: %v\n", originalNsIfaces)

}
