package main

import (
	"fmt"
	"net"
	"os"

	"github.com/vishvananda/netns"
)

func main() {

	//Checking if the user is root
	if os.Getuid() != 0 {
		fmt.Println("Error: user must be root!")
		return
	}

	originalNs, err := netns.Get()
	if err != nil {
		fmt.Println("Error getting original namespace:", err)
		return
	}

	defer originalNs.Close()

	newNs, err := netns.GetFromName("netns0")
	if err != nil {
		fmt.Println("Error getting new namespace:", err)
		return
	}

	defer newNs.Close()

	netns.Set(newNs)

	newNsIfaces, _ := net.Interfaces()
	fmt.Printf("New Namespace Interfaces: %v\n", newNsIfaces)

	netns.Set(originalNs)

	originalNsIfaces, _ := net.Interfaces()
	fmt.Printf("Root Namespace Interfaces: %v\n", originalNsIfaces)

}
