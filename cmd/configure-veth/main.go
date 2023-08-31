package main

import (
	"fmt"
	"os"
	"runtime"

	user "github.com/narharim/container-networking/pkg/utils"
	networking "github.com/narharim/container-networking/pkg/veth"
	"github.com/vishvananda/netns"
)

const (
	newNamespace = "netns0"
	vethName1    = "veth0"
	vethName2    = "ceth0"
	ipAddress1   = "172.18.0.11/16"
	ipAddress2   = "172.18.0.10/16"
)

func main() {
	//https://stackoverflow.com/questions/27629380/how-to-exit-a-go-program-honoring-deferred-calls
	defer os.Exit(0)

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

	if err := configureVethPair(newNs); err != nil {
		fmt.Println("Error configuring veth pair:", err)
		return
	}

	fmt.Println("Configuration completed successfully!")
}

func configureVethPair(newNs netns.NsHandle) error {
	vethPair := networking.VethPair{
		Name1: vethName1,
		Name2: vethName2,
	}

	if err := vethPair.Create(); err != nil {
		return fmt.Errorf("failed to create veth pair: %w", err)
	}

	if err := vethPair.MoveOneToNamespace(newNs); err != nil {
		return fmt.Errorf("failed to move veth pair to namespace: %w", err)
	}

	if err := networking.SetLinkUp(vethPair.Name1); err != nil {
		return fmt.Errorf("failed to set link up for %s: %w", vethPair.Name1, err)
	}

	if err := networking.ConfigureIP(vethPair.Name1, ipAddress1); err != nil {
		return fmt.Errorf("failed to configure IP for %s: %w", vethPair.Name1, err)
	}

	if err := netns.Set(newNs); err != nil {
		return fmt.Errorf("failed to set new namespace:%w", err)
	}

	if err := networking.SetLinkUp(vethPair.Name2); err != nil {
		return fmt.Errorf("failed to set link up for %s: %w", vethPair.Name2, err)
	}

	if err := networking.ConfigureIP(vethPair.Name2, ipAddress2); err != nil {
		return fmt.Errorf("failed to configure IP for %s: %w", vethPair.Name2, err)
	}

	return nil
}
