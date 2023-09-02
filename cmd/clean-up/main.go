package main

import (
	"flag"
	"fmt"

	"github.com/vishvananda/netns"
)

func main() {
	flag.Parse()

	args := flag.Args()

	if len(args) == 0 {
		fmt.Println("No arguments provided.")
		return
	}
	for _, ns := range args {
		if err := netns.DeleteNamed(ns); err != nil {
			fmt.Printf("unable to delete the ns: %s\n", err)
		}
	}

}
