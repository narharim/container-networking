package main

import (
	"fmt"
	"os"
)

func CheckRootPrivileges() {
	if os.Getuid() != 0 {
		fmt.Println("Error: user must be root!")
		os.Exit(1)
	}
}
