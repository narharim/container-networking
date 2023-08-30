package main

import (
	"fmt"
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

	//Getting a new namespace netns0
	newNs, err := netns.GetFromName("netns0")
	if err != nil {
		fmt.Println("Error creating new namespace:", err)
		return
	}
	defer newNs.Close()

	//Getting the veth0
	vlink, err := netlink.LinkByName("veth0")
	if err != nil {
		fmt.Printf("Failed to find interface %s: %v\n", "veth0", err)
		os.Exit(1)
	}

	//Setting the link up
	if err := netlink.LinkSetUp(vlink); err != nil {
		fmt.Printf("Failed to move interface to namespace: %v\n", err)
		os.Exit(1)
	}

	//Adding address to veth0
	vaddr, _ := netlink.ParseAddr("172.18.0.11/16")
	netlink.AddrAdd(vlink, vaddr)

	netns.Set(newNs)

	//Getting the virtual ethernet device which is created
	clink, err := netlink.LinkByName("ceth0")
	if err != nil {
		fmt.Printf("Failed to find interface %s: %v\n", "ceth0", err)
		os.Exit(1)
	}

	//Setting the link up
	if err := netlink.LinkSetUp(clink); err != nil {
		fmt.Printf("Failed to set up interface : %v\n", err)
		os.Exit(1)
	}

	caddr, _ := netlink.ParseAddr("172.18.0.10/16")
	netlink.AddrAdd(clink, caddr)

}
