package main

import (
	"flag"
	"fmt"
	"os"

	network "github.com/narharim/container-networking/pkg/network"
	user "github.com/narharim/container-networking/pkg/utils"
)

var (
	veth       string
	bridgeName string
)

func main() {
	flag.StringVar(&veth, "veth", "", "virtual ethernet device")
	flag.StringVar(&bridgeName, "bridge", "", "new bridge network")

	flag.Parse()

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n", os.Args[0])
		flag.PrintDefaults()
		return
	}

	if veth == "" || bridgeName == "" {
		flag.Usage()
		return
	}

	user.CheckRootPrivileges()

	bridge := network.Bridge{
		Name: bridgeName,
	}

	if err := bridge.ConfigureBridge(); err != nil {
		fmt.Println("Error configuring bridge:", err)
		return
	}

	if err := network.AttachToBridge(veth, bridge); err != nil {
		fmt.Println("Error configuring bridge:", err)
		return
	}
	fmt.Println("Configuration completed successfully!")
}
