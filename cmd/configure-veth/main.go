package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	network "github.com/narharim/container-networking/pkg/network"
	user "github.com/narharim/container-networking/pkg/utils"
	"github.com/vishvananda/netns"
)

var (
	newNamespace   string
	veth           string
	peerVeth       string
	vethIpAddr     string
	peerVethIpAddr string
	assignAddress  bool
)

func main() {

	flag.StringVar(&newNamespace, "namespace", "", "new namespace name")
	flag.StringVar(&veth, "veth", "", "virtual ethernet device")
	flag.StringVar(&peerVeth, "pveth", "", "peer virtual ethernet device")
	flag.StringVar(&vethIpAddr, "veth-ip", "", "virtual ethernet device ip address")
	flag.StringVar(&peerVethIpAddr, "pveth-ip", "", "peer virtual ip address")
	flag.BoolVar(&assignAddress, "addr", true, "assign ip address to veth")

	flag.Parse()

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n", os.Args[0])
		flag.PrintDefaults()
		return
	}

	if assignAddress {
		if newNamespace == "" || veth == "" || peerVeth == "" || vethIpAddr == "" || peerVethIpAddr == "" {
			flag.Usage()
			return
		}
	}

	if !assignAddress {
		if newNamespace == "" || veth == "" || peerVeth == "" || peerVethIpAddr == "" {
			flag.Usage()
			return
		}
	}

	user.CheckRootPrivileges()

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	rootNs, err := netns.Get()
	if err != nil {
		fmt.Println("unable to get root namespace:", err)
		return
	}

	defer rootNs.Close()

	newNs, err := netns.NewNamed(newNamespace)
	if err != nil {
		fmt.Println("unable to create namespace:", err)
		return
	}

	defer newNs.Close()

	fmt.Println("Created new network namespace:", newNs)

	netns.Set(rootNs)

	vethPair := &network.VethPair{
		Name1: veth,
		Name2: peerVeth,
	}

	if err := vethPair.ConfigureVethPair(newNs, assignAddress, vethIpAddr, peerVethIpAddr); err != nil {
		fmt.Println("Error configuring veth pair:", err)
		return
	}

	fmt.Println("Configuration completed successfully!")
}
