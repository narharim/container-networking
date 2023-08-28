package main

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/vishvananda/netlink"
	"github.com/vishvananda/netns"
	"golang.org/x/sys/unix"
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

	//Printing the process id that which runs this program
	//See the information about of namespace created in /proc/<pid>/task/<pid>/ns/ in linux machine
	fmt.Printf("%d %d\n", os.Getpid(), unix.Gettid())

	//Getting the existing namespace which is root namespace
	originalNs, err := netns.Get()
	if err != nil {
		fmt.Println("Error getting original namespace:", err)
		return
	}
	defer originalNs.Close()

	//Creating a new named namespace netns0
	newNs, err := netns.NewNamed("netns0")
	if err != nil {
		fmt.Println("Error creating new namespace:", err)
		return
	}
	defer newNs.Close()

	fmt.Println("Created new network namespace:", newNs)

	//falling back to original namespace
	netns.Set(originalNs)
	fmt.Println("Switched back to the original network namespace")

	//Creating pair of virtual ethernet devices which will acts a tunnel to interact with new namespace
	veth := &netlink.Veth{
		LinkAttrs: netlink.LinkAttrs{Name: "veth0"},
		PeerName:  "ceth0",
	}

	if err := netlink.LinkAdd(veth); err != nil {
		fmt.Println("Failed to create veth pair:", err)
		return
	}

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
	addr, _ := netlink.ParseAddr("172.18.0.11/16")
	netlink.AddrAdd(vlink, addr)

	time.Sleep(time.Minute * 2)
	netns.DeleteNamed("netns0")

}
