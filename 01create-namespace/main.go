package main

import (
	"fmt"
	"os"
	"runtime"

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
}
