package main

import (
	"fmt"
	"net"
	"os"
	"runtime"

	"github.com/vishvananda/netlink"
	"github.com/vishvananda/netns"
)

func main() {

	//Checking if the user is root
	if os.Getuid() != 0 {
		fmt.Println("Error: user must be root!")
		return
	}

	// Lock the OS Thread so we don't accidentally switch namespaces
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	originalNs, err := netns.Get()
	if err != nil {
		fmt.Println("Error getting original namespace:", err)
		return
	}

	defer originalNs.Close()

	//Getting the created namespace
	newNs, err := netns.GetFromName("netns0")
	if err != nil {
		fmt.Println("Error getting new namespace:", err)
		return
	}
	defer newNs.Close()

	//Getting the virtual ethernet device which is created
	clink, err := netlink.LinkByName("ceth0")
	if err != nil {
		fmt.Printf("Failed to find interface %s: %v\n", "ceth0", err)
		os.Exit(1)
	}

	//Sending the device from root namespace to new namespace
	if err = netlink.LinkSetNsFd(clink, int(newNs)); err != nil {
		fmt.Printf("Failed to move interface to namespace: %v\n", err)
		os.Exit(1)
	}

	netns.Set(newNs)

	newNsIfaces, _ := net.Interfaces()
	fmt.Printf("New Namespace Interfaces: %v\n", newNsIfaces)

	netns.Set(originalNs)

	originalNsIfaces, _ := net.Interfaces()
	fmt.Printf("Root Namespace Interfaces: %v\n", originalNsIfaces)
}
